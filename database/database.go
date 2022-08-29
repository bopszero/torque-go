package database

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

var DbConnectionMap map[string]Database

func Init() error {
	dbConnMap := viper.GetStringMapString(config.KeyDbConnectionMap)
	if _, hasDefault := dbConnMap[AliasDefault]; !hasDefault {
		return fmt.Errorf("`default` connection is missing in `DB_CONNECTION_MAP`")
	}

	dbContainer := make(map[string]Database)
	for alias, connStr := range dbConnMap {
		db, err := initDbFromConnectionString(connStr)
		if err != nil {
			return err
		}
		dbContainer[alias] = db
	}
	if _, hasMaster := dbConnMap[AliasMainMaster]; !hasMaster {
		dbContainer[AliasMainMaster] = dbContainer[AliasDefault]
	}
	if _, hasSlave := dbConnMap[AliasMainSlave]; !hasSlave {
		dbContainer[AliasMainSlave] = dbContainer[AliasDefault]
	}
	DbConnectionMap = dbContainer
	config.CmdRegisterRootDefer(Close)
	return nil
}

func initDbFromConnectionString(dsn string) (_ Database, err error) {
	if tz := config.TimeZone; tz != "" {
		tzStr := url.Values{"loc": []string{tz}}.Encode()
		dsn = dsn + "&" + tzStr
	}
	var (
		dsnParts  = strings.Split(dsn, "://")
		dsnScheme = dsnParts[0]
		dsnURI    = dsnParts[1]
		dialector gorm.Dialector
	)
	switch dsnScheme {
	case "mysql":
		dialector = mysql.Open(dsnURI)
		break
	default:
		err = fmt.Errorf("Database DSN scheme `%s` hasn't been supported yet", dsnScheme)
		return
	}

	dbConf := gorm.Config{
		Logger: gormlogger.Discard,
	}
	if config.Debug {
		dbConf.Logger = gormlogger.Default
	}
	db, err := gorm.Open(dialector, &dbConf)
	if err != nil {
		return
	}
	if config.Debug {
		db = db.Debug()
	}

	poolOpts, err := GetDbPoolOptions()
	if err != nil {
		return
	}
	dbPool, err := db.DB()
	if err != nil {
		return
	}
	if poolOpts.PoolSize > 0 {
		dbPool.SetMaxIdleConns(poolOpts.PoolSize)
	}
	if poolOpts.MaxLifetime.Duration > 0 {
		dbPool.SetConnMaxLifetime(poolOpts.MaxLifetime.Duration)
	}
	if poolOpts.MaxSize > 0 {
		dbPool.SetMaxOpenConns(poolOpts.MaxSize)
	}

	return Database{db}, nil
}

func Close() {
	closeErrorMap := make(map[string]error)
	for alias, db := range DbConnectionMap {
		sqlDB, err := db.DB.DB()
		if err != nil {
			continue
		}
		if err := sqlDB.Close(); err != nil {
			closeErrorMap[alias] = err
		}
	}

	if len(closeErrorMap) == 0 {
		return
	}

	logger := comlogging.GetLogger()
	for alias, err := range closeErrorMap {
		logger.
			WithError(err).
			WithField("alias", alias).
			Error("close database connection failed")
	}
}

func GetDb(alias string) (db Database, err error) {
	db, ok := DbConnectionMap[alias]
	if !ok {
		err = fmt.Errorf("cannot get database connection `%v`", alias)
	}
	return
}

func GetDbF(alias string) Database {
	db, err := GetDb(alias)
	comutils.PanicOnError(err)
	return db
}

func GetDbDefault() Database {
	return GetDbF(AliasDefault)
}

func GetDbMaster() Database {
	return GetDbF(AliasMainMaster)
}

func GetDbSlave() Database {
	return GetDbF(AliasMainSlave)
}

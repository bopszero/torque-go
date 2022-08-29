package thirdpartymod

import (
	"net/url"
	"strings"

	pool "github.com/jolestar/go-commons-pool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
	"gopkg.in/gomail.v2"
)

var (
	vEmailSystemClientSingleton = comtypes.NewSingleton(func() interface{} {
		client, err := NewEmailServiceClient(viper.GetString(config.KeyEmailDSN))
		comutils.PanicOnError(err)
		return client
	})
	vEmailConnPoolConfig pool.ObjectPoolConfig
)

func init() {
	conf := pool.NewDefaultPoolConfig()
	conf.MaxTotal = 10
	conf.MaxIdle = 5
	conf.MaxWaitMillis = 3000
	conf.MinEvictableIdleTimeMillis = 15000
	conf.TimeBetweenEvictionRunsMillis = 15000
	vEmailConnPoolConfig = *conf
}

type EmailServiceClient struct {
	dialer   *gomail.Dialer
	connPool *pool.ObjectPool

	fromName  string
	fromEmail string
}

func NewEmailServiceClient(dsn string) (client EmailServiceClient, err error) {
	emailURL, err := url.Parse(viper.GetString(config.KeyEmailDSN))
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	var (
		urlData   = emailURL.Query()
		fromName  = urlData.Get("from_name")
		fromEmail = urlData.Get("from_address")
	)
	var (
		port        = comutils.ParseIntF(emailURL.Port())
		password, _ = emailURL.User.Password()
		dialer      = gomail.NewDialer(
			emailURL.Hostname(), port,
			emailURL.User.Username(), password)

		connPoolFactory = pool.NewPooledObjectFactory(
			func() (interface{}, error) { return dialer.Dial() },
			func(obj *pool.PooledObject) error { return obj.Object.(gomail.SendCloser).Close() },
			nil, nil, nil,
		)
	)

	client = EmailServiceClient{
		dialer:    dialer,
		connPool:  pool.NewObjectPool(connPoolFactory, &vEmailConnPoolConfig),
		fromName:  fromName,
		fromEmail: fromEmail,
	}
	config.CmdRegisterRootDefer(func() { client.connPool.Close() })

	return client, nil
}

func GetEmailServiceSystemClient() EmailServiceClient {
	return vEmailSystemClientSingleton.Get().(EmailServiceClient)
}

func (this EmailServiceClient) borrowConnection() (conn gomail.SendCloser, err error) {
	connObj, err := this.connPool.BorrowObject()
	if err != nil {
		err = utils.WrapError(err)
		return
	}
	conn, ok := connObj.(gomail.SendCloser)
	if !ok {
		err = utils.WrapError(constants.ErrorInvalidData)
		return
	}
	return conn, nil
}

func (this EmailServiceClient) Send(ctx comcontext.Context, msg *gomail.Message) error {
	if len(msg.GetHeader(constants.EmailHeaderFrom)) == 0 {
		msg.SetAddressHeader(constants.EmailHeaderFrom, this.fromEmail, this.fromName)
	}
	connection, err := this.borrowConnection()
	if err != nil {
		return err
	}
	defer this.connPool.ReturnObject(connection)

	var (
		subject   = strings.Join(msg.GetHeader(constants.EmailHeaderSubject), ", ")
		toAddress = strings.Join(msg.GetHeader(constants.EmailHeaderTo), ", ")
	)
	logEntry := comlogging.GetLogger().
		WithContext(ctx).
		WithFields(logrus.Fields{
			"subject": subject,
			"to":      toAddress,
			"from":    strings.Join(msg.GetHeader(constants.EmailHeaderFrom), ", "),
		})
	if err := gomail.Send(connection, msg); err != nil {
		logEntry.
			WithError(err).
			Errorf("send email failed | to=%v,subject=%s", toAddress, subject)
		return err
	}

	logEntry.Infof("send email ok | to=%v,subject=%s", toAddress, subject)
	return nil
}

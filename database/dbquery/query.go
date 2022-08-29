package dbquery

import (
	"fmt"
	"strings"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func JoinExpr(values ...interface{}) string {
	strValues := make([]string, len(values))
	for i, value := range values {
		strValues[i] = comutils.Stringify(value)
	}
	return strings.Join(strValues, ", ")
}

func Or(conditions ...string) string {
	for i := range conditions {
		conditions[i] = fmt.Sprintf("(%s)", conditions[i])
	}
	return strings.Join(conditions, " OR ")
}

func In(colName string, values interface{}) (string, interface{}) {
	return fmt.Sprintf("%s IN (?)", colName), values
}

func NotIn(colName string, values interface{}) (string, interface{}) {
	return fmt.Sprintf("%s NOT IN (?)", colName), values
}

func Compare(colName string, operator string, value interface{}) (string, interface{}) {
	return fmt.Sprintf("%s %s ?", colName, operator), value
}

func Gt(colName string, value interface{}) (string, interface{}) {
	return Compare(colName, ">", value)
}

func Gte(colName string, value interface{}) (string, interface{}) {
	return Compare(colName, ">=", value)
}

func Equal(colName string, value interface{}) (string, interface{}) {
	return Compare(colName, "=", value)
}

func NotEqual(colName string, value interface{}) (string, interface{}) {
	return Compare(colName, "!=", value)
}

func Lte(colName string, value interface{}) (string, interface{}) {
	return Compare(colName, "<=", value)
}

func Lt(colName string, value interface{}) (string, interface{}) {
	return Compare(colName, "<", value)
}

func Between(colName string, fromValue interface{}, toValue interface{}) (string, interface{}, interface{}) {
	return fmt.Sprintf("%s BETWEEN ? AND ?", colName), fromValue, toValue
}

func SelectForUpdate(db *gorm.DB) *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}

func OrderAsc(colName string) string {
	return fmt.Sprintf("%s ASC", colName)
}

func OrderDesc(colName string) string {
	return fmt.Sprintf("%s DESC", colName)
}

func MaxAlias(colName string, alias string) string {
	return fmt.Sprintf("MAX(%s) AS %s", colName, alias)
}

func Max(colName string) string {
	return MaxAlias(colName, "max_"+colName)
}

func SumAlias(colName string, alias string) string {
	return fmt.Sprintf("SUM(%s) AS %s", colName, alias)
}

func Sum(colName string) string {
	return SumAlias(colName, "sum_"+colName)
}

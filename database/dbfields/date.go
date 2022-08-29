package dbfields

import (
	"database/sql/driver"
	"time"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

const (
	DateFieldFormat     = "2006-01-02"
	DateFieldEmptyValue = "1000-01-01"
)

type DateField struct {
	valueStr string
}

func NewDateFieldFromString(date string) (field DateField, err error) {
	if _, err = comutils.TimeParse(DateFieldFormat, date); err != nil {
		err = utils.WrapError(err)
		return
	}
	field.valueStr = date
	return
}

func NewDateFieldFromStringF(date string) DateField {
	field, err := NewDateFieldFromString(date)
	comutils.PanicOnError(err)
	return field
}

func NewDateFieldFromTime(dateTime time.Time) (field DateField) {
	field.valueStr = dateTime.Format(DateFieldFormat)
	return
}

func (this DateField) String() string {
	return this.valueStr
}

func (this DateField) Time() (time.Time, error) {
	return time.Parse(DateFieldFormat, this.valueStr)
}

func (this DateField) TimeF() time.Time {
	timeVal, err := this.Time()
	comutils.PanicOnError(err)
	return timeVal
}

func (this DateField) Value() (driver.Value, error) {
	if this.valueStr == "" {
		return DateFieldEmptyValue, nil
	}

	if _, err := time.Parse(DateFieldFormat, this.valueStr); err != nil {
		return nil, err
	}

	return this.valueStr, nil
}

func (this *DateField) Scan(input interface{}) error {
	dateTime := input.(time.Time)
	this.valueStr = dateTime.Format(DateFieldFormat)
	return nil
}

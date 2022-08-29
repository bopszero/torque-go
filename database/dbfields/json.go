package dbfields

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type JsonField meta.O

func (this JsonField) Value() (driver.Value, error) {
	bytes, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	valueStr := string(bytes)
	if valueStr == "{}" {
		valueStr = ""
	}

	return valueStr, err
}

func (this *JsonField) Scan(input interface{}) error {
	inputBytes := input.([]byte)
	if len(inputBytes) == 0 {
		inputBytes = []byte("{}")
	}

	return json.Unmarshal(inputBytes, this)
}

type JsonListField []interface{}

func (this JsonListField) Value() (driver.Value, error) {
	bytes, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	valueStr := string(bytes)
	if valueStr == "[]" {
		valueStr = ""
	}

	return valueStr, err
}

func (this *JsonListField) Scan(input interface{}) error {
	inputBytes := input.([]byte)
	if len(inputBytes) == 0 {
		inputBytes = []byte("[]")
	}

	return json.Unmarshal(inputBytes, this)
}

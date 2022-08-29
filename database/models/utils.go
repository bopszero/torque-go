package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

func NewBool(value bool) sql.NullBool {
	return sql.NullBool{
		Bool:  value,
		Valid: true,
	}
}

func NewString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func NewStringNull() sql.NullString {
	return sql.NullString{
		String: "-",
		Valid:  false,
	}
}

func NewBytes(data []byte) sql.NullString {
	return sql.NullString{
		String: string(data),
		Valid:  true,
	}
}

func NewInt32(value int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: value,
		Valid: true,
	}
}

func NewInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: value,
		Valid: true,
	}
}

func NewUInt64(value uint64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(value),
		Valid: true,
	}
}

func toJson(obj interface{}) (driver.Value, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	valueStr := string(bytes)
	if valueStr == "{}" {
		valueStr = ""
	}

	return valueStr, err
}

func fromJson(input interface{}, obj interface{}) error {
	inputBytes := input.([]byte)
	if len(inputBytes) == 0 {
		inputBytes = []byte("{}")
	}

	return json.Unmarshal(inputBytes, obj)
}

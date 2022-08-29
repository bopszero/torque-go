package dbfields

import (
	"database/sql"
	"strings"

	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

type BoolField struct {
	sql.NullBool
}

func (this BoolField) MarshalJSON() ([]byte, error) {
	var jsonText string
	if !this.Valid {
		jsonText = "null"
	} else if this.Bool {
		jsonText = "true"
	} else {
		jsonText = "false"
	}
	return []byte(jsonText), nil
}

func (this *BoolField) UnmarshalJSON(data []byte) error {
	text := strings.ToLower(string(data))
	switch text {
	case "true":
		this.Valid = true
		this.Bool = true
		break
	case "false":
		this.Valid = true
		this.Bool = false
		break
	case "null":
		this.Valid = false
		break
	default:
		return utils.IssueErrorf("BoolField cannot `UnmarshalJSON` for value `%s`", text)
	}
	return nil
}

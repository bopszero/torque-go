package utils

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"time"
)

func TimeParseDate(dateStr string) (time.Time, error) {
	return time.Parse(constants.DateFormatISO, dateStr)
}

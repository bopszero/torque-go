package responses

import (
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

type ApiResponse struct {
	Code    meta.ErrorCode `json:"code"`
	Message string         `json:"message,omitempty"`
	Errors  []string       `json:"errors,omitempty"`
	Data    interface{}    `json:"data"`
}

package constants

import "gitlab.com/snap-clickstaff/torque-go/config"

const (
	ContentTypeJSON   = "application/json"
	ContentTypeBinary = "application/octet-stream"
	ContentTypeCSV    = "text/csv"
)

var RequestIpHeaders = []string{
	config.HttpHeaderCfConnectingIP,
	config.HttpHeaderXForwardedFor,
}

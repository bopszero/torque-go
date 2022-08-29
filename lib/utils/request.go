package utils

import (
	"net"
	"net/http"
	"strings"

	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
)

func GetRequestUserIP(request *http.Request) string {
	if request == nil {
		return ""
	}

	for _, header := range constants.RequestIpHeaders {
		headerValue := request.Header.Get(header)
		if headerValue != "" {
			return strings.Split(headerValue, ",")[0]
		}
	}

	return request.RemoteAddr
}

func IsRequestTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	switch errT := err.(type) {
	case *comlogging.SentryError:
		err = errT.GetError()
		break
	}

	if networkErr, ok := err.(net.Error); ok {
		return networkErr.Timeout()
	}

	if strings.Contains(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
		return true
	}

	return false
}

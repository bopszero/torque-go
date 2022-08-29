package middleware

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/torque-go/api/apiutils"
	"gitlab.com/snap-clickstaff/torque-go/config"
)

type (
	LogRequestOptions struct {
		RequestLimit  int
		LogResponse   bool
		ResponseLimit int
		HeaderPrefix  string
	}

	logResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

var (
	LogRequestDefaultOptions = LogRequestOptions{
		RequestLimit:  500,
		LogResponse:   true,
		ResponseLimit: 500,
	}
	LogRequestDefaultMiddleware = NewRequestLogger(LogRequestDefaultOptions)
)

func NewRequestLogger(options LogRequestOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			startTime := time.Now()

			requestBody := readContextRequestBody(c)

			responseBodyBuffer := new(bytes.Buffer)
			if options.LogResponse {
				response := c.Response()
				multiWriter := io.MultiWriter(response.Writer, responseBodyBuffer)
				response.Writer = &logResponseWriter{Writer: multiWriter, ResponseWriter: response.Writer}
			} else {
				responseBodyBuffer.WriteString("--ignored--")
			}

			defer func() {
				log(c, options, startTime, requestBody, responseBodyBuffer, err)
			}()
			err = next(c)

			return err
		}
	}
}

func log(
	c echo.Context, options LogRequestOptions, startTime time.Time,
	requestBody string, repsonseBodyBuffer *bytes.Buffer, err error,
) {
	var (
		ctx      = apiutils.EchoWrapContext(c)
		headers  map[string]string
		request  = c.Request()
		response = c.Response()

		endTime      = time.Now()
		latency      = endTime.Sub(startTime)
		responseCode = response.Header().Get(config.HttpHeaderApiResponseCode)
	)
	if options.HeaderPrefix != "" {
		headers = make(map[string]string)
		lowerPrefix := strings.ToLower(options.HeaderPrefix)
		for name, values := range request.Header {
			lowerName := strings.ToLower(name)
			if strings.HasPrefix(lowerName, lowerPrefix) {
				headers[lowerName] = strings.Join(values, "\n")
			}
		}
	}
	if requestBody == "" && len(request.URL.Query()) > 0 {
		requestBody = "query:" + request.URL.Query().Encode()
	}
	var (
		requestBodyLog  = truncateStringDisplay(requestBody, options.RequestLimit)
		responseBodyLog = truncateStringDisplay(
			string(repsonseBodyBuffer.Bytes()),
			options.ResponseLimit,
		)
		uid, _   = apiutils.GetContextUID(ctx)
		logEntry = comlogging.GetLogger().
				GenHttpInboundEnty(
				uint64(uid), request.Method, latency, request.URL.EscapedPath(), c.RealIP(),
				requestBodyLog, headers, response.Status, responseCode, responseBodyLog,
			).
			WithContext(ctx)
	)

	if err != nil {
		logEntry.
			WithError(err).
			Infof("http request error | err=%s", err.Error())
	} else {
		logEntry.Info("http request")
	}
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *logResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *logResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

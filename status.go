package hproblem

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusContinue                      statusError = 100 // RFC 9110, 15.2.1
	StatusSwitchingProtocols            statusError = 101 // RFC 9110, 15.2.2
	StatusProcessing                    statusError = 102 // RFC 2518, 10.1
	StatusEarlyHints                    statusError = 103 // RFC 8297
	StatusOK                            statusError = 200 // RFC 9110, 15.3.1
	StatusCreated                       statusError = 201 // RFC 9110, 15.3.2
	StatusAccepted                      statusError = 202 // RFC 9110, 15.3.3
	StatusNonAuthoritativeInfo          statusError = 203 // RFC 9110, 15.3.4
	StatusNoContent                     statusError = 204 // RFC 9110, 15.3.5
	StatusResetContent                  statusError = 205 // RFC 9110, 15.3.6
	StatusPartialContent                statusError = 206 // RFC 9110, 15.3.7
	StatusMultiStatus                   statusError = 207 // RFC 4918, 11.1
	StatusAlreadyReported               statusError = 208 // RFC 5842, 7.1
	StatusIMUsed                        statusError = 226 // RFC 3229, 10.4.1
	StatusMultipleChoices               statusError = 300 // RFC 9110, 15.4.1
	StatusMovedPermanently              statusError = 301 // RFC 9110, 15.4.2
	StatusFound                         statusError = 302 // RFC 9110, 15.4.3
	StatusSeeOther                      statusError = 303 // RFC 9110, 15.4.4
	StatusNotModified                   statusError = 304 // RFC 9110, 15.4.5
	StatusUseProxy                      statusError = 305 // RFC 9110, 15.4.6
	StatusTemporaryRedirect             statusError = 307 // RFC 9110, 15.4.8
	StatusPermanentRedirect             statusError = 308 // RFC 9110, 15.4.9
	StatusBadRequest                    statusError = 400 // RFC 9110, 15.5.1
	StatusUnauthorized                  statusError = 401 // RFC 9110, 15.5.2
	StatusPaymentRequired               statusError = 402 // RFC 9110, 15.5.3
	StatusForbidden                     statusError = 403 // RFC 9110, 15.5.4
	StatusNotFound                      statusError = 404 // RFC 9110, 15.5.5
	StatusMethodNotAllowed              statusError = 405 // RFC 9110, 15.5.6
	StatusNotAcceptable                 statusError = 406 // RFC 9110, 15.5.7
	StatusProxyAuthRequired             statusError = 407 // RFC 9110, 15.5.8
	StatusRequestTimeout                statusError = 408 // RFC 9110, 15.5.9
	StatusConflict                      statusError = 409 // RFC 9110, 15.5.10
	StatusGone                          statusError = 410 // RFC 9110, 15.5.11
	StatusLengthRequired                statusError = 411 // RFC 9110, 15.5.12
	StatusPreconditionFailed            statusError = 412 // RFC 9110, 15.5.13
	StatusRequestEntityTooLarge         statusError = 413 // RFC 9110, 15.5.14
	StatusRequestURITooLong             statusError = 414 // RFC 9110, 15.5.15
	StatusUnsupportedMediaType          statusError = 415 // RFC 9110, 15.5.16
	StatusRequestedRangeNotSatisfiable  statusError = 416 // RFC 9110, 15.5.17
	StatusExpectationFailed             statusError = 417 // RFC 9110, 15.5.18
	StatusTeapot                        statusError = 418 // RFC 9110, 15.5.19 (Unused)
	StatusMisdirectedRequest            statusError = 421 // RFC 9110, 15.5.20
	StatusUnprocessableEntity           statusError = 422 // RFC 9110, 15.5.21
	StatusLocked                        statusError = 423 // RFC 4918, 11.3
	StatusFailedDependency              statusError = 424 // RFC 4918, 11.4
	StatusTooEarly                      statusError = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired               statusError = 426 // RFC 9110, 15.5.22
	StatusPreconditionRequired          statusError = 428 // RFC 6585, 3
	StatusTooManyRequests               statusError = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge   statusError = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons    statusError = 451 // RFC 7725, 3
	StatusInternalServerError           statusError = 500 // RFC 9110, 15.6.1
	StatusNotImplemented                statusError = 501 // RFC 9110, 15.6.2
	StatusBadGateway                    statusError = 502 // RFC 9110, 15.6.3
	StatusServiceUnavailable            statusError = 503 // RFC 9110, 15.6.4
	StatusGatewayTimeout                statusError = 504 // RFC 9110, 15.6.5
	StatusHTTPVersionNotSupported       statusError = 505 // RFC 9110, 15.6.6
	StatusVariantAlsoNegotiates         statusError = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           statusError = 507 // RFC 4918, 11.5
	StatusLoopDetected                  statusError = 508 // RFC 5842, 7.2
	StatusNotExtended                   statusError = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired statusError = 511 // RFC 6585, 6
)

// statusJSONS prebuilds the json responses for the standard status codes.
var statusJSONS map[int][]byte = func() map[int][]byte {
	jsons := make(map[int][]byte)
	for code := 100; code < 600; code++ {
		text := http.StatusText(code)
		if text != "" {
			json := fmt.Sprintf(`{"detail":"%s","status":%d,"title":"%s"}`, text, code, text)
			jsons[code] = []byte(json)
		}
	}
	return jsons
}()

type statusError int

func (err statusError) Error() string {
	return http.StatusText(err.StatusCode())
}

func (err statusError) StatusCode() int {
	return int(err)
}

func (err statusError) MarshalJSON() ([]byte, error) {
	return statusJSONS[int(err)], nil
}

func (err statusError) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	title := err.Error()
	return e.Encode(struct {
		Detail  string   `xml:"detail"`
		Status  int      `xml:"status"`
		Title   string   `xml:"title"`
		XMLName xml.Name `xml:"urn:ietf:rfc:7807 problem"`
	}{
		Detail: title,
		Status: err.StatusCode(),
		Title:  title,
	})
}

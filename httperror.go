// Package hproblem provides a standard interface for handling
// API error responses in web applications.
// It implements RFC 7807 (Problem Details for HTTP APIs)
// which specifies a way to carry machine-readable
// details of errors in a HTTP response to avoid the need
// to define new error response formats for HTTP APIs.
package hproblem

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
)

type httpError struct {
	error
	statusCode int
}

func (err *httpError) StatusCode() int {
	return err.statusCode
}

func (err *httpError) Unwrap() error {
	return err.error
}

func (err *httpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Detail string `json:"detail"`
		Status int    `json:"status"`
		Title  string `json:"title"`
	}{
		Detail: err.Error(),
		Status: err.statusCode,
		Title:  http.StatusText(err.statusCode),
	})
}

func (err *httpError) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(struct {
		Detail  string   `xml:"detail"`
		Status  int      `xml:"status"`
		Title   string   `xml:"title"`
		XMLName xml.Name `xml:"urn:ietf:rfc:7807 problem"`
	}{
		Detail: err.Error(),
		Status: err.statusCode,
		Title:  http.StatusText(err.statusCode),
	})
}

func serveJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func serveXML(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/problem+xml; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_, _ = io.WriteString(w, xml.Header)
	_ = xml.NewEncoder(w).Encode(v)
}

// Wrap associates an error with a status code.
func Wrap(statusCode int, err error) error {
	return &httpError{err, statusCode}
}

// Errorf is a shorthand for Wrap(fmt.Errorf(...), statusCode).
func Errorf(statusCode int, format string, a ...interface{}) error {
	return Wrap(statusCode, fmt.Errorf(format, a...))
}

// StatusCode reports the HTTP status code associated with err
// if it implements the StatusCode() int method,
// 504 Gateway Timeout if it implements Timeout() bool,
// 503 Service Unavailable if it implements Temporary() bool,
// 500 Internal Server Error otherwise, or 200 OK if err is nil.
// StatusCode will unwrap err to find the most precise status code.
func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	for ; err != nil; err = errors.Unwrap(err) {
		if sc, ok := err.(interface{ StatusCode() int }); ok { //nolint
			return sc.StatusCode()
		} else if to, ok := err.(interface{ Timeout() bool }); ok && to.Timeout() { //nolint
			return http.StatusGatewayTimeout
		} else if te, ok := err.(interface{ Temporary() bool }); ok && te.Temporary() { //nolint
			return http.StatusServiceUnavailable
		}
	}

	return http.StatusInternalServerError
}

// ServeError replies to the request by rendering err.
// If err implements http.Handler, its ServeHTTP method is called.
// Otherwise, err is rendered as JSON, XML or plain text depending on the
// request's Accept header.
// If err is nil, it will be rendered as StatusOK.
func ServeError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = StatusOK
	}

	if h, ok := err.(http.Handler); ok { //nolint
		h.ServeHTTP(w, r)
		return
	}

	for _, accept := range r.Header["Accept"] {
		if ok, _ := path.Match("*/*json", accept); ok {
			serveJSON(w, err, StatusCode(err))
			return
		} else if ok, _ := path.Match("*/*xml", accept); ok {
			serveXML(w, err, StatusCode(err))
			return
		}
	}

	http.Error(w, err.Error(), StatusCode(err))
}

// MethodNotAllowed replies to the request with StatusMethodNotAllowed.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	ServeError(w, r, StatusMethodNotAllowed)
}

// NotFound replies to the request with StatusNotFound.
func NotFound(w http.ResponseWriter, r *http.Request) {
	ServeError(w, r, StatusNotFound)
}

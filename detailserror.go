package hproblem

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"unicode"
)

// DetailsError implements the RFC 7807 model.
// See: https://datatracker.ietf.org/doc/html/rfc7807
//
// Additional fields can be added by embedding it inside another struct.
//
//	type TraceDetailsError struct {
//	    *hproblem.DetailsError
//	    TraceID string `json:"trace_id" xml:"trace_id"`
//	}
//
//	hproblem.ServeError(w, r, TraceDetailsError{})
type DetailsError struct {
	// A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty" xml:"detail,omitempty"`

	// A URI reference that identifies the specific occurrence of the problem.
	// It may or may not yield further information if dereferenced.
	Instance string `json:"instance,omitempty" xml:"instance,omitempty"`

	// The HTTP status code ([RFC7231], Section 6)
	// generated by the origin server for this occurrence of the problem.
	Status int `json:"status,omitempty" xml:"status,omitempty"`

	// A short, human-readable summary of the problem
	// type. It SHOULD NOT change from occurrence to occurrence of the
	// problem, except for purposes of localization (e.g., using
	// proactive content negotiation; see [RFC7231], Section 3.4).
	Title string `json:"title,omitempty" xml:"title,omitempty"`

	// A URI reference [RFC3986] that identifies the
	// problem type. This specification encourages that, when
	// dereferenced, it provide human-readable documentation for the
	// problem type (e.g., using HTML [W3C.REC-html5-20141028]). When
	// this member is not present, its value is assumed to be
	// "about:blank".
	Type string `json:"type,omitempty" xml:"type,omitempty"`

	// XMLName is needed to marshal to XML.
	XMLName xml.Name `json:"-" xml:"urn:ietf:rfc:7807 problem"`

	wrappedError error
}

// Error implements the error interface and returns the Detail field.
func (details *DetailsError) Error() string { return details.Detail }

// StatusCode implements the interface used by StatusCode and returns the Status field.
func (details *DetailsError) StatusCode() int { return details.Status }

// Unwrap implements the interface used by errors.Unwrap() and returns the wrapped error.
func (details *DetailsError) Unwrap() error { return details.wrappedError }

// NewDetailsError returns a new DetailsError with the
// Detail, Status and Title fields set according to err.
func NewDetailsError(err error) *DetailsError {
	var detail string
	if err != nil {
		detail = err.Error()
	}

	statusCode := StatusCode(err)

	return &DetailsError{
		Detail:       detail,
		Status:       statusCode,
		Title:        http.StatusText(statusCode),
		wrappedError: err,
	}
}

var ErrInvalidEncoding = errors.New("hproblem: invalid details error encoding")

// Unmarshal parses a JSON or XML encoded details error.
// Returns ErrInvalidEncoding if the encoding is invalid.
func (details *DetailsError) Unmarshal(data []byte) error {
	data = bytes.TrimLeftFunc(data, unicode.IsSpace)
	if len(data) == 0 {
		return ErrInvalidEncoding
	}

	switch data[0] {
	case '{':
		return json.Unmarshal(data, details)
	case '<':
		return xml.Unmarshal(data, details)
	default:
		return ErrInvalidEncoding
	}
}

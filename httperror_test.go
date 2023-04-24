package hproblem

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type customError struct{}

func (err customError) StatusCode() int {
	return 123
}

func (err customError) Error() string {
	return "hello"
}

func (err customError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(err.StatusCode())
	fmt.Fprintf(w, "hello")
}

func TestServeError(t *testing.T) {
	t.Run("XML", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", "application/xml")
		NotFound(w, r)
		if w.Result().StatusCode != http.StatusNotFound {
			t.Fatal()
		} else if w.Body.String() != xml.Header+`<problem xmlns="urn:ietf:rfc:7807"><detail>Not Found</detail><status>404</status><title>Not Found</title></problem>` {
			t.Fatal(w.Body.String())
		}
	})

	t.Run("JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", "application/json")
		MethodNotAllowed(w, r)
		if w.Result().StatusCode != http.StatusMethodNotAllowed {
			t.Fatal()
		} else if w.Body.String() != `{"detail":"Method Not Allowed","status":405,"title":"Method Not Allowed"}`+"\n" {
			t.Fatal()
		}
	})

	t.Run("Text", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ServeError(w, r, StatusBadRequest)
		if w.Result().StatusCode != http.StatusBadRequest {
			t.Fatal()
		} else if w.Body.String() != "Bad Request\n" {
			t.Fatal()
		}
	})

	t.Run("nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ServeError(w, r, nil)
		if w.Result().StatusCode != http.StatusOK {
			t.Fatal()
		}
		if w.Body.String() != http.StatusText(http.StatusOK)+"\n" {
			t.Fatal()
		}
	})

	t.Run("custom", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ServeError(w, r, customError{})
		if w.Result().StatusCode != 123 {
			t.Fatal(w.Result().StatusCode)
		}
		if w.Body.String() != "hello" {
			t.Fatal()
		}
	})
}

func TestUnwrap(t *testing.T) {
	err := errors.New("")
	err2 := Wrap(444, err)
	if !errors.Is(errors.Unwrap(err2), err) {
		t.Fatal()
	}
}

type mockTemporaryError struct{}

func (e mockTemporaryError) Error() string   { return "" }
func (e mockTemporaryError) Temporary() bool { return true }

func TestStatusCode(t *testing.T) {
	for _, testCase := range []struct {
		Err      error
		Expected int
	}{
		{nil, http.StatusOK},
		{context.DeadlineExceeded, http.StatusGatewayTimeout},
		{fmt.Errorf("error: %w", mockTemporaryError{}), http.StatusServiceUnavailable},
		{Wrap(http.StatusBadRequest, errors.New("bla")), http.StatusBadRequest},
		{errors.New("bla"), http.StatusInternalServerError},
	} {
		if StatusCode(testCase.Err) != testCase.Expected {
			t.Error(testCase.Err, StatusCode(testCase.Err), testCase.Expected)
		}
	}
}

func TestErrorMarshalJSON(t *testing.T) {
	err := Wrap(http.StatusTeapot, errors.New("guru meditation"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "application/json")
	ServeError(w, r, err)
	if w.Result().StatusCode != StatusCode(err) {
		t.Fatal()
	}
	if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/problem+json") {
		t.Fatal()
	}
}

func TestErrorMarshalXML(t *testing.T) {
	err := Errorf(http.StatusTeapot, "guru meditation")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "application/xml")
	ServeError(w, r, err)
	if w.Result().StatusCode != StatusCode(err) {
		t.Fatal()
	}
	if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/problem+xml") {
		t.Fatal()
	}
}

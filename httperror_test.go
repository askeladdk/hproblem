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
		{Wrap(errors.New("bla"), http.StatusBadRequest), http.StatusBadRequest},
		{errors.New("bla"), http.StatusInternalServerError},
	} {
		if StatusCode(testCase.Err) != testCase.Expected {
			t.Error(testCase.Err, StatusCode(testCase.Err), testCase.Expected)
		}
	}
}

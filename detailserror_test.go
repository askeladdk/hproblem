package hproblem

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	_ error = &DetailsError{}
	_ error = &struct{ DetailsError }{}
	_ error = &struct{ *DetailsError }{}
)

type testEmbeddedDetails struct {
	*DetailsError
	ID string `json:"id" xml:"id"`
}

func TestProblemDetailsError(t *testing.T) {
	detail := testEmbeddedDetails{
		DetailsError: NewDetailsError(Wrap(http.StatusBadRequest, errors.New("error"))),
		ID:           "myid",
	}

	if detail.Error() != "error" {
		t.Fatal()
	}

	if detail.StatusCode() != http.StatusBadRequest {
		t.Fatal()
	}

	t.Run("JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", "application/json")
		ServeError(w, r, &detail)
		b := w.Body.String()
		if b != `{"detail":"error","status":400,"title":"Bad Request","id":"myid"}`+"\n" {
			t.Fatal(b)
		}
	})

	t.Run("XML", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Accept", "text/xml")
		ServeError(w, r, &detail)
		b := w.Body.String()
		if b != xml.Header+`<problem xmlns="urn:ietf:rfc:7807"><detail>error</detail><status>400</status><title>Bad Request</title><id>myid</id></problem>` {
			t.Fatal(b)
		}
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		x := []byte(`{"detail":"error","status":400,"title":"Bad Request"}`)
		var dt DetailsError
		if err := dt.Unmarshal(x); err != nil {
			t.Fatal(err)
		}
		if dt.Status != 400 {
			t.Fatal()
		}
		if dt.Title != "Bad Request" {
			t.Fatal()
		}
	})

	t.Run("xml", func(t *testing.T) {
		x := []byte(xml.Header + `<problem xmlns="urn:ietf:rfc:7807"><detail>error</detail><status>400</status><title>Bad Request</title></problem>`)
		var dt DetailsError
		if err := dt.Unmarshal(x); err != nil {
			t.Fatal(err)
		}
		if dt.Status != 400 {
			t.Fatal()
		}
		if dt.Title != "Bad Request" {
			t.Fatal()
		}
	})

	t.Run("error1", func(t *testing.T) {
		var dt DetailsError
		if err := dt.Unmarshal([]byte{}); !errors.Is(err, ErrInvalidEncoding) {
			t.Fatal()
		}
	})

	t.Run("error2", func(t *testing.T) {
		var dt DetailsError
		if err := dt.Unmarshal([]byte(`  []`)); !errors.Is(err, ErrInvalidEncoding) {
			t.Fatal()
		}
	})
}

func TestDetailsUnwrap(t *testing.T) {
	err := errors.New("")
	err2 := NewDetailsError(err)
	if !errors.Is(errors.Unwrap(err2), err) {
		t.Fatal()
	}
}

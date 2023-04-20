package hproblem_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/askeladdk/hproblem"
)

func ExampleServeError_json() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "application/json")

	hproblem.ServeError(w, r, hproblem.StatusBadRequest)

	fmt.Println(w.Result().Status)
	_, _ = io.Copy(os.Stdout, w.Body)

	// Output:
	// 400 Bad Request
	// {"detail":"Bad Request","status":400,"title":"Bad Request"}
}

func ExampleServeError_xml() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept", "text/xml")

	hproblem.ServeError(w, r, hproblem.StatusBadRequest)

	fmt.Println(w.Result().Status)
	_, _ = io.Copy(os.Stdout, w.Body)

	// Output:
	// 400 Bad Request
	// <?xml version="1.0" encoding="UTF-8"?>
	// <problem xmlns="urn:ietf:rfc:7807"><detail>Bad Request</detail><status>400</status><title>Bad Request</title></problem>
}

func ExampleServeError_text() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	hproblem.ServeError(w, r, hproblem.StatusBadRequest)

	fmt.Println(w.Result().Status)
	_, _ = io.Copy(os.Stdout, w.Body)

	// Output:
	// 400 Bad Request
	// Bad Request
}

func ExampleStatusCode() {
	fmt.Println(hproblem.StatusCode(nil))
	fmt.Println(hproblem.StatusCode(io.EOF))
	fmt.Println(hproblem.StatusCode(hproblem.Wrap(io.EOF, http.StatusBadRequest)))
	// Output:
	// 200
	// 500
	// 400
}

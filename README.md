# hproblem - Error responses for HTTP APIs in Go

[![GoDoc](https://godoc.org/github.com/askeladdk/hproblem?status.png)](https://godoc.org/github.com/askeladdk/hproblem)
[![Go Report Card](https://goreportcard.com/badge/github.com/askeladdk/hproblem)](https://goreportcard.com/report/github.com/askeladdk/hproblem)
[![Coverage Status](https://coveralls.io/repos/github/askeladdk/hproblem/badge.svg?branch=master)](https://coveralls.io/github/askeladdk/hproblem?branch=master)

## Overview

Package hproblem provides a standard interface for handling API error responses in web applications. It implements [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) (Problem Details for HTTP APIs) which specifies a way to carry machine-readable details of errors in a HTTP response to avoid the need to define new error response formats for HTTP APIs.

## Install

```
go get -u github.com/askeladdk/hproblem
```

## Quickstart

The two basic functions are `Wrap` and `ServeError`. Wrap associates an error with a status code. `ServeError` replies to a request by marshaling the error to JSON, XML or plain text depending on the request's Accept header. Use it instead of `http.Error`. `ServeError` also accepts errors that implement the `http.Handler` interface, in which case the error is in charge of marshaling itself.

```go
func endpoint(w http.ResponseWriter, r *http.Request) {
    hproblem.ServeError(w, r, hproblem.Wrap(http.StatusBadRequest, io.EOF))
}
```

Use `Errorf` as a shorthand for `Wrap(statusCode, fmt.Errorf(...))`.

```go
err = hproblem.Errorf(http.StatusBadRequest, "package: error: %w", err)
```

Use the `DetailsError` type directly if you need more control.

```go
var err error = &hproblem.DetailsError{
    Detail: "This is not the Jedi that you are looking for",
    Instance: "/jedi/obi-wan",
    Status: http.StatusNotFound,
    Title: "Jedi Mind Trick",
}
```

Embed `DetailsError` inside another type to add custom fields and use `NewDetailsError` to initialize it.

```go
type TraceError struct {
    *hproblem.DetailsError
    TraceID string `json:"trace_id" xml:"trace_id"`
}

var err error = &TraceError{
    DetailsError: hproblem.NewDetailsError(hproblem.Wrap(http.StatusBadRequest, io.EOF)),
    TraceID: "42",
}
```

Use the predefined `Status*` errors to serve HTTP status codes without needing to wrap. This is convenient in cases where it is not needed to attach extra information to an error. Every status code present in the `http` package has an equivalent error in `hproblem`. Handlers `MethodNotFound` and `NotFound` are also provided.

```go
hproblem.ServeError(w, r, hproblem.StatusForbidden)
```

Read the rest of the [documentation on pkg.go.dev](https://pkg.go.dev/github.com/askeladdk/hproblem). It's easy-peasy!

## License

Package hproblem is released under the terms of the ISC license.

go-e5e
======

[![PkgGoDev](https://pkg.go.dev/badge/anexia-it/go-e5e)](https://pkg.go.dev/anexia-it/go-e5e)
[![Build Status](https://travis-ci.org/anexia-it/go-e5e.svg?branch=master)](https://travis-ci.org/anexia-it/go-e5e)
[![codecov](https://codecov.io/gh/anexia-it/go-e5e/branch/master/graph/badge.svg)](https://codecov.io/gh/anexia-it/go-e5e)
[![Go Report Card](https://goreportcard.com/badge/github.com/anexia-it/go-e5e)](https://goreportcard.com/report/github.com/anexia-it/go-e5e)

go-e5e is a support library to help Go developers develop Anexia e5e functions.

# Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/anexia-it/go-e5e
```

# Getting started

```go
package main

import (
	"runtime"

	e5e "github.com/anexia-it/go-e5e"
)

type entrypoints struct{}

func (f *entrypoints) MyEntrypoint(event e5e.Event, context e5e.Context) (*e5e.Return, error) {
	return &e5e.Return{
		Status: 200,
		ResponseHeaders: map[string]string{
			"x-custom-response-header": "This is a custom response header",
		},
		Data: map[string]interface{}{
			"version": runtime.Version(),
		},
	}, nil
}

func main() {
	if err := e5e.Start(&entrypoints{}); err != nil {
		panic(err)
	}
}
```

# List of developers

* Andreas Stocker <AStocker@anexia-it.com>, Lead Developer
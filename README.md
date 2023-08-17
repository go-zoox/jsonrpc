# JSONRPC Client/Server
> According to the [JSON-RPC 2.0 specification](http://www.jsonrpc.org/specification),
> JSON-RPC is a lightweight remote procedure call (RPC) protocol.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/jsonrpc)](https://pkg.go.dev/github.com/go-zoox/jsonrpc)
[![Build Status](https://github.com/go-zoox/jsonrpc/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/jsonrpc/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/jsonrpc)](https://goreportcard.com/report/github.com/go-zoox/jsonrpc)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/jsonrpc/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/jsonrpc?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/jsonrpc.svg)](https://github.com/go-zoox/jsonrpc/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/jsonrpc.svg?label=Release)](https://github.com/go-zoox/jsonrpc/releases)

## Installation
To install the package, run:
```bash
go get github.com/go-zoox/jsonrpc
```

## Getting Started

```go
// server
package main

import (
	"context"

	"github.com/go-zoox/jsonrpc"
	"github.com/go-zoox/jsonrpc/server"
	"github.com/go-zoox/logger"
)

func main() {
	s := server.New()

	s.Register("echo", func(ctx context.Context, params jsonrpc.Params) (jsonrpc.Result, error) {
		logger.Info("params: %s", params)

		return jsonrpc.Result{
			"name": params.Get("name"),
			"age":  18,
		}, nil
	})

	s.Run()
}
```

```go
// client
package main

import (
	"context"

	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/jsonrpc"
	"github.com/go-zoox/jsonrpc/client"
	"github.com/go-zoox/logger"
)

func main() {
	c := client.New("http://localhost:8080")

	r, err := c.Call(context.Background(), "echo", jsonrpc.Params{
		"name": "zero",
	})
	if err != nil {
		logger.Errorf("failed to call: %s", err)
		return
	}

	logger.Info("result: %d", cast.ToInt64(r.Get("age")))
}
```

## License
GoZoox is released under the [MIT License](./LICENSE).

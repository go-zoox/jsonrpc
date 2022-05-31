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
func main() {
  s := jsonrpc.NewServer(8080)

	s.Register("echo", func(params gjson.Result) Result {
		logger.Info("params: %s", params.String())

		return Result{
			"name": params.Get("name").String(),
			"age":  18,
		}
	})

	s.Start()
}
```

```go
// client
func main() {
  c := jsonrpc.NewClient("http://localhost:8080/")

	r, err := c.Invoke("echo", map[string]string{
		"name": "zero",
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("result: %d", r.Get("age").Int())
}
```

## License
GoZoox is released under the [MIT License](./LICENSE).
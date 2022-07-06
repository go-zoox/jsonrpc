package jsonrpc

import (
	"github.com/go-zoox/jsonrpc/client"
	"github.com/go-zoox/jsonrpc/server"
)

// NewClient creates a JSONRPC client.
func NewClient(uri string, cfg ...*client.Config) *client.Client {
	return client.New(uri, cfg...)
}

// NewServer creates a JSONRPC server.
func NewServer(cfg ...*server.Config) *server.Server {
	return server.New(cfg...)
}

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

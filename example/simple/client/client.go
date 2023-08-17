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

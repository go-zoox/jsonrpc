package server

import (
	"testing"

	"github.com/go-zoox/logger"
	"github.com/tidwall/gjson"
)

func TestServer(t *testing.T) {
	s := New()

	s.Register("echo", func(ctx *Context, params gjson.Result) Result {
		logger.Info("params: %s", params.String())

		return Result{
			"name": params.Get("name").String(),
			"age":  18,
		}
	})

	s.Start(":8080")
}

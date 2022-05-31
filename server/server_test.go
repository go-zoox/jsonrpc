package server

import (
	"testing"

	"github.com/go-zoox/logger"
	"github.com/tidwall/gjson"
)

func TestServer(t *testing.T) {
	s := New(8080)

	s.Register("echo", func(params gjson.Result) Result {
		logger.Info("params: %s", params.String())

		return Result{
			"name": params.Get("name").String(),
			"age":  18,
		}
	})

	s.Start()
}

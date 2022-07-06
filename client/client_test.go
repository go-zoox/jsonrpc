package client

import (
	"testing"

	"github.com/go-zoox/logger"
)

func TestClient(t *testing.T) {
	c := New("http://localhost:8080")

	r, err := c.Invoke("echo", map[string]string{
		"name": "zero",
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.Info("result: %d", r.Get("age").Int())

	// if r.Get("message").String() != "Hello World" {
	// 	t.Error("Expected 'Hello World', got", r.Get("message").String())
	// }
}

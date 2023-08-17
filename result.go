package jsonrpc

import (
	"github.com/go-zoox/core-utils/object"
	"github.com/go-zoox/tag"
	"github.com/go-zoox/tag/datasource"
)

// Result is a map of result.
type Result map[string]any

// Bind binds the jsonrpc params into the given struct.
func (p Result) Bind(obj any) error {
	return tag.New("json", datasource.NewMapDataSource(p)).Decode(obj)
}

// Get returns the value of the given key.
func (p Result) Get(key string) any {
	return object.Get(p, key)
}

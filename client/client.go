package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-zoox/fetch"
	"github.com/go-zoox/jsonrpc"
	"github.com/go-zoox/uuid"
)

// Client is a JSON-RPC client.
type Client interface {
	Call(ctx context.Context, method string, params jsonrpc.Params) (jsonrpc.Result, error)
}

type client struct {
	server string
	path   string
}

// New creates a new JSON-RPC client.
func New(server string, path ...string) Client {
	pathX := jsonrpc.JSONRPCDefaultPath
	if len(path) > 0 {
		if path[0] != "" {
			pathX = path[0]
		}
	}

	return &client{
		server: server,
		path:   pathX,
	}
}

// Call calls a JSON-RPC method.
func (c *client) Call(ctx context.Context, method string, params jsonrpc.Params) (jsonrpc.Result, error) {
	response, err := fetch.Post(c.server+c.path, &fetch.Config{
		Body: map[string]any{
			"jsonrpc": "2.0",
			"method":  method,
			"params":  params,
			"id":      uuid.V4(),
		},
	})
	if err != nil {
		return nil, err
	}

	var res jsonrpc.Response
	err = json.Unmarshal(response.Body, &res)
	if err != nil {
		return nil, err
	}

	if res.JSONRPC != jsonrpc.JSONRPCVersion {
		return nil, fmt.Errorf("invalid jsonrpc version: %s", res.JSONRPC)
	}

	if res.Error != nil && res.Error.Code != 0 {
		return nil, fmt.Errorf("[%d] %s", res.Error.Code, res.Error.Message)
	}

	return res.Result, nil
}

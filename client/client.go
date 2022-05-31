package client

import (
	"fmt"

	"github.com/go-zoox/fetch"
	"github.com/tidwall/gjson"
)

// Client is a JSONRPC client.
type Client struct {
	ServerURI string
	Count     int
	//
	Config *Config
}

// Config is a JSONRPC client configuration.
type Config struct {
	Headers map[string]string
	Query   map[string]string
}

// New creates a JSONRPC client.
func New(uri string, cfg ...*Config) *Client {
	var config *Config
	if len(cfg) > 0 && cfg[0] != nil {
		config = cfg[0]
	} else {
		config = &Config{}
	}

	return &Client{
		ServerURI: uri,
		Config:    config,
	}
}

// Invoke invokes a JSONRPC method.
func (c *Client) Invoke(method string, params any) (gjson.Result, error) {
	c.Count += 1
	response, err := fetch.Post(c.ServerURI, &fetch.Config{
		Headers: c.Config.Headers,
		Query:   c.Config.Query,
		Body: map[string]any{
			"jsonrpc": "2.0",
			"method":  method,
			"params":  params,
			"id":      c.Count,
		},
	})
	if err != nil {
		return gjson.Result{}, err
	}

	if response.Status != 200 {
		err := response.Get("error")
		code := err.Get("code").Int()
		message := err.Get("message").String()
		return gjson.Result{}, fmt.Errorf("[%d] %s", code, message)
	}

	return response.Get("result"), nil
}

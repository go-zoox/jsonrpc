package server

import (
	"context"
	"encoding/json"

	"github.com/go-zoox/jsonrpc"
	"github.com/go-zoox/logger"
)

// HandlerFunc is a handler function.
type HandlerFunc func(ctx context.Context, params jsonrpc.Params) (jsonrpc.Result, error)

// Server is a JSON-RPC server.
type Server interface {
	Register(method string, handler HandlerFunc)
	Invoke(ctx context.Context, body []byte) ([]byte, error)
	//
	Path() string
	Run(addr ...string) error
}

type server struct {
	path string
	//
	methods map[string]HandlerFunc
}

// New creates a new JSON-RPC server.
func New(path ...string) Server {
	pathX := jsonrpc.JSONRPCDefaultPath
	if len(path) > 0 {
		if path[0] != "" {
			pathX = path[0]
		}
	}

	return &server{
		path:    pathX,
		methods: make(map[string]HandlerFunc),
	}
}

func (s *server) Register(method string, handler HandlerFunc) {
	s.methods[method] = handler
}

func (s *server) Invoke(ctx context.Context, body []byte) ([]byte, error) {
	response := &jsonrpc.Response{
		JSONRPC: jsonrpc.JSONRPCVersion,
	}

	var request jsonrpc.Request
	err := json.Unmarshal(body, &request)
	if err != nil {
		logger.Info("jsonrpc: invalid request: %s(%s)", err, string(body))

		response.Error = &jsonrpc.Error{
			Code:    -32700,
			Message: "Parse error",
		}

		return json.Marshal(response)
	}

	if request.JSONRPC != "2.0" {
		response.Error = &jsonrpc.Error{
			Code:    -32600,
			Message: "Invalid Request (invlid JSON-RPC version)",
		}
		return json.Marshal(response)
	}

	if request.Method == "" {
		response.Error = &jsonrpc.Error{
			Code:    -32600,
			Message: "Invalid Request (method is required)",
		}
		return json.Marshal(response)
	}

	if request.ID == "" {
		response.Error = &jsonrpc.Error{
			Code:    -32600,
			Message: "Invalid Request (id is required)",
		}
		return json.Marshal(response)
	}

	// fmt.Println("request.ID", request.ID)
	// fmt.Println("request.Method", request.Method)
	// fmt.Println("request.Params", request.Params)

	logger.Infof("[jsonrpc][id: %s][method: %s]", request.ID, request.Method)

	response.ID = request.ID

	handler, ok := s.methods[request.Method]
	if !ok {
		response.Error = &jsonrpc.Error{
			Code:    -32601,
			Message: "Method not found",
		}

		return json.Marshal(response)
	}

	result, err := handler(ctx, request.Params)
	if err != nil {
		response.Error = &jsonrpc.Error{
			Code:    -32603,
			Message: err.Error(),
		}

		return json.Marshal(response)
	}

	response.Result = result
	return json.Marshal(response)
}

// Path returns the path of the server.
func (s *server) Path() string {
	return s.path
}

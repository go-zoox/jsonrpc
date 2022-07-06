package server

import (
	"io"
	"net/http"

	"github.com/go-zoox/encoding/json"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/safe"
	"github.com/tidwall/gjson"
)

// Server is a JSONRPC server.
type Server struct {
	Path string
	//
	methods map[string]Method
}

// Config is a JSONRPC server configuration.
type Config struct {
	Path string
}

// Method is a JSONRPC method.
type Method func(ctx *Context, params Params) Result

// Context is a JSONRPC context.
type Context struct {
	Request *http.Request
}

// Params is a JSONRPC params.
type Params = gjson.Result

// Result is a JSONRPC result.
type Result map[string]any

// New creates a new JSONRPC server.
func New(cfg ...*Config) *Server {
	path := ""
	if len(cfg) > 0 && cfg[0] != nil {
		if cfg[0].Path != "" {
			path = cfg[0].Path
		}
	}

	return &Server{
		Path:    path,
		methods: map[string]Method{},
	}
}

// Start starts the JSONRPC server.
func (s *Server) Start(addr string) {
	logger.Info("Starting JSONRPC server at: %s", addr)
	if err := http.ListenAndServe(addr, s); err != nil {
		logger.Error("Failed to start JSONRPC server: %s", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Path != "" && r.URL.Path != s.Path {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := gjson.Parse(string(bytes))
	if !body.Get("jsonrpc").Exists() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Encode(map[string]any{
			"jsonrpc": "2.0",
			"error": map[string]any{
				"code":    -32600,
				"message": "Invalid Request",
			},
			"id": nil,
		})

		w.Write(response)
		return
	}

	jsonrpc := body.Get("jsonrpc").String()
	if jsonrpc != "2.0" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Encode(map[string]any{
			"jsonrpc": "2.0",
			"error": map[string]any{
				"code":    -32600,
				"message": "Invlid JSONRPC Version",
			},
			"id": nil,
		})

		w.Write(response)
		return
	}

	method := body.Get("method").String()
	params := body.Get("params")
	id := body.Get("id").Int()

	if !s.has(method) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Encode(map[string]any{
			"jsonrpc": "2.0",
			"error": map[string]any{
				"code":    -32601,
				"message": "Method not found",
			},
			"id": nil,
		})
		w.Write(response)
		return
	}

	var result map[string]any
	context := &Context{
		Request: r,
	}

	err = safe.Do(func() error {
		result = s.invoke(context, method, params)
		return nil
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Encode(map[string]any{
			"jsonrpc": "2.0",
			"error": map[string]any{
				"code":    -32000,
				"message": err.Error(),
			},
			"id": nil,
		})

		w.Write(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Encode(map[string]any{
		"jsonrpc": "2.0",
		"result":  result,
		"id":      id,
	})

	w.Write(response)
}

func (s *Server) has(method string) bool {
	_, ok := s.methods[method]
	return ok
}

func (s *Server) invoke(ctx *Context, method string, params gjson.Result) Result {
	return s.methods[method](ctx, params)
}

// Register registers a method.
func (s *Server) Register(method string, fn Method) {
	s.methods[method] = fn
}

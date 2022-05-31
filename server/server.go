package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-zoox/encoding/json"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/safe"
	"github.com/tidwall/gjson"
)

// Server is a JSONRPC server.
type Server struct {
	Port int
	Path string
	//
	methods map[string]Method
}

// Config is a JSONRPC server configuration.
type Config struct {
	Path string
}

// Method is a JSONRPC method.
type Method func(params gjson.Result) Result

// Result is a JSONRPC result.
type Result map[string]any

// New creates a new JSONRPC server.
func New(port int, cfg ...*Config) *Server {
	path := "/"
	if len(cfg) > 0 && cfg[0] != nil {
		if cfg[0].Path != "" {
			path = cfg[0].Path
		}
	}

	return &Server{
		Port:    port,
		Path:    path,
		methods: map[string]Method{},
	}
}

// Start starts the JSONRPC server.
func (s *Server) Start() {
	http.HandleFunc(s.Path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

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
		err = safe.Do(func() error {
			result = s.invoke(method, params)
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
	})

	logger.Info("Starting JSONRPC server at port: %d", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

func (s *Server) has(method string) bool {
	_, ok := s.methods[method]
	return ok
}

func (s *Server) invoke(method string, params gjson.Result) Result {
	return s.methods[method](params)
}

// Register registers a method.
func (s *Server) Register(method string, fn Method) {
	s.methods[method] = fn
}

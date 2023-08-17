package server

import (
	"io"
	"net/http"

	"github.com/go-zoox/logger"
)

// ServeHTTP implements the http.Handler interface.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.path != "" && r.URL.Path != s.path {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	response, err := s.Invoke(r.Context(), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Run runs the JSONRPC server.
func (s *server) Run(addr ...string) error {
	addrX := ":8080"
	if len(addr) > 0 {
		if addr[0] != "" {
			addrX = addr[0]
		}
	}

	logger.Info("Run JSON-RPC server at: %s", addrX)
	return http.ListenAndServe(addrX, s)
}

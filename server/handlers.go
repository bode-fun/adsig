package server

import (
	"fmt"
	"net/http"
)

func (s *Server) registerHandlers() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

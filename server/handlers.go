package server

import (
	"fmt"
	"net/http"

	"git.bode.fun/adsig/signature"
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerHandlers() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/signature", signature.RegisterHandlers())
	})
}

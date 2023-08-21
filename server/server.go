package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Type Definition
// ------------------------------------------------------------------------

var _ http.Handler = (*Server)(nil)

type Server struct {
	router chi.Router
}

func New() *Server {
	router := chi.NewRouter()

	srv := &Server{
		router: router,
	}

	srv.useMiddleware()
	srv.registerHandlers()

	return srv
}

// Public Methods
// ------------------------------------------------------------------------

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

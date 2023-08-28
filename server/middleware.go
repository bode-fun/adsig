package server

import "github.com/go-chi/chi/v5/middleware"

func (s *Server) useMiddleware() {
	// FIXME: When needed, add the following CORS middleware: "github.com/rs/cors"
	// FIXME: Add security headers via "github.com/unrolled/secure"
	// FIXME: Maybe add UBER's rate-limit middleware to avoid flooding the AD
	// FIXME: When adding a ui, add CSRF protection via "github.com/justinas/nosurf"
	// FIXME: Add Gzip via github.com/klauspost/compress

	// Give each request a unique ID
	s.router.Use(middleware.RequestID)

	// Get the client's ip-address, even when proxied
	s.router.Use(middleware.RealIP)

	// Remove multiple slashes from the requested resource path
	s.router.Use(middleware.CleanPath)

	// Remove any trailing slash from the requested resource path
	s.router.Use(middleware.StripSlashes)

	// Paths are clean and ready to be logged :)

	// Log every incoming request
	// INFO: The Log middleware needs Recover middleware to be registered
	s.router.Use(middleware.Logger)

	// A panic should not quit the program
	s.router.Use(middleware.Recoverer)
}

package signature

import (
	"fmt"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(db *badger.DB) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/{account}", handleGetForAccount(db))
	}
}

func handleGetForAccount(db *badger.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "account")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, email)
	}
}

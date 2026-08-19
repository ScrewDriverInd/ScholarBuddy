//go:build production

package httpapi

import "github.com/go-chi/chi/v5"

// Production builds have no docs routes or embedded OpenAPI spec.
func registerDocs(_ chi.Router) {}

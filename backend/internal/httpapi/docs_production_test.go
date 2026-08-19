//go:build production

package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scholarbuddy/backend/internal/auth"
	"github.com/scholarbuddy/backend/internal/config"
)

func TestProductionDocsAreUnavailable(t *testing.T) {
	cfg := config.Config{AuthJWTSecret: "this-is-a-test-secret-that-is-long-enough", AdminUserID: "admin", AdminPassword: "password", AdminKey: "key"}
	h := New(nil, auth.NewService(nil, "", cfg.AuthJWTSecret), cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("docs route should be absent in production, got %d", w.Code)
	}
}

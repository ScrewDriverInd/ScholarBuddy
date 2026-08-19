//go:build !production

package httpapi

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ScrewDriverInd/ScholarBuddy/internal/auth"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/config"
)

func TestDevelopmentDocsAreAvailable(t *testing.T) {
	cfg := config.Config{AuthJWTSecret: "this-is-a-test-secret-that-is-long-enough", AdminUserID: "admin", AdminPassword: "password", AdminKey: "key"}
	h := New(nil, auth.NewService(nil, "", cfg.AuthJWTSecret), cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil))
	if w.Code != http.StatusOK || !bytes.Contains(w.Body.Bytes(), []byte("openapi: 3.1.0")) {
		t.Fatalf("development spec unavailable: status=%d body=%s", w.Code, w.Body.String())
	}
}

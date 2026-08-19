package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scholarbuddy/backend/internal/auth"
	"github.com/scholarbuddy/backend/internal/config"
)

func TestAdminLogin(t *testing.T) {
	cfg := config.Config{AuthJWTSecret: "this-is-a-test-secret-that-is-long-enough", AdminUserID: "admin", AdminPassword: "password", AdminKey: "key"}
	h := New(nil, auth.NewService(nil, "", cfg.AuthJWTSecret), cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))

	for _, tc := range []struct {
		body string
		want int
	}{
		{`{"user_id":"admin","password":"password","admin_key":"key"}`, http.StatusOK},
		{`{"user_id":"admin","password":"wrong","admin_key":"key"}`, http.StatusUnauthorized},
	} {
		r := httptest.NewRequest(http.MethodPost, "/abbujaan/login", bytes.NewBufferString(tc.body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("status = %d, want %d; body=%s", w.Code, tc.want, w.Body.String())
		}
		if tc.want == http.StatusOK {
			var response struct {
				Data struct {
					Token string `json:"access_token"`
				} `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Data.Token == "" {
				t.Fatalf("expected access token: %v", err)
			}
		}
	}
}

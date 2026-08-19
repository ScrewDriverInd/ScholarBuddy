package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ScrewDriverInd/ScholarBuddy/internal/auth"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/config"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/opportunity"
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

func TestOpportunityFilter(t *testing.T) {
	tests := []struct {
		value string
		want  opportunity.Type
		valid bool
	}{
		{"", "", true},
		{"all", "", true},
		{"research", opportunity.ResearchExtra, true},
		{"scholarship", opportunity.Scholarship, true},
		{"research_extra", "", false},
		{"unknown", "", false},
	}
	for _, tt := range tests {
		got, valid := opportunityFilter(tt.value)
		if got != tt.want || valid != tt.valid {
			t.Errorf("opportunityFilter(%q) = (%q, %t), want (%q, %t)", tt.value, got, valid, tt.want, tt.valid)
		}
	}
}

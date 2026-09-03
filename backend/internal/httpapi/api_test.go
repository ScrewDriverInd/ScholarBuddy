package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ScrewDriverInd/ScholarBuddy/internal/auth"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/config"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/opportunity"
)

func TestLegacyAdminLoginRouteIsUnavailable(t *testing.T) {
	cfg := config.Config{AuthJWTSecret: "this-is-a-test-secret-that-is-long-enough"}
	h := New(nil, auth.NewService(nil, "", cfg.AuthJWTSecret), cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/abbujaan/login", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("legacy admin login route status = %d, want 404", w.Code)
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
		{"research", opportunity.Research, true},
		{"extras", opportunity.Extras, true},
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

func TestWelcomeRoute(t *testing.T) {
	cfg := config.Config{AuthJWTSecret: "this-is-a-test-secret-that-is-long-enough"}
	h := New(nil, auth.NewService(nil, "", cfg.AuthJWTSecret), cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Welcome to ScholarBuddy") {
		t.Fatalf("welcome response = %d %s", w.Code, w.Body.String())
	}
}

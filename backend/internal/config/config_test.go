package config

import "testing"

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgresql://localhost/scholarbuddy")
	t.Setenv("SUPABASE_URL", "https://example.supabase.co")
	t.Setenv("SUPABASE_JWKS_URL", "")
	t.Setenv("AUTH_JWT_SECRET", "this-secret-is-longer-than-thirty-two-characters")
	t.Setenv("ADMIN_USER_ID", "admin")
	t.Setenv("ADMIN_PASSWORD", "password")
	t.Setenv("ADMIN_KEY", "admin-key")
}

func TestLoadDerivesJWKSURLAndParsesCORS(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://scholarbuddy.example")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got, want := cfg.SupabaseJWKSURL, "https://example.supabase.co/auth/v1/.well-known/jwks.json"; got != want {
		t.Errorf("SupabaseJWKSURL = %q, want %q", got, want)
	}
	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Errorf("CORSAllowedOrigins = %#v, want two origins", cfg.CORSAllowedOrigins)
	}
}

func TestLoadRejectsShortInternalJWTSecret(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("AUTH_JWT_SECRET", "too-short")
	if _, err := Load(); err == nil {
		t.Fatal("Load() succeeded with a short AUTH_JWT_SECRET")
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() succeeded without DATABASE_URL")
	}
}

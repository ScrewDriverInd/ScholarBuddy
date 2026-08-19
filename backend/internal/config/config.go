package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr           string
	DatabaseURL        string
	CORSAllowedOrigins []string
	SupabaseURL        string
	SupabaseJWKSURL    string
	AuthJWTSecret      string
	AdminUserID        string
	AdminPassword      string
	AdminKey           string
}

func Load() (Config, error) {
	c := Config{
		HTTPAddr:    env("HTTP_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"), SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseJWKSURL: os.Getenv("SUPABASE_JWKS_URL"), AuthJWTSecret: os.Getenv("AUTH_JWT_SECRET"),
		AdminUserID: os.Getenv("ADMIN_USER_ID"), AdminPassword: os.Getenv("ADMIN_PASSWORD"), AdminKey: os.Getenv("ADMIN_KEY"),
	}
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		c.CORSAllowedOrigins = strings.Split(origins, ",")
	}
	if c.SupabaseJWKSURL == "" && c.SupabaseURL != "" {
		c.SupabaseJWKSURL = strings.TrimRight(c.SupabaseURL, "/") + "/auth/v1/.well-known/jwks.json"
	}
	for name, value := range map[string]string{"DATABASE_URL": c.DatabaseURL, "SUPABASE_JWKS_URL": c.SupabaseJWKSURL, "AUTH_JWT_SECRET": c.AuthJWTSecret, "ADMIN_USER_ID": c.AdminUserID, "ADMIN_PASSWORD": c.AdminPassword, "ADMIN_KEY": c.AdminKey} {
		if value == "" {
			return Config{}, fmt.Errorf("%s is required", name)
		}
	}
	if len(c.AuthJWTSecret) < 32 {
		return Config{}, fmt.Errorf("AUTH_JWT_SECRET must be at least 32 characters")
	}
	return c, nil
}
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

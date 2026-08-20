package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"golang.org/x/crypto/bcrypt"
)

type Principal struct {
	ID    uuid.UUID
	Email string
	Roles []string
}

func (p Principal) HasRole(role string) bool {
	for _, value := range p.Roles {
		if value == role {
			return true
		}
	}
	return false
}

type contextKey struct{}

func FromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(contextKey{}).(Principal)
	return p, ok
}

type Service struct {
	db                      *pgxpool.Pool
	jwksURL, internalSecret string
	httpClient              *http.Client
}

func NewService(db *pgxpool.Pool, jwksURL, internalSecret string) *Service {
	return &Service{db: db, jwksURL: jwksURL, internalSecret: internalSecret, httpClient: &http.Client{Timeout: 5 * time.Second}}
}
func (s *Service) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := s.supabasePrincipal(r.Context(), bearer(r))
		if err != nil {
			unauthorized(w)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, p)))
	})
}
func (s *Service) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := s.internalPrincipal(bearer(r))
		if err != nil || !p.HasRole("ROLE_ADMIN") {
			unauthorized(w)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, p)))
	})
}
func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
}
func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"valid authentication is required"}}`))
}

func (s *Service) supabasePrincipal(ctx context.Context, raw string) (Principal, error) {
	if raw == "" {
		return Principal{}, errors.New("missing token")
	}
	claims := jwt.MapClaims{}
	token, _, err := new(jwt.Parser).ParseUnverified(raw, claims)
	if err != nil {
		return Principal{}, err
	}
	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		return Principal{}, errors.New("token kid missing")
	}
	set, err := jwk.Fetch(ctx, s.jwksURL, jwk.WithHTTPClient(s.httpClient))
	if err != nil {
		return Principal{}, err
	}
	key, ok := set.LookupKeyID(kid)
	if !ok {
		return Principal{}, errors.New("signing key not found")
	}
	var rawKey any
	if err := jwk.Export(key, &rawKey); err != nil {
		return Principal{}, err
	}
	parsed, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) { return rawKey, nil }, jwt.WithValidMethods([]string{"RS256", "ES256", "EdDSA"}))
	if err != nil || !parsed.Valid {
		return Principal{}, errors.New("invalid token")
	}
	sub, _ := claims["sub"].(string)
	id, err := uuid.Parse(sub)
	if err != nil {
		return Principal{}, errors.New("invalid subject")
	}
	email, _ := claims["email"].(string)
	name, _ := claims["user_metadata"].(map[string]any)
	fullName := ""
	if name != nil {
		fullName, _ = name["full_name"].(string)
	}
	// Google OAuth accounts are always created with ROLE_USER. Existing database
	// roles are preserved, including ROLE_ADMIN for administrator accounts.
	_, err = s.db.Exec(ctx, "INSERT INTO users (id,email,full_name,roles) VALUES ($1,$2,$3,ARRAY['ROLE_USER'::user_role]) ON CONFLICT (id) DO UPDATE SET email=EXCLUDED.email,full_name=EXCLUDED.full_name,updated_at=now()", id, email, fullName)
	if err != nil {
		return Principal{}, err
	}
	return Principal{ID: id, Email: email, Roles: []string{"ROLE_USER"}}, nil
}
func (s *Service) LoginAdmin(ctx context.Context, username, password string) (string, error) {
	var id uuid.UUID
	var passwordHash string
	if err := s.db.QueryRow(ctx, "SELECT id,password_hash FROM users WHERE username=$1 AND 'ROLE_ADMIN'::user_role = ANY(roles)", username).Scan(&id, &passwordHash); err != nil {
		return "", err
	}
	if passwordHash == "" || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"role": "ROLE_ADMIN", "sub": id.String(), "exp": time.Now().Add(8 * time.Hour).Unix(), "iat": time.Now().Unix()}).SignedString([]byte(s.internalSecret))
}
func (s *Service) internalPrincipal(raw string) (Principal, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.internalSecret), nil
	})
	if err != nil || !t.Valid {
		return Principal{}, errors.New("invalid token")
	}
	role, _ := claims["role"].(string)
	return Principal{Roles: []string{role}}, nil
}

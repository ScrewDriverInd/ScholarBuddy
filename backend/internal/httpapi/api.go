package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/ScrewDriverInd/ScholarBuddy/internal/auth"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/config"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/opportunity"
	"github.com/ScrewDriverInd/ScholarBuddy/internal/platform/httperr"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type API struct {
	opportunities *opportunity.Store
	auth          *auth.Service
	cfg           config.Config
	log           *slog.Logger
}

func New(store *opportunity.Store, authService *auth.Service, cfg config.Config, log *slog.Logger) http.Handler {
	a := &API{opportunities: store, auth: authService, cfg: cfg, log: log}
	r := chi.NewRouter()
	r.Use(a.requestID, a.recover, a.cors)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	registerDocs(r)
	r.Get("/", a.list)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/opportunities", a.list)
		r.Get("/opportunities/{id}", a.get)
		r.With(a.auth.UserMiddleware).Get("/me", a.me)
		r.With(a.auth.UserMiddleware).Post("/opportunities", a.create)
		r.With(a.auth.UserMiddleware).Patch("/opportunities/{id}", a.update)
		r.With(a.auth.AdminMiddleware).Get("/admin/opportunities", a.adminList)
	})
	r.Post("/abbujaan/login", a.adminLogin)
	return r
}
func (a *API) list(w http.ResponseWriter, r *http.Request) {
	filter, valid := opportunityFilter(r.URL.Query().Get("type"))
	if !valid {
		a.error(w, r, http.StatusBadRequest, "invalid_type", "type is invalid")
		return
	}
	page, perPage := pagination(r)
	result, err := a.opportunities.List(r.Context(), filter, page, perPage)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	respond(w, http.StatusOK, result)
}

// opportunityFilter accepts only the public opportunity enum values.
func opportunityFilter(value string) (opportunity.Type, bool) {
	switch value {
	case "", "all":
		return "", true
	case string(opportunity.Scholarship):
		return opportunity.Scholarship, true
	case string(opportunity.Hackathon):
		return opportunity.Hackathon, true
	case string(opportunity.Internship):
		return opportunity.Internship, true
	case string(opportunity.Research):
		return opportunity.Research, true
	case string(opportunity.Extras):
		return opportunity.Extras, true
	default:
		return "", false
	}
}
func (a *API) get(w http.ResponseWriter, r *http.Request) {
	id, err := opportunityDisplayID(chi.URLParam(r, "id"))
	if err != nil {
		a.error(w, r, 400, "invalid_id", "opportunity ID must be a positive number")
		return
	}
	o, err := a.opportunities.Get(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		a.error(w, r, 404, "not_found", "opportunity was not found")
		return
	}
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	respond(w, 200, o)
}
func (a *API) me(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.FromContext(r.Context())
	respond(w, 200, map[string]any{"id": p.ID, "email": p.Email, "role": p.Role})
}
func (a *API) create(w http.ResponseWriter, r *http.Request) {
	in, ok := a.input(w, r)
	if !ok {
		return
	}
	p, _ := auth.FromContext(r.Context())
	o, err := a.opportunities.Create(r.Context(), in, p.ID)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	respond(w, 201, o)
}
func (a *API) update(w http.ResponseWriter, r *http.Request) {
	id, err := opportunityDisplayID(chi.URLParam(r, "id"))
	if err != nil {
		a.error(w, r, 400, "invalid_id", "opportunity ID must be a positive number")
		return
	}
	in, ok := a.input(w, r)
	if !ok {
		return
	}
	p, _ := auth.FromContext(r.Context())
	o, err := a.opportunities.Update(r.Context(), id, p.ID, in)
	if errors.Is(err, pgx.ErrNoRows) {
		a.error(w, r, 403, "forbidden", "you can only update opportunities you created")
		return
	}
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	respond(w, 200, o)
}

func opportunityDisplayID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid display ID")
	}
	return id, nil
}
func (a *API) input(w http.ResponseWriter, r *http.Request) (opportunity.Input, bool) {
	var in opportunity.Input
	if err := decode(r, &in); err != nil {
		a.error(w, r, 400, "invalid_json", "request body must be valid JSON")
		return in, false
	}
	if err := in.Validate(); err != nil {
		a.error(w, r, 400, "validation_error", err.Error())
		return in, false
	}
	return in, true
}
func (a *API) adminList(w http.ResponseWriter, r *http.Request) {
	page, perPage := pagination(r)
	result, err := a.opportunities.List(r.Context(), "", page, perPage)
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	respond(w, 200, result)
}
func (a *API) adminLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
		AdminKey string `json:"admin_key"`
	}
	if err := decode(r, &req); err != nil {
		a.error(w, r, 400, "invalid_json", "request body must be valid JSON")
		return
	}
	if !auth.CredentialsMatch(a.cfg.AdminUserID, a.cfg.AdminPassword, a.cfg.AdminKey, req.UserID, req.Password, req.AdminKey) {
		a.error(w, r, 401, "invalid_credentials", "admin credentials are invalid")
		return
	}
	token, err := a.auth.IssueAdminToken()
	if err != nil {
		a.serverError(w, r, err)
		return
	}
	respond(w, 200, map[string]any{"access_token": token, "token_type": "Bearer", "role": "ROLE_ADMIN"})
}
func decode(r *http.Request, target any) error {
	d := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	d.DisallowUnknownFields()
	return d.Decode(target)
}
func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
}
func pagination(r *http.Request) (int, int) {
	page, perPage := 1, 20
	if v, _ := strconv.Atoi(r.URL.Query().Get("page")); v > 0 {
		page = v
	}
	if v, _ := strconv.Atoi(r.URL.Query().Get("per_page")); v > 0 && v <= 100 {
		perPage = v
	}
	return page, perPage
}
func (a *API) error(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	httperr.Write(w, status, code, msg, r.Header.Get("X-Request-ID"))
}
func (a *API) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.log.Error("request failed", "error", err, "request_id", r.Header.Get("X-Request-ID"))
	a.error(w, r, 500, "internal_error", "an unexpected error occurred")
}
func (a *API) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id := uuid.NewString()
			r.Header.Set("X-Request-ID", id)
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func (a *API) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				a.error(w, r, 500, "internal_error", "an unexpected error occurred")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func (a *API) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		for _, allowed := range a.cfg.CORSAllowedOrigins {
			if strings.TrimSpace(allowed) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

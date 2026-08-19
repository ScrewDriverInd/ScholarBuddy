# ScholarBuddy API

Go backend for opportunities (scholarships, hackathons, internships, research/extra).

## Quick start

1. Copy `.env.example` to `.env` and fill in the Supabase Postgres connection string, project credentials, and admin secrets.
2. Create the PostgreSQL database and run `make migrate-up` (requires `migrate`).
3. Run `make test`, then `make run`.

Swagger UI is available in a normal development build at `http://localhost:8080/docs`. Build production with `go build -tags production ./cmd/api`; that binary contains neither `/docs` nor `/docs/openapi.yaml`.

## API

All responses are JSON envelopes: `{ "data": ... }` or `{ "error": { "code", "message", "request_id" } }`.

| Method | Path | Authentication | Purpose |
|---|---|---|---|
| GET | `/healthz` | No | Service health |
| GET | `/` or `/api/v1/opportunities?type=scholarship` | No | Paginated opportunity list; `type` may be `all`, `scholarship`, `hackathon`, `internship`, or `research` |
| GET | `/api/v1/opportunities/{id}` | No | Full opportunity details |
| GET | `/api/v1/me` | Supabase bearer token | Provision/read current user |
| POST | `/api/v1/opportunities` | Supabase bearer token | Create an opportunity |
| PATCH | `/api/v1/opportunities/{id}` | Supabase bearer token | Update own opportunity only |
| GET | `/api/v1/admin/opportunities` | Admin session token | Admin dashboard list |
| POST | `/abbujaan/login` | No | Admin login using `user_id`, `password`, `admin_key` |

For Google login, the frontend should call Supabase OAuth and send the returned access token as `Authorization: Bearer <token>` to protected routes. New accounts are always provisioned as `ROLE_USER`.

Create/update request fields: `title`, `description`, `type`, `eligibility`, `steps`, `benefits`, `link`, `referral`. `type` must use the values above.

# ScholarBuddy API

Go backend for opportunities (scholarships, hackathons, internships, research, and extras).

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
| GET | `/` or `/api/v1/opportunities?type=scholarship` | No | Paginated opportunity list; `type` may be `all`, `scholarship`, `hackathon`, `internship`, `research`, or `extras` |
| GET | `/api/v1/opportunities/{id}` | No | Full opportunity details |
| GET | `/api/v1/me` | Supabase bearer token | Provision/read current user |
| POST | `/api/v1/opportunities` | Supabase bearer token | Create an opportunity |
| PATCH | `/api/v1/opportunities/{id}` | Supabase bearer token | Update own opportunity only |
| GET | `/api/v1/abbujaan` or `/api/v1/abbujaan/opportunities` | Admin session token | Pending opportunities for the admin dashboard |
| PATCH | `/api/v1/abbujaan/opportunities/{id}/approve` | Admin session token | Approve a pending opportunity |
| DELETE | `/api/v1/abbujaan/opportunities/{id}` | Admin session token | Delete an opportunity |
| POST | `/api/v1/abbujaan/login` | No | Admin login using `username` and `password` |

For Google login, the frontend should call Supabase OAuth and send the returned access token as `Authorization: Bearer <token>` to protected routes. New accounts are always provisioned as `ROLE_USER`. Created and updated opportunities remain pending until approved by an administrator; only approved opportunities appear publicly.

Create/update request fields: `title`, `description`, `types`, `eligibility`, `steps`, `benefits`, `link`, `referral`. `types` is a required, non-empty array using any combination of `scholarship`, `hackathon`, `internship`, `research`, and `extras`.

Opportunity responses expose a stable numeric `id` for tables and routes. UUID primary keys remain internal to the backend.

## Create an administrator

Run the following in Supabase SQL Editor, replacing the username, email, and password. It uses bcrypt through Supabase's `pgcrypto` extension and assigns both roles:

```sql
INSERT INTO users (id, email, full_name, username, password_hash, roles)
VALUES (
  gen_random_uuid(),
  'admin@example.com',
  'ScholarBuddy Admin',
  'admin',
  crypt('replace-with-a-strong-password', gen_salt('bf')),
  ARRAY['ROLE_USER'::user_role, 'ROLE_ADMIN'::user_role]
);
```

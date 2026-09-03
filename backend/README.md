# ScholarBuddy API

Go backend for opportunities (scholarships, hackathons, internships, research, and extras).

## Quick start

1. Copy `.env.example` to `.env` and fill in the Supabase Postgres connection string, project credentials, and admin secrets.
2. Create the PostgreSQL database and run `make migrate-up`.
3. Run `make test`, then `make run`.

## Reset an existing development database

The project now has one clean initial migration. To discard the previous ScholarBuddy schema, opportunity data, custom users, and migration history, export `DATABASE_URL` and run:

```bash
make db-reset
make migrate-up
```

`db-reset` is destructive. It does not remove Supabase Auth users or Google OAuth settings.

The built-in Go migrator applies every pending `*.up.sql` migration with `make migrate-up`; `make migrate-down` rolls back only the most recently applied migration.

Swagger UI is available in a normal development build at `http://localhost:8080/docs`. Build production with `go build -tags production ./cmd/api`; that binary contains neither `/docs` nor `/docs/openapi.yaml`.

## API

All responses are JSON envelopes: `{ "data": ... }` or `{ "error": { "code", "message", "request_id" } }`.

| Method | Path | Authentication | Purpose |
|---|---|---|---|
| GET | `/healthz` | No | Service health |
| GET | `/` | No | `{ "message": "Welcome to ScholarBuddy" }` |
| GET | `/api/v1/opportunities?type=scholarship` | No | Approved opportunities, ranked by click count; `type` may be `all`, `scholarship`, `hackathon`, `internship`, `research`, or `extras` |
| GET | `/api/v1/opportunities/{id}` | No | Full opportunity details |
| GET | `/api/v1/me` | Supabase bearer token | Provision/read current user |
| POST | `/api/v1/opportunities` | Supabase bearer token | Create an opportunity |
| PATCH | `/api/v1/opportunities/{id}` | Supabase bearer token | Update own opportunity only |
| GET | `/api/v1/abbujaan` | Admin session token | Pending opportunities for the admin dashboard |
| PATCH | `/api/v1/abbujaan/opportunities/{id}/approve` | Admin session token | Approve a pending opportunity |
| DELETE | `/api/v1/abbujaan/opportunities/{id}` | Admin session token | Delete an opportunity |
| POST | `/api/v1/abbujaan/login` | No | Admin login using `username` and `password` |

For Google login, the frontend should call Supabase OAuth and send the returned access token as `Authorization: Bearer <token>` to protected routes. New accounts are always provisioned as `ROLE_USER`. Created and updated opportunities remain pending until approved by an administrator; only approved opportunities appear publicly.

Create/update request fields: `title`, `description`, `types`, `eligibility`, `steps`, `benefits`, `link`, `referral`. `types` is a required, non-empty array using any combination of `scholarship`, `hackathon`, `internship`, `research`, and `extras`.

Opportunity responses expose their UUID `id`. The frontend renders the row number it wants to display; UUIDs are used only for API routes and data identity. Opening a public opportunity increments its `click_count`, and public lists are ranked by that count.

## Authentication

### User Authentication (Google OAuth via Supabase)
Users log in with Google OAuth through Supabase. When a user logs in:
- Supabase handles the OAuth flow and stores the user in their auth system
- Our backend receives the Supabase JWT token
- We automatically store the user in our `users` table with:
  - `id`: Supabase user UUID
  - `email`: Gmail address
  - `full_name`: Name from Google account
  - `username`: Auto-generated from email (part before @)
  - `roles`: `['ROLE_USER']` by default

### Admin Authentication (Username & Password)
Admins use a separate authentication system with username and password stored in the `admins` table:
- Passwords are hashed using argon2id algorithm
- Admin login returns a JWT token valid for 8 hours
- Admin tokens include the `ROLE_ADMIN` role

To create an admin, hash the password with argon2id and insert into the `admins` table:
```sql
INSERT INTO admins (username, password_hash) VALUES ('admin', '$argon2id$v=19$m=65536,t=3,p=4$...');
```

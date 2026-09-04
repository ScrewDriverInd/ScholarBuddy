# ScholarBuddy Backend - Spring Boot 4

Production-grade Spring Boot 4 backend with Google OAuth2 authentication.

## Features

- **Google OAuth2 only** - Single sign-in for all users
- **Role-based access** - ROLE_USER (all), ROLE_ADMIN (privileged)
- **Opportunity workflow** - Create → Pending → Admin Approval → Public
- **Click tracking** - Records user engagement
- **Production-ready** - Validation, error handling, migrations, health checks

## Tech Stack

- Java 25
- Spring Boot 4.1.1
- Spring Security with OAuth2 Client
- Spring Data JPA
- PostgreSQL
- Flyway
- Lombok
- SpringDoc OpenAPI

## Quick Start

1. **Create PostgreSQL database** (e.g., Neon)

2. **Set up Google OAuth2**
   - [Google Cloud Console](https://console.cloud.google.com/)
   - Create OAuth 2.0 credentials
   - Authorized redirect URI: `http://localhost:8080/login/oauth2/code/google`

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your credentials
   ```

4. **Run**
   ```bash
   ./gradlew bootRun
   ```

5. **Create first admin**
   - Sign in with Google at http://localhost:8080/oauth2/authorization/google
   - Run in PostgreSQL:
     ```sql
     SELECT id FROM users WHERE email = 'your-email@example.com';
     INSERT INTO user_roles (user_id, role) VALUES ('your-user-id', 'ROLE_ADMIN');
     ```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL JDBC URL | `jdbc:postgresql://ep-xxx.neon.tech/scholarbuddy` |
| `GOOGLE_CLIENT_ID` | Google OAuth2 client ID | `xxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Google OAuth2 client secret | `GOCSPX-xxx` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated origins | `http://localhost:5173` |
| `HTTP_ADDR` | Server port (optional) | `8080` |

## API Endpoints

### Public
- `GET /` - Welcome message
- `GET /actuator/health` - Health check
- `GET /api/v1/opportunities` - List approved (with `?type`, `?page`, `?per_page`)
- `GET /api/v1/opportunities/{id}` - Get opportunity (records click)

### Authenticated
- `GET /api/v1/user/me` - Current user info
- `POST /api/v1/opportunities` - Create opportunity (pending)
- `PATCH /api/v1/opportunities/{id}` - Update own opportunity (resets to pending)
- `POST /logout` - Logout

### Admin
- `GET /api/v1/abbujaan/opportunities` - List pending opportunities
- `PATCH /api/v1/abbujaan/opportunities/{id}/approve` - Approve opportunity
- `DELETE /api/v1/abbujaan/opportunities/{id}` - Delete any opportunity
- `POST /api/v1/abbujaan/users/{id}/admin` - Grant ROLE_ADMIN to user

## Authentication

**Google OAuth2 Flow:**
1. Frontend redirects to `/oauth2/authorization/google`
2. User authenticates with Google
3. Backend exchanges code for user info
4. User created/updated in database with ROLE_USER
5. Session cookie set
6. User redirected to frontend

**Session-based** - Cookies, not tokens. Sessions are httpOnly and secure in production.

## Database Schema

```sql
-- Enums
CREATE TYPE user_role AS ENUM ('ROLE_USER', 'ROLE_ADMIN');
CREATE TYPE opportunity_type AS ENUM ('SCHOLARSHIP', 'HACKATHON', 'INTERNSHIP', 'RESEARCH', 'EXTRAS');
CREATE TYPE approval_status AS ENUM ('PENDING', 'APPROVED');

-- Users (from Google OAuth)
CREATE TABLE users (
    id UUID PRIMARY KEY,                    -- from Google 'sub'
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL GENERATED AS (split_part(email, '@', 1)) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- User roles (many-to-many)
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role user_role NOT NULL,
    PRIMARY KEY (user_id, role)
);

-- Opportunities
CREATE TABLE opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 250),
    description TEXT NOT NULL CHECK (char_length(description) BETWEEN 1 AND 10000),
    eligibility TEXT NOT NULL DEFAULT '',
    steps TEXT NOT NULL DEFAULT '',
    benefits TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    referral TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id),
    approval_status approval_status NOT NULL DEFAULT 'PENDING',
    click_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Opportunity types (many-to-many)
CREATE TABLE opportunity_types (
    opportunity_id UUID NOT NULL REFERENCES opportunities(id) ON DELETE CASCADE,
    type opportunity_type NOT NULL,
    PRIMARY KEY (opportunity_id, type)
);
```

## Project Structure

```
src/main/java/com/libreturtle/scholarbuddy/
├── config/
│   ├── CorsConfig.java           # CORS configuration
│   ├── RequestIdFilter.java      # X-Request-ID tracking
│   └── SecurityConfig.java       # Spring Security + OAuth2
├── controller/
│   ├── AdminController.java      # /api/v1/abbujaan/*
│   ├── OpportunityController.java # /api/v1/opportunities/*
│   ├── RootController.java       # /
│   └── UserController.java       # /api/v1/user/*
├── dto/
│   ├── ApiResponseBody.java      # {"data": ...} wrapper (record)
│   ├── ErrorResponse.java        # {"error": {...}} (record)
│   ├── OpportunityRequest.java   # Create/update DTO (record)
│   ├── OpportunityResponse.java  # Opportunity DTO (record)
│   ├── PageRequest.java          # Pagination validation (record)
│   ├── PageResponse.java         # Paginated response (record)
│   └── UserResponse.java         # User DTO (record)
├── exception/
│   ├── ApiException.java         # Custom exceptions
│   └── GlobalExceptionHandler.java # @ControllerAdvice
├── model/
│   ├── ApprovalStatus.java       # PENDING, APPROVED
│   ├── Opportunity.java          # JPA entity
│   ├── OpportunityType.java      # SCHOLARSHIP, HACKATHON, etc.
│   ├── User.java                 # JPA entity
│   └── UserRole.java             # ROLE_USER, ROLE_ADMIN
├── repository/
│   ├── OpportunityRepository.java
│   └── UserRepository.java
├── security/
│   ├── CustomOAuth2User.java     # OAuth2User wrapper
│   ├── CustomOAuth2UserService.java # Loads/creates users
│   └── SecurityUtils.java        # Get current user
├── service/
│   ├── OpportunityService.java   # Business logic
│   └── UserService.java          # User management
└── validation/
    ├── ValidUrl.java             # Custom @ValidUrl annotation
    └── UrlValidator.java         # URL validation logic

src/main/resources/
├── application.yaml              # Configuration
└── db/migration/
    └── V1__initial_schema.sql    # Flyway migration
```

## Response Format

**Success:**
```json
{
  "data": { ... }
}
```

**Error:**
```json
{
  "error": {
    "code": "validation_error",
    "message": "title is required",
    "requestId": "uuid"
  }
}
```

## Development

```bash
# Build
./gradlew build

# Run tests (requires database)
./gradlew test

# Clean build
./gradlew clean build

# Run with profile
./gradlew bootRun --args='--spring.profiles.active=dev'
```

## API Documentation

Once running:
- Swagger UI: http://localhost:8080/swagger-ui.html
- OpenAPI spec: http://localhost:8080/v3/api-docs

## Production Deployment

1. Set environment variables
2. Use HTTPS (Spring Security will set secure cookies)
3. Configure CORS for your frontend domain
4. Set up database backups
5. Monitor `/actuator/health`
6. Review security headers and CSP

## Standards Applied

✅ **Java conventions** - UPPERCASE enums, records for DTOs  
✅ **Validation** - Bean Validation on all request DTOs  
✅ **Error handling** - Consistent format, no internal details leaked  
✅ **Security** - OAuth2, CSRF disabled (session + CORS), role-based access  
✅ **Database** - Flyway migrations, constraints, indexes  
✅ **API design** - RESTful, proper status codes, pagination  
✅ **Code quality** - Lombok, layered architecture, clear separation of concerns  

## Changes from Go Backend

**Removed:**
- Supabase JWT validation
- `/api/v1/abbujaan/login` endpoint (admin login)
- Separate `admins` table

**Added:**
- Google OAuth2 via Spring Security
- Session-based authentication
- `/api/v1/user/me` endpoint
- `POST /api/v1/abbujaan/users/{id}/admin` endpoint

**Updated:**
- Enums to UPPERCASE (SCHOLARSHIP not scholarship)
- Health endpoint: `/actuator/health` (Spring standard)
- Logout: `POST /logout` (Spring standard)

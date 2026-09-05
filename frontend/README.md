# ScholarBuddy Frontend

React frontend for ScholarBuddy with Google OAuth2 authentication via Spring Boot backend.

## Features

- Browse approved opportunities (public)
- Filter by type (scholarship, hackathon, internship, research, extras)
- Google OAuth2 login
- Admin dashboard for approving/deleting opportunities
- Modular component architecture
- HackerNews-inspired design

## Tech Stack

- React 19
- Vite 8
- CSS (no frameworks)
- Session-based authentication (cookies)

## Project Structure

```
src/
├── components/
│   ├── Header.jsx          # Navigation header with login/logout
│   ├── Home.jsx            # Main opportunity list with filters
│   ├── Detail.jsx          # Opportunity detail view
│   └── AdminDashboard.jsx  # Admin panel for pending opportunities
├── api.js                  # API client with fetch wrapper
├── useAuth.js              # Authentication hook
├── App.jsx                 # Main app component
├── App.css                 # Application styles
├── index.css               # Global styles
└── main.jsx                # Entry point
```

## Setup

The frontend requires Node.js 20.19 or newer (Node.js 26 is recommended).

1. **Install dependencies**
   ```bash
   npm install
   ```

2. **Configure environment** (optional)
   ```bash
   cp .env.example .env
   # VITE_API_BASE is empty by default (uses proxy)
   ```

3. **Run development server**
   ```bash
   npm run dev
   ```

4. **Build for production**
   ```bash
   npm run build
   ```

## Development

The dev server runs on `http://localhost:5173` and proxies API requests to the backend at `http://localhost:8080`.

### Proxy Configuration

Vite is configured to proxy these paths to the backend:
- `/api/*` - API endpoints
- `/oauth2/*` - OAuth2 login flow
- `/logout` - Logout endpoint

### Authentication Flow

1. User clicks "login with google"
2. Frontend redirects to `/oauth2/authorization/google`
3. Spring backend handles OAuth2 flow with Google
4. User is redirected back with session cookie
5. Frontend calls `/api/v1/user/me` to get user info
6. Session cookie is included automatically in all requests

### API Endpoints Used

**Public:**
- `GET /api/v1/opportunities` - List approved opportunities
- `GET /api/v1/opportunities/{id}` - Get opportunity detail

**Authenticated:**
- `GET /api/v1/user/me` - Get current user
- `POST /logout` - Logout

**Admin:**
- `GET /api/v1/abbujaan/opportunities` - List pending opportunities
- `PATCH /api/v1/abbujaan/opportunities/{id}/approve` - Approve
- `DELETE /api/v1/abbujaan/opportunities/{id}` - Delete

## Components

### Header
Navigation bar with app name, tagline, and login/logout button.

### Home
Main view showing:
- Filter chips for opportunity types
- Paginated table of opportunities
- Click tracking on opportunity links

### Detail
Shows full opportunity information:
- Title, description, types
- Eligibility, steps, benefits
- External link to opportunity

### AdminDashboard
Admin-only view for:
- Listing pending opportunities
- Approving opportunities
- Deleting opportunities

## Hooks

### useAuth
Custom hook that provides:
- `user` - Current user object (null if not logged in)
- `loading` - Auth check loading state
- `login()` - Redirect to Google OAuth
- `logout()` - Logout and clear session
- `isAdmin` - Boolean if user has ROLE_ADMIN
- `checkAuth()` - Manually refresh auth state

## API Client

`api.js` provides:
- `apiFetch(path, options)` - Fetch wrapper with credentials
- `unwrapData(payload)` - Extract data from `{data: ...}` envelope
- `apiURL(path)` - Build full API URL

All requests include `credentials: "include"` for session cookies.

## Changes from Original

**Removed:**
- Supabase client (`@supabase/supabase-js`)
- Admin username/password login
- Bearer token authentication
- `/abbujaan/login` endpoint

**Added:**
- Google OAuth2 via Spring backend
- Session-based authentication
- `/api/v1/user/me` endpoint for user info
- Modular component structure
- `useAuth` hook for auth state

**Updated:**
- Enum values to uppercase (SCHOLARSHIP not scholarship)
- Response format handling (`data.clickCount` not `data.click_count`)
- Admin access check (via roles, not separate token)

## Notes

- Admin access requires `ROLE_ADMIN` in user roles
- First admin must be created via database (see backend README)
- Session cookies are httpOnly and secure in production
- OAuth redirect URI: `http://localhost:8080/login/oauth2/code/google`

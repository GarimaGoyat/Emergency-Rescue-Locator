# Emergency Rescue Locator

A full-stack emergency SOS and real-time location tracking web application. Users can send emergency alerts with live GPS tracking; admins monitor active emergencies on a map-based dashboard.

## Tech Stack

| Layer | Technologies |
|-------|-------------|
| Frontend | React, Vite, Tailwind CSS, React Router, Axios, Leaflet.js |
| Backend | Go, Gin, JWT, GORM |
| Database | PostgreSQL |
| Deployment | Vercel (frontend), Render (backend), Neon (PostgreSQL) |

## Project Structure

```
Emergency-Rescue-Locator/
├── backend/
│   ├── cmd/server/main.go          # Application entry point
│   ├── internal/
│   │   ├── config/                 # Environment configuration
│   │   ├── controllers/            # HTTP request handlers
│   │   ├── database/               # DB connection & auto-migrate
│   │   ├── middleware/             # Auth, CORS, rate limiting
│   │   ├── models/                 # GORM data models
│   │   ├── repositories/           # Database access layer
│   │   ├── routes/                 # API route definitions
│   │   └── services/               # Business logic
│   ├── migrations/                 # SQL migration files
│   ├── scripts/seed_admin.sql      # Admin user seed script
│   ├── Dockerfile
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── components/             # Reusable UI components
│   │   ├── context/                # Auth & notification state
│   │   ├── pages/                  # Route pages
│   │   └── services/api.js         # Axios API client
│   ├── vercel.json
│   └── .env.example
├── docker-compose.yml
├── render.yaml
└── README.md
```

## Features

- **Authentication** — Register, login, JWT tokens, protected routes
- **Emergency SOS** — Large confirm-to-send SOS button with location capture
- **Real-Time Tracking** — Browser geolocation updates every 10 seconds
- **Admin Dashboard** — View/search/resolve emergencies, statistics cards
- **Maps** — Leaflet.js map with clickable emergency markers
- **Notifications** — In-app toast alerts for success, errors, and status updates
- **Security** — bcrypt passwords, JWT auth, input validation, rate limiting

## API Routes

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Register user |
| POST | `/api/v1/auth/login` | Login user |

### Protected (Bearer JWT)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/auth/profile` | Get current user profile |
| POST | `/api/v1/emergencies` | Create SOS emergency |
| GET | `/api/v1/emergencies/active` | Get user's active emergency |
| GET | `/api/v1/emergencies/:id` | Get emergency details |
| POST | `/api/v1/emergencies/:id/cancel` | Cancel emergency |
| POST | `/api/v1/emergencies/:id/location` | Send location update |
| GET | `/api/v1/emergencies/:id/location/latest` | Get latest location |
| GET | `/api/v1/emergencies/:id/location/history` | Get location history |

### Admin (JWT + admin role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/emergencies` | List active emergencies |
| GET | `/api/v1/admin/emergencies/search?q=&status=` | Search emergencies |
| GET | `/api/v1/admin/emergencies/:id` | Full emergency details |
| POST | `/api/v1/admin/emergencies/:id/resolve` | Mark as resolved |
| GET | `/api/v1/admin/stats` | Dashboard statistics |

## Database Schema

### users
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| email | VARCHAR(255) | Unique email |
| password | VARCHAR(255) | bcrypt hash |
| name | VARCHAR(255) | Full name |
| phone | VARCHAR(20) | Phone number |
| role | VARCHAR(20) | `user` or `admin` |

### emergencies
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | FK → users |
| status | VARCHAR(20) | active, resolved, cancelled |
| description | TEXT | Emergency description |
| latitude | DOUBLE | Initial latitude |
| longitude | DOUBLE | Initial longitude |
| address | TEXT | Location address |
| resolved_at | TIMESTAMPTZ | Resolution timestamp |

### location_updates
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| emergency_id | UUID | FK → emergencies |
| latitude | DOUBLE | GPS latitude |
| longitude | DOUBLE | GPS longitude |
| accuracy | DOUBLE | GPS accuracy (meters) |
| recorded_at | TIMESTAMPTZ | When location was recorded |

## Local Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 16+ (or Docker)

### Option 1: Docker Compose (recommended)

```bash
# Start PostgreSQL + backend
docker-compose up -d

# Frontend (separate terminal)
cd frontend
cp .env.example .env
npm install
npm run dev
```

Backend: http://localhost:8080  
Frontend: http://localhost:5173

### Option 2: Manual setup

**1. Database**

```bash
# Create database
createdb emergency_rescue

# Run migrations (optional — GORM auto-migrates on startup)
psql emergency_rescue < backend/migrations/001_create_users.sql
psql emergency_rescue < backend/migrations/002_create_emergencies.sql
psql emergency_rescue < backend/migrations/003_create_location_updates.sql
```

**2. Backend**

```bash
cd backend
cp .env.example .env
# Edit .env with your DATABASE_URL and JWT_SECRET

go mod tidy
go run ./cmd/server
```

**3. Frontend**

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

### Create Admin User

Register a normal account, then promote to admin:

```sql
UPDATE users SET role = 'admin' WHERE email = 'your@email.com';
```

Or run the seed script (default: `admin@rescue.local` / `admin123`):

```bash
psql $DATABASE_URL -f backend/scripts/seed_admin.sql
```

## Deployment

### PostgreSQL on Neon

1. Create a project at [neon.tech](https://neon.tech)
2. Copy the connection string (with `?sslmode=require`)
3. Run migration SQL files in the Neon SQL Editor

### Backend on Render

1. Connect your GitHub repo to Render
2. Use the included `render.yaml` or create a Web Service:
   - **Runtime:** Docker
   - **Dockerfile Path:** `backend/Dockerfile`
   - **Health Check Path:** `/health`
3. Set environment variables:
   - `DATABASE_URL` — Neon connection string
   - `JWT_SECRET` — long random secret
   - `CORS_ORIGINS` — your Vercel frontend URL
   - `GIN_MODE=release`

### Frontend on Vercel

1. Import the repo on [vercel.com](https://vercel.com)
2. Set **Root Directory** to `frontend`
3. Set environment variable:
   - `VITE_API_URL` — `https://your-api.onrender.com/api/v1`
4. Deploy — `vercel.json` handles SPA routing

## Environment Variables

### Backend (`backend/.env`)
| Variable | Description | Example |
|----------|-------------|---------|
| PORT | Server port | `8080` |
| DATABASE_URL | PostgreSQL connection | `postgres://...` |
| JWT_SECRET | JWT signing key | random 32+ chars |
| JWT_EXPIRY_HOURS | Token lifetime | `24` |
| CORS_ORIGINS | Allowed origins (comma-separated) | `http://localhost:5173` |
| RATE_LIMIT_REQUESTS | Max requests per period | `100` |
| RATE_LIMIT_DURATION | Rate limit window | `1m` |

### Frontend (`frontend/.env`)
| Variable | Description | Example |
|----------|-------------|---------|
| VITE_API_URL | Backend API base URL | `http://localhost:8080/api/v1` |

## Architecture

```
┌─────────────┐     REST/JWT      ┌─────────────┐
│   React     │ ◄──────────────► │  Gin API    │
│  Frontend   │                   │  (Go)       │
└─────────────┘                   └──────┬──────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    │                    │                    │
              Controllers           Services            Repositories
                    │                    │                    │
                    └────────────────────┴────────────────────┘
                                         │
                                   PostgreSQL
```

**Request flow:** Route → Middleware (CORS, rate limit, JWT) → Controller → Service → Repository → Database

## License

MIT

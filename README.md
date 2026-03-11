# Study Platform API

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18.1-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=flat-square&logo=swagger&logoColor=black)

REST API for managing courses, students, and teachers built with Go.

## Tech Stack

- **Go** — standard `net/http`
- **PostgreSQL** — raw SQL queries via `database/sql` + `lib/pq`
- **JWT** — access/refresh token authentication with role-based authorization
- **Swagger** — auto-generated API documentation
- **Goose** — database migrations (auto-applied on startup)

## Key Features

- **Custom Error Handling:** Centralized middleware that maps domain-specific errors to standard HTTP status codes.
- **Transaction Management:** Context-based transaction injection for complex cross-repository flows.
- **Graceful Shutdown:** Ensures all active HTTP requests and database connections are cleanly closed on termination.
- **Custom Request Validation:** Human-readable validation errors dynamically generated from struct tags.

## Project Structure

```text
cmd/study-platform/     — application entrypoint
internal/
  app/                  — application setup, DI, HTTP server
  apperror/             — custom error types and codes
  config/               — configuration loading
  entity/               — domain entities
  handler/              — HTTP handlers
  middleware/           — auth, role-based access, logging
  repository/           — database repositories
  service/              — business logic layer
  util/                 — helpers
pkg/
  closer/               — graceful shutdown
  database/             — DB interface, transaction manager
  errwrap/              — error wrapping
  httpresponse/         — JSON response helpers
  logger/               — structured logging
  validator/            — request validation
migrations/             — goose SQL migrations
docs/                   — swagger generated docs
```

## Getting Started

### 1. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your values, or use the defaults as-is for local development.

### 2. Run with Docker

```bash
docker compose up --build
```

This starts both the API and PostgreSQL. Migrations are applied automatically on startup.

### 3. Open Swagger UI

```text
http://localhost:3800/swagger/index.html
```

_(Note: If you changed the application port in your `.env` file, make sure to replace `3800` in the URL with your custom port)._

## API Endpoints

> **Auth column legend:**
>
> - 🔓 **Public** — no authentication required
> - 🔒 **Any** — JWT required (any authenticated user)
> - 🔒 **Student** — JWT required + Student role
> - 🔒 **Teacher** — JWT required + Teacher role

### Authentication

| Method | Path             | Auth      | Description                                               |
| ------ | ---------------- | --------- | --------------------------------------------------------- |
| `POST` | `/auth/register` | 🔓 Public | Register new user (creates student profile automatically) |
| `POST` | `/auth/login`    | 🔓 Public | Login, returns access + refresh tokens                    |
| `POST` | `/auth/refresh`  | 🔓 Public | Get new token pair using refresh token                    |

### Users

| Method   | Path        | Auth   | Description                                                   |
| -------- | ----------- | ------ | ------------------------------------------------------------- |
| `PUT`    | `/users/me` | 🔒 Any | Update current user (name, email, password)                   |
| `DELETE` | `/users/me` | 🔒 Any | Delete current user (cascades to student/teacher/enrollments) |

### Students

| Method   | Path                              | Auth       | Description                    |
| -------- | --------------------------------- | ---------- | ------------------------------ |
| `GET`    | `/students`                       | 🔓 Public  | List all students              |
| `GET`    | `/students/{id}`                  | 🔓 Public  | Get student by ID              |
| `PUT`    | `/students/me`                    | 🔒 Student | Update current student profile |
| `POST`   | `/students/me/courses/{courseId}` | 🔒 Student | Enroll in a course             |
| `DELETE` | `/students/me/courses/{courseId}` | 🔒 Student | Unenroll from a course         |

### Teachers

| Method | Path             | Auth       | Description                                |
| ------ | ---------------- | ---------- | ------------------------------------------ |
| `GET`  | `/teachers`      | 🔓 Public  | List all teachers                          |
| `GET`  | `/teachers/{id}` | 🔓 Public  | Get teacher by ID                          |
| `POST` | `/teachers`      | 🔒 Any     | Become a teacher (creates teacher profile) |
| `PUT`  | `/teachers/me`   | 🔒 Teacher | Update current teacher profile             |

### Courses

| Method   | Path            | Auth       | Description       |
| -------- | --------------- | ---------- | ----------------- |
| `GET`    | `/courses`      | 🔓 Public  | List all courses  |
| `GET`    | `/courses/{id}` | 🔓 Public  | Get course by ID  |
| `POST`   | `/courses`      | 🔒 Teacher | Create a course   |
| `PUT`    | `/courses/{id}` | 🔒 Teacher | Update own course |
| `DELETE` | `/courses/{id}` | 🔒 Teacher | Delete own course |

## Authentication

The API uses **JWT** with access and refresh tokens.

- **Access token** — short-lived, sent in `Authorization: Bearer <token>` header
- **Refresh token** — long-lived, used to obtain a new token pair via `POST /auth/refresh`
- **Roles** — `Student` (assigned on registration), `Teacher` (assigned when creating teacher profile)

JWT payload contains `userId` and `roles[]`. Role-based middleware restricts access to protected endpoints.

## Architecture

```text
Client → net/http Router → Middleware (Auth, Logging) → Handler → Service → Repository → PostgreSQL
```

- **Handlers** — parse requests, extract path variables/context, validate input, return JSON.
- **Services** — handle core business logic and authorization checks.
- **Repositories** — execute raw SQL queries via `database/sql`.

## Testing

Unit tests for the HTTP layer (Handlers) and Business logic layer (Services) are implemented using the **Table-Driven Tests** pattern. External dependencies are isolated using mocks generated via `testify/mock`.

Run the tests with coverage:

```bash
go test -v -cover ./...
```

| Layer    | Packages                             | Coverage  |
| -------- | ------------------------------------ | --------- |
| Handlers | auth, course, student, teacher, user | 86% – 95% |
| Services | course, teacher                      | 100%      |

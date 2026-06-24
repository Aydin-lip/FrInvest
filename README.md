# Recruitment API

A production-ready REST API built with Go, Gin, GORM, and MySQL.

## Tech Stack

- **Language:** Go 1.21
- **Framework:** Gin
- **ORM:** GORM
- **Database:** MySQL
- **Auth:** JWT (24h expiry)
- **Email:** SMTP with HTML templates

## Project Structure

```
recruitment-api/
├── cmd/
│   └── main.go                  # Entry point, DI wiring, routes
├── config/
│   └── config.go                # Env loading, DB connection
├── internal/
│   ├── controller/
│   │   ├── auth_controller.go
│   │   └── user_controller.go
│   ├── dto/
│   │   └── auth.go              # Request/Response DTOs
│   ├── email/
│   │   └── mailer.go            # SMTP mailer
│   ├── middleware/
│   │   └── auth_middleware.go   # JWT middleware
│   ├── models/
│   │   ├── user.go
│   │   └── verification_code.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── verification_repository.go
│   └── service/
│       ├── auth_service.go
│       ├── jwt_service.go
│       └── user_service.go
├── templates/
│   ├── verification_code.html
│   └── welcome.html
├── .env
├── go.mod
└── README.md
```

## Setup

### 1. Create MySQL database

```sql
CREATE DATABASE app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. Configure environment variables

Edit `.env`:

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_NAME=app_db
DB_USER=root
DB_PASSWORD=password
JWT_SECRET=super_secret_key
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=test@gmail.com
SMTP_PASSWORD=password
```

### 3. Install dependencies & run

```bash
go mod tidy
go run cmd/main.go
```

Tables are created automatically via GORM AutoMigrate.

---

## API Reference

### Public Endpoints

#### POST /api/auth/send-code
```json
{ "email": "user@example.com" }
```
Sends a 6-digit code valid for 2 minutes. Returns error if an active code already exists.

#### POST /api/auth/register
```json
{
  "firstName": "Ali",
  "lastName": "Ahmadi",
  "phoneNumber": "09123456789",
  "email": "ali@gmail.com",
  "code": "123456"
}
```
Registers user, marks code as used, returns JWT + user info.

#### POST /api/auth/login
```json
{ "email": "ali@gmail.com", "code": "123456" }
```
Login via email code, returns JWT + user info.

---

### Protected Endpoints (Bearer token required)

#### GET /api/users
Returns list of all active users (firstName, lastName, email, status).

#### PUT /api/users/status
```json
{ "userId": 1, "status": 3 }
```
Updates user status (0–4).

#### GET /api/users/statistics
Returns percentage breakdown of users by status (sum = 100%).

---

## User Status Values

| Value | Meaning     |
|-------|-------------|
| 0     | New         |
| 1     | Reviewed    |
| 2     | Interviewed |
| 3     | Offer Sent  |
| 4     | Rejected    |

## Phone Number Format

Must match: `^09\d{9}$`  
Example: `09123456789`

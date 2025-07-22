# Task Management API

A RESTful API built with Go, Gin, and MongoDB for managing users and tasks with role-based access control. Users can register, log in, and view tasks, while admins can manage users and tasks. The API uses JWT for authentication and bcrypt for password hashing.

---

## Features
- **User Management**: Register, log in, promote users to admin (admin-only), and fetch user details.
- **Task Management**: Create, read, update, and delete tasks with title, description, due date, and status (`pending`, `completed`, `missed`).
- **Role-Based Access**: Public routes for auth, authenticated routes for task viewing, and admin routes for management.
- **Security**: JWT tokens in secure cookies, bcrypt-hashed passwords.

---

## Technologies
- Go (v1.16+)
- Gin (HTTP framework)
- MongoDB (NoSQL database)
- Bcrypt (password hashing)
- JWT (authentication)
- godotenv (environment variables)

---

## Setup

### Environment Variables
Create a `.env` file in the project root:
```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=taskdb
JWT_SECRET=your-secure-jwt-secret
PORT=8080
```
- **MONGO_URI**: MongoDB connection string.
- **DB_NAME**: Database name (e.g., `taskdb`).
- **JWT_SECRET**: Unique secret for JWT signing.
- **PORT**: Server port (defaults to `8080`).

**Note**: Do not commit `.env` to version control.

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd task-management-api
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run main.go
   ```
   The API will be available at `http://localhost:8080` (or your specified port).

---

## Usage
- **Register**: `POST /auth/register` with `{ "username": "user", "password": "pass" }`. The first user is an admin.
- **Login**: `POST /auth/login` to get a JWT cookie.
- **Tasks**: Use `GET /tasks` (authenticated) or `POST /tasks` (admin) to manage tasks.
- **Users**: Use `GET /users` or `PATCH /promote/:id` (admin) to manage users.
- See [Documentation.md](Documentation.md) for full API details.

---

## Project Structure
```
task-manager/
├── Delivery/
│   ├── main.go                    # App entry point
│   ├── controllers/
│   │   └── controller.go          # HTTP handlers
│   └── routers/
│       └── router.go              # Route definitions
├── Domain/
│   └── domain.go                  # Models and interfaces
├── Infrastructure/
│   ├── auth_middleWare.go         # Authentication middleware
│   ├── jwt_service.go             # JWT handling
│   └── password_service.go        # Password hashing
├── Repositories/
│   ├── task_repository.go         # Task MongoDB operations
│   └── user_repository.go         # User MongoDB operations
├── Usecases/
│   ├── task_usecases.go           # Task business logic
│   └── user_usecases.go           # User business logic
├── api/
│   └── api_documentation.go       # API documentation utilities
└── .env                           # Environment variables
```

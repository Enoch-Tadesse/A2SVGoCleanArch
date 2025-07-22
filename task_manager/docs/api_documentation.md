# Task Management API Documentation

## Overview
The **Task Management API** is a RESTful service built with Go, Gin, and MongoDB. It enables user authentication, task management, and role-based access control. Users can register, log in, and manage tasks (create, read, update, delete), while administrators have additional privileges to manage users and tasks. The API uses JWT for secure authentication and bcrypt for password hashing.

This documentation covers setup, architecture, API endpoints, and usage details.

---

## Table of Contents
1. [Features](#features)
2. [Technologies](#technologies)
3. [Setup Instructions](#setup-instructions)
   - [Environment Variables](#environment-variables)
   - [Installation](#installation)
4. [Project Structure](#project-structure)
5. [API Endpoints](#api-endpoints)
   - [Public Routes](#public-routes)
   - [Authenticated Routes](#authenticated-routes)
   - [Admin Routes](#admin-routes)
6. [Authentication](#authentication)
7. [Error Handling](#error-handling)

---

## Features
- **User Management**:
  - Register users with unique usernames and hashed passwords.
  - Log in with JWT-based authentication (stored in secure cookies).
  - Promote users to admin (admin-only).
  - Retrieve user details by ID or list all users (admin-only).
- **Task Management**:
  - Create, read, update, and delete tasks with fields for title, description, due date, and status (`pending`, `completed`, `missed`).
  - Regular users can view tasks; admins can manage all tasks.
- **Role-Based Access Control**:
  - Public routes for registration and login.
  - Authenticated routes for task retrieval.
  - Admin-only routes for user and task management.
- **Security**:
  - Passwords hashed with bcrypt.
  - JWT tokens with expiration checks.
  - Secure cookies with `SameSite` and `HttpOnly` flags.

---

## Technologies
- **Go**: Backend programming language (v1.16+).
- **Gin**: HTTP web framework for routing and middleware.
- **MongoDB**: NoSQL database for storing users and tasks.
- **Bcrypt**: Password hashing library.
- **JWT (golang-jwt)**: Token-based authentication.
- **godotenv**: Environment variable management.

---

## Setup Instructions

### Environment Variables
Create a `.env` file in the project root with the following variables:
```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=taskdb
JWT_SECRET=your-secure-jwt-secret
PORT=8080
```
- **MONGO_URI**: MongoDB connection string (e.g., local or MongoDB Atlas).
- **DB_NAME**: Name of the MongoDB database (e.g., `taskdb`).
- **JWT_SECRET**: A secure, unique string for signing JWT tokens.
- **PORT**: Server port (defaults to `8080` if unset).

**Note**: Ensure the `.env` file is not committed to version control for security.

### Installation
1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   cd task-management-api
   ```

2. **Install Dependencies**:
   Install required Go modules:
   ```bash
   go mod tidy
   ```

3. **Run the Application**:
   Start the server:
   ```bash
   go run main.go
   ```
   The server will run at `http://localhost:8080` (or the port specified in `.env`).

4. **Verify Setup**:
   - Ensure MongoDB is running and accessible.
   - Test the `/auth/register` endpoint to create a user.
   - The first registered user is automatically an admin.

---

## Project Structure
The project follows a clean architecture with separation of concerns:

```
task-manager/
├── Delivery/
│   ├── main.go                    # Application entry point, initializes server and dependencies
│   ├── controllers/
│   │   └── controller.go          # HTTP handlers for task and user operations
│   └── routers/
│       └── router.go              # Route definitions and middleware setup
├── Domain/
│   └── domain.go                  # Models and interfaces for tasks and users
├── Infrastructure/
│   ├── auth_middleWare.go         # Authentication and authorization middleware
│   ├── jwt_service.go             # JWT generation and validation
│   └── password_service.go        # Password hashing and comparison
├── Repositories/
│   ├── task_repository.go         # MongoDB operations for tasks
│   └── user_repository.go         # MongoDB operations for users
├── Usecases/
│   ├── task_usecases.go           # Business logic for task operations
│   └── user_usecases.go           # Business logic for user operations
├── api/
│   └── api_documentation.go       # API documentation (optional or generated)
└── .env                           # Environment variables (not committed)
```

- **Domain**: Defines `Task` and `User` models and interfaces (`TaskRepository`, `UserRepository`, `TaskUsecase`, `UserUsecase`).
- **Infrastructure**: Implements utilities (JWT generation/validation, password hashing, middleware).
- **Repositories**: Handles MongoDB operations for users and tasks.
- **Usecases**: Contains business logic for task and user operations.
- **Delivery**: Manages HTTP routes and request handlers.
- **api**: Contains API documentation or related utilities.

---

## API Endpoints

### Public Routes
No authentication required.

#### `POST /auth/register`
Registers a new user. The first user is automatically an admin.

**Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response**:
- **201 Created**: `{ "message": "user created successfully", "data": { user } }`
- **400 Bad Request**: Invalid body or username exists.
- **500 Internal Server Error**: Server failure.

#### `POST /auth/login`
Authenticates a user and sets a JWT cookie.

**Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response**:
- **200 OK**: Returns user data, sets `Authentication` cookie.
- **400 Bad Request**: Invalid credentials or body.
- **500 Internal Server Error**: Server failure.

### Authenticated Routes
Requires a valid JWT cookie (`Authentication`).

#### `GET /tasks`
Fetches all tasks.

**Response**:
- **200 OK**: Array of tasks.
- **400 Bad Request**: No tasks exist.
- **500 Internal Server Error**: Server failure.

#### `GET /tasks/:id`
Fetches a task by ID.

**Response**:
- **200 OK**: Task object.
- **400 Bad Request**: Invalid ID or task not found.
- **500 Internal Server Error**: Server failure.

### Admin Routes
Requires a valid JWT cookie and admin privileges.

#### `GET /users`
Fetches all users.

**Response**:
- **200 OK**: Array of users.
- **400 Bad Request**: No users exist.
- **500 Internal Server Error**: Server failure.

#### `GET /users/:id`
Fetches a user by ID.

**Response**:
- **200 OK**: User object.
- **400 Bad Request**: Invalid ID or user not found.
- **500 Internal Server Error**: Server failure.

#### `PATCH /promote/:id`
Promotes a user to admin.

**Response**:
- **200 OK**: `{ "error": "user updated successfully" }`
- **400 Bad Request**: Invalid ID or user not found.
- **500 Internal Server Error**: Server failure.

#### `POST /tasks`
Creates a new task.

**Request Body**:
```json
{
  "title": "string",
  "description": "string",
  "due_date": "2025-12-31T23:59:59Z",
  "status": "pending"
}
```

**Response**:
- **201 Created**: Task object.
- **400 Bad Request**: Invalid body or past due date.
- **500 Internal Server Error**: Server failure.

#### `DELETE /tasks/:id`
Deletes a task by ID.

**Response**:
- **200 OK**: `{ "message": "task delete successfully" }`
- **400 Bad Request**: Invalid ID or task not found.
- **500 Internal Server Error**: Server failure.

#### `PUT /tasks/:id`
Updates a task by ID.

**Request Body**:
```json
{
  "title": "string",
  "description": "string",
  "due_date": "2025-12-31T23:59:59Z",
  "status": "pending|completed|missed"
}
```

**Response**:
- **200 OK**: Updated task or `{ "message": "no changes were made", "data": task }`.
- **400 Bad Request**: Invalid ID, body, status, or past due date.
- **500 Internal Server Error**: Server failure.

---

## Authentication
- **JWT Tokens**: Generated on login, stored in an `Authentication` cookie (24-hour expiry, `HttpOnly`, `Secure`, `SameSite=Lax`).
- **Middleware**:
  - `AuthenticationMiddleware`: Verifies JWT and sets user context.
  - `AuthorizationMiddleware`: Ensures the user is an admin for protected routes.
- **Usage**:
  - Include the `Authentication` cookie in requests to authenticated/admin routes.
  - Obtain the cookie via `/auth/login`.

---

## Error Handling
The API returns standardized JSON error responses:

```json
{
  "error": "error message"
}
```

Common errors:
- **400 Bad Request**: Invalid input (e.g., missing fields, invalid ID, past due date).
- **401 Unauthorized**: Missing/invalid JWT or non-admin access to admin routes.
- **500 Internal Server Error**: Database or server issues.

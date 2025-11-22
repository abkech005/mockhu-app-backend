# 🚀 Beginner's Guide to Building APIs

A complete step-by-step guide for creating REST APIs in Go.

---

## 📚 Table of Contents

1. [What is an API?](#what-is-an-api)
2. [API Development Flow](#api-development-flow)
3. [Step-by-Step Process](#step-by-step-process)
4. [Real Example: Building Login API](#real-example-building-login-api)
5. [Best Practices](#best-practices)
6. [Common Patterns](#common-patterns)

---

## 🤔 What is an API?

**API** = Application Programming Interface

Think of it as a **waiter in a restaurant**:
- **You (Client)**: Order food
- **Waiter (API)**: Takes your order to the kitchen
- **Kitchen (Backend)**: Prepares the food
- **Waiter (API)**: Brings food back to you

```
Client (Mobile/Web) → API → Server → Database
                      ↓
              Response comes back
```

---

## 🔄 API Development Flow

```
1. PLAN        → What endpoints do you need?
2. DESIGN      → Define request/response structure (DTOs)
3. DATABASE    → Design database tables (if needed)
4. IMPLEMENT   → Write the code
5. TEST        → Test with curl/Postman
6. DOCUMENT    → Write API docs
7. DEPLOY      → Put it on a server
```

---

## 📝 Step-by-Step Process

### **Step 1: Plan Your API**

Ask yourself:
- What features do I need? (Login, Signup, Posts, etc.)
- What data will I send/receive?
- Who can access what? (Public vs Protected)

**Example Planning:**
```
Features needed:
- User Signup
- User Login
- Get User Profile
- Update Profile
```

---

### **Step 2: Design Your Endpoints**

Follow REST conventions:

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/v1/auth/signup` | Create new user |
| POST | `/v1/auth/login` | Login user |
| GET | `/v1/users/:id` | Get user details |
| PUT | `/v1/users/:id` | Update user |
| DELETE | `/v1/users/:id` | Delete user |

**REST Conventions:**
- **GET** = Read/Fetch data
- **POST** = Create new data
- **PUT/PATCH** = Update data
- **DELETE** = Remove data

---

### **Step 3: Define Data Structures (DTOs)**

**DTO** = Data Transfer Object (what data goes in/out)

```go
// What client SENDS
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// What server RETURNS
type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}
```

---

### **Step 4: Create Project Structure**

```
myapi/
├── cmd/api/
│   └── main.go              # Entry point
├── internal/
│   └── app/
│       ├── auth/            # Auth feature
│       │   ├── dto.go       # Data structures
│       │   ├── handler.go   # HTTP handlers
│       │   └── routes.go    # Route definitions
│       └── user/            # User feature
│           ├── dto.go
│           ├── handler.go
│           └── routes.go
└── go.mod
```

---

### **Step 5: Write the Code**

#### **5.1 Create DTOs (Data Structures)**

`internal/app/auth/dto.go`
```go
package auth

type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  UserInfo `json:"user"`
}

type UserInfo struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}
```

---

#### **5.2 Create Handler (Business Logic)**

`internal/app/auth/handler.go`
```go
package auth

import "github.com/gofiber/fiber/v2"

type Handler struct {
    // Add dependencies (database, services, etc.)
}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) Login(c *fiber.Ctx) error {
    // Step 1: Parse incoming request
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    // Step 2: Validate (check if fields are correct)
    if req.Email == "" || req.Password == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Email and password required",
        })
    }

    // Step 3: Business logic
    // - Check if user exists in database
    // - Verify password
    // - Generate JWT token
    
    // For now, return dummy response
    return c.JSON(LoginResponse{
        Token: "dummy-jwt-token",
        User: UserInfo{
            ID:    "user-123",
            Email: req.Email,
            Name:  "John Doe",
        },
    })
}
```

---

#### **5.3 Create Routes**

`internal/app/auth/routes.go`
```go
package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
    handler := NewHandler()
    
    auth := app.Group("/v1/auth")
    
    auth.Post("/login", handler.Login)
    auth.Post("/signup", handler.Signup)
    auth.Post("/logout", handler.Logout)
}
```

---

#### **5.4 Main Application**

`cmd/api/main.go`
```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "myapi/internal/app/auth"
)

func main() {
    // Create Fiber app
    app := fiber.New()
    
    // Register routes from different domains
    auth.RegisterRoutes(app)
    
    // Start server
    log.Println("Server starting on :8080")
    app.Listen(":8080")
}
```

---

### **Step 6: Test Your API**

#### **Using curl:**
```bash
# Start server
go run cmd/api/main.go

# Test login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### **Using Postman:**
1. Open Postman
2. Create POST request to `http://localhost:8080/v1/auth/login`
3. Add JSON body
4. Send request

---

## 🎯 Real Example: Building Login API

Let's build a **complete login endpoint** from scratch:

### **Step 1: Define What You Need**
```
Feature: User Login
- User sends: email + password
- Server returns: JWT token + user info
- If wrong password: return error
```

### **Step 2: Create DTO**
```go
// What comes IN
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// What goes OUT (success)
type LoginResponse struct {
    Token string   `json:"token"`
    User  UserInfo `json:"user"`
}

// What goes OUT (error)
type ErrorResponse struct {
    Error string `json:"error"`
}
```

### **Step 3: Write Handler**
```go
func (h *Handler) Login(c *fiber.Ctx) error {
    // 1. Parse request
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{
            Error: "Invalid JSON",
        })
    }
    
    // 2. Validate
    if req.Email == "" {
        return c.Status(400).JSON(ErrorResponse{
            Error: "Email is required",
        })
    }
    
    // 3. Check database (pseudo-code)
    user := database.FindUserByEmail(req.Email)
    if user == nil {
        return c.Status(401).JSON(ErrorResponse{
            Error: "Invalid credentials",
        })
    }
    
    // 4. Verify password
    if !checkPassword(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(ErrorResponse{
            Error: "Invalid credentials",
        })
    }
    
    // 5. Generate token
    token := generateJWT(user.ID)
    
    // 6. Return success
    return c.JSON(LoginResponse{
        Token: token,
        User: UserInfo{
            ID:    user.ID,
            Email: user.Email,
            Name:  user.Name,
        },
    })
}
```

### **Step 4: Register Route**
```go
auth := app.Group("/v1/auth")
auth.Post("/login", handler.Login)
```

### **Step 5: Test**
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}'
```

---

## ✅ Best Practices

### **1. Use Proper HTTP Status Codes**
```go
200 - OK (success)
201 - Created (new resource)
400 - Bad Request (invalid input)
401 - Unauthorized (not logged in)
403 - Forbidden (no permission)
404 - Not Found
500 - Internal Server Error
```

### **2. Consistent Response Format**
```json
// Success
{
  "data": { ... },
  "message": "Success"
}

// Error
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Email is required"
  }
}
```

### **3. Validate Input**
```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### **4. Handle Errors Gracefully**
```go
if err != nil {
    return c.Status(500).JSON(fiber.Map{
        "error": "Something went wrong",
    })
}
```

### **5. Use Middleware**
```go
// Logging
app.Use(logger.New())

// Authentication check
app.Use("/v1/protected", authMiddleware)

// CORS
app.Use(cors.New())
```

---

## 🔄 Common Patterns

### **Pattern 1: CRUD Operations**

```go
// CREATE
POST   /v1/posts          → Create new post
// READ
GET    /v1/posts          → Get all posts
GET    /v1/posts/:id      → Get single post
// UPDATE
PUT    /v1/posts/:id      → Update post
// DELETE
DELETE /v1/posts/:id      → Delete post
```

### **Pattern 2: Authentication Flow**

```
1. User signs up      → POST /v1/auth/signup
2. User logs in       → POST /v1/auth/login (get token)
3. Use token          → Add header: Authorization: Bearer <token>
4. Access protected   → GET /v1/profile (with token)
5. Refresh token      → POST /v1/auth/refresh
6. Logout             → POST /v1/auth/logout
```

### **Pattern 3: Error Handling**

```go
func (h *Handler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    
    user, err := h.db.FindUser(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(404).JSON(fiber.Map{
                "error": "User not found",
            })
        }
        return c.Status(500).JSON(fiber.Map{
            "error": "Internal error",
        })
    }
    
    return c.JSON(user)
}
```

---

## 🎓 Learning Path

### **Week 1: Basics**
- Understand HTTP methods (GET, POST, PUT, DELETE)
- Learn JSON format
- Create simple endpoints

### **Week 2: Structure**
- Organize code into handlers/routes
- Use DTOs for requests/responses
- Add validation

### **Week 3: Database**
- Connect to PostgreSQL
- Create models
- Write queries

### **Week 4: Authentication**
- Implement JWT tokens
- Create middleware
- Protect routes

### **Week 5: Advanced**
- Error handling
- Logging
- Testing
- Documentation

---

## 📊 Your Current Project Structure

This is what you have now:

```
internal/app/
├── auth/               # Authentication feature
│   ├── dto.go         # LoginRequest, SignupRequest, etc.
│   ├── handler.go     # Login(), Signup(), Verify()
│   └── routes.go      # Register auth routes
│
├── onboarding/        # User onboarding
│   ├── dto.go
│   ├── handler.go
│   └── routes.go
│
└── upload/            # File uploads
    ├── dto.go
    ├── handler.go
    └── routes.go
```

**Each feature is self-contained!** ✅

---

## 🚀 Next Steps for You

1. ✅ **You have:** Dummy endpoints returning fake data
2. ⏳ **Next:** Connect to database (PostgreSQL)
3. ⏳ **Then:** Add real authentication (JWT tokens)
4. ⏳ **Then:** Add validation
5. ⏳ **Finally:** Deploy to production

---

## 🎯 Quick Reference

### **Creating a New Endpoint (5 Steps)**

1. **Plan:** What does it do?
2. **DTO:** Define request/response structure
3. **Handler:** Write the logic
4. **Routes:** Register the endpoint
5. **Test:** Use curl/Postman

### **Example: Create "Get Profile" Endpoint**

```go
// 1. DTO
type ProfileResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 2. Handler
func (h *Handler) GetProfile(c *fiber.Ctx) error {
    userID := c.Params("id")
    
    // Get from database (dummy for now)
    return c.JSON(ProfileResponse{
        ID:    userID,
        Name:  "John Doe",
        Email: "john@example.com",
    })
}

// 3. Route
app.Get("/v1/users/:id", handler.GetProfile)

// 4. Test
curl http://localhost:8085/v1/users/123
```

---

## 💡 Tips for Beginners

1. **Start Simple** - Don't add everything at once
2. **Test Often** - Test after each change
3. **Read Errors** - Error messages tell you what's wrong
4. **Use Examples** - Copy patterns from working code
5. **Ask Questions** - Don't hesitate to ask when stuck

---

**You're on the right path! Keep building! 🚀**


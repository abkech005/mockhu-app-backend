# Domain-Driven Architecture 🏗️

Your app is now organized by **business domains** instead of technical layers!

---

## 📁 New Structure

```
internal/app/
├── auth/                  # Authentication domain
│   ├── dto.go            # Auth-specific DTOs
│   ├── handler.go        # Auth HTTP handlers
│   ├── routes.go         # Auth route registration
│   └── (future)
│       ├── service.go    # Business logic
│       └── repository.go # Database operations
│
├── onboarding/           # Onboarding domain
│   ├── dto.go
│   ├── handler.go
│   ├── routes.go
│   └── (future) service.go, repository.go
│
└── upload/               # Upload domain
    ├── dto.go
    ├── handler.go
    ├── routes.go
    └── (future) storage.go

cmd/api/
└── main.go               # Clean entry point
```

---

## ✅ Benefits of Domain-Driven Design

### 1. **Self-Contained Domains**
Each domain has everything it needs in one folder:
- DTOs
- Handlers
- Routes
- Business logic (future)
- Database queries (future)

### 2. **Easy to Scale**
Add a new feature? Just create a new domain folder!

```bash
internal/app/
└── post/          # New feature
    ├── dto.go
    ├── handler.go
    ├── routes.go
    └── service.go
```

### 3. **Team-Friendly**
Different teams can work on different domains without conflicts:
- Team A: `auth/`
- Team B: `onboarding/`
- Team C: `post/`

### 4. **Clear Ownership**
Want to know where login logic is? Look in `auth/` folder!

---

## 🔄 How It Works

### 1. main.go (Entry Point)
```go
func main() {
    router := setupRouter()
    // Start server...
}

func setupRouter() *gin.Engine {
    r := gin.Default()
    
    // Each domain registers its own routes
    auth.RegisterRoutes(r)
    onboarding.RegisterRoutes(r)
    upload.RegisterRoutes(r)
    
    return r
}
```

### 2. Domain Routes (e.g., auth/routes.go)
```go
func RegisterRoutes(r *gin.Engine) {
    handler := NewHandler()
    
    auth := r.Group("/v1/auth")
    {
        auth.POST("/signup", handler.Signup)
        auth.POST("/login", handler.Login)
        // ... all auth routes
    }
}
```

### 3. Handler (e.g., auth/handler.go)
```go
func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest  // From auth/dto.go
    c.ShouldBindJSON(&req)
    
    // Business logic here
    
    c.JSON(200, LoginResponse{...})
}
```

---

## 📊 Comparison: Before vs After

### Before (Layer-Based)
```
internal/
├── transport/
│   ├── dtos/
│   │   ├── auth.go      # All DTOs mixed
│   │   ├── onboard.go
│   │   └── upload.go
│   └── http/
│       ├── handlers/
│       │   ├── auth_handler.go
│       │   ├── onboard_handler.go
│       │   └── upload_handler.go
│       └── router.go     # All routes in one file
```

**Problem:** Hard to find related code. Auth logic spread across multiple folders.

---

### After (Domain-Based) ✅
```
internal/app/
├── auth/            # Everything auth-related
│   ├── dto.go
│   ├── handler.go
│   └── routes.go
├── onboarding/      # Everything onboarding-related
│   ├── dto.go
│   ├── handler.go
│   └── routes.go
└── upload/          # Everything upload-related
    ├── dto.go
    ├── handler.go
    └── routes.go
```

**Benefit:** All auth logic in one place. Easy to find and modify!

---

## 🚀 Adding a New Feature

Want to add a "Post" feature? Here's how:

### Step 1: Create Domain Folder
```bash
mkdir -p internal/app/post
```

### Step 2: Create Files
```
internal/app/post/
├── dto.go        # Post DTOs
├── handler.go    # Post handlers
└── routes.go     # Post routes
```

### Step 3: Register Routes in main.go
```go
import "github.com/mockhu-app-backend/internal/app/post"

func setupRouter() *gin.Engine {
    r := gin.Default()
    
    auth.RegisterRoutes(r)
    onboarding.RegisterRoutes(r)
    upload.RegisterRoutes(r)
    post.RegisterRoutes(r)  // ← Add this line
    
    return r
}
```

**That's it!** Your new feature is live. 🎉

---

## 🔮 Future Additions (Per Domain)

Each domain can grow independently:

```
internal/app/auth/
├── dto.go         ✅ Created
├── handler.go     ✅ Created
├── routes.go      ✅ Created
├── service.go     ⏳ Add business logic
├── repository.go  ⏳ Add database queries
├── middleware.go  ⏳ Add auth-specific middleware
└── validator.go   ⏳ Add custom validation
```

---

## 📖 Learning Resources

### Domain-Driven Design Concepts

1. **Domain** = A business capability (auth, posts, payments)
2. **Bounded Context** = Each domain is isolated
3. **Aggregate** = Domain entity + business rules
4. **Repository** = Data access for a domain

### Your Current Domains

| Domain | Purpose | Endpoints |
|--------|---------|-----------|
| **auth** | Authentication & authorization | 6 endpoints |
| **onboarding** | User onboarding flow | 3 endpoints |
| **upload** | File uploads | 1 endpoint |

---

## 🧪 Testing the API

All endpoints still work the same!

```bash
# Health check
curl http://localhost:8082/health

# Login (auth domain)
curl -X POST http://localhost:8082/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"user@example.com","password":"test123"}'

# Onboard basic (onboarding domain)
curl -X POST http://localhost:8082/v1/onboard/basic \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","dob":"1999-05-12"}'
```

---

## 🎯 Key Takeaways

1. ✅ **Scalable** - Easy to add 100+ domains
2. ✅ **Organized** - Related code stays together
3. ✅ **Team-Friendly** - Multiple teams can work in parallel
4. ✅ **Maintainable** - Easy to find and modify features
5. ✅ **Clean main.go** - Just registers domains

---

## Next Steps

1. Add **service layer** for business logic
2. Add **repository layer** for database access
3. Add **middleware** (auth, logging, rate limiting)
4. Add **tests** for each domain
5. Add **validation** for DTOs

---

Congrats! You now have a professional, scalable architecture! 🚀


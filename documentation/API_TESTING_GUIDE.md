# API Testing Guide 🚀

## What We Built

You now have a working REST API with **10 endpoints** returning dummy data!

---

## Go Concepts You Learned

### 1. **Structs**
```go
type AuthHandler struct {
    // Fields go here
}
```
- Like classes in other languages
- Group related data and functions together

### 2. **Methods (Receiver Functions)**
```go
func (h *AuthHandler) HandleLogin(c *gin.Context) {
    // h is the receiver - access to struct fields
}
```
- `(h *AuthHandler)` makes this function belong to AuthHandler
- `*` means pointer (efficient, can modify struct)

### 3. **Error Handling**
```go
if err := c.ShouldBindJSON(&req); err != nil {
    // Handle error
    return
}
```
- No try/catch in Go - explicitly check errors
- Return early if error occurs

### 4. **Goroutines**
```go
go func() {
    // This runs in parallel
}()
```
- Lightweight threads
- Used in main.go to run server while listening for shutdown signals

---

## Testing Your API

### 1. **Health Check**
```bash
curl http://localhost:8081/health
```

### 2. **Login** (POST /v1/auth/login)
```bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user@example.com",
    "password": "test123"
  }'
```

### 3. **Signup** (POST /v1/auth/signup)
```bash
curl -X POST http://localhost:8081/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "method": "email",
    "email": "newuser@example.com",
    "password": "Secret123!"
  }'
```

### 4. **Onboard Basic** (POST /v1/onboard/basic)
```bash
curl -X POST http://localhost:8081/v1/onboard/basic \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "dob": "1999-05-12"
  }'
```

### 5. **Verify OTP** (POST /v1/auth/verify)
```bash
curl -X POST http://localhost:8081/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "method": "email",
    "code": "123456"
  }'
```

---

## All Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/v1/onboard/basic` | Submit name & DOB |
| POST | `/v1/onboard/profile` | Submit username & avatar |
| POST | `/v1/onboard/interests` | Save interests |
| POST | `/v1/auth/signup` | Create account |
| POST | `/v1/auth/verify` | Verify OTP |
| POST | `/v1/auth/login` | Login |
| POST | `/v1/auth/refresh` | Refresh token |
| POST | `/v1/auth/logout` | Logout |
| POST | `/v1/auth/resend` | Resend OTP |
| POST | `/v1/upload/avatar` | Upload avatar |

---

## Project Structure

```
cmd/api/
  main.go                      # Entry point - starts server

internal/transport/
  dtos/                        # Data Transfer Objects (API contracts)
    auth.go                    # Auth DTOs
    onboard.go                 # Onboarding DTOs
    upload.go                  # Upload DTOs
    common.go                  # Shared DTOs & errors
  
  http/
    router.go                  # URL routing
    handlers/
      auth_handler.go          # Auth endpoints
      onboard_handler.go       # Onboarding endpoints
      upload_handler.go        # Upload endpoints
```

---

## Next Steps (Learning Path)

### Step 1: Add Validation ✅
Currently, validation is basic. You can add:
- Email format validation
- Password strength checks
- Custom validation rules

### Step 2: Connect Database
Replace dummy responses with real data from Postgres.

### Step 3: Add JWT Authentication
Generate real JWT tokens instead of dummy strings.

### Step 4: Add Middleware
- Authentication middleware (check JWT tokens)
- Logging middleware
- Rate limiting

### Step 5: Add Tests
Write unit tests for your handlers.

---

## Common Commands

```bash
# Run server
make run

# Stop server
pkill -f "cmd/api/main.go"

# Check if server is running
curl http://localhost:8081/health

# View server logs
# (logs appear in terminal where you ran 'make run')
```

---

## Understanding the Flow

```
1. Client sends HTTP request
   ↓
2. Gin router matches URL to handler
   ↓
3. Handler parses JSON into DTO
   ↓
4. Handler validates input (binding tags)
   ↓
5. Handler processes (dummy logic for now)
   ↓
6. Handler returns JSON response
```

---

## Questions to Explore

1. **What happens if I send invalid JSON?**
   Try it! Send `{"invalid": }` and see the error response.

2. **What if I skip a required field?**
   Remove `"password"` from login and see what happens.

3. **How do I change the port?**
   Edit `main.go` line: `Addr: ":8081"` → `Addr: ":3000"`

4. **Can I add a new endpoint?**
   Yes! Add to handler, then add route in `router.go`.

---

Happy coding! 🎉

If you have questions about any Go concepts, feel free to ask!


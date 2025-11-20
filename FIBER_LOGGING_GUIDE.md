# 🔥 Fiber Logging Guide

Your Fiber app now has **automatic request logging** + **panic recovery**!

---

## ✅ What's Added

### **1. Logger Middleware**
Logs every HTTP request automatically:

```
[2025-11-21 02:50:15] 200 - GET /health (234μs)
[2025-11-21 02:50:16] 200 - POST /v1/auth/login (1.2ms)
[2025-11-21 02:50:17] 404 - GET /notfound (156μs)
```

### **2. Recovery Middleware**
Catches panics and prevents server crashes:

```
[Recovery] panic recovered: runtime error: invalid memory address
```

---

## 📝 Current Configuration

```go
// Logger format
app.Use(logger.New(logger.Config{
    Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
    TimeFormat: "2006-01-02 15:04:05",
    TimeZone:   "Local",
}))

// Recovery from panics
app.Use(recover.New())
```

---

## 🎨 Custom Log Formats

### **Minimal**
```go
app.Use(logger.New(logger.Config{
    Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
}))
```

**Output:**
```
15:04:05 | 200 | 234μs | GET /health
```

---

### **Detailed (Production)**
```go
app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
    TimeFormat: "2006-01-02 15:04:05",
}))
```

**Output:**
```
[2025-11-21 02:50:15] 200 - GET /health - 127.0.0.1 - 234μs
```

---

### **JSON Format (for log aggregators)**
```go
app.Use(logger.New(logger.Config{
    Format: `{"time":"${time}","status":${status},"method":"${method}","path":"${path}","latency":"${latency}","ip":"${ip}"}` + "\n",
}))
```

**Output:**
```json
{"time":"2025-11-21 02:50:15","status":200,"method":"GET","path":"/health","latency":"234μs","ip":"127.0.0.1"}
```

---

## 📊 Available Log Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `${time}` | Request timestamp | `15:04:05` |
| `${status}` | HTTP status code | `200` |
| `${method}` | HTTP method | `GET` |
| `${path}` | Request path | `/v1/auth/login` |
| `${latency}` | Request duration | `1.234ms` |
| `${ip}` | Client IP address | `127.0.0.1` |
| `${ua}` | User agent | `curl/7.64.1` |
| `${error}` | Error message (if any) | `not found` |
| `${body}` | Request body | `{"email":"..."}` |
| `${queryParams}` | Query parameters | `?page=1&limit=10` |
| `${header:name}` | Request header | `${header:Authorization}` |
| `${resBody}` | Response body | `{"status":"ok"}` |

---

## 🎯 Different Logging Strategies

### **1. Console Only (Development)**
```go
app.Use(logger.New()) // Uses default format
```

### **2. File Logging**
```go
file, _ := os.OpenFile("./logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

app.Use(logger.New(logger.Config{
    Output: file,
}))
```

### **3. Both Console + File**
```go
file, _ := os.OpenFile("./logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
multiWriter := io.MultiWriter(os.Stdout, file)

app.Use(logger.New(logger.Config{
    Output: multiWriter,
}))
```

### **4. Conditional Logging (Skip Health Checks)**
```go
app.Use(logger.New(logger.Config{
    Next: func(c *fiber.Ctx) bool {
        return c.Path() == "/health" // Skip logging /health
    },
}))
```

---

## 🔥 Advanced: Custom Logger

### **Log Only Errors (4xx, 5xx)**
```go
app.Use(logger.New(logger.Config{
    Next: func(c *fiber.Ctx) bool {
        return c.Response().StatusCode() < 400 // Skip 2xx, 3xx
    },
}))
```

### **Custom Colors**
```go
import "github.com/gofiber/fiber/v2/middleware/logger"

app.Use(logger.New(logger.Config{
    Format: "${time} | ${status} | ${latency} | ${method} ${path}\n",
    CustomTags: map[string]logger.LogFunc{
        "status": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
            status := c.Response().StatusCode()
            var color string
            
            if status >= 200 && status < 300 {
                color = "\033[32m" // Green
            } else if status >= 400 && status < 500 {
                color = "\033[33m" // Yellow
            } else if status >= 500 {
                color = "\033[31m" // Red
            }
            
            return output.WriteString(fmt.Sprintf("%s%d\033[0m", color, status))
        },
    },
}))
```

---

## 📁 Structured Logging with External Libraries

### **Using Zerolog**
```bash
go get github.com/rs/zerolog/log
```

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// Setup
zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// In handlers
func (h *Handler) Login(c *fiber.Ctx) error {
    log.Info().
        Str("method", c.Method()).
        Str("path", c.Path()).
        Str("ip", c.IP()).
        Msg("Login attempt")
    
    // ... handler logic
}
```

### **Using Zap**
```bash
go get go.uber.org/zap
```

```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("request received",
    zap.String("method", "POST"),
    zap.String("path", "/v1/auth/login"),
    zap.Int("status", 200),
)
```

---

## 🚨 Error Logging Examples

### **Log Errors to File**
```go
// Create logs directory
os.MkdirAll("./logs", os.ModePerm)

errorFile, _ := os.OpenFile("./logs/errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

// In handlers
func (h *Handler) Login(c *fiber.Ctx) error {
    if err != nil {
        errorFile.WriteString(fmt.Sprintf("[%s] ERROR: %v\n", time.Now().Format(time.RFC3339), err))
        return c.Status(500).JSON(fiber.Map{"error": "Internal error"})
    }
}
```

### **Custom Error Logger Middleware**
```go
app.Use(func(c *fiber.Ctx) error {
    err := c.Next()
    
    if err != nil {
        log.Printf("[ERROR] %s %s - %v", c.Method(), c.Path(), err)
        return err
    }
    
    return nil
})
```

---

## 📊 Log Rotation (Production)

```bash
go get gopkg.in/natefinch/lumberjack.v2
```

```go
import "gopkg.in/natefinch/lumberjack.v2"

app.Use(logger.New(logger.Config{
    Output: &lumberjack.Logger{
        Filename:   "./logs/app.log",
        MaxSize:    10,    // megabytes
        MaxBackups: 3,     // number of backups
        MaxAge:     28,    // days
        Compress:   true,  // compress old logs
    },
}))
```

---

## 🎯 Recommended Setup

### **Development**
```go
app.Use(logger.New(logger.Config{
    Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
    TimeFormat: "15:04:05",
}))
```

### **Production**
```go
// File logging with rotation
logFile := &lumberjack.Logger{
    Filename:   "./logs/app.log",
    MaxSize:    100, // MB
    MaxBackups: 5,
    MaxAge:     30, // days
    Compress:   true,
}

app.Use(logger.New(logger.Config{
    Format: `{"time":"${time}","status":${status},"method":"${method}","path":"${path}","ip":"${ip}","latency":"${latency}"}` + "\n",
    Output: logFile,
}))
```

---

## 🔍 Monitoring & Debugging

### **Request ID Tracking**
```go
import "github.com/gofiber/fiber/v2/middleware/requestid"

app.Use(requestid.New())

app.Use(logger.New(logger.Config{
    Format: "[${time}] ${locals:requestid} ${status} - ${method} ${path}\n",
}))
```

**Output:**
```
[15:04:05] abc-123-def-456 200 - GET /health
```

---

## 📝 Current Logs You'll See

When you run `make run`, you'll see:

```bash
Server starting on :8084

# Every request is logged:
[2025-11-21 02:50:15] 200 - GET /health (234μs)
[2025-11-21 02:50:16] 200 - POST /v1/auth/login (1.2ms)
[2025-11-21 02:50:17] 200 - POST /v1/onboard/basic (856μs)
[2025-11-21 02:50:18] 404 - GET /notfound (123μs)
```

---

## 🎉 Benefits

✅ **See all requests** in real-time  
✅ **Track performance** (latency)  
✅ **Debug issues** easily  
✅ **Monitor API usage**  
✅ **Automatic panic recovery** (no crashes)  

Your Fiber app is now production-ready with logging! 🔥


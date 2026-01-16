# Backend Architecture Improvement Proposal (Stdlib Only)

## Philosophy: Zero External Dependencies

Using only Go's standard library forces better design decisions and keeps the codebase lean, maintainable, and free from dependency hell.

---

## Proposed Improvements (Stdlib Only)

### 1. Service Layer Pattern

**Extract business logic from handlers and stores into a dedicated service layer.**

```
internal/
  domain/              # Business entities
    user.go
    vacation.go
    errors.go
  service/             # Business logic
    user_service.go
    vacation_service.go
  store/               # Data access (current structure)
    user/
    vacation/
  handler/             # HTTP layer
    api/
    web/
```

**Example Service**:
```go
// internal/service/vacation_service.go
type VacationService struct {
    vacationStore types.VacationStore
    userStore     types.UserStore
}

func (s *VacationService) CreateVacation(ctx context.Context, req *types.Vacation) (int, error) {
    // Business logic here
    user, err := s.userStore.FindById(req.PersonId)
    if err != nil {
        return 0, fmt.Errorf("user not found: %w", err)
    }
    
    if err := s.validateLeaveBalance(user, req.DaysUsed); err != nil {
        return 0, err
    }
    
    return s.vacationStore.CreateVacation(req)
}
```

**Benefits**: Business logic reusable, testable, and independent of HTTP.

---

### 2. Structured Error Handling

**Use custom error types with error wrapping (stdlib `errors` package).**

```go
// internal/domain/errors.go
type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// Predefined errors
var (
    ErrUserNotFound      = &AppError{Code: "USER_NOT_FOUND", Message: "user not found"}
    ErrInsufficientLeave = &AppError{Code: "INSUFFICIENT_LEAVE", Message: "insufficient leave balance"}
    ErrInvalidInput      = &AppError{Code: "INVALID_INPUT", Message: "invalid input"}
)

// Helper to create error with context
func NewAppError(code, message string, err error) *AppError {
    return &AppError{Code: code, Message: message, Err: err}
}
```

**Usage**:
```go
if user == nil {
    return 0, ErrUserNotFound
}

// With context
if days < 0 {
    return 0, NewAppError("INVALID_INPUT", "days cannot be negative", nil)
}
```

---

### 3. Input Validation (Stdlib Only)

**Manual validation with clear error messages.**

```go
// internal/domain/validation.go
type ValidationError struct {
    Field   string
    Message string
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
    var sb strings.Builder
    for _, err := range v {
        sb.WriteString(fmt.Sprintf("%s: %s; ", err.Field, err.Message))
    }
    return sb.String()
}

// Validators
func ValidateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return &ValidationError{Field: "email", Message: "invalid email format"}
    }
    return nil
}

func ValidateUser(u *types.User) error {
    var errs ValidationErrors
    
    if len(u.FirstName) < 2 {
        errs = append(errs, ValidationError{"first_name", "must be at least 2 characters"})
    }
    if u.Age < 18 || u.Age > 100 {
        errs = append(errs, ValidationError{"age", "must be between 18 and 100"})
    }
    if err := ValidateEmail(u.Email); err != nil {
        errs = append(errs, err.(ValidationError))
    }
    
    if len(errs) > 0 {
        return errs
    }
    return nil
}
```

---

### 4. Context Propagation

**Use `context.Context` throughout the stack.**

```go
// Update all store methods
type UserStore interface {
    FindById(ctx context.Context, id int) (*User, error)
    CreateUser(ctx context.Context, user *User) (int, error)
    // ...
}

// In handlers, create context with timeout
func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    id := extractID(r)
    user, err := h.userStore.FindById(ctx, id)
    // ...
}
```

**Benefits**: Timeout propagation, request cancellation, tracing support.

---

### 5. Enhanced Middleware (Stdlib)

```go
// middleware/recovery.go
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// middleware/request_id.go
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateID() // Use crypto/rand
        }
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// middleware/cors.go
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

### 6. Structured Logging (Stdlib `log/slog`)

**Go 1.21+ includes `log/slog` in stdlib.**

```go
// configs/logger.go
func NewLogger() *slog.Logger {
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}

// Usage
logger.Info("user created",
    "user_id", userID,
    "email", user.Email,
    "request_id", ctx.Value("request_id"),
)

logger.Error("database error",
    "error", err,
    "operation", "create_user",
)
```

---

### 7. Response Helpers (Stdlib)

```go
// pkg/response/response.go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{Success: true, Data: data})
}

func Error(w http.ResponseWriter, status int, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    var appErr *AppError
    if errors.As(err, &appErr) {
        json.NewEncoder(w).Encode(Response{
            Success: false,
            Error:   &ErrorInfo{Code: appErr.Code, Message: appErr.Message},
        })
        return
    }
    
    json.NewEncoder(w).Encode(Response{
        Success: false,
        Error:   &ErrorInfo{Code: "INTERNAL_ERROR", Message: err.Error()},
    })
}
```

---

### 8. Configuration (Stdlib)

```go
// configs/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
}

type ServerConfig struct {
    Host string
    Port string
}

type DatabaseConfig struct {
    User     string
    Password string
    Address  string
    Name     string
}

func Load() (*Config, error) {
    cfg := &Config{
        Server: ServerConfig{
            Host: getEnv("HOST", "localhost"),
            Port: getEnv("PORT", ":8080"),
        },
        Database: DatabaseConfig{
            User:     getEnv("DB_USER", "portal"),
            Password: getEnv("DB_PASSWORD", "password123"),
            Address:  getEnv("DB_ADDRESS", "127.0.0.1:3307"),
            Name:     getEnv("DB_NAME", "vacation_tool"),
        },
    }
    
    return cfg, cfg.Validate()
}

func (c *Config) Validate() error {
    if c.Database.Address == "" {
        return errors.New("DB_ADDRESS is required")
    }
    return nil
}
```

---

### 9. Testing Strategy (Stdlib)

```go
// service/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Mock store
    mockStore := &MockUserStore{
        createFunc: func(ctx context.Context, u *types.User) (int, error) {
            return 1, nil
        },
    }
    
    service := &UserService{store: mockStore}
    
    user := &types.User{FirstName: "John", LastName: "Doe", Age: 30, Email: "john@example.com"}
    id, err := service.CreateUser(context.Background(), user)
    
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if id != 1 {
        t.Errorf("expected id 1, got %d", id)
    }
}

// handler/user_handler_test.go
func TestHandleGetUser(t *testing.T) {
    handler := &Handler{/* ... */}
    
    req := httptest.NewRequest("GET", "/users/1", nil)
    rec := httptest.NewRecorder()
    
    handler.HandleGetUser(rec, req)
    
    if rec.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rec.Code)
    }
}
```

---

### 10. Health Checks (Stdlib)

```go
// handler/health.go
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()
    
    // Check database
    if err := h.db.PingContext(ctx); err != nil {
        response.JSON(w, http.StatusServiceUnavailable, map[string]string{
            "status": "unhealthy",
            "error":  err.Error(),
        })
        return
    }
    
    response.JSON(w, http.StatusOK, map[string]string{
        "status": "healthy",
        "uptime": time.Since(startTime).String(),
    })
}
```

---

## Implementation Priority

### Phase 1: Quick Wins (1-2 days)
1. ✅ Add recovery middleware
2. ✅ Add request ID middleware
3. ✅ Implement structured errors
4. ✅ Add response helpers
5. ✅ Add health check endpoint

### Phase 2: Core Refactoring (3-5 days)
1. ✅ Extract service layer
2. ✅ Add input validation
3. ✅ Add context propagation
4. ✅ Implement structured logging (slog)

### Phase 3: Testing & Polish (2-3 days)
1. ✅ Write unit tests for services
2. ✅ Write integration tests for stores
3. ✅ Write HTTP tests for handlers
4. ✅ Add documentation

---

## Benefits of Stdlib-Only Approach

✅ **Zero dependency management** - No version conflicts, no supply chain risks
✅ **Faster builds** - No external dependencies to download
✅ **Smaller binaries** - Only what you use
✅ **Long-term stability** - Stdlib is backward compatible
✅ **Better understanding** - Forces you to understand the patterns
✅ **Easier onboarding** - New developers only need to know Go

---

## Summary

All proposed improvements use **only Go's standard library**:
- `context` - Request context and cancellation
- `errors` - Error wrapping and type checking
- `log/slog` - Structured logging (Go 1.21+)
- `net/http` - HTTP server and testing
- `encoding/json` - JSON handling
- `testing` - Test framework
- `database/sql` - Database access

The architecture becomes cleaner, more maintainable, and completely self-contained.

**Ready to implement?** I recommend starting with Phase 1 (middleware + errors + responses) as it's non-breaking and provides immediate value.

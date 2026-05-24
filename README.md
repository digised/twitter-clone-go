# Clean Architecture

> What each file does, why it exists, and how they work together

---

## model.go — Defines what your data looks like

Before you can store or send any data, you need to describe its shape. model.go is where you define that shape. Think of it as a blueprint: a User has an ID, a username, an email, and a creation date. Every other file references this blueprint, but this file does not depend on anyone else. It sits at the very bottom of the chain.

**What is inside:**
- The main structs: User, Tweet, Follow, etc.
- Field tags that tell Go how to read/write JSON and database columns
- No logic here — only data shapes

```go
type User struct {
    ID        uuid.UUID `json:"id"       db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"-"        db:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

> **Note:** `json:"-"` on Email means it will never appear in API responses — a small tag that prevents leaking private data.

---

## dto.go — Defines what goes in and what comes out of the API

DTO stands for Data Transfer Object. It is the contract between your API and the outside world. When a user signs up, what fields must they send? When you reply, what fields do you include? By keeping this separate from model.go, you can change your database schema without breaking the API, and change the API without touching the database. They are independent on purpose.

**What is inside:**
- Request structs: RegisterRequest, LoginRequest — what the client sends
- Response structs: UserResponse — what the server sends back (no passwords!)
- Validation tags that automatically reject bad input before it reaches your logic

```go
type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
}
```

> **Note:** The validate tags are read by go-playground/validator. If a field is missing or wrong, the request is rejected automatically — no manual if-checks needed.

---

## repository.go — The only place that talks to the database

All SQL queries live here and only here. No other file is allowed to write a SELECT or INSERT. This is called the Repository Pattern. The big benefit: if you ever switch from PostgreSQL to a different database, you only rewrite this one file. The rest of the code does not care what database you are using — it just calls repository functions and gets data back.

**What is inside:**
- An interface that lists what actions are possible: FindByID, Create, Update, Delete
- A struct that actually runs the SQL using sqlx or pgx
- No business rules here — just fetch, save, update, delete

```go
type UserRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
}

type userRepo struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepo{db: db}
}
```

> **Note:** Defining an interface (not just a struct) makes it easy to swap out the real database for a fake one during testing. No real database needed to run tests.

---

## service.go — Where all the decisions are made

This is the brain of the domain. It answers questions like: Is this email already taken? Should I hash the password before saving? Do I need to send a notification after this action? These are business rules — they do not belong in handler.go (which is about HTTP) or repository.go (which is about data). They belong here, where they can be tested and changed in isolation.

**What is inside:**
- Calls the repository to read or write data
- Enforces rules: duplicate checks, permission checks, etc.
- Coordinates between multiple repositories when needed
- Publishes events to a queue (like Kafka) when something important happens

```go
func (s *userService) Register(
    ctx context.Context, req RegisterRequest,
) (*UserResponse, error) {
    existing, _ := s.repo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }
    hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
    user := &User{Email: req.Email, PasswordHash: string(hashed)}
    s.repo.Create(ctx, user)
    return toResponse(user), nil
}
```

> **Note:** The service never imports a database driver. It only knows about the repository interface. This keeps business logic clean and portable.

---

## handler.go — The front door, handles HTTP requests

handler.go is the first thing that runs when a request arrives. Its job is simple: read the request, call the service, send back a response. It does not contain any database queries. It does not contain any business rules. If you ever move from a REST API to gRPC or GraphQL, this is the only file you would need to replace.

**What is inside:**
- Reads JSON from the request body and fills a DTO struct
- Returns 400 if the input is invalid, before calling any service
- Calls the correct service method and waits for the result
- Sends back the right HTTP status code: 201 Created, 200 OK, 404 Not Found, etc.

```go
func (h *UserHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    res, err := h.service.Register(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, res)
}
```

> **Note:** handler.go is not allowed to touch the database directly. Ever. That line must never be crossed.

---

## service_test.go — Proof that your logic actually works

Because service.go only talks to an interface (not a real database), you can create a fake repository in tests that behaves however you want. Want to test what happens when the email already exists? Tell the fake repository to return a user. No database setup, no Docker, no SQL — just fast, reliable tests that run in seconds and catch bugs before they reach production.

**What is inside:**
- A mock repository that returns whatever the test tells it to
- Happy path tests: everything works as expected
- Error path tests: duplicate email, wrong password, missing fields
- Table-driven tests: test 10 cases with one loop instead of 10 functions

```go
func TestRegister_DuplicateEmail(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockRepo.On("FindByEmail", mock.Anything, "a@b.com").
        Return(&User{}, nil) // pretend user already exists

    svc := NewUserService(mockRepo)
    _, err := svc.Register(ctx, RegisterRequest{Email: "a@b.com"})

    assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}
```

> **Note:** Run all tests with: `go test ./...` — This is also what your CI pipeline runs on every push.

---

## How a request flows through all the layers

When a user sends `POST /users/register`, here is the exact path the code takes:

1. **handler.go** — Receives the HTTP request, parses JSON, validates the DTO
2. **service.go** — Checks business rules: is this email already taken?
3. **repository.go** — Runs INSERT query, saves the new user in PostgreSQL
4. **service.go** — Publishes a "user.registered" event to Kafka (optional)
5. **handler.go** — Returns HTTP 201 Created with a UserResponse body

---

## Quick reference

| File | Layer | Main responsibility | Touches DB? |
|---|---|---|---|
| `model.go` | Domain | Defines data shapes (structs) | No |
| `dto.go` | API | Input/output contracts for the API | No |
| `repository.go` | Data | All SQL queries, nothing else | Yes |
| `service.go` | Logic | Business rules and coordination | No |
| `handler.go` | HTTP | Parse requests, send responses | No |
| `service_test.go` | Test | Tests for service logic (no real DB) | No |

---

> This structure follows Clean Architecture principles used in large Go codebases at companies like Uber, Cloudflare, and others. Each file has one job. Change one thing in one place. Test each part in isolation.

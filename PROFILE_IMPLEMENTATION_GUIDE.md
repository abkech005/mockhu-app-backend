# User Profiles Implementation Guide - Step by Step

**Follow this guide to implement the User Profiles & Following feature yourself.**

---

## 📋 DAY 1: Database Migration & Setup (4-5 hours)

### Step 1: Create Migration Files

#### 1.1 Create UP Migration
Create file: `migrations/000014_add_profile_privacy_fields.up.sql`

Add the following columns to the `users` table:
- `bio TEXT` - User biography (max 500 chars)
- `institution_id UUID` - Link to institution (for future)
- `who_can_message VARCHAR(20) DEFAULT 'everyone'` - Privacy: who can send DMs
- `who_can_see_posts VARCHAR(20) DEFAULT 'everyone'` - Privacy: who can see posts
- `show_followers_list BOOLEAN DEFAULT true` - Privacy: show follower list
- `show_following_list BOOLEAN DEFAULT true` - Privacy: show following list

Add these constraints:
- Check `who_can_message` is one of: 'everyone', 'followers', 'none'
- Check `who_can_see_posts` is one of: 'everyone', 'followers', 'none'
- Check `bio` length is max 500 characters

Add these indexes:
- Unique index on `LOWER(username)` for case-insensitive uniqueness
- Index on `institution_id` for faster joins

#### 1.2 Create DOWN Migration
Create file: `migrations/000014_add_profile_privacy_fields.down.sql`

Should drop:
- All indexes created above
- All constraints created above
- All columns added above

---

### Step 2: Update User Model

Open: `internal/app/auth/model.go`

Add these new fields to the `User` struct:
```go
// Profile fields
Bio           string  `json:"bio,omitempty"`
InstitutionID *string `json:"institution_id,omitempty"`

// Privacy settings
WhoCanMessage      string `json:"who_can_message"`
WhoCanSeePosts     string `json:"who_can_see_posts"`
ShowFollowersList  bool   `json:"show_followers_list"`
ShowFollowingList  bool   `json:"show_following_list"`
```

**Tip:** Add these after `AvatarURL` and before `IsActive`

---

### Step 3: Create Profile DTOs

Open: `internal/app/auth/dto.go`

Add at the end of the file:

#### 3.1 ProfileStats struct
```go
type ProfileStats struct {
    PostsCount     int `json:"posts_count"`
    FollowersCount int `json:"followers_count"`
    FollowingCount int `json:"following_count"`
}
```

#### 3.2 InstitutionInfo struct
```go
type InstitutionInfo struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

#### 3.3 ProfileResponse (for GET /v1/users/:userId/profile)
```go
type ProfileResponse struct {
    ID                     string           `json:"id"`
    Username               string           `json:"username"`
    FirstName              string           `json:"first_name"`
    LastName               string           `json:"last_name"`
    AvatarURL              string           `json:"avatar_url,omitempty"`
    Bio                    string           `json:"bio,omitempty"`
    Institution            *InstitutionInfo `json:"institution,omitempty"`
    Stats                  ProfileStats     `json:"stats"`
    IsFollowing            bool             `json:"is_following"`
    IsFollowedBy           bool             `json:"is_followed_by"`
    MutualConnectionsCount int              `json:"mutual_connections_count"`
    CreatedAt              string           `json:"created_at"`
}
```

#### 3.4 PrivacySettings struct
```go
type PrivacySettings struct {
    WhoCanMessage     string `json:"who_can_message"`
    WhoCanSeePosts    string `json:"who_can_see_posts"`
    ShowFollowersList bool   `json:"show_followers_list"`
    ShowFollowingList bool   `json:"show_following_list"`
}
```

#### 3.5 OwnProfileResponse (for GET /v1/users/me/profile)
```go
type OwnProfileResponse struct {
    ID                  string           `json:"id"`
    Username            string           `json:"username"`
    FirstName           string           `json:"first_name"`
    LastName            string           `json:"last_name"`
    Email               string           `json:"email"`
    Phone               string           `json:"phone"`
    DateOfBirth         string           `json:"date_of_birth"`
    AvatarURL           string           `json:"avatar_url,omitempty"`
    Bio                 string           `json:"bio,omitempty"`
    Institution         *InstitutionInfo `json:"institution,omitempty"`
    Stats               ProfileStats     `json:"stats"`
    PrivacySettings     PrivacySettings  `json:"privacy_settings"`
    EmailVerified       bool             `json:"email_verified"`
    PhoneVerified       bool             `json:"phone_verified"`
    OnboardingCompleted bool             `json:"onboarding_completed"`
    CreatedAt           string           `json:"created_at"`
}
```

#### 3.6 UpdateProfileRequest (for PUT /v1/users/me/profile)
```go
type UpdateProfileRequest struct {
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
    Username  string `json:"username,omitempty"`
    Bio       string `json:"bio,omitempty"`
}
```

#### 3.7 UpdatePrivacyRequest (for PUT /v1/users/me/privacy)
```go
type UpdatePrivacyRequest struct {
    WhoCanMessage     string `json:"who_can_message,omitempty"`
    WhoCanSeePosts    string `json:"who_can_see_posts,omitempty"`
    ShowFollowersList *bool  `json:"show_followers_list,omitempty"`
    ShowFollowingList *bool  `json:"show_following_list,omitempty"`
}
```

**Note:** Use `*bool` (pointer) so we can distinguish between false and not provided.

#### 3.8 MutualConnectionUser struct
```go
type MutualConnectionUser struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    AvatarURL string `json:"avatar_url,omitempty"`
}
```

#### 3.9 MutualConnectionsResponse
```go
type MutualConnectionsResponse struct {
    MutualConnections []*MutualConnectionUser `json:"mutual_connections"`
    Pagination        PaginationInfo          `json:"pagination"`
}
```

#### 3.10 MutualConnectionsCountResponse
```go
type MutualConnectionsCountResponse struct {
    UserID                 string `json:"user_id"`
    MutualConnectionsCount int    `json:"mutual_connections_count"`
}
```

#### 3.11 PaginationInfo struct (if not exists)
```go
type PaginationInfo struct {
    Page       int `json:"page"`
    TotalPages int `json:"total_pages"`
    TotalItems int `json:"total_items"`
    Limit      int `json:"limit"`
}
```

#### 3.12 AvatarUploadResponse
```go
type AvatarUploadResponse struct {
    AvatarURL string `json:"avatar_url"`
    Message   string `json:"message"`
}
```

---

### Step 4: Test Migration

#### 4.1 Run Migration
```bash
make migrate
# or
migrate -path migrations -database "postgresql://user:pass@localhost:5432/mockhu?sslmode=disable" up
```

#### 4.2 Verify in Database
```bash
# Connect to database
psql -d mockhu

# Check if columns exist
\d users

# You should see all 6 new columns
```

#### 4.3 Test Rollback
```bash
make migrate-down
# or
migrate -path migrations -database "..." down 1
```

#### 4.4 Run Migration Again
```bash
make migrate
```

---

### Step 5: Verify Code Compiles

```bash
cd /Users/abkech/Documents/mockhu-app-backend
go build ./cmd/api/main.go
```

**If there are errors:**
- Check User struct field names match exactly
- Check all DTOs are properly formatted
- Run `go fmt ./...` to format code

---

## ✅ Day 1 Completion Checklist

- [ ] Created `migrations/000014_add_profile_privacy_fields.up.sql`
- [ ] Created `migrations/000014_add_profile_privacy_fields.down.sql`
- [ ] Updated `User` struct in `model.go` with 6 new fields
- [ ] Added 12 new DTOs to `dto.go`
- [ ] Migration runs successfully (UP)
- [ ] Migration rollback works (DOWN)
- [ ] Verified columns exist in database
- [ ] Code compiles without errors
- [ ] Ran `go fmt ./...` to format code

---

## 🎯 What You've Accomplished

After completing Day 1, you have:
1. ✅ Database schema updated with profile and privacy fields
2. ✅ User model updated to include new fields
3. ✅ All DTOs created for API requests/responses
4. ✅ Migration tested (up and down)

---

## 📋 Next Steps (Day 2)

Tomorrow you'll implement:
- Repository methods for profile operations
- Service layer for business logic
- 2 Handler endpoints:
  - `GET /v1/users/:userId/profile` (public profile)
  - `GET /v1/users/me/profile` (own profile)
- Routes registration
- Testing

---

## 🆘 Troubleshooting

### Migration Fails
**Error:** Column already exists
- **Solution:** The migration ran partially. Use `migrate force 13` then `migrate up`

**Error:** Dirty database version
- **Solution:** `migrate -path migrations -database "..." force 14`

### Code Won't Compile
**Error:** Undefined field
- **Solution:** Check field names match exactly (case-sensitive)

**Error:** Duplicate declaration
- **Solution:** Make sure you didn't add DTOs twice

### Database Connection Issues
- **Solution:** Ensure Docker is running: `docker-compose up -d`
- Check connection string in Makefile or env vars

---

## 📝 Quick Reference

### File Locations
- Migrations: `migrations/000014_*`
- User model: `internal/app/auth/model.go`
- DTOs: `internal/app/auth/dto.go`

### Commands
```bash
# Run migration
make migrate

# Rollback migration
make migrate-down

# Check migration status
migrate -path migrations -database "..." version

# Build code
go build ./cmd/api/main.go

# Format code
go fmt ./...
```

---

**Time for Day 1:** 4-5 hours  
**Status:** Ready to start!  

Good luck! 🚀


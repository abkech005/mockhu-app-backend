# User Profile Feature - Implementation Checklist

**Complete this checklist at your own pace.**

---

## 📋 PHASE 1: Database Setup

### Migration Files
- [ ] Create `migrations/000014_add_profile_privacy_fields.up.sql`
  - [ ] Add `bio TEXT` column
  - [ ] Add `institution_id UUID` column
  - [ ] Add `who_can_message VARCHAR(20) DEFAULT 'everyone'` column
  - [ ] Add `who_can_see_posts VARCHAR(20) DEFAULT 'everyone'` column
  - [ ] Add `show_followers_list BOOLEAN DEFAULT true` column
  - [ ] Add `show_following_list BOOLEAN DEFAULT true` column
  - [ ] Add CHECK constraint on `who_can_message` (everyone/followers/none)
  - [ ] Add CHECK constraint on `who_can_see_posts` (everyone/followers/none)
  - [ ] Add CHECK constraint on `bio` length (max 500 chars)
  - [ ] Create unique index on `LOWER(username)`
  - [ ] Create index on `institution_id`

- [ ] Create `migrations/000014_add_profile_privacy_fields.down.sql`
  - [ ] Drop all indexes
  - [ ] Drop all constraints
  - [ ] Drop all columns

### Test Migration
- [ ] Run migration: `make migrate`
- [ ] Verify columns exist in database (check with `\d users`)
- [ ] Test rollback: `make migrate-down`
- [ ] Run migration again: `make migrate`

---

## 📋 PHASE 2: Models & DTOs

### Update User Model
- [ ] Open `internal/app/auth/model.go`
- [ ] Add `Bio string` field
- [ ] Add `InstitutionID *string` field
- [ ] Add `WhoCanMessage string` field
- [ ] Add `WhoCanSeePosts string` field
- [ ] Add `ShowFollowersList bool` field
- [ ] Add `ShowFollowingList bool` field

### Create DTOs
- [ ] Open `internal/app/auth/dto.go`
- [ ] Add `ProfileStats` struct
- [ ] Add `InstitutionInfo` struct
- [ ] Add `ProfileResponse` struct (for public profile)
- [ ] Add `OwnProfileResponse` struct (for own profile)
- [ ] Add `UpdateProfileRequest` struct
- [ ] Add `PrivacySettings` struct
- [ ] Add `UpdatePrivacyRequest` struct
- [ ] Add `MutualConnectionUser` struct
- [ ] Add `MutualConnectionsResponse` struct
- [ ] Add `MutualConnectionsCountResponse` struct
- [ ] Add `PaginationInfo` struct (if not exists)
- [ ] Add `AvatarUploadResponse` struct

### Verify Code
- [ ] Run `go build ./cmd/api/main.go`
- [ ] Fix any compilation errors
- [ ] Run `go fmt ./...`

---

## 📋 PHASE 3: Profile Viewing (2 endpoints)

### Repository Layer
- [ ] Open `internal/app/auth/repository.go`
- [ ] Add `GetProfileByID(ctx, userID) (*User, error)` method to interface
- [ ] Open `internal/app/auth/repository_postgres.go`
- [ ] Implement `GetProfileByID` method
  - [ ] SELECT with new profile fields
  - [ ] Handle NULL values properly

### Service Layer - Get User Profile
- [ ] Open `internal/app/auth/service.go`
- [ ] Add `GetUserProfile(ctx, userID, currentUserID) (*ProfileResponse, error)` method
- [ ] Implement logic:
  - [ ] Get user from repository
  - [ ] Get posts count (query posts table)
  - [ ] Get followers count (use follow repo)
  - [ ] Get following count (use follow repo)
  - [ ] Check if current user is following target
  - [ ] Check if target is following current user
  - [ ] Calculate mutual connections count (if authenticated)
  - [ ] Build and return ProfileResponse

### Service Layer - Get Own Profile
- [ ] Add `GetOwnProfile(ctx, userID) (*OwnProfileResponse, error)` method
- [ ] Implement logic:
  - [ ] Get user from repository
  - [ ] Get stats (same as above)
  - [ ] Build and return OwnProfileResponse with all private fields

### Handler Layer
- [ ] Open `internal/app/auth/handler.go`
- [ ] Add `GetUserProfile` handler function
  - [ ] Get `userId` from URL params
  - [ ] Get current user ID from JWT (optional, use `c.Locals("user_id")`)
  - [ ] Call `service.GetUserProfile`
  - [ ] Handle errors (404 if user not found)
  - [ ] Return JSON response

- [ ] Add `GetOwnProfile` handler function
  - [ ] Get current user ID from JWT (required)
  - [ ] Call `service.GetOwnProfile`
  - [ ] Handle errors
  - [ ] Return JSON response

### Routes
- [ ] Open `internal/app/auth/routes.go`
- [ ] Add `GET /v1/users/:userId/profile` route (public, no auth)
- [ ] Add `GET /v1/users/me/profile` route (protected, requires auth)

### Testing
- [ ] Test GET user profile without authentication
- [ ] Test GET user profile with authentication (should show mutual connections)
- [ ] Test GET own profile (should show private fields)
- [ ] Test with non-existent user ID (should return 404)

---

## 📋 PHASE 4: Update Profile (1 endpoint)

### Repository Layer
- [ ] Open `internal/app/auth/repository.go`
- [ ] Add `UpdateProfile(ctx, userID, updates) error` method to interface
- [ ] Add `CheckUsernameExists(ctx, username, excludeUserID) (bool, error)` method
- [ ] Open `internal/app/auth/repository_postgres.go`
- [ ] Implement `UpdateProfile` method
  - [ ] Build dynamic UPDATE query
  - [ ] Update only provided fields
- [ ] Implement `CheckUsernameExists` method
  - [ ] Query with LOWER(username) for case-insensitive check
  - [ ] Exclude current user from check

### Service Layer
- [ ] Open `internal/app/auth/service.go`
- [ ] Add `UpdateProfile(ctx, userID, req) (*User, error)` method
- [ ] Add validation logic:
  - [ ] First name: 1-50 characters (if provided)
  - [ ] Last name: 1-50 characters (if provided)
  - [ ] Bio: max 500 characters (if provided)
  - [ ] Username: 3-30 chars, alphanumeric + underscore (if provided)
- [ ] Check username uniqueness (case-insensitive)
- [ ] Sanitize bio (strip HTML tags using bluemonday)
- [ ] Call repository to update
- [ ] Return updated user

### Handler Layer
- [ ] Open `internal/app/auth/handler.go`
- [ ] Add `UpdateProfile` handler function
  - [ ] Get current user ID from JWT
  - [ ] Parse `UpdateProfileRequest` from body
  - [ ] Call `service.UpdateProfile`
  - [ ] Handle errors (400 for validation, 409 for duplicate username)
  - [ ] Return updated profile

### Routes
- [ ] Open `internal/app/auth/routes.go`
- [ ] Add `PUT /v1/users/me/profile` route (protected)

### Testing
- [ ] Test update profile with valid data (all fields)
- [ ] Test update only bio
- [ ] Test update only username
- [ ] Test update only first/last name
- [ ] Test with duplicate username (should return 409)
- [ ] Test with too long bio (should return 400)
- [ ] Test with invalid username format (should return 400)

---

## 📋 PHASE 5: Avatar Upload & Delete (2 endpoints)

### Setup Image Processing
- [ ] Install imaging library: `go get github.com/disintegration/imaging`
- [ ] Install mimetype library: `go get github.com/gabriel-vasile/mimetype`
- [ ] Run `go mod tidy`

### Create Avatar Service
- [ ] Create directory: `internal/pkg/avatar/`
- [ ] Create `internal/pkg/avatar/avatar.go`
- [ ] Add `ProcessAvatar(fileBytes, filename) (string, error)` function
  - [ ] Validate file type (JPEG, PNG, WebP)
  - [ ] Validate file size (max 5MB)
  - [ ] Decode image
  - [ ] Resize to 400x400 (square, centered crop)
  - [ ] Generate unique filename (UUID + extension)
  - [ ] Save to storage directory (`storage/avatars/`)
  - [ ] Return file URL/path

- [ ] Add `DeleteAvatar(avatarURL) error` function
  - [ ] Extract filename from URL
  - [ ] Delete file from storage

### Storage Setup
- [ ] Create directory: `storage/avatars/`
- [ ] Add `storage/` to `.gitignore`

### Repository Layer
- [ ] Open `internal/app/auth/repository.go`
- [ ] Add `UpdateAvatar(ctx, userID, avatarURL) error` method to interface
- [ ] Open `internal/app/auth/repository_postgres.go`
- [ ] Implement `UpdateAvatar` method

### Service Layer
- [ ] Open `internal/app/auth/service.go`
- [ ] Add `UploadAvatar(ctx, userID, fileBytes, filename) (string, error)` method
  - [ ] Get current avatar URL from user
  - [ ] Process new avatar (resize, save)
  - [ ] Update database with new URL
  - [ ] Delete old avatar file (if exists)
  - [ ] Return new avatar URL

- [ ] Add `DeleteAvatar(ctx, userID) error` method
  - [ ] Get current avatar URL
  - [ ] Delete file from storage
  - [ ] Set avatar_url to NULL in database

### Handler Layer
- [ ] Open `internal/app/auth/handler.go`
- [ ] Add `UploadAvatar` handler function
  - [ ] Get current user ID from JWT
  - [ ] Parse multipart form
  - [ ] Get file from form field "avatar"
  - [ ] Read file bytes
  - [ ] Call `service.UploadAvatar`
  - [ ] Handle errors (400 for invalid file)
  - [ ] Return avatar URL

- [ ] Add `DeleteAvatar` handler function
  - [ ] Get current user ID from JWT
  - [ ] Call `service.DeleteAvatar`
  - [ ] Return success message

### Routes
- [ ] Open `internal/app/auth/routes.go`
- [ ] Add `POST /v1/users/me/avatar` route (protected)
- [ ] Add `DELETE /v1/users/me/avatar` route (protected)

### Testing
- [ ] Test upload JPEG image
- [ ] Test upload PNG image
- [ ] Test upload WebP image
- [ ] Test upload invalid file type (should fail)
- [ ] Test upload file > 5MB (should fail)
- [ ] Test upload replaces old avatar
- [ ] Verify old avatar file is deleted
- [ ] Test delete avatar
- [ ] Verify avatar URL is NULL after delete

---

## 📋 PHASE 6: Privacy Settings (2 endpoints)

### Repository Layer
- [ ] Open `internal/app/auth/repository.go`
- [ ] Add `GetPrivacySettings(ctx, userID) (*PrivacySettings, error)` method
- [ ] Add `UpdatePrivacySettings(ctx, userID, settings) error` method
- [ ] Open `internal/app/auth/repository_postgres.go`
- [ ] Implement `GetPrivacySettings` method
- [ ] Implement `UpdatePrivacySettings` method (dynamic UPDATE)

### Service Layer
- [ ] Open `internal/app/auth/service.go`
- [ ] Add `GetPrivacySettings(ctx, userID) (*PrivacySettings, error)` method
- [ ] Add `UpdatePrivacySettings(ctx, userID, req) (*PrivacySettings, error)` method
- [ ] Add validation:
  - [ ] `who_can_message` must be: everyone/followers/none
  - [ ] `who_can_see_posts` must be: everyone/followers/none

### Handler Layer
- [ ] Open `internal/app/auth/handler.go`
- [ ] Add `GetPrivacySettings` handler
  - [ ] Get current user ID from JWT
  - [ ] Call service
  - [ ] Return settings

- [ ] Add `UpdatePrivacySettings` handler
  - [ ] Get current user ID from JWT
  - [ ] Parse request body
  - [ ] Call service
  - [ ] Handle validation errors
  - [ ] Return updated settings

### Routes
- [ ] Open `internal/app/auth/routes.go`
- [ ] Add `GET /v1/users/me/privacy` route (protected)
- [ ] Add `PUT /v1/users/me/privacy` route (protected)

### Testing
- [ ] Test get privacy settings
- [ ] Test update all privacy settings
- [ ] Test update only who_can_message
- [ ] Test update only who_can_see_posts
- [ ] Test update only visibility settings
- [ ] Test with invalid values (should return 400)

---

## 📋 PHASE 7: Mutual Connections (2 endpoints)

### Repository Layer
- [ ] Open `internal/app/follow/repository.go` (or create in auth)
- [ ] Add `GetMutualConnections(ctx, user1ID, user2ID, limit, offset) ([]*User, error)` method
- [ ] Add `GetMutualConnectionsCount(ctx, user1ID, user2ID) (int, error)` method
- [ ] Open repository implementation file
- [ ] Implement `GetMutualConnections` method
  - [ ] SQL: Find users that BOTH users follow
  - [ ] Use INTERSECT or subqueries
  - [ ] Join with users table for user info
  - [ ] Add pagination (LIMIT/OFFSET)
- [ ] Implement `GetMutualConnectionsCount` method
  - [ ] Similar query but COUNT only

### Service Layer
- [ ] Open appropriate service file
- [ ] Add `GetMutualConnections(ctx, currentUserID, targetUserID, page, limit)` method
  - [ ] Validate pagination (page >= 1, limit 1-50)
  - [ ] Calculate offset
  - [ ] Call repository
  - [ ] Convert to response format
  - [ ] Add pagination info

- [ ] Add `GetMutualConnectionsCount(ctx, currentUserID, targetUserID) (int, error)` method
  - [ ] Call repository
  - [ ] Return count

### Handler Layer
- [ ] Add `GetMutualConnections` handler
  - [ ] Get current user ID from JWT
  - [ ] Get target user ID from URL params
  - [ ] Parse page/limit from query params
  - [ ] Call service
  - [ ] Return response with pagination

- [ ] Add `GetMutualConnectionsCount` handler
  - [ ] Get current user ID from JWT
  - [ ] Get target user ID from URL params
  - [ ] Call service
  - [ ] Return count

### Routes
- [ ] Add `GET /v1/users/:userId/mutual-connections` route (protected)
- [ ] Add `GET /v1/users/:userId/mutual-connections/count` route (protected)

### Testing
- [ ] Test mutual connections with users who have mutual friends
- [ ] Test mutual connections with users who have no mutual friends
- [ ] Test pagination (multiple pages)
- [ ] Test mutual connections count accuracy
- [ ] Test with self (should be 0 or empty)

---

## 📋 PHASE 8: Privacy Enforcement

### Update Follow Service
- [ ] Open `internal/app/follow/service.go`
- [ ] Update `GetFollowers` method
  - [ ] Check if owner's `show_followers_list` is false
  - [ ] If false and viewer is not owner, return empty or forbidden
  
- [ ] Update `GetFollowing` method
  - [ ] Check if owner's `show_following_list` is false
  - [ ] If false and viewer is not owner, return empty or forbidden

### Update Post Service (Prepare for Future)
- [ ] Open `internal/app/post/service.go`
- [ ] Add helper function `CanViewPosts(viewer, owner *User) bool`
  - [ ] If viewer is owner, return true
  - [ ] Check owner's `who_can_see_posts` setting
  - [ ] Return true/false based on privacy
- [ ] Add TODO comments where this check will be used in feed

### Testing
- [ ] Test follower list hidden when show_followers_list = false
- [ ] Test following list hidden when show_following_list = false
- [ ] Test owner can still see own lists
- [ ] Test privacy checks work correctly

---

## 📋 PHASE 9: Integration & Testing

### Create Test Script
- [ ] Create `test_profile.sh`
- [ ] Add test for: GET user profile (public)
- [ ] Add test for: GET user profile (authenticated)
- [ ] Add test for: GET own profile
- [ ] Add test for: Update profile (valid data)
- [ ] Add test for: Update profile (duplicate username)
- [ ] Add test for: Upload avatar (JPEG)
- [ ] Add test for: Upload avatar (invalid file)
- [ ] Add test for: Delete avatar
- [ ] Add test for: Get privacy settings
- [ ] Add test for: Update privacy settings
- [ ] Add test for: Get mutual connections
- [ ] Add test for: Get mutual connections count
- [ ] Add test for: Privacy enforcement
- [ ] Make script executable: `chmod +x test_profile.sh`

### Run All Tests
- [ ] Run test script: `./test_profile.sh`
- [ ] Fix any failing tests
- [ ] Verify all endpoints work correctly

---

## 📋 PHASE 10: Documentation

### Postman Collection
- [ ] Create `postman/Mockhu_Profile_API.postman_collection.json`
- [ ] Add "Profile Viewing" folder
  - [ ] GET User Profile
  - [ ] GET Own Profile
- [ ] Add "Profile Management" folder
  - [ ] PUT Update Profile
  - [ ] POST Upload Avatar
  - [ ] DELETE Avatar
- [ ] Add "Privacy Settings" folder
  - [ ] GET Privacy Settings
  - [ ] PUT Privacy Settings
- [ ] Add "Mutual Connections" folder
  - [ ] GET Mutual Connections List
  - [ ] GET Mutual Connections Count
- [ ] Add test scripts to each request
- [ ] Add example responses
- [ ] Test all endpoints in Postman

### Postman Environment
- [ ] Create `postman/Mockhu_Profile_API_Environment.postman_environment.json`
- [ ] Add `base_url` variable
- [ ] Add `access_token` variable
- [ ] Add `user_id` variable
- [ ] Add `target_user_id` variable

### API Documentation
- [ ] Create `postman/PROFILE_API_DOCUMENTATION.md`
- [ ] Document all 9 endpoints
- [ ] Add request/response examples
- [ ] Document privacy rules
- [ ] Add business rules
- [ ] Add error codes
- [ ] Add usage examples

### Update Project Docs
- [ ] Update `SECTIONS_STATUS.md`
  - [ ] Mark User Profiles as complete
  - [ ] Add endpoint count (9 endpoints)
  - [ ] Add completion date
- [ ] Update `MVP1_WORK_CHECKLIST.md`
  - [ ] Check off User Profiles feature
  - [ ] Update progress percentage

---

## 📋 PHASE 11: Code Quality & Cleanup

### Code Quality
- [ ] Run `go fmt ./...` to format all code
- [ ] Run `go vet ./...` to check for issues
- [ ] Check for any linter errors
- [ ] Add code comments where needed
- [ ] Remove any debug/test code
- [ ] Remove any TODO comments (or track them)

### Error Handling
- [ ] Verify all repository errors are handled
- [ ] Verify all service errors are handled
- [ ] Verify all handler errors return proper status codes
- [ ] Add proper error messages for users

### Performance Check
- [ ] Profile queries are using indexes
- [ ] No N+1 query problems
- [ ] Pagination works efficiently
- [ ] Image processing is reasonably fast

---

## 📋 PHASE 12: Final Review & Push

### Final Testing
- [ ] Test all 9 endpoints one more time
- [ ] Test error cases
- [ ] Test with different user types
- [ ] Verify privacy enforcement works

### Git Commit
- [ ] Review all changes: `git status`
- [ ] Stage migration files: `git add migrations/`
- [ ] Stage model changes: `git add internal/app/auth/model.go`
- [ ] Stage DTO changes: `git add internal/app/auth/dto.go`
- [ ] Stage repository changes: `git add internal/app/auth/repository*.go`
- [ ] Stage service changes: `git add internal/app/auth/service.go`
- [ ] Stage handler changes: `git add internal/app/auth/handler.go`
- [ ] Stage route changes: `git add internal/app/auth/routes.go`
- [ ] Stage avatar service: `git add internal/pkg/avatar/`
- [ ] Stage follow service updates: `git add internal/app/follow/`
- [ ] Stage documentation: `git add postman/ *.md`
- [ ] Stage test script: `git add test_profile.sh`
- [ ] Commit: `git commit -m "feat: implement user profiles & privacy settings"`

### Push to Repository
- [ ] Push to remote: `git push origin main`
- [ ] Verify push was successful

---

## ✅ COMPLETION SUMMARY

### What You've Built:
- [x] 9 new API endpoints for user profiles
- [x] Profile viewing (public + private)
- [x] Profile editing (name, username, bio)
- [x] Avatar upload/delete with image processing
- [x] Privacy settings (4 privacy controls)
- [x] Mutual connections feature
- [x] Privacy enforcement in Follow system
- [x] Complete test suite
- [x] Postman collection & documentation

### Deliverables:
- [x] Database migration (000014)
- [x] Updated User model with 6 new fields
- [x] 12 new DTOs
- [x] Repository methods (profile operations)
- [x] Service layer (business logic + validation)
- [x] Handler layer (HTTP handlers)
- [x] Routes registration
- [x] Avatar processing service
- [x] Test script
- [x] Postman collection
- [x] API documentation

---

## 🎯 Next Feature

After completing this feature, you're ready for:
- **Student Verification System**
- **Institution System**
- **Direct Messaging**

---

**Total Items:** 200+ checkboxes  
**Feature:** User Profiles & Following Enhancement  
**Status:** Ready to implement at your pace!

Good luck! 🚀


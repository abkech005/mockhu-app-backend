# User Profiles & Following System - Implementation Plan

## 📋 Overview

**Feature:** User Profiles & Following (MVP1 - Priority 1, Feature #4)  
**Timeline:** 1 week  
**Dependencies:** Auth System ✅, Follow System ✅

---

## ✅ Already Completed

### From Section 1: Follow System
- ✅ Follow/unfollow users
- ✅ Get followers list (paginated)
- ✅ Get following list (paginated)
- ✅ Check if user is following another user
- ✅ Get follow statistics (follower count, following count)

### From Auth System
- ✅ User model exists with basic fields:
  - ID, email, phone, username, password
  - first_name, last_name, date_of_birth
  - avatar_url, is_active, created_at, updated_at

---

## 🔨 Features To Build

### 1. View User Profiles ✅
### 2. Mutual Connections Display 🔨
### 3. Edit Profile (bio, avatar, institution) 🔨
### 4. Privacy Settings 🔨

---

## 📊 Database Schema Changes

### 1. Extend Users Table

```sql
-- Migration: 000014_add_profile_fields.up.sql

-- Add profile fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS institution_id UUID REFERENCES institutions(id);

-- Add privacy settings
ALTER TABLE users ADD COLUMN IF NOT EXISTS who_can_message VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN IF NOT EXISTS who_can_see_posts VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_followers_list BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_following_list BOOLEAN DEFAULT true;

-- Add constraints
ALTER TABLE users ADD CONSTRAINT valid_message_privacy 
    CHECK (who_can_message IN ('everyone', 'followers', 'none'));

ALTER TABLE users ADD CONSTRAINT valid_posts_privacy 
    CHECK (who_can_see_posts IN ('everyone', 'followers', 'none'));

-- Add bio length constraint (optional)
ALTER TABLE users ADD CONSTRAINT bio_length_check 
    CHECK (LENGTH(bio) <= 500);

-- Index for username uniqueness
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique 
    ON users(LOWER(username)) WHERE username IS NOT NULL;

-- Index for institution lookup
CREATE INDEX IF NOT EXISTS idx_users_institution 
    ON users(institution_id) WHERE institution_id IS NOT NULL;
```

### 2. Down Migration

```sql
-- Migration: 000014_add_profile_fields.down.sql

-- Remove constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS valid_message_privacy;
ALTER TABLE users DROP CONSTRAINT IF EXISTS valid_posts_privacy;
ALTER TABLE users DROP CONSTRAINT IF EXISTS bio_length_check;

-- Remove indexes
DROP INDEX IF EXISTS idx_users_username_unique;
DROP INDEX IF EXISTS idx_users_institution;

-- Remove columns
ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS institution_id;
ALTER TABLE users DROP COLUMN IF EXISTS who_can_message;
ALTER TABLE users DROP COLUMN IF EXISTS who_can_see_posts;
ALTER TABLE users DROP COLUMN IF EXISTS show_followers_list;
ALTER TABLE users DROP COLUMN IF EXISTS show_following_list;
```

---

## 📡 API Endpoints

### Profile Viewing (2 endpoints)

#### 1. Get User Profile (Public)
**Endpoint:** `GET /v1/users/:userId/profile`

**Path Parameters:**
- `userId` (string, required): The user's ID or username

**Response:** `200 OK`
```json
{
    "id": "user-uuid",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "avatar_url": "https://storage.../avatar.jpg",
    "bio": "Computer Science student at Stanford",
    "institution": {
        "id": "inst-uuid",
        "name": "Stanford University",
        "type": "university"
    },
    "stats": {
        "posts_count": 42,
        "followers_count": 150,
        "following_count": 89
    },
    "is_following": true,
    "is_followed_by": false,
    "mutual_connections_count": 12,
    "created_at": "2024-01-15T10:00:00Z"
}
```

**Privacy Considerations:**
- If `show_followers_list` is false and not own profile, hide follower count or show as "hidden"
- If `show_following_list` is false and not own profile, hide following count
- Bio is always public if set
- Stats may be hidden based on privacy settings

#### 2. Get Own Profile
**Endpoint:** `GET /v1/users/me/profile`

**Response:** `200 OK`
```json
{
    "id": "user-uuid",
    "username": "johndoe",
    "email": "john@stanford.edu",
    "phone": "+1234567890",
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "2000-05-15",
    "avatar_url": "https://storage.../avatar.jpg",
    "bio": "Computer Science student",
    "institution": {
        "id": "inst-uuid",
        "name": "Stanford University",
        "type": "university"
    },
    "privacy_settings": {
        "who_can_message": "followers",
        "who_can_see_posts": "everyone",
        "show_followers_list": true,
        "show_following_list": true
    },
    "stats": {
        "posts_count": 42,
        "followers_count": 150,
        "following_count": 89,
        "profile_views": 523
    },
    "email_verified": true,
    "phone_verified": false,
    "is_active": true,
    "created_at": "2024-01-15T10:00:00Z"
}
```

**Note:** Returns all fields including private data

---

### Profile Management (3 endpoints)

#### 3. Update Profile
**Endpoint:** `PUT /v1/users/me/profile`

**Request Body:**
```json
{
    "first_name": "John",
    "last_name": "Doe",
    "bio": "Updated bio text...",
    "username": "johndoe",
    "institution_id": "inst-uuid"
}
```

**Validation Rules:**
- `first_name`: 1-50 characters, required
- `last_name`: 1-50 characters, required
- `bio`: max 500 characters, optional
- `username`: 3-30 characters, alphanumeric + underscore, unique
- `institution_id`: must exist in institutions table

**Response:** `200 OK`
```json
{
    "id": "user-uuid",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "bio": "Updated bio text...",
    "institution_id": "inst-uuid",
    "updated_at": "2024-11-26T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Validation errors
- `409 Conflict`: Username already taken
- `401 Unauthorized`: Not authenticated

#### 4. Upload Avatar
**Endpoint:** `POST /v1/users/me/avatar`

**Request:** `multipart/form-data`
- `avatar` (file, required): Image file

**Accepted Formats:** JPEG, PNG, WebP  
**Max Size:** 5MB  
**Processing:**
- Resize to 400x400 (square crop)
- Compress for web
- Generate unique filename
- Store in S3 or local storage
- Delete old avatar if exists

**Response:** `200 OK`
```json
{
    "avatar_url": "https://storage.../avatars/user-uuid-123456.jpg",
    "updated_at": "2024-11-26T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid file type, file too large
- `413 Payload Too Large`: File > 5MB
- `401 Unauthorized`: Not authenticated

#### 5. Delete Avatar
**Endpoint:** `DELETE /v1/users/me/avatar`

**Response:** `200 OK`
```json
{
    "message": "avatar deleted successfully",
    "avatar_url": null
}
```

**Business Logic:**
- Delete file from storage
- Set avatar_url to NULL in database
- Return default avatar URL (optional)

---

### Privacy Settings (2 endpoints)

#### 6. Get Privacy Settings
**Endpoint:** `GET /v1/users/me/privacy`

**Response:** `200 OK`
```json
{
    "who_can_message": "followers",
    "who_can_see_posts": "everyone",
    "show_followers_list": true,
    "show_following_list": true
}
```

#### 7. Update Privacy Settings
**Endpoint:** `PUT /v1/users/me/privacy`

**Request Body:**
```json
{
    "who_can_message": "followers",
    "who_can_see_posts": "everyone",
    "show_followers_list": false,
    "show_following_list": true
}
```

**Privacy Options:**

**`who_can_message`:**
- `everyone`: Anyone can send messages
- `followers`: Only followers can send messages
- `none`: No one can send messages (DMs disabled)

**`who_can_see_posts`:**
- `everyone`: Posts are public
- `followers`: Only followers see posts
- `none`: Only user can see own posts (private account)

**`show_followers_list`:**
- `true`: Followers list is visible to everyone
- `false`: Followers list only visible to user

**`show_following_list`:**
- `true`: Following list is visible to everyone
- `false`: Following list only visible to user

**Response:** `200 OK`
```json
{
    "who_can_message": "followers",
    "who_can_see_posts": "everyone",
    "show_followers_list": false,
    "show_following_list": true,
    "updated_at": "2024-11-26T10:00:00Z"
}
```

---

### Mutual Connections (2 endpoints)

#### 8. Get Mutual Connections
**Endpoint:** `GET /v1/users/:userId/mutual-connections`

**Path Parameters:**
- `userId` (string, required): The target user's ID

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 20, max: 50)

**Response:** `200 OK`
```json
{
    "mutual_connections": [
        {
            "id": "user-uuid-1",
            "username": "janedoe",
            "first_name": "Jane",
            "last_name": "Doe",
            "avatar_url": "https://storage.../avatar.jpg",
            "institution": {
                "id": "inst-uuid",
                "name": "Stanford University"
            }
        },
        {
            "id": "user-uuid-2",
            "username": "bobsmith",
            "first_name": "Bob",
            "last_name": "Smith",
            "avatar_url": "https://storage.../avatar2.jpg",
            "institution": {
                "id": "inst-uuid",
                "name": "Stanford University"
            }
        }
    ],
    "pagination": {
        "page": 1,
        "total_pages": 1,
        "total_items": 12,
        "limit": 20
    }
}
```

**Algorithm:**
```sql
-- Find users that both current user AND target user are following
SELECT u.* 
FROM users u
WHERE u.id IN (
    SELECT followed_id FROM user_follows WHERE follower_id = $1 -- current user
)
AND u.id IN (
    SELECT followed_id FROM user_follows WHERE follower_id = $2 -- target user
)
AND u.id NOT IN ($1, $2) -- exclude both users
ORDER BY u.first_name, u.last_name
LIMIT $3 OFFSET $4;
```

#### 9. Get Mutual Connections Count
**Endpoint:** `GET /v1/users/:userId/mutual-connections/count`

**Response:** `200 OK`
```json
{
    "user_id": "user-uuid",
    "mutual_connections_count": 12
}
```

**Business Logic:**
- Count users followed by both current user and target user
- Cache this count for performance
- Return 0 if not authenticated

---

## 🏗️ Directory Structure

```
internal/app/profile/
├── model.go              # Profile domain model (if needed)
├── dto.go                # Request/Response DTOs
├── repository.go         # Repository interface
├── repository_postgres.go # PostgreSQL implementation
├── service.go            # Business logic
├── handler.go            # HTTP handlers
└── routes.go             # Route registration
```

**Note:** May extend `internal/app/auth` instead of creating new package, since it's user-related.

---

## 🔐 Privacy Enforcement

### Service Layer Checks

#### When viewing posts:
```go
func CanViewPosts(viewer *User, postOwner *User) bool {
    // Owner can always see own posts
    if viewer.ID == postOwner.ID {
        return true
    }
    
    // Check privacy setting
    switch postOwner.WhoCanSeePosts {
    case "everyone":
        return true
    case "followers":
        return isFollowing(viewer.ID, postOwner.ID)
    case "none":
        return false
    default:
        return false
    }
}
```

#### When sending messages:
```go
func CanSendMessage(sender *User, recipient *User) error {
    // Check if recipient allows messages from sender
    switch recipient.WhoCanMessage {
    case "everyone":
        return nil
    case "followers":
        if !isFollowing(sender.ID, recipient.ID) {
            return ErrCannotMessage
        }
        return nil
    case "none":
        return ErrMessagingDisabled
    default:
        return ErrCannotMessage
    }
}
```

#### When viewing followers/following lists:
```go
func CanViewFollowersList(viewer *User, profileOwner *User) bool {
    // Owner can always see own list
    if viewer.ID == profileOwner.ID {
        return true
    }
    
    // Check privacy setting
    return profileOwner.ShowFollowersList
}

func CanViewFollowingList(viewer *User, profileOwner *User) bool {
    // Owner can always see own list
    if viewer.ID == profileOwner.ID {
        return true
    }
    
    // Check privacy setting
    return profileOwner.ShowFollowingList
}
```

---

## 🚀 Implementation Order

### Phase 1: Database Migration (0.5 day)
1. Create migration file with all new fields
2. Test migration up/down
3. Update User model in auth package

### Phase 2: Profile Viewing (1 day)
1. Create profile DTOs
2. Implement GetUserProfile endpoint
3. Implement GetOwnProfile endpoint
4. Add stats calculation (posts, followers, following counts)
5. Test profile viewing

### Phase 3: Profile Management (1.5 days)
1. Implement UpdateProfile endpoint
2. Add validation logic
3. Implement UploadAvatar endpoint
4. Integrate image processing library
5. Implement DeleteAvatar endpoint
6. Test profile updates and avatar operations

### Phase 4: Privacy Settings (0.5 day)
1. Implement GetPrivacySettings endpoint
2. Implement UpdatePrivacySettings endpoint
3. Add validation for privacy values
4. Test privacy settings updates

### Phase 5: Mutual Connections (1 day)
1. Write mutual connections SQL query
2. Optimize query with indexes
3. Implement GetMutualConnections endpoint
4. Implement GetMutualConnectionsCount endpoint
5. Test mutual connections accuracy
6. Performance testing

### Phase 6: Privacy Enforcement (1 day)
1. Add privacy checks to Follow service
2. Add privacy checks to Post service (for future)
3. Update existing handlers to respect privacy
4. Test privacy enforcement across systems

### Phase 7: Testing & Documentation (0.5 day)
1. Integration testing all endpoints
2. Create Postman collection
3. Write API documentation
4. Update SECTIONS_STATUS.md

**Total: 6-7 days**

---

## 📊 API Endpoints Summary

| Category | Method | Endpoint | Auth | Description |
|----------|--------|----------|------|-------------|
| **Profile View** | GET | `/v1/users/:userId/profile` | Optional | Get user profile (public) |
| | GET | `/v1/users/me/profile` | ✅ Required | Get own profile (private) |
| **Profile Edit** | PUT | `/v1/users/me/profile` | ✅ Required | Update profile |
| | POST | `/v1/users/me/avatar` | ✅ Required | Upload avatar |
| | DELETE | `/v1/users/me/avatar` | ✅ Required | Delete avatar |
| **Privacy** | GET | `/v1/users/me/privacy` | ✅ Required | Get privacy settings |
| | PUT | `/v1/users/me/privacy` | ✅ Required | Update privacy settings |
| **Mutual** | GET | `/v1/users/:userId/mutual-connections` | ✅ Required | Get mutual connections list |
| | GET | `/v1/users/:userId/mutual-connections/count` | ✅ Required | Get mutual count |

**Total: 9 new endpoints**

---

## 🧪 Testing Requirements

### Unit Tests
- Profile update validation
- Username uniqueness check
- Bio length validation
- Privacy settings validation
- Mutual connections calculation

### Integration Tests
- Get user profile (authenticated vs public view)
- Update profile with valid/invalid data
- Upload avatar (various file types, sizes)
- Privacy settings prevent unauthorized access
- Mutual connections pagination

### Edge Cases
- Update username to existing username (should fail)
- Upload non-image file as avatar (should fail)
- Set invalid privacy value (should fail)
- View profile of user who blocked you
- Mutual connections when one user has no follows

---

## 📝 Business Rules

### Profile Visibility
1. **Public fields:** username, name, avatar, bio, institution, stats
2. **Private fields:** email, phone, DOB, privacy settings
3. **Conditional fields:** Stats may be hidden based on privacy

### Username Rules
1. Must be unique (case-insensitive)
2. 3-30 characters
3. Alphanumeric and underscore only
4. Cannot start with underscore
5. Cannot change more than once per 30 days (optional)

### Bio Rules
1. Max 500 characters
2. Allow emojis and special characters
3. Strip HTML/script tags (sanitize input)
4. Optional field

### Avatar Rules
1. Max file size: 5MB
2. Accepted formats: JPEG, PNG, WebP
3. Auto-resize to 400x400 (square crop)
4. Delete old avatar when uploading new
5. Generate unique filename to prevent collisions

### Privacy Settings
1. Default to most open (`everyone`)
2. Changes take effect immediately
3. `none` option makes content visible only to user
4. Privacy settings don't affect profile view itself

### Mutual Connections
1. Only visible to authenticated users
2. Show users followed by BOTH parties
3. Exclude the two users themselves
4. Paginated (max 50 per page)
5. Sorted alphabetically

---

## 🔄 Integration Points

### With Follow System
- Check follow status for privacy enforcement
- Get follower/following counts for stats
- Calculate mutual connections

### With Post System
- Count user's posts for stats
- Filter posts based on `who_can_see_posts` setting
- Hide posts from users based on privacy

### With Messaging System (Future)
- Check `who_can_message` before allowing DM
- Show/hide message button based on privacy
- Block messages from users not meeting criteria

### With Institution System (Future)
- Link profile to institution
- Display institution name in profile
- Filter users by institution

---

## 📚 Technology Stack

### Go Libraries

```go
// Image processing
"github.com/disintegration/imaging"
// or
"github.com/h2non/bimg" // requires libvips

// HTML sanitization
"github.com/microcosm-cc/bluemonday"

// File upload validation
"github.com/gabriel-vasile/mimetype"
```

### Storage Options

**Avatar Storage:**
1. **AWS S3** (Recommended for production)
   - Scalable, reliable
   - CDN integration
   - Cost-effective

2. **Local Storage** (Good for development)
   - Simple, no external dependencies
   - Not scalable
   - Backup required

**Decision:** Start with local, plan for S3 migration

---

## ⚠️ Considerations & Risks

### Technical Risks
1. **Avatar storage:** Need S3 or local storage strategy
2. **Image processing:** May need system dependencies (libvips)
3. **Username changes:** Need strategy to prevent abuse
4. **Privacy complexity:** Cross-system checks may impact performance

### Business Risks
1. **Username squatting:** Popular usernames taken early
2. **Privacy defaults:** Too open may hurt privacy, too closed may hurt engagement
3. **Mutual connections:** May expose social graph unintentionally

### Performance Considerations
1. **Mutual connections query** - May be slow with large networks (add indexes)
2. **Profile stats** - Calculate on demand or cache?
3. **Privacy checks** - Add latency to requests (optimize with joins)
4. **Avatar uploads** - Async processing for large images?

---

## 🎯 Success Metrics

### Completion Criteria
- All 9 endpoints implemented and tested
- Privacy settings working across systems
- Avatar upload/delete functional
- Mutual connections accurate
- Profile updates validated correctly
- Postman collection complete

### Performance Targets
- Profile GET: < 100ms
- Profile UPDATE: < 200ms
- Avatar upload: < 2s for 5MB image
- Mutual connections: < 300ms

### Quality Metrics
- Username uniqueness: 100% (no duplicates)
- Avatar upload success rate: > 95%
- Privacy enforcement: 100% (no leaks)
- Mutual connections accuracy: 100%

---

## 📦 Deliverables

### Code
- [ ] Migration: `000014_add_profile_fields`
- [ ] Profile service (or extend auth service)
- [ ] Profile handler with 9 endpoints
- [ ] Image upload & processing
- [ ] Privacy enforcement logic
- [ ] Routes registration

### Testing
- [ ] Unit tests for profile service
- [ ] Integration tests for all endpoints
- [ ] Privacy enforcement tests
- [ ] Test script: `test_profile.sh`

### Documentation
- [ ] Postman collection: `Mockhu_Profile_API.postman_collection.json`
- [ ] Postman environment: `Mockhu_Profile_API_Environment.postman_environment.json`
- [ ] API documentation: `PROFILE_API_DOCUMENTATION.md`
- [ ] Update SECTIONS_STATUS.md

---

## 🎯 Next Steps After Completion

1. **Update Follow System** - Add privacy checks to follower lists
2. **Update Post System** - Add privacy checks to post visibility
3. **Prepare for Messaging** - Privacy checks are foundation for DM
4. **Add Profile Views Counter** - Track who viewed profile (optional)

---

**Estimated Effort:** 6-7 days  
**Priority:** High (MVP1 - Priority 1)  
**Complexity:** Medium (Image processing, privacy logic)  
**Status:** Planning Complete, Ready for Implementation

---

## 📌 Notes

### Regarding Follow System
- Follow/unfollow endpoints already exist ✅
- Follower/following lists already exist ✅
- This plan focuses on NEW features:
  - Profile viewing & management
  - Mutual connections
  - Privacy settings

### Regarding Institution System
- `institution_id` field added to users table
- Full Institution System to be built separately
- For now, institution_id can be NULL or reference future table

---

**Ready to start implementation?** Let's build the User Profiles & Following system! 🚀


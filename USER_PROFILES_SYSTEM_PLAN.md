# User Profiles & Following System - Implementation Plan

## 📋 Overview

**Feature:** User Profiles & Following (MVP1 - Priority 1, Feature #4)  
**Timeline:** 1 week  
**Dependencies:** Auth System, Follow System (already complete)

---

## ✅ Already Completed

### Follow System (Section 1) ✅
- Follow/unfollow users
- Get followers list (paginated)
- Get following list (paginated)
- Check if user is following another user
- Get follow statistics (follower count, following count)

### Auth System ✅
- User signup with email/phone
- Basic user model exists with fields:
  - ID, email, phone, username, password
  - first_name, last_name, date_of_birth
  - avatar_url, institution_id
  - is_active, created_at, updated_at

---

## 🔨 Features To Build

### 1. View User Profiles

#### **1.1 Get User Profile (Public)**
**Endpoint:** `GET /v1/users/:userId/profile`

**Response Data:**
```json
{
    "id": "uuid",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "avatar_url": "https://...",
    "bio": "Computer Science student...",
    "institution": {
        "id": "uuid",
        "name": "Stanford University",
        "type": "university"
    },
    "stats": {
        "posts_count": 42,
        "followers_count": 150,
        "following_count": 89
    },
    "is_following": true,
    "is_followed_by_me": true,
    "mutual_connections_count": 12,
    "is_verified": true,
    "created_at": "2024-01-15T10:00:00Z"
}
```

**Privacy Considerations:**
- Public profile fields: username, name, avatar, institution, bio, verification status
- Private fields based on privacy settings: posts visibility, contact info

#### **1.2 Get Own Profile**
**Endpoint:** `GET /v1/users/me/profile`

**Response includes additional private data:**
```json
{
    // All public fields +
    "email": "john@stanford.edu",
    "phone": "+1234567890",
    "date_of_birth": "2000-05-15",
    "privacy_settings": {
        "who_can_message": "everyone|followers|none",
        "who_can_see_posts": "everyone|followers|none",
        "show_followers_list": true,
        "show_following_list": true
    }
}
```

---

### 2. Mutual Connections Display

#### **2.1 Get Mutual Connections**
**Endpoint:** `GET /v1/users/:userId/mutual-connections`

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 50)

**Response:**
```json
{
    "mutual_connections": [
        {
            "id": "uuid",
            "username": "janedoe",
            "first_name": "Jane",
            "last_name": "Doe",
            "avatar_url": "https://...",
            "institution_name": "Stanford University"
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

**Business Logic:**
- Find users that both the authenticated user AND the target user are following
- SQL: Users who appear in both followers lists

#### **2.2 Get Mutual Connections Count**
**Endpoint:** `GET /v1/users/:userId/mutual-connections/count`

**Response:**
```json
{
    "user_id": "uuid",
    "mutual_connections_count": 12
}
```

---

### 3. Edit Profile

#### **3.1 Update Profile**
**Endpoint:** `PUT /v1/users/me/profile`

**Request Body:**
```json
{
    "first_name": "John",
    "last_name": "Doe",
    "bio": "Updated bio text...",
    "username": "newusername",
    "institution_id": "uuid",
    "date_of_birth": "2000-05-15"
}
```

**Validation Rules:**
- `first_name`: 1-50 characters
- `last_name`: 1-50 characters
- `bio`: max 500 characters
- `username`: 3-30 characters, alphanumeric + underscore, unique
- `date_of_birth`: valid date, user must be 13+ years old

**Response:** Updated profile object

#### **3.2 Update Avatar**
**Endpoint:** `POST /v1/users/me/avatar`

**Request:** Multipart form-data with image file

**Processing:**
- Accept: JPEG, PNG, WebP
- Max size: 5MB
- Resize to: 400x400 (square crop)
- Store in: S3/local storage
- Update `avatar_url` in database

**Response:**
```json
{
    "avatar_url": "https://storage.../avatar.jpg",
    "updated_at": "2024-11-26T10:00:00Z"
}
```

#### **3.3 Delete Avatar**
**Endpoint:** `DELETE /v1/users/me/avatar`

**Response:**
```json
{
    "message": "avatar deleted successfully"
}
```

---

### 4. Privacy Settings

#### **4.1 Database Schema Extension**

**Add to `users` table:**
```sql
ALTER TABLE users ADD COLUMN bio TEXT;
ALTER TABLE users ADD COLUMN who_can_message VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN who_can_see_posts VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN show_followers_list BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN show_following_list BOOLEAN DEFAULT true;

-- Constraints
ALTER TABLE users ADD CONSTRAINT valid_message_privacy 
    CHECK (who_can_message IN ('everyone', 'followers', 'none'));

ALTER TABLE users ADD CONSTRAINT valid_posts_privacy 
    CHECK (who_can_see_posts IN ('everyone', 'followers', 'none'));
```

#### **4.2 Get Privacy Settings**
**Endpoint:** `GET /v1/users/me/privacy`

**Response:**
```json
{
    "who_can_message": "followers",
    "who_can_see_posts": "everyone",
    "show_followers_list": true,
    "show_following_list": false
}
```

#### **4.3 Update Privacy Settings**
**Endpoint:** `PUT /v1/users/me/privacy`

**Request Body:**
```json
{
    "who_can_message": "followers",
    "who_can_see_posts": "everyone",
    "show_followers_list": true,
    "show_following_list": false
}
```

**Privacy Options:**
- `who_can_message`:
  - `everyone`: Anyone can message
  - `followers`: Only followers can message
  - `none`: No one can message

- `who_can_see_posts`:
  - `everyone`: Posts are public
  - `followers`: Only followers can see posts
  - `none`: Posts are private (only user can see)

**Response:** Updated privacy settings

---

## 🏗️ Architecture & Implementation

### Directory Structure

```
internal/app/profile/
├── model.go              # Profile domain model
├── dto.go                # Request/Response DTOs
├── repository.go         # Repository interface
├── repository_postgres.go # PostgreSQL implementation
├── service.go            # Business logic
├── handler.go            # HTTP handlers
└── routes.go             # Route registration
```

### Database Changes

#### Migration: `000014_add_profile_fields.up.sql`
```sql
-- Add profile and privacy fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS who_can_message VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN IF NOT EXISTS who_can_see_posts VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_followers_list BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_following_list BOOLEAN DEFAULT true;

-- Constraints
ALTER TABLE users ADD CONSTRAINT valid_message_privacy 
    CHECK (who_can_message IN ('everyone', 'followers', 'none'));

ALTER TABLE users ADD CONSTRAINT valid_posts_privacy 
    CHECK (who_can_see_posts IN ('everyone', 'followers', 'none'));

-- Index for username uniqueness check
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique ON users(username) WHERE username IS NOT NULL;
```

---

## 📡 API Endpoints Summary

### Profile Management (7 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/v1/users/:userId/profile` | Optional | Get user profile (public) |
| GET | `/v1/users/me/profile` | ✅ Required | Get own profile |
| PUT | `/v1/users/me/profile` | ✅ Required | Update profile |
| POST | `/v1/users/me/avatar` | ✅ Required | Upload avatar |
| DELETE | `/v1/users/me/avatar` | ✅ Required | Delete avatar |
| GET | `/v1/users/me/privacy` | ✅ Required | Get privacy settings |
| PUT | `/v1/users/me/privacy` | ✅ Required | Update privacy settings |

### Mutual Connections (2 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/v1/users/:userId/mutual-connections` | ✅ Required | Get mutual connections list |
| GET | `/v1/users/:userId/mutual-connections/count` | ✅ Required | Get mutual connections count |

**Total: 9 new endpoints**

---

## 🔐 Privacy Enforcement

### Service Layer Checks

#### When viewing posts:
```
if post.user.who_can_see_posts == "followers":
    if !isFollowing(currentUser, post.user):
        return 403 Forbidden

if post.user.who_can_see_posts == "none":
    if currentUser.id != post.user.id:
        return 403 Forbidden
```

#### When sending messages:
```
if recipient.who_can_message == "followers":
    if !isFollowing(sender, recipient):
        return 403 Forbidden

if recipient.who_can_message == "none":
    return 403 Forbidden
```

#### When viewing followers/following lists:
```
if !user.show_followers_list:
    if currentUser.id != user.id:
        return 403 Forbidden (or return empty list)

if !user.show_following_list:
    if currentUser.id != user.id:
        return 403 Forbidden (or return empty list)
```

---

## 🧪 Testing Requirements

### Unit Tests
- Profile CRUD operations
- Privacy settings validation
- Username uniqueness check
- Bio length validation
- Mutual connections calculation

### Integration Tests
- Get profile with auth token
- Get profile without auth (public view)
- Update profile with valid data
- Update profile with invalid data (username taken, bio too long)
- Privacy enforcement (posts, messages, lists)
- Mutual connections pagination

### Manual Testing
- Avatar upload and deletion
- Profile updates reflect immediately
- Privacy settings prevent unauthorized access
- Mutual connections are accurate

---

## 📝 Business Rules

### Profile Visibility
1. Public profiles: username, name, avatar, bio, institution, stats, verification status
2. Private profiles: Based on `who_can_see_posts` setting
3. Own profile: All fields visible including email, phone, privacy settings

### Username Rules
1. Must be unique across all users
2. 3-30 characters
3. Alphanumeric and underscore only
4. Cannot be changed more than once per month (optional enhancement)

### Bio Rules
1. Max 500 characters
2. Allow emojis and special characters
3. No HTML/script tags (sanitize input)

### Avatar Rules
1. Max file size: 5MB
2. Accepted formats: JPEG, PNG, WebP
3. Auto-resize to 400x400
4. Delete old avatar when uploading new one

### Mutual Connections
1. Only visible to authenticated users
2. Show users that both parties follow
3. Paginated (max 50 per page)
4. Sorted by most recent follow

### Privacy Settings
1. Default to most open (`everyone`)
2. Changes take effect immediately
3. `none` option hides content from everyone except user
4. Privacy settings don't affect profile view (only posts/messages)

---

## 🔄 Integration Points

### Follow System
- Check if users follow each other (for mutual connections)
- Privacy enforcement for followers-only content

### Posts System
- Filter posts based on `who_can_see_posts` setting
- Show/hide posts in feed based on privacy

### Messaging System (Future)
- Check `who_can_message` before allowing DM
- Show/hide message button based on privacy

### Institution System (Future)
- Link user profile to institution
- Show institution name and type in profile

---

## 📚 Dependencies

### Existing Systems to Use
- Auth System: User authentication and authorization
- Follow System: Follow relationships, follower/following lists
- Upload System: Avatar upload (may need enhancement)

### New Dependencies
- Image processing library (for avatar resizing)
  - Option 1: `github.com/disintegration/imaging`
  - Option 2: `github.com/h2non/bimg` (requires libvips)

---

## 📊 Database Queries Optimization

### Indexes Needed
```sql
-- For username lookups
CREATE INDEX idx_users_username ON users(username);

-- For mutual connections query (if not exists)
CREATE INDEX idx_user_follows_follower ON user_follows(follower_id);
CREATE INDEX idx_user_follows_followed ON user_follows(followed_id);
```

### Complex Queries

#### Mutual Connections:
```sql
SELECT DISTINCT u.*
FROM users u
JOIN user_follows uf1 ON u.id = uf1.followed_id
JOIN user_follows uf2 ON u.id = uf2.followed_id
WHERE uf1.follower_id = $1  -- current user
  AND uf2.follower_id = $2  -- target user
  AND u.id != $1 AND u.id != $2
ORDER BY uf1.created_at DESC
LIMIT $3 OFFSET $4;
```

---

## 🚀 Implementation Order

### Phase 1: Core Profile (2 days)
1. Create profile package structure
2. Add bio field to users table
3. Implement Get Profile endpoints (public + own)
4. Implement Update Profile endpoint
5. Test profile CRUD operations

### Phase 2: Privacy Settings (1 day)
1. Add privacy fields to users table
2. Implement Get/Update Privacy Settings endpoints
3. Test privacy settings validation

### Phase 3: Avatar Management (1 day)
1. Integrate image processing library
2. Implement avatar upload endpoint
3. Implement avatar delete endpoint
4. Test avatar operations

### Phase 4: Mutual Connections (1 day)
1. Implement mutual connections query
2. Create mutual connections endpoints
3. Test mutual connections accuracy
4. Optimize query performance

### Phase 5: Privacy Enforcement (1 day)
1. Add privacy checks to Posts service
2. Add privacy checks to Follow service (lists visibility)
3. Prepare for Messaging system integration
4. Test privacy enforcement across systems

### Phase 6: Testing & Documentation (1 day)
1. Integration testing all endpoints
2. Create Postman collection
3. Write API documentation
4. Update SECTIONS_STATUS.md

---

## 📦 Deliverables

### Code
- [ ] Profile package (model, dto, repository, service, handler, routes)
- [ ] Migration: `000014_add_profile_fields`
- [ ] Updated Follow service with privacy checks
- [ ] Updated Posts service with privacy checks
- [ ] Avatar upload/delete functionality

### Testing
- [ ] Unit tests for profile service
- [ ] Integration tests for all endpoints
- [ ] Test script: `test_profile.sh`

### Documentation
- [ ] Postman collection: `Mockhu_Profile_API.postman_collection.json`
- [ ] Postman environment: `Mockhu_Profile_API_Environment.postman_environment.json`
- [ ] API documentation: `PROFILE_API_DOCUMENTATION.md`

---

## ⚠️ Considerations & Risks

### Technical Risks
1. **Avatar storage**: Need to decide on S3 vs local storage
2. **Image processing**: May need additional system dependencies
3. **Privacy complexity**: Cross-system privacy checks may impact performance

### Business Risks
1. **Username changes**: Need policy on how often users can change usernames
2. **Privacy defaults**: Too restrictive may hurt engagement, too open may hurt privacy
3. **Mutual connections**: May expose social graph, need to consider privacy implications

### Performance Considerations
1. Mutual connections query may be slow with large follower counts (needs optimization)
2. Profile stats calculation (posts count) may need caching
3. Privacy checks on every request may add latency

---

## 🎯 Success Metrics

### Completion Criteria
- All 9 endpoints implemented and tested
- Privacy settings working across systems
- Avatar upload/delete functional
- Mutual connections accurate and performant
- Postman collection and documentation complete

### Performance Targets
- Profile GET requests: < 100ms
- Profile UPDATE requests: < 200ms
- Avatar upload: < 2s for 5MB image
- Mutual connections query: < 300ms

---

**Estimated Effort:** 5-7 days  
**Priority:** High (MVP1 - Priority 1)  
**Status:** Planning Complete, Ready for Implementation


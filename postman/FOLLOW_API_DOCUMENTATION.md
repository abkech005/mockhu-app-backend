# 📚 Mockhu Follow API Documentation

**Version:** 1.0  
**Base URL:** `http://localhost:8085`  
**Authentication:** JWT Bearer Token (except `/follow-stats`)

---

## 🔐 Authentication

All endpoints (except `/follow-stats`) require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

**Get Token:**
```bash
POST /v1/auth/login
{
  "identifier": "email@example.com",  # or phone number
  "password": "your-password"
}
```

---

## 📋 Endpoints Overview

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/v1/users/:userId/follow` | ✅ | Follow a user |
| DELETE | `/v1/users/:userId/follow` | ✅ | Unfollow a user |
| GET | `/v1/users/:userId/is-following` | ✅ | Check if following |
| GET | `/v1/users/:userId/followers` | ✅ | Get followers list |
| GET | `/v1/users/:userId/following` | ✅ | Get following list |
| GET | `/v1/users/:userId/follow-stats` | ❌ | Get follow statistics |

---

## 📖 Endpoint Details

### 1. Follow User

**Endpoint:** `POST /v1/users/:userId/follow`

**Description:** Creates a follow relationship between the authenticated user and the target user.

**Authentication:** Required

**Path Parameters:**
- `userId` (string, required) - UUID of user to follow

**Request Example:**
```bash
curl -X POST http://localhost:8085/v1/users/a9bcd0ce-a4e1-400d-a4cc-46ad69f32858/follow \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Success Response (200):**
```json
{
  "message": "followed successfully",
  "is_following": true
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `cannot follow yourself` | Attempting to follow your own account |
| 401 | `unauthorized` | Missing or invalid JWT token |
| 404 | `user not found` | Target user doesn't exist |
| 500 | `failed to follow user` | Internal server error |

---

### 2. Unfollow User

**Endpoint:** `DELETE /v1/users/:userId/follow`

**Description:** Removes the follow relationship between the authenticated user and the target user.

**Authentication:** Required

**Path Parameters:**
- `userId` (string, required) - UUID of user to unfollow

**Request Example:**
```bash
curl -X DELETE http://localhost:8085/v1/users/a9bcd0ce-a4e1-400d-a4cc-46ad69f32858/follow \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Success Response (200):**
```json
{
  "message": "unfollowed successfully",
  "is_following": false
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `cannot unfollow yourself` | Attempting to unfollow your own account |
| 401 | `unauthorized` | Missing or invalid JWT token |
| 500 | `failed to unfollow user` | Internal server error |

---

### 3. Check if Following

**Endpoint:** `GET /v1/users/:userId/is-following`

**Description:** Checks if the authenticated user is following the target user.

**Authentication:** Required

**Path Parameters:**
- `userId` (string, required) - UUID of user to check

**Request Example:**
```bash
curl -X GET http://localhost:8085/v1/users/a9bcd0ce-a4e1-400d-a4cc-46ad69f32858/is-following \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Success Response (200):**
```json
{
  "is_following": true
}
```

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 401 | `unauthorized` | Missing or invalid JWT token |
| 500 | `failed to check follow status` | Internal server error |

---

### 4. Get Followers

**Endpoint:** `GET /v1/users/:userId/followers`

**Description:** Returns a paginated list of users who follow the target user. Includes whether the current user follows each follower (for "Follow back" functionality).

**Authentication:** Required

**Path Parameters:**
- `userId` (string, required) - UUID of user to get followers for

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1, min: 1)
- `limit` (integer, optional) - Results per page (default: 20, min: 1, max: 100)

**Request Example:**
```bash
curl -X GET "http://localhost:8085/v1/users/a9bcd0ce-a4e1-400d-a4cc-46ad69f32858/followers?page=1&limit=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Success Response (200):**
```json
{
  "users": [
    {
      "id": "152e7cd0-e870-4a96-a9d3-08d980897292",
      "username": "john_doe",
      "first_name": "John",
      "last_name": "Doe",
      "avatar_url": "https://cdn.mockhu.com/avatars/john.jpg",
      "is_followed_by_me": false
    }
  ],
  "total_count": 1,
  "page": 1,
  "limit": 20
}
```

**Response Fields:**
- `users` (array) - List of follower user objects
  - `id` (string) - User UUID
  - `username` (string) - Username
  - `first_name` (string) - First name
  - `last_name` (string) - Last name
  - `avatar_url` (string) - Profile picture URL
  - `is_followed_by_me` (boolean) - Whether current user follows this follower
- `total_count` (integer) - Total number of followers
- `page` (integer) - Current page number
- `limit` (integer) - Results per page

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `user ID is required` | Missing userId parameter |
| 401 | `unauthorized` | Missing or invalid JWT token |
| 500 | `failed to get followers` | Internal server error |

---

### 5. Get Following

**Endpoint:** `GET /v1/users/:userId/following`

**Description:** Returns a paginated list of users that the target user follows.

**Authentication:** Required

**Path Parameters:**
- `userId` (string, required) - UUID of user to get following list for

**Query Parameters:**
- `page` (integer, optional) - Page number (default: 1, min: 1)
- `limit` (integer, optional) - Results per page (default: 20, min: 1, max: 100)

**Request Example:**
```bash
curl -X GET "http://localhost:8085/v1/users/152e7cd0-e870-4a96-a9d3-08d980897292/following?page=1&limit=20" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Success Response (200):**
```json
{
  "users": [
    {
      "id": "a9bcd0ce-a4e1-400d-a4cc-46ad69f32858",
      "username": "jane_smith",
      "first_name": "Jane",
      "last_name": "Smith",
      "avatar_url": "https://cdn.mockhu.com/avatars/jane.jpg",
      "is_followed_by_me": true
    }
  ],
  "total_count": 1,
  "page": 1,
  "limit": 20
}
```

**Response Fields:** Same as Get Followers endpoint

**Error Responses:** Same as Get Followers endpoint

---

### 6. Get Follow Stats

**Endpoint:** `GET /v1/users/:userId/follow-stats`

**Description:** Returns follower and following counts for a user. **Public endpoint** - no authentication required.

**Authentication:** Not required

**Path Parameters:**
- `userId` (string, required) - UUID of user to get stats for

**Request Example:**
```bash
curl -X GET http://localhost:8085/v1/users/a9bcd0ce-a4e1-400d-a4cc-46ad69f32858/follow-stats
```

**Success Response (200):**
```json
{
  "user_id": "a9bcd0ce-a4e1-400d-a4cc-46ad69f32858",
  "follower_count": 5,
  "following_count": 12
}
```

**Response Fields:**
- `user_id` (string) - User UUID
- `follower_count` (integer) - Number of followers
- `following_count` (integer) - Number of users being followed

**Error Responses:**

| Status | Error | Description |
|--------|-------|-------------|
| 400 | `user ID is required` | Missing userId parameter |
| 500 | `failed to get follow stats` | Internal server error |

---

## 🔄 Common Workflows

### Workflow 1: Follow a User
```
1. POST /v1/users/:userId/follow
   → Returns: {"message": "followed successfully", "is_following": true}

2. GET /v1/users/:userId/is-following
   → Returns: {"is_following": true}

3. GET /v1/users/:userId/follow-stats
   → Returns: {"follower_count": X, "following_count": Y}
```

### Workflow 2: Check Follow Status Before Following
```
1. GET /v1/users/:userId/is-following
   → Returns: {"is_following": false}

2. POST /v1/users/:userId/follow
   → Returns: {"message": "followed successfully", "is_following": true}
```

### Workflow 3: Get User's Social Graph
```
1. GET /v1/users/:userId/followers?page=1&limit=20
   → Returns: List of followers

2. GET /v1/users/:userId/following?page=1&limit=20
   → Returns: List of following

3. GET /v1/users/:userId/follow-stats
   → Returns: Counts
```

---

## ⚠️ Error Handling

### Common Errors

**401 Unauthorized:**
```json
{
  "error": "unauthorized"
}
```
**Solution:** Include valid JWT token in Authorization header

**400 Bad Request:**
```json
{
  "error": "cannot follow yourself"
}
```
**Solution:** Don't attempt to follow your own user ID

**404 Not Found:**
```json
{
  "error": "user not found"
}
```
**Solution:** Verify the userId exists in the system

**500 Internal Server Error:**
```json
{
  "error": "failed to follow user"
}
```
**Solution:** Check server logs, may indicate database issue

---

## 📝 Notes

1. **Idempotency:** Following the same user multiple times is safe - subsequent calls return success without creating duplicates.

2. **Self-Follow Prevention:** The system prevents users from following themselves (returns 400 error).

3. **Pagination:** Followers and Following lists support pagination. Use `page` and `limit` query parameters.

4. **Public Stats:** Follow stats endpoint is public (no auth required) for displaying counts on public profiles.

5. **Follow Back:** The `is_followed_by_me` field in followers/following lists helps implement "Follow back" functionality.

---

## 🧪 Testing

### Quick Test Sequence

```bash
# 1. Login to get token
TOKEN=$(curl -s -X POST http://localhost:8085/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier": "test@example.com", "password": "test123"}' \
  | jq -r '.access_token')

# 2. Follow a user
curl -X POST http://localhost:8085/v1/users/TARGET_USER_ID/follow \
  -H "Authorization: Bearer $TOKEN"

# 3. Check if following
curl http://localhost:8085/v1/users/TARGET_USER_ID/is-following \
  -H "Authorization: Bearer $TOKEN"

# 4. Get followers
curl "http://localhost:8085/v1/users/TARGET_USER_ID/followers?page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN"

# 5. Get follow stats (no auth)
curl http://localhost:8085/v1/users/TARGET_USER_ID/follow-stats

# 6. Unfollow
curl -X DELETE http://localhost:8085/v1/users/TARGET_USER_ID/follow \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📦 Postman Collection

Import the collection file:
- **Collection:** `Mockhu_Follow_API.postman_collection.json`
- **Environment:** `Mockhu_Follow_API_Environment.postman_environment.json`

**Setup:**
1. Import both files into Postman
2. Set environment variables:
   - `base_url`: `http://localhost:8085`
   - `access_token`: (get from login)
   - `target_user_id`: (UUID of user to test with)

---

**Last Updated:** November 26, 2025  
**API Version:** 1.0


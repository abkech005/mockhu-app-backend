# Mockhu Complete API Documentation

Complete API documentation for all Mockhu Backend endpoints.

## 📋 Table of Contents

1. [Authentication](#authentication)
2. [Interests](#interests)
3. [Onboarding](#onboarding)
4. [Upload](#upload)
5. [Follow](#follow)
6. [Posts](#posts)

---

## 🔐 Authentication

### POST /v1/auth/signup
Create a new user account.

**Request Body:**
```json
{
  "method": "email",
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "user_id": "uuid",
  "verification_needed": true,
  "verification_channel": "email",
  "verification_code": "123456"
}
```

### POST /v1/auth/verify
Verify email/phone with verification code.

**Request Body:**
```json
{
  "user_id": "uuid",
  "method": "email",
  "code": "123456"
}
```

**Response (200):**
```json
{
  "access_token": "jwt-token",
  "refresh_token": "refresh-token",
  "expires_in": 900
}
```

### POST /v1/auth/login
Login with email/phone and password.

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "access_token": "jwt-token",
  "refresh_token": "refresh-token",
  "expires_in": 900
}
```

### POST /v1/auth/refresh
Refresh access token.

**Request Body:**
```json
{
  "refresh_token": "refresh-token"
}
```

### POST /v1/auth/logout
Logout user.

**Request Body:**
```json
{
  "refresh_token": "refresh-token"
}
```

---

## 🎯 Interests

### GET /v1/interests/
Get all interests. Optional category filter.

**Query Parameters:**
- `category` (optional): Filter by category (e.g., "tech", "arts")

**Response (200):**
```json
{
  "interests": [
    {
      "id": "uuid",
      "name": "Technology",
      "slug": "technology",
      "category": "tech",
      "icon": "💻"
    }
  ]
}
```

### GET /v1/interests/categories
Get all interest categories with counts.

**Response (200):**
```json
{
  "categories": [
    {
      "category": "tech",
      "count": 10
    }
  ]
}
```

### GET /v1/users/:id/interests
Get interests for a specific user.

**Response (200):**
```json
{
  "interests": [...]
}
```

---

## 👤 Onboarding

### POST /v1/onboarding/complete
Complete user onboarding with profile information.

**Headers:**
- `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe",
  "avatar_url": "https://example.com/avatar.jpg",
  "interests": ["technology", "programming"]
}
```

**Response (200):**
```json
{
  "message": "onboarding completed successfully"
}
```

### GET /v1/onboarding/status/:user_id
Get onboarding status for a user.

**Response (200):**
```json
{
  "is_complete": true,
  "completed_at": "2025-11-26T10:00:00Z"
}
```

---

## 📤 Upload

### POST /v1/upload/avatar
Upload user avatar image.

**Headers:**
- `Authorization: Bearer <token>`

**Request:**
- `multipart/form-data`
- Field: `avatar` (file)

**Response (200):**
```json
{
  "avatar_url": "https://example.com/avatars/uuid.jpg"
}
```

---

## 👥 Follow

### POST /v1/users/:userId/follow
Follow a user.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "followed successfully",
  "is_following": true
}
```

### DELETE /v1/users/:userId/follow
Unfollow a user.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "unfollowed successfully",
  "is_following": false
}
```

### GET /v1/users/:userId/is-following
Check if current user is following target user.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "is_following": true
}
```

### GET /v1/users/:userId/followers
Get list of followers for a user.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 50)

**Response (200):**
```json
{
  "users": [...],
  "pagination": {
    "page": 1,
    "total_pages": 1,
    "total_items": 10,
    "limit": 20
  }
}
```

### GET /v1/users/:userId/following
Get list of users that a user is following.

**Headers:**
- `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 50)

**Response (200):**
```json
{
  "users": [...],
  "pagination": {...}
}
```

### GET /v1/users/:userId/follow-stats
Get follow statistics for a user (public endpoint).

**Response (200):**
```json
{
  "followers_count": 100,
  "following_count": 50
}
```

---

## 📝 Posts

### POST /v1/posts
Create a new post.

**Headers:**
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "content": "This is my first post!",
  "images": ["https://example.com/image1.jpg"],
  "is_anonymous": false
}
```

**Constraints:**
- `content`: 1-5000 characters
- `images`: Max 10 images
- `is_anonymous`: Boolean

**Response (201):**
```json
{
  "id": "uuid",
  "author": {
    "id": "uuid",
    "username": "johndoe",
    "first_name": "John",
    "avatar_url": "https://example.com/avatar.jpg"
  },
  "content": "This is my first post!",
  "images": ["https://example.com/image1.jpg"],
  "reactions": {
    "fire_count": 0,
    "is_fired_by_me": false,
    "recent_users": []
  },
  "created_at": "2025-11-26T10:00:00Z"
}
```

### GET /v1/posts/:postId
Get a single post by ID (public endpoint).

**Response (200):**
```json
{
  "id": "uuid",
  "author": {...},
  "content": "...",
  "images": [...],
  "reactions": {...},
  "created_at": "..."
}
```

**Note:** This endpoint increments the view count.

### GET /v1/users/:userId/posts
Get all posts by a specific user.

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 50)

**Headers (optional):**
- `Authorization: Bearer <token>` - For reaction info

**Response (200):**
```json
{
  "posts": [...],
  "pagination": {
    "page": 1,
    "total_pages": 1,
    "total_items": 10,
    "limit": 20
  }
}
```

### GET /v1/feed
Get feed of posts from users you follow.

**Headers:**
- `Authorization: Bearer <token>` (required)

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 50)

**Response (200):**
```json
{
  "posts": [...],
  "pagination": {...}
}
```

### POST /v1/posts/:postId/reactions
Toggle fire reaction on a post.

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "post_id": "uuid",
  "fire_count": 5,
  "is_fired_by_me": true
}
```

**Note:** First call adds reaction, second call removes it.

### DELETE /v1/posts/:postId
Delete a post (soft delete).

**Headers:**
- `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "post deleted successfully"
}
```

**Note:** Only post owner can delete. Post is soft-deleted (is_active = false).

---

## 🔑 Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## 📊 Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "error message"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

---

## 🚀 Quick Start

1. **Import Postman Collection:**
   - Import `Mockhu_Complete_API.postman_collection.json`
   - Import `Mockhu_Complete_API_Environment.postman_environment.json`

2. **Set Environment:**
   - Select "Mockhu - Complete API Environment"
   - Update `base_url` if needed (default: `http://localhost:8085`)

3. **Test Flow:**
   - Signup → Verify → Login (auto-sets `access_token`)
   - Complete Onboarding
   - Create Post (auto-sets `post_id`)
   - Follow Users
   - Get Feed

---

## 📝 Notes

- All timestamps are in ISO 8601 format (UTC)
- Pagination: `page` starts at 1, `limit` max is 50
- Post content: 1-5000 characters
- Images: Max 10 per post
- Reactions: Currently only "fire" type supported
- Soft delete: Deleted posts are marked `is_active = false`


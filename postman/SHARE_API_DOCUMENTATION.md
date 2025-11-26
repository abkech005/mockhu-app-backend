# Share System API Documentation

Complete API documentation for the Share System (Section 4) of the Mockhu backend.

## Base URL

```
http://localhost:8085
```

## Authentication

Protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

---

## Public Endpoints (No Authentication Required)

### 1. Get Share

Get a single share by its ID.

**Endpoint:** `GET /v1/shares/:shareId`

**Path Parameters:**
- `shareId` (string, required): The ID of the share

**Response:** `200 OK`
```json
{
    "id": "share-uuid",
    "post_id": "post-uuid",
    "user": {
        "id": "user-uuid",
        "username": "johndoe",
        "first_name": "John",
        "avatar_url": "https://example.com/avatar.jpg"
    },
    "shared_to_type": "timeline",
    "created_at": "2025-11-26T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Share ID is required
- `404 Not Found`: Share not found
- `500 Internal Server Error`: Failed to get share

---

### 2. Get Post Shares

Get all shares for a specific post with pagination.

**Endpoint:** `GET /v1/posts/:postId/shares`

**Path Parameters:**
- `postId` (string, required): The ID of the post

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 20, max: 50)

**Example:** `GET /v1/posts/123/shares?page=1&limit=20`

**Response:** `200 OK`
```json
{
    "shares": [
        {
            "id": "share-uuid-1",
            "post_id": "post-uuid",
            "user": {
                "id": "user-uuid-1",
                "username": "johndoe",
                "first_name": "John",
                "avatar_url": "https://example.com/avatar1.jpg"
            },
            "shared_to_type": "timeline",
            "created_at": "2025-11-26T10:30:00Z"
        },
        {
            "id": "share-uuid-2",
            "post_id": "post-uuid",
            "user": {
                "id": "user-uuid-2",
                "username": "janedoe",
                "first_name": "Jane",
                "avatar_url": "https://example.com/avatar2.jpg"
            },
            "shared_to_type": "dm",
            "created_at": "2025-11-26T09:15:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "total_pages": 1,
        "total_items": 2,
        "limit": 20
    }
}
```

**Error Responses:**
- `400 Bad Request`: Post ID is required
- `500 Internal Server Error`: Failed to get shares

---

### 3. Get Share Count

Get the total count of shares for a specific post.

**Endpoint:** `GET /v1/posts/:postId/shares/count`

**Path Parameters:**
- `postId` (string, required): The ID of the post

**Response:** `200 OK`
```json
{
    "post_id": "post-uuid",
    "count": 42
}
```

**Error Responses:**
- `400 Bad Request`: Post ID is required
- `500 Internal Server Error`: Failed to get share count

---

### 4. Get User Shares

Get all shares made by a specific user with pagination.

**Endpoint:** `GET /v1/users/:userId/shares`

**Path Parameters:**
- `userId` (string, required): The ID of the user

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 20, max: 50)

**Example:** `GET /v1/users/123/shares?page=1&limit=20`

**Response:** `200 OK`
```json
{
    "shares": [
        {
            "id": "share-uuid-1",
            "post_id": "post-uuid-1",
            "user": {
                "id": "user-uuid",
                "username": "johndoe",
                "first_name": "John",
                "avatar_url": "https://example.com/avatar.jpg"
            },
            "shared_to_type": "timeline",
            "created_at": "2025-11-26T10:30:00Z"
        }
    ],
    "pagination": {
        "page": 1,
        "total_pages": 1,
        "total_items": 1,
        "limit": 20
    }
}
```

**Error Responses:**
- `400 Bad Request`: User ID is required
- `500 Internal Server Error`: Failed to get shares

---

## Protected Endpoints (Authentication Required)

### 5. Create Share

Share a post. A user can only share a post once.

**Endpoint:** `POST /v1/posts/:postId/shares`

**Headers:**
- `Authorization: Bearer <access_token>` (required)
- `Content-Type: application/json` (required)

**Path Parameters:**
- `postId` (string, required): The ID of the post to share

**Request Body:**
```json
{
    "shared_to_type": "timeline"
}
```

**Share Types:**
- `timeline` (default): Share to user's timeline
- `dm`: Share via direct message
- `external`: Share externally

**Response:** `201 Created`
```json
{
    "id": "share-uuid",
    "post_id": "post-uuid",
    "user": {
        "id": "user-uuid",
        "username": "johndoe",
        "first_name": "John",
        "avatar_url": "https://example.com/avatar.jpg"
    },
    "shared_to_type": "timeline",
    "created_at": "2025-11-26T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: 
  - Invalid request body
  - Post ID is required
  - Invalid share type (must be 'timeline', 'dm', or 'external')
- `401 Unauthorized`: Missing or invalid authorization token
- `404 Not Found`: Post not found
- `409 Conflict`: Post already shared by user (duplicate share attempt)
- `500 Internal Server Error`: Failed to create share

---

### 6. Delete Share

Delete a share. Only the owner of the share can delete it.

**Endpoint:** `DELETE /v1/shares/:shareId`

**Headers:**
- `Authorization: Bearer <access_token>` (required)

**Path Parameters:**
- `shareId` (string, required): The ID of the share to delete

**Response:** `200 OK`
```json
{
    "message": "share deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Share ID is required
- `401 Unauthorized`: Missing or invalid authorization token
- `403 Forbidden`: Unauthorized to delete this share (not the owner)
- `404 Not Found`: Share not found
- `500 Internal Server Error`: Failed to delete share

---

## Share Types

| Type | Description |
|------|-------------|
| `timeline` | Share to user's timeline (default) |
| `dm` | Share via direct message |
| `external` | Share externally |

---

## Business Rules

1. **Duplicate Prevention**: A user can only share a post once. Attempting to share the same post again will return a `409 Conflict` error.

2. **Ownership**: Only the user who created a share can delete it. Attempting to delete someone else's share will return a `403 Forbidden` error.

3. **Post Validation**: The post must exist before it can be shared. Sharing a non-existent post will return a `404 Not Found` error.

4. **Share Type Validation**: Only `timeline`, `dm`, and `external` are valid share types. Invalid types will return a `400 Bad Request` error.

---

## Example Workflow

1. **Get a post to share:**
   ```bash
   GET /v1/posts?page=1&limit=1
   ```

2. **Create a share:**
   ```bash
   POST /v1/posts/{postId}/shares
   {
       "shared_to_type": "timeline"
   }
   ```

3. **Get share count:**
   ```bash
   GET /v1/posts/{postId}/shares/count
   ```

4. **Get all shares for the post:**
   ```bash
   GET /v1/posts/{postId}/shares?page=1&limit=20
   ```

5. **Get a specific share:**
   ```bash
   GET /v1/shares/{shareId}
   ```

6. **Delete the share:**
   ```bash
   DELETE /v1/shares/{shareId}
   ```

---

## Postman Collection

Import the Postman collection from:
- `postman/Mockhu_Share_API.postman_collection.json`
- `postman/Mockhu_Share_API_Environment.postman_environment.json`

The collection includes:
- All 6 endpoints
- Pre-configured tests
- Example requests and responses
- Environment variables

---

## Error Codes Summary

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (invalid input) |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (not authorized for action) |
| 404 | Not Found (resource doesn't exist) |
| 409 | Conflict (duplicate share) |
| 500 | Internal Server Error |

---

## Notes

- All timestamps are in ISO 8601 format (RFC3339)
- Pagination defaults: page=1, limit=20, max limit=50
- Share IDs are UUIDs
- User information is automatically enriched in responses
- Shares are ordered by `created_at DESC` (newest first)


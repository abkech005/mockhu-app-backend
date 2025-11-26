# Comments System - Single-Level Implementation Summary

## ✅ **Single-Level Comment System (ENFORCED)**

### What is Single-Level?

**Single-Level means:**
- ✅ **Comment on Post** → Top-level comment (`parent_comment_id = NULL`)
- ✅ **Comment on Comment** → Reply to top-level comment (`parent_comment_id = comment_id`)
- ❌ **Comment on Reply** → **NOT ALLOWED** (replies cannot have replies)

### Structure

```
Post
├── Comment 1 (top-level)
│   ├── Reply 1.1 ✅ (to Comment 1)
│   ├── Reply 1.2 ✅ (to Comment 1)
│   └── Reply 1.3 ✅ (to Comment 1)
├── Comment 2 (top-level)
│   └── Reply 2.1 ✅ (to Comment 2)
└── Comment 3 (top-level)
    └── (no replies)
```

**NOT ALLOWED:**
```
Post
└── Comment 1
    └── Reply 1.1
        └── Reply 1.1.1 ❌ (REJECTED - cannot reply to reply)
```

### Implementation

**Validation Added:**
```go
// In CreateComment service
if parent.ParentCommentID != nil && *parent.ParentCommentID != "" {
    return nil, errors.New("cannot reply to a reply - only top-level comments can have replies")
}
```

**Error Response:**
```json
{
  "error": "cannot reply to a reply - only top-level comments can have replies"
}
```
**Status:** `400 Bad Request`

### API Endpoints

1. **POST /v1/posts/:postId/comments**
   - Create top-level comment (no `parent_comment_id`)
   - Create reply (with `parent_comment_id` = top-level comment)
   - ❌ Reject reply to reply

2. **GET /v1/posts/:postId/comments**
   - Returns top-level comments
   - Each comment includes its replies (one level only)

3. **GET /v1/comments/:commentId**
   - Returns comment with up to 5 replies
   - Replies don't show nested replies (empty array)

### Database Schema

```sql
parent_comment_id UUID REFERENCES post_comments(id)
```

- Database **supports** multi-level (self-referential)
- Service **enforces** single-level (validation)

### Response Structure

```json
{
  "id": "comment-uuid",
  "content": "Top-level comment",
  "parent_comment_id": null,
  "replies": [
    {
      "id": "reply-uuid",
      "content": "Reply",
      "parent_comment_id": "comment-uuid",
      "replies": [],  // Always empty - no nested replies
      "reply_count": 0
    }
  ],
  "reply_count": 1
}
```

### Testing

**Valid Operations:**
- ✅ Comment on post
- ✅ Reply to comment

**Invalid Operations:**
- ❌ Reply to reply (returns 400 error)

### Files Modified

1. `internal/app/comment/service.go` - Added validation
2. `internal/app/comment/handler.go` - Added error handling
3. Documentation updated

---

## ✅ **Status: Single-Level Enforced**

The system now **enforces single-level comments**:
- Comments can be made on posts
- Replies can be made to comments
- Replies **cannot** be made to replies (validation prevents it)



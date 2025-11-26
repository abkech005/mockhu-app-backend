# Comments System - Analysis

## ✅ **Single-Level Comment System (ENFORCED)**

### Current Implementation

**Structure:**
```
Post
├── Comment 1 (top-level, parent_comment_id = NULL)
│   ├── Reply 1.1 (parent_comment_id = Comment 1)
│   ├── Reply 1.2 (parent_comment_id = Comment 1)
│   └── Reply 1.3 (parent_comment_id = Comment 1)
├── Comment 2 (top-level, parent_comment_id = NULL)
│   └── Reply 2.1 (parent_comment_id = Comment 2)
└── Comment 3 (top-level, parent_comment_id = NULL)
```

**Rules:**
- ✅ **Top-level comments**: `parent_comment_id = NULL` (direct comment on post)
- ✅ **Replies**: `parent_comment_id = top-level comment ID` (reply to comment)
- ❌ **Replies to replies**: **REJECTED** - Validation prevents nested replies

### Validation Logic

**In `CreateComment` service:**
```go
// ENFORCE SINGLE-LEVEL: Reject replies to replies
if parent.ParentCommentID != nil && *parent.ParentCommentID != "" {
    return nil, errors.New("cannot reply to a reply - only top-level comments can have replies")
}
```

**What this means:**
- ✅ Can comment on a post (top-level)
- ✅ Can reply to a top-level comment
- ❌ **Cannot reply to a reply** (enforced at service level)

### Database Schema

```sql
parent_comment_id UUID REFERENCES post_comments(id) ON DELETE CASCADE
```

- Database **supports** multi-level (self-referential)
- Service **enforces** single-level (validation prevents nesting)

### API Response Structure

```json
{
  "id": "comment-uuid",
  "content": "Top-level comment",
  "parent_comment_id": null,
  "replies": [
    {
      "id": "reply-uuid",
      "content": "Reply to comment",
      "parent_comment_id": "comment-uuid",
      "replies": [],  // Empty - replies don't have replies
      "reply_count": 0
    }
  ],
  "reply_count": 1
}
```

### Error Handling

**Attempting to reply to a reply:**
```json
{
  "error": "cannot reply to a reply - only top-level comments can have replies"
}
```
**Status Code:** `400 Bad Request`

---

## 📊 Summary

**System Type:** **Single-Level Comment System**

**Allowed:**
- ✅ Comment on post (top-level)
- ✅ Reply to comment (one level only)

**Not Allowed:**
- ❌ Reply to reply (enforced by validation)

**Implementation:**
- Database: Supports multi-level (unlimited nesting possible)
- Service: **Enforces single-level** (validation prevents replies to replies)
- Response: Shows one level (comment → reply)


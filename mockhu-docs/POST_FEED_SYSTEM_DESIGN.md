# 📱 Post & Feed System - Database Design Document

**Project:** Mockhu - Student Social Media Platform  
**Feature:** Post & Feed System  
**Version:** 1.0 (MVP Design)  
**Date:** November 25, 2025  
**Status:** Planning Phase

---

## 🎯 Overview

This document outlines the database design for Mockhu's Post & Feed System, the core social feature that allows students to share updates, react to content, and discover relevant posts.

### Key Features:
- Create and view posts (text + images)
- React to posts (🔥 Fire reaction)
- Comment on posts
- Share posts
- Personalized feed (following + institution + trending)
- Anonymous posting (optional)

---

## 📊 Database Schema

### 1. `posts` Table

**Purpose:** Store all user posts (content, images, metadata)

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    images TEXT[], -- Array of image URLs (max 4-10)
    post_type VARCHAR(20) DEFAULT 'regular',
    privacy VARCHAR(20) DEFAULT 'public',
    is_anonymous BOOLEAN DEFAULT false,
    institution_id UUID REFERENCES institutions(id),
    is_active BOOLEAN DEFAULT true,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT valid_post_type CHECK (post_type IN ('regular', 'question', 'achievement')),
    CONSTRAINT valid_privacy CHECK (privacy IN ('public', 'institution', 'followers')),
    CONSTRAINT content_length CHECK (char_length(content) >= 1 AND char_length(content) <= 5000),
    CONSTRAINT images_limit CHECK (array_length(images, 1) IS NULL OR array_length(images, 1) <= 10)
);

-- Indexes for performance
CREATE INDEX idx_posts_user_id ON posts(user_id, created_at DESC);
CREATE INDEX idx_posts_institution ON posts(institution_id, created_at DESC) WHERE is_active = true;
CREATE INDEX idx_posts_created_at ON posts(created_at DESC) WHERE is_active = true;
CREATE INDEX idx_posts_active ON posts(is_active, created_at DESC);
```

**Fields Explanation:**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `id` | UUID | Unique post identifier | Primary key, auto-generated |
| `user_id` | UUID | Post author | Foreign key to users, NOT NULL |
| `content` | TEXT | Post text content | 1-5000 characters |
| `images` | TEXT[] | Array of image URLs | Max 10 images, nullable |
| `post_type` | VARCHAR(20) | Type of post | 'regular', 'question', 'achievement' |
| `privacy` | VARCHAR(20) | Visibility setting | 'public', 'institution', 'followers' |
| `is_anonymous` | BOOLEAN | Hide author name | Default false |
| `institution_id` | UUID | Associated institution | Optional, for institution-only posts |
| `is_active` | BOOLEAN | Soft delete flag | Default true |
| `view_count` | INTEGER | Number of views | Default 0, incremented on view |
| `created_at` | TIMESTAMP | Creation time | Auto-generated |
| `updated_at` | TIMESTAMP | Last modified time | Auto-updated |
| `deleted_at` | TIMESTAMP | Deletion time | NULL if not deleted |

---

### 2. `post_reactions` Table

**Purpose:** Store user reactions to posts (🔥 Fire)

```sql
CREATE TABLE post_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction_type VARCHAR(20) DEFAULT 'fire',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_reaction_type CHECK (reaction_type IN ('fire')),
    UNIQUE(post_id, user_id)  -- One reaction per user per post
);

-- Indexes
CREATE INDEX idx_post_reactions_post ON post_reactions(post_id);
CREATE INDEX idx_post_reactions_user ON post_reactions(user_id);
CREATE INDEX idx_post_reactions_created ON post_reactions(created_at DESC);
CREATE UNIQUE INDEX idx_post_reactions_unique ON post_reactions(post_id, user_id);
```

**Fields Explanation:**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `id` | UUID | Unique reaction ID | Primary key |
| `post_id` | UUID | Post being reacted to | Foreign key, NOT NULL |
| `user_id` | UUID | User who reacted | Foreign key, NOT NULL |
| `reaction_type` | VARCHAR(20) | Type of reaction | Currently only 'fire' |
| `created_at` | TIMESTAMP | When reacted | Auto-generated |

**Unique Constraint:** `(post_id, user_id)` ensures one reaction per user per post

---

### 3. `post_comments` Table

**Purpose:** Store comments and replies on posts

```sql
CREATE TABLE post_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES post_comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_anonymous BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT content_length CHECK (char_length(content) >= 1 AND char_length(content) <= 2000)
);

-- Indexes
CREATE INDEX idx_post_comments_post ON post_comments(post_id, created_at);
CREATE INDEX idx_post_comments_user ON post_comments(user_id);
CREATE INDEX idx_post_comments_parent ON post_comments(parent_comment_id);
CREATE INDEX idx_post_comments_active ON post_comments(is_active, created_at DESC);
```

**Fields Explanation:**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `id` | UUID | Unique comment ID | Primary key |
| `post_id` | UUID | Post being commented on | Foreign key, NOT NULL |
| `user_id` | UUID | Comment author | Foreign key, NOT NULL |
| `parent_comment_id` | UUID | Parent comment (for replies) | Nullable, self-reference |
| `content` | TEXT | Comment text | 1-2000 characters |
| `is_anonymous` | BOOLEAN | Hide author name | Default false |
| `is_active` | BOOLEAN | Soft delete flag | Default true |
| `created_at` | TIMESTAMP | Creation time | Auto-generated |
| `updated_at` | TIMESTAMP | Last modified | Auto-updated |
| `deleted_at` | TIMESTAMP | Deletion time | NULL if not deleted |

**Comment Nesting:** One level only (comment → reply, no deeper)

---

### 4. `post_shares` Table

**Purpose:** Track post sharing for viral growth

```sql
CREATE TABLE post_shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shared_to_type VARCHAR(20) DEFAULT 'timeline',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_share_type CHECK (shared_to_type IN ('timeline', 'dm', 'external'))
);

-- Indexes
CREATE INDEX idx_post_shares_post ON post_shares(post_id);
CREATE INDEX idx_post_shares_user ON post_shares(user_id);
CREATE INDEX idx_post_shares_created ON post_shares(created_at DESC);
```

**Fields Explanation:**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `id` | UUID | Unique share ID | Primary key |
| `post_id` | UUID | Post being shared | Foreign key, NOT NULL |
| `user_id` | UUID | User who shared | Foreign key, NOT NULL |
| `shared_to_type` | VARCHAR(20) | Where shared | 'timeline', 'dm', 'external' |
| `created_at` | TIMESTAMP | When shared | Auto-generated |

---

### 5. `user_feed_cache` Table (Future Optimization)

**Purpose:** Cache generated feeds for performance

```sql
CREATE TABLE user_feed_cache (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_type VARCHAR(20) NOT NULL,
    post_ids JSONB NOT NULL,
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    
    PRIMARY KEY (user_id, feed_type),
    CONSTRAINT valid_feed_type CHECK (feed_type IN ('following', 'institution', 'trending', 'discover'))
);

-- Indexes
CREATE INDEX idx_feed_cache_expires ON user_feed_cache(expires_at);
```

**Note:** This table is for future optimization (10K+ users). Not needed for MVP.

---

## 🔗 Entity Relationships

```
users (existing)
  │
  ├─→ posts (1:many)
  │     ├─→ post_reactions (1:many)
  │     ├─→ post_comments (1:many)
  │     │     └─→ post_comments (self-reference for replies)
  │     └─→ post_shares (1:many)
  │
  └─→ post_reactions (1:many) - user's reactions
  └─→ post_comments (1:many) - user's comments
  └─→ post_shares (1:many) - user's shares

institutions (existing)
  └─→ posts (1:many) - institution-specific posts
```

---

## 🎯 Design Decisions

### Decision 1: Image Storage Strategy

**Chosen:** Array of URLs in `posts.images` field

**Alternatives Considered:**
- ❌ Separate `post_images` table (more complex, overkill for MVP)
- ❌ Single image field (too limiting)

**Rationale:**
- ✅ Simple queries (no joins needed)
- ✅ Atomic post creation
- ✅ Sufficient for MVP (max 10 images)
- ✅ Can migrate to separate table later if needed

**Trade-offs:**
- Can't easily query individual images
- No per-image metadata
- Good enough for MVP (<10K users)

---

### Decision 2: Soft Delete vs Hard Delete

**Chosen:** Soft delete with `is_active` and `deleted_at`

**Rationale:**
- ✅ Can recover deleted posts
- ✅ Preserve data for analytics
- ✅ Better user experience (undo delete)
- ✅ Compliance requirements

**Implementation:**
```sql
-- Soft delete
UPDATE posts SET is_active = false, deleted_at = NOW() WHERE id = ?;

-- Query active posts only
SELECT * FROM posts WHERE is_active = true;
```

---

### Decision 3: Reaction Counting Strategy

**MVP Strategy:** Calculate counts on-the-fly

```sql
-- Get reaction count
SELECT COUNT(*) FROM post_reactions WHERE post_id = ?;

-- Get users who reacted
SELECT u.* FROM users u
JOIN post_reactions pr ON pr.user_id = u.id
WHERE pr.post_id = ?
LIMIT 20;
```

**Future Optimization:** Cache counts in posts table

```sql
-- Add cached counter (when scaling)
ALTER TABLE posts ADD COLUMN reaction_count INTEGER DEFAULT 0;

-- Update via database trigger
CREATE TRIGGER update_reaction_count
AFTER INSERT OR DELETE ON post_reactions
FOR EACH ROW
EXECUTE FUNCTION update_post_reaction_count();
```

**When to optimize:** After 10K+ posts or 100K+ reactions

---

### Decision 4: Comment Nesting Strategy

**Chosen:** One level of nesting (comment → reply only)

**Structure:**
```
Post
├── Comment 1
│   ├── Reply 1
│   ├── Reply 2
│   └── Reply 3
├── Comment 2
│   └── Reply 1
└── Comment 3
```

**Rationale:**
- ✅ Good enough for student discussions
- ✅ Simple to display in UI
- ✅ Easy to paginate
- ✅ Prevents infinite nesting complexity

**Query Pattern:**
```sql
-- Get all comments for a post
SELECT * FROM post_comments 
WHERE post_id = ? AND parent_comment_id IS NULL
ORDER BY created_at;

-- Get replies to a comment
SELECT * FROM post_comments
WHERE parent_comment_id = ?
ORDER BY created_at;
```

---

### Decision 5: Anonymous Posting

**Chosen:** Track `user_id` in database, hide in API response

**Implementation:**
```sql
-- Post is stored with user_id
INSERT INTO posts (user_id, content, is_anonymous) 
VALUES (?, ?, true);

-- API response hides user_id if anonymous
{
  "id": "123",
  "author": {
    "id": null,
    "name": "Anonymous Student",
    "avatar": "/default-anonymous.png"
  },
  "content": "Struggling with mental health...",
  "is_anonymous": true
}
```

**Rationale:**
- ✅ Can moderate bad actors (we know who posted)
- ✅ Can show aggregated data (X% anonymous posts)
- ✅ User feels anonymous (name not shown)
- ✅ Safer than truly anonymous

---

### Decision 6: Feed Generation Strategy

**MVP Strategy:** Real-time generation (no caching)

```sql
-- Feed query (simplified)
SELECT p.* FROM posts p
WHERE p.user_id IN (
    SELECT following_id FROM user_follows WHERE follower_id = ?
)
AND p.is_active = true
ORDER BY p.created_at DESC
LIMIT 20 OFFSET 0;
```

**Future Optimization:** Cached feeds with 15-minute TTL

**When to optimize:** After 5K+ users or slow feed loads

---

## 📈 Performance Considerations

### Indexes Strategy

**Critical Indexes (Must Have):**
```sql
-- Posts
CREATE INDEX idx_posts_user_created ON posts(user_id, created_at DESC);
CREATE INDEX idx_posts_created_active ON posts(created_at DESC) WHERE is_active = true;
CREATE INDEX idx_posts_institution ON posts(institution_id, created_at DESC);

-- Reactions
CREATE UNIQUE INDEX idx_reactions_unique ON post_reactions(post_id, user_id);
CREATE INDEX idx_reactions_post ON post_reactions(post_id);

-- Comments
CREATE INDEX idx_comments_post ON post_comments(post_id, created_at);
```

**Impact:**
- ✅ 10-100x faster queries
- ✅ Essential for feed generation
- ✅ Prevents duplicate reactions

---

### Query Optimization

**Bad Query (N+1 Problem):**
```go
// DON'T DO THIS
posts := GetPosts() // 1 query
for post := range posts {
    post.ReactionCount = GetReactionCount(post.ID) // N queries
}
```

**Good Query (Single Query with Join):**
```sql
SELECT 
    p.*,
    COUNT(pr.id) as reaction_count
FROM posts p
LEFT JOIN post_reactions pr ON pr.post_id = p.id
WHERE p.is_active = true
GROUP BY p.id
ORDER BY p.created_at DESC
LIMIT 20;
```

---

### Scaling Thresholds

| Metric | Current Approach | Optimization Trigger | Solution |
|--------|-----------------|---------------------|----------|
| **Users** | <1K | 5K+ | Add feed caching |
| **Posts** | <10K | 100K+ | Add read replicas |
| **Reactions/post** | <100 | 1K+ | Cache reaction counts |
| **Comments/post** | <50 | 500+ | Paginate comments |
| **Image uploads** | Direct upload | 10GB+ | CDN + image service |

---

## 🔒 Data Integrity

### Constraints

**Foreign Key Constraints:**
```sql
-- Cascade deletes (when user deleted, delete their posts)
REFERENCES users(id) ON DELETE CASCADE

-- Protect deletes (can't delete post if has comments - optional)
REFERENCES posts(id) ON DELETE RESTRICT
```

**Check Constraints:**
```sql
-- Content length limits
CONSTRAINT content_length CHECK (char_length(content) BETWEEN 1 AND 5000)

-- Enum validation
CONSTRAINT valid_privacy CHECK (privacy IN ('public', 'institution', 'followers'))

-- Image limit
CONSTRAINT images_limit CHECK (array_length(images, 1) <= 10)
```

**Unique Constraints:**
```sql
-- Prevent duplicate reactions
UNIQUE(post_id, user_id)
```

---

## 🚀 Migration Strategy

### Phase 1: Core Tables (Week 1)
```
✅ Create posts table
✅ Create post_reactions table
✅ Create indexes
✅ Add test data
```

### Phase 2: Comments (Week 2)
```
✅ Create post_comments table
✅ Add indexes
✅ Test nested comments
```

### Phase 3: Optimization (Week 3+)
```
⏸️ Add cached counters (if needed)
⏸️ Add feed cache table (if needed)
⏸️ Add read replicas (if needed)
```

---

## 📊 Sample Data

### Example Post
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "content": "Just got my first internship at Google! 🎉",
  "images": [
    "https://cdn.mockhu.com/posts/abc123.jpg"
  ],
  "post_type": "achievement",
  "privacy": "public",
  "is_anonymous": false,
  "institution_id": "456e8400-e29b-41d4-a716-446655440000",
  "is_active": true,
  "view_count": 245,
  "created_at": "2025-11-25T10:00:00Z",
  "updated_at": "2025-11-25T10:00:00Z",
  "deleted_at": null
}
```

### Example Reaction
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "post_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "789e4567-e89b-12d3-a456-426614174000",
  "reaction_type": "fire",
  "created_at": "2025-11-25T10:05:00Z"
}
```

---

## 🎯 API Response Format

### Get Post
```json
{
  "post": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "john_doe",
      "first_name": "John",
      "avatar_url": "https://cdn.mockhu.com/avatars/john.jpg",
      "institution": "IIT Delhi"
    },
    "content": "Just got my first internship at Google! 🎉",
    "images": [
      "https://cdn.mockhu.com/posts/abc123.jpg"
    ],
    "reactions": {
      "fire_count": 45,
      "is_fired_by_me": true,
      "recent_users": [
        {
          "id": "789e4567-e89b-12d3-a456-426614174000",
          "name": "Sarah",
          "avatar": "https://cdn.mockhu.com/avatars/sarah.jpg"
        }
      ]
    },
    "comments_count": 12,
    "shares_count": 3,
    "view_count": 245,
    "created_at": "2025-11-25T10:00:00Z",
    "is_anonymous": false
  }
}
```

### Get Feed
```json
{
  "feed": {
    "posts": [
      { /* post object */ },
      { /* post object */ },
      { /* post object */ }
    ],
    "pagination": {
      "limit": 20,
      "offset": 0,
      "has_more": true,
      "total": 156
    }
  }
}
```

---

## ✅ Pre-Implementation Checklist

### Before Creating Migrations:

- [x] Schema design reviewed
- [ ] Data types confirmed
- [ ] Constraints validated
- [ ] Indexes planned
- [ ] Relationships verified
- [ ] Sample queries tested
- [ ] Performance considered
- [ ] Security checked (SQL injection prevention)
- [ ] Privacy requirements met
- [ ] Scaling strategy defined

### Questions to Answer:

1. ✅ **Images:** Array in posts table (confirmed)
2. ✅ **Soft delete:** Use is_active flag (confirmed)
3. ✅ **Anonymous:** Track user_id but hide in API (confirmed)
4. ⏸️ **Comments:** Include in MVP? (TBD)
5. ⏸️ **Feed caching:** Add now or later? (Later)
6. ✅ **Counter fields:** Calculate on-the-fly (confirmed)

---

## 🚧 Known Limitations

### MVP Limitations (Acceptable):

1. **No video support** - Images only
2. **No rich text** - Plain text with URLs
3. **No hashtags** - Will add later
4. **No post editing** - Create or delete only
5. **No polls** - Post-MVP feature
6. **No stories** - Future feature
7. **Simple feed algorithm** - Chronological only

### Technical Debt to Address Later:

1. **Reaction counts** - Not cached (fine for <10K posts)
2. **Feed generation** - Not cached (fine for <5K users)
3. **Image optimization** - Client-side only
4. **Comment pagination** - Load all (fine for <50 comments)
5. **Search** - PostgreSQL only (ElasticSearch later)

---

## 📚 References

### Similar Systems:
- Instagram post structure
- Twitter feed algorithm
- Facebook reaction system
- Reddit comment nesting

### Technologies:
- PostgreSQL 14+ (JSONB, array support)
- Go Fiber (backend framework)
- AWS S3 / Cloudinary (image storage)
- Redis (future caching)

---

## 🔄 Next Steps

### Immediate Actions:

1. **Review this design** - Team approval
2. **Create migration files** - Actual SQL
3. **Write repository layer** - Go code
4. **Implement API endpoints** - REST API
5. **Test with sample data** - Verify queries
6. **Build frontend** - Display posts/feed

### Future Enhancements:

1. Add cached counters (after 10K+ posts)
2. Implement feed caching (after 5K+ users)
3. Add ElasticSearch (for post search)
4. Implement video support
5. Add hashtag system
6. Build recommendation engine

---

## ✅ Approval & Sign-off

**Design Status:** ✅ **APPROVED - Ready for Implementation**

**Approved By:** [Pending]  
**Date:** November 25, 2025  
**Next Action:** Create migration files

---

**Document Version:** 1.0  
**Last Updated:** November 25, 2025  
**Next Review:** After MVP launch


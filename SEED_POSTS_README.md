# Seed Posts - Documentation

This guide explains how to seed sample posts into your database for testing and development.

## 📋 Options

There are two ways to seed posts:

### Option 1: Migration (Recommended for Production)

Use the migration file to seed posts as part of your database setup.

```bash
# Run migration
make migrate

# Or manually
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/mockhu_db?sslmode=disable" up
```

**Note:** The migration will only create posts if users already exist in the database.

### Option 2: Go Script (Recommended for Development)

Use the Go script for more control and flexibility.

```bash
# Run seed script
make seed-posts

# Or manually
go run scripts/seed_posts.go
```

## 🎯 What Gets Seeded

### Posts
- **15+ diverse post contents** covering:
  - Technology & coding
  - Lifestyle & photography
  - Food & travel
  - Anonymous posts
- **3-5 posts per user** (if multiple users exist)
- **Random view counts** (0-200)
- **Random creation times** (within last 3 days)
- **Some posts with images** (using Unsplash URLs)

### Reactions
- **30% of posts get reactions**
- **1-5 reactions per post**
- **Random users** react to posts (not their own)

## 📊 Requirements

1. **Users must exist** in the database first
2. **Database connection** via `DATABASE_URL` environment variable
3. **Migrations run** (posts and reactions tables must exist)

## 🔧 Environment Setup

The seed script uses the same `DATABASE_URL` environment variable as the main app:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mockhu_db?sslmode=disable"
```

Or set it in your `.env` file.

## 🚀 Quick Start

```bash
# 1. Ensure database is running
docker-compose up -d

# 2. Run migrations (if not already done)
make migrate

# 3. Create some users (via API or manually)
# Use the signup endpoint to create test users

# 4. Seed posts
make seed-posts
```

## 📝 Sample Output

```
✅ Database connected
Found user: user1@test.com (26b32915-7228-46c9-b798-d9b8e6ba4601)
Found user: user2@test.com (299c7646-8994-4921-87e2-4a30f23e21aa)
Found 2 users. Creating posts...
Created post 632da781-a5b5-48a5-b438-493eb319fd51 for user 26b32915-7228-46c9-b798-d9b8e6ba4601
Created post d9d2afa7-fe20-4c30-a601-0a30516d7084 for user 299c7646-8994-4921-87e2-4a30f23e21aa
...
✅ Successfully created 8 posts with reactions!
```

## 🧹 Cleanup

### Remove Seeded Posts (Migration)

```bash
# Rollback migration
make migrate-down

# Or manually
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/mockhu_db?sslmode=disable" down 1
```

### Remove Seeded Posts (Manual SQL)

```sql
-- Delete reactions first
DELETE FROM post_reactions WHERE post_id IN (
    SELECT id FROM posts WHERE created_at >= NOW() - INTERVAL '3 days'
);

-- Delete posts
DELETE FROM posts WHERE created_at >= NOW() - INTERVAL '3 days';
```

## ⚠️ Notes

- **Migration approach**: Posts are created based on existing users at migration time
- **Script approach**: More flexible, can be run multiple times
- **Reactions**: Only added if multiple users exist
- **Timestamps**: Posts are backdated to simulate a realistic timeline
- **Images**: Uses Unsplash placeholder URLs (replace with real URLs in production)

## 🎨 Customization

To customize the seeded posts, edit:
- **Migration**: `migrations/000011_seed_posts.up.sql`
- **Script**: `scripts/seed_posts.go` - modify the `postContents` array

## 🔍 Verify Seeded Data

```sql
-- Count posts
SELECT COUNT(*) FROM posts;

-- View recent posts
SELECT id, user_id, LEFT(content, 50) as preview, view_count, created_at 
FROM posts 
ORDER BY created_at DESC 
LIMIT 10;

-- Count reactions
SELECT COUNT(*) FROM post_reactions;

-- Posts with most reactions
SELECT p.id, p.content, COUNT(pr.id) as reaction_count
FROM posts p
LEFT JOIN post_reactions pr ON p.id = pr.post_id
GROUP BY p.id, p.content
ORDER BY reaction_count DESC
LIMIT 10;
```


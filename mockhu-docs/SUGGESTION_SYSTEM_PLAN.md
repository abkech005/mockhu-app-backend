# 🎯 User Suggestion System - Implementation Plan

## Overview
Build a recommendation system that suggests users, pages, and groups to follow based on:
- User's selected interests
- User's contacts (phone/email)
- Trending content and popularity
- Social graph (mutual connections)

---

## 📊 Phase 1: Database Schema Design

### 1.1 Pages Table
```sql
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    category VARCHAR(100),
    avatar_url TEXT,
    cover_url TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_follower_count CHECK (follower_count >= 0),
    CONSTRAINT valid_post_count CHECK (post_count >= 0)
);

CREATE INDEX idx_pages_category ON pages(category);
CREATE INDEX idx_pages_owner_id ON pages(owner_id);
CREATE INDEX idx_pages_follower_count ON pages(follower_count DESC);
CREATE INDEX idx_pages_slug ON pages(slug);
```

### 1.2 Groups Table
```sql
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    category VARCHAR(100),
    avatar_url TEXT,
    cover_url TEXT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    member_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    privacy VARCHAR(20) DEFAULT 'public', -- public, private, secret
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_member_count CHECK (member_count >= 0),
    CONSTRAINT valid_post_count CHECK (post_count >= 0),
    CONSTRAINT valid_privacy CHECK (privacy IN ('public', 'private', 'secret'))
);

CREATE INDEX idx_groups_category ON groups(category);
CREATE INDEX idx_groups_created_by ON groups(created_by);
CREATE INDEX idx_groups_member_count ON groups(member_count DESC);
CREATE INDEX idx_groups_privacy ON groups(privacy);
CREATE INDEX idx_groups_slug ON groups(slug);
```

### 1.3 Page Interests (Junction Table)
```sql
CREATE TABLE page_interests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    interest_category VARCHAR(100) NOT NULL,
    interest_subcategory VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(page_id, interest_category, interest_subcategory)
);

CREATE INDEX idx_page_interests_page_id ON page_interests(page_id);
CREATE INDEX idx_page_interests_category ON page_interests(interest_category);
CREATE INDEX idx_page_interests_subcategory ON page_interests(interest_subcategory);
```

### 1.4 Group Interests (Junction Table)
```sql
CREATE TABLE group_interests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    interest_category VARCHAR(100) NOT NULL,
    interest_subcategory VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(group_id, interest_category, interest_subcategory)
);

CREATE INDEX idx_group_interests_group_id ON group_interests(group_id);
CREATE INDEX idx_group_interests_category ON group_interests(interest_category);
CREATE INDEX idx_group_interests_subcategory ON group_interests(interest_subcategory);
```

### 1.5 User Follows (Users following other users)
```sql
CREATE TABLE user_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(follower_id, following_id),
    CONSTRAINT no_self_follow CHECK (follower_id != following_id)
);

CREATE INDEX idx_user_follows_follower ON user_follows(follower_id);
CREATE INDEX idx_user_follows_following ON user_follows(following_id);
CREATE INDEX idx_user_follows_created_at ON user_follows(created_at DESC);
```

### 1.6 Page Follows
```sql
CREATE TABLE page_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, page_id)
);

CREATE INDEX idx_page_follows_user ON page_follows(user_id);
CREATE INDEX idx_page_follows_page ON page_follows(page_id);
CREATE INDEX idx_page_follows_created_at ON page_follows(created_at DESC);
```

### 1.7 Group Members
```sql
CREATE TABLE group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', -- admin, moderator, member
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, group_id),
    CONSTRAINT valid_role CHECK (role IN ('admin', 'moderator', 'member'))
);

CREATE INDEX idx_group_members_user ON group_members(user_id);
CREATE INDEX idx_group_members_group ON group_members(group_id);
CREATE INDEX idx_group_members_role ON group_members(role);
CREATE INDEX idx_group_members_joined_at ON group_members(joined_at DESC);
```

### 1.8 User Contacts (For Contact-based Suggestions)
```sql
CREATE TABLE user_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_type VARCHAR(20) NOT NULL, -- email, phone
    contact_value VARCHAR(255) NOT NULL,
    contact_name VARCHAR(255),
    matched_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_invited BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_contact_type CHECK (contact_type IN ('email', 'phone'))
);

CREATE INDEX idx_user_contacts_user_id ON user_contacts(user_id);
CREATE INDEX idx_user_contacts_value ON user_contacts(contact_value);
CREATE INDEX idx_user_contacts_matched ON user_contacts(matched_user_id);
CREATE INDEX idx_user_contacts_type ON user_contacts(contact_type);
```

### 1.9 Trending Scores (Cached trending data)
```sql
CREATE TABLE trending_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(20) NOT NULL, -- user, page, group
    entity_id UUID NOT NULL,
    score DECIMAL(10,2) DEFAULT 0,
    period VARCHAR(20) NOT NULL, -- daily, weekly, monthly
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    
    UNIQUE(entity_type, entity_id, period),
    CONSTRAINT valid_entity_type CHECK (entity_type IN ('user', 'page', 'group')),
    CONSTRAINT valid_period CHECK (period IN ('daily', 'weekly', 'monthly')),
    CONSTRAINT valid_score CHECK (score >= 0)
);

CREATE INDEX idx_trending_scores_entity ON trending_scores(entity_type, entity_id);
CREATE INDEX idx_trending_scores_period ON trending_scores(period);
CREATE INDEX idx_trending_scores_score ON trending_scores(score DESC);
CREATE INDEX idx_trending_scores_expires ON trending_scores(expires_at);
```

### 1.10 User Activity Stats (For engagement metrics)
```sql
CREATE TABLE user_activity_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    follower_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    engagement_score DECIMAL(10,2) DEFAULT 0,
    last_active_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT valid_follower_count CHECK (follower_count >= 0),
    CONSTRAINT valid_following_count CHECK (following_count >= 0),
    CONSTRAINT valid_post_count CHECK (post_count >= 0),
    CONSTRAINT valid_engagement_score CHECK (engagement_score >= 0)
);

CREATE INDEX idx_user_activity_follower_count ON user_activity_stats(follower_count DESC);
CREATE INDEX idx_user_activity_engagement ON user_activity_stats(engagement_score DESC);
CREATE INDEX idx_user_activity_last_active ON user_activity_stats(last_active_at DESC);
```

---

## 🎯 Phase 2: Recommendation Algorithm Design

### 2.1 Interest-Based Suggestions

#### Algorithm:
1. Get user's selected interests
2. Find users/pages/groups with matching interests
3. Score by:
   - Number of matching interests (weight: 40%)
   - Follower/member count (weight: 30%)
   - Engagement score (weight: 20%)
   - Recency (weight: 10%)

```go
type InterestMatch struct {
    EntityID          string
    EntityType        string  // user, page, group
    MatchingInterests []string
    MatchScore        float64
    FollowerCount     int
    EngagementScore   float64
}

func CalculateInterestScore(userInterests, entityInterests []string, 
                            followerCount int, engagementScore float64) float64 {
    // Jaccard similarity for interests
    matchCount := len(intersect(userInterests, entityInterests))
    totalCount := len(union(userInterests, entityInterests))
    interestScore := float64(matchCount) / float64(totalCount) * 40.0
    
    // Normalize follower count (log scale)
    followerScore := math.Log10(float64(followerCount + 1)) / 6.0 * 30.0
    
    // Engagement score (already 0-100)
    engagementWeight := engagementScore / 100.0 * 20.0
    
    // Recency bonus (new accounts get 10% boost)
    recencyBonus := 10.0
    
    return interestScore + followerScore + engagementWeight + recencyBonus
}
```

### 2.2 Contact-Based Suggestions

#### Algorithm:
1. Match user's contacts (email/phone) with registered users
2. Find mutual connections (friends of friends)
3. Score by:
   - Direct contact match (weight: 50%)
   - Mutual connections count (weight: 30%)
   - Shared interests (weight: 20%)

```go
type ContactMatch struct {
    UserID            string
    MatchType         string  // direct, mutual
    MutualCount       int
    SharedInterests   []string
    MatchScore        float64
}

func CalculateContactScore(matchType string, mutualCount int, 
                           sharedInterests int) float64 {
    baseScore := 0.0
    if matchType == "direct" {
        baseScore = 50.0
    }
    
    // Mutual connections (logarithmic scale)
    mutualScore := math.Log10(float64(mutualCount + 1)) / 4.0 * 30.0
    
    // Shared interests
    interestScore := float64(sharedInterests) / 10.0 * 20.0
    
    return baseScore + mutualScore + interestScore
}
```

### 2.3 Trending-Based Suggestions

#### Algorithm:
1. Calculate trending score based on:
   - New followers in last 24h/7d/30d
   - Engagement rate (likes, comments, shares)
   - Growth velocity
   - Content freshness

```go
type TrendingScore struct {
    EntityID         string
    EntityType       string
    DailyScore       float64
    WeeklyScore      float64
    MonthlyScore     float64
    GrowthVelocity   float64
}

func CalculateTrendingScore(entity Entity, period string) float64 {
    // New followers/members
    newFollowers := entity.GetNewFollowers(period)
    followerScore := float64(newFollowers) / float64(entity.TotalFollowers + 1) * 40.0
    
    // Engagement rate
    engagementRate := entity.GetEngagementRate(period)
    engagementScore := engagementRate * 30.0
    
    // Growth velocity (acceleration)
    velocity := entity.GetGrowthVelocity()
    velocityScore := math.Min(velocity * 20.0, 20.0)
    
    // Content freshness
    freshnessScore := entity.GetContentFreshnessScore() * 10.0
    
    return followerScore + engagementScore + velocityScore + freshnessScore
}
```

### 2.4 Hybrid Recommendation

Combine all three approaches with configurable weights:

```go
type SuggestionConfig struct {
    InterestWeight  float64  // 0.40
    ContactWeight   float64  // 0.35
    TrendingWeight  float64  // 0.25
}

func GenerateFinalSuggestions(userID string, config SuggestionConfig) []Suggestion {
    interestSuggestions := GetInterestBasedSuggestions(userID)
    contactSuggestions := GetContactBasedSuggestions(userID)
    trendingSuggestions := GetTrendingSuggestions()
    
    // Merge and deduplicate
    merged := mergeAndDeduplicate(interestSuggestions, contactSuggestions, trendingSuggestions)
    
    // Calculate final score
    for _, suggestion := range merged {
        suggestion.FinalScore = 
            suggestion.InterestScore * config.InterestWeight +
            suggestion.ContactScore * config.ContactWeight +
            suggestion.TrendingScore * config.TrendingWeight
    }
    
    // Sort by final score
    sort.Slice(merged, func(i, j int) bool {
        return merged[i].FinalScore > merged[j].FinalScore
    })
    
    return merged
}
```

---

## 🎨 Phase 3: API Design

### 3.1 Get Suggestions for New User
```
GET /v1/suggestions/onboarding
Authorization: Bearer <token>

Query Parameters:
- type: string (all, users, pages, groups)
- limit: int (default: 20, max: 100)
- offset: int (default: 0)
- source: string (interests, contacts, trending, all)

Response:
{
  "suggestions": {
    "users": [
      {
        "id": "uuid",
        "username": "johndoe",
        "first_name": "John",
        "last_name": "Doe",
        "avatar_url": "https://...",
        "follower_count": 1234,
        "is_verified": true,
        "match_reason": "3 shared interests: Technology, Gaming, Music",
        "match_score": 85.5,
        "mutual_connections": 5,
        "is_contact": false
      }
    ],
    "pages": [
      {
        "id": "uuid",
        "name": "Tech News Daily",
        "slug": "tech-news-daily",
        "description": "Latest tech news...",
        "avatar_url": "https://...",
        "category": "Technology",
        "follower_count": 50000,
        "is_verified": true,
        "match_reason": "Matches your interest in Technology",
        "match_score": 92.3,
        "trending_rank": 15
      }
    ],
    "groups": [
      {
        "id": "uuid",
        "name": "Gaming Community",
        "slug": "gaming-community",
        "description": "Connect with gamers...",
        "avatar_url": "https://...",
        "category": "Gaming",
        "member_count": 15000,
        "privacy": "public",
        "match_reason": "Popular in your area + Gaming interest",
        "match_score": 88.7,
        "trending_rank": 8
      }
    ]
  },
  "meta": {
    "total_suggestions": 45,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

### 3.2 Upload Contacts
```
POST /v1/suggestions/upload-contacts
Authorization: Bearer <token>

Request:
{
  "contacts": [
    {
      "type": "email",
      "value": "friend@example.com",
      "name": "Friend Name"
    },
    {
      "type": "phone",
      "value": "+1234567890",
      "name": "Another Friend"
    }
  ]
}

Response:
{
  "message": "Contacts uploaded successfully",
  "total_uploaded": 150,
  "matched_users": 23,
  "processing": true
}
```

### 3.3 Follow Multiple Users/Pages/Groups
```
POST /v1/suggestions/bulk-follow
Authorization: Bearer <token>

Request:
{
  "follows": [
    {"type": "user", "id": "uuid1"},
    {"type": "page", "id": "uuid2"},
    {"type": "group", "id": "uuid3"}
  ]
}

Response:
{
  "message": "Followed 3 items successfully",
  "success_count": 3,
  "failed_count": 0,
  "results": [
    {"type": "user", "id": "uuid1", "success": true},
    {"type": "page", "id": "uuid2", "success": true},
    {"type": "group", "id": "uuid3", "success": true}
  ]
}
```

### 3.4 Get Trending
```
GET /v1/suggestions/trending
Authorization: Bearer <token>

Query Parameters:
- type: string (users, pages, groups)
- period: string (daily, weekly, monthly)
- limit: int (default: 20)
- category: string (optional filter by category)

Response:
{
  "trending": [
    {
      "id": "uuid",
      "type": "page",
      "name": "Viral Page",
      "trending_score": 95.8,
      "growth_rate": "+250%",
      "new_followers": 10000,
      "period": "daily"
    }
  ]
}
```

### 3.5 Dismiss Suggestion
```
POST /v1/suggestions/:id/dismiss
Authorization: Bearer <token>

Response:
{
  "message": "Suggestion dismissed"
}
```

---

## 🏗️ Phase 4: Implementation Roadmap

### Sprint 1: Foundation (1-2 weeks)
- [ ] Create database migrations for all tables
- [ ] Set up basic models and repositories
- [ ] Implement user_follows, page_follows, group_members
- [ ] Add follower count tracking

### Sprint 2: Contact Matching (1 week)
- [ ] Implement contact upload API
- [ ] Build contact matching algorithm
- [ ] Create background job for contact processing
- [ ] Add privacy settings for contact visibility

### Sprint 3: Interest-Based Suggestions (1-2 weeks)
- [ ] Implement interest matching algorithm
- [ ] Create suggestion scoring system
- [ ] Build suggestion cache layer (Redis)
- [ ] Implement GET /v1/suggestions/onboarding API

### Sprint 4: Trending System (1 week)
- [ ] Build trending score calculation
- [ ] Create background job for trending updates (cron)
- [ ] Implement trending API endpoints
- [ ] Add trending badges/indicators

### Sprint 5: Optimization & Caching (1 week)
- [ ] Add Redis caching for suggestions
- [ ] Implement pagination and infinite scroll
- [ ] Optimize database queries
- [ ] Add rate limiting

### Sprint 6: Testing & Polish (1 week)
- [ ] Unit tests for recommendation algorithms
- [ ] Integration tests for APIs
- [ ] Load testing with realistic data
- [ ] UI/UX improvements

---

## 🔧 Technology Stack

### Backend
- **Go Fiber**: API endpoints
- **PostgreSQL**: Primary database
- **Redis**: Caching layer for suggestions
- **Background Jobs**: Trending calculations, contact matching
- **Elasticsearch** (optional): Full-text search for users/pages/groups

### Algorithms
- **Jaccard Similarity**: Interest matching
- **Collaborative Filtering**: User-based recommendations
- **Graph Analysis**: Mutual connections
- **Time-decay Functions**: Trending calculations

### Performance
- **Caching Strategy**: 
  - Suggestion results: 15 minutes TTL
  - Trending data: 1 hour TTL
  - User stats: 5 minutes TTL
- **Background Jobs**: Run every hour for trending updates
- **Database Indexing**: All foreign keys and scoring columns

---

## 📊 Phase 5: Metrics & Analytics

### Track:
1. **Suggestion Quality**
   - Click-through rate (CTR)
   - Follow rate per suggestion
   - Time to first follow

2. **User Engagement**
   - Average suggestions viewed
   - Bulk follow adoption rate
   - Contact upload rate

3. **System Performance**
   - API response times
   - Cache hit rate
   - Background job duration

---

## 🎯 Success Metrics

- **70%+ follow rate** from suggestions
- **5+ follows** per new user during onboarding
- **<200ms** API response time for suggestions
- **90%+ cache hit rate** for repeated requests

---

## 🚀 Future Enhancements

1. **Machine Learning**
   - Train ML model on user behavior
   - Personalized ranking
   - A/B testing different algorithms

2. **Real-time Updates**
   - WebSocket for live trending updates
   - Push notifications for contact joins

3. **Advanced Features**
   - "People you may know" sidebar
   - Weekly digest of new suggestions
   - Suggestion quality feedback loop

4. **Social Features**
   - See what your friends are following
   - Group suggestions based on friend activity
   - Collaborative interests

---

## 📝 Next Steps

1. Review and approve this plan
2. Set up development environment
3. Start with Sprint 1: Database migrations
4. Build MVP with basic interest-based suggestions
5. Iterate based on user feedback

---

**Estimated Total Time**: 6-8 weeks for full implementation
**MVP Time**: 2-3 weeks (interest-based suggestions only)


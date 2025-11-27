# MVP1 Work Checklist - Mockhu Backend

**Last Updated:** 2024-11-26  
**Total Estimated Time:** ~8-10 weeks remaining

---

## ✅ COMPLETED WORK

### Core Systems (Pre-existing) ✅
- [x] User authentication (email/phone)
- [x] Email/phone verification codes
- [x] JWT token authentication
- [x] User profile creation (basic)
- [x] Interest selection system
- [x] Onboarding flow
- [x] File upload (avatar)

### Section 1: Follow System ✅
- [x] Follow/unfollow users
- [x] Get followers list (paginated)
- [x] Get following list (paginated)
- [x] Check follow status
- [x] Get follow statistics
- [x] Postman collection
- [x] Documentation
- [x] Tests

### Section 2: Posts System ✅
- [x] Create posts (text + images)
- [x] Get single post
- [x] Get user posts
- [x] Delete posts
- [x] Toggle reactions (like/unlike)
- [x] Get user feed
- [x] Post seeding
- [x] Postman collection
- [x] Tests

### Section 3: Comments System ✅
- [x] Create comments on posts
- [x] Get single comment
- [x] Get post comments (paginated)
- [x] Update comments
- [x] Delete comments
- [x] Single-level comment enforcement
- [x] Postman collection
- [x] Documentation
- [x] Tests

### Section 4: Shares System ✅
- [x] Share posts (timeline, dm, external)
- [x] Get single share
- [x] Get post shares (paginated)
- [x] Get share count
- [x] Get user shares (paginated)
- [x] Delete shares
- [x] Duplicate share prevention
- [x] Postman collection
- [x] Documentation
- [x] Tests

**Total Completed: ~35 API endpoints across 4 sections**

---

## 🔨 PRIORITY 1: CRITICAL PATH (Must Have for Launch)

### Feature 4: User Profiles & Following Enhancement 📋 Planned
**Time:** 1 week | **Status:** Planning Complete | **Plan:** `USER_PROFILE_SYSTEM_PLAN.md`

#### Database
- [ ] Create migration `000014_add_profile_fields.up.sql`
  - [ ] Add bio field (TEXT)
  - [ ] Add institution_id field (UUID)
  - [ ] Add privacy settings (who_can_message, who_can_see_posts)
  - [ ] Add visibility settings (show_followers_list, show_following_list)
  - [ ] Add constraints and indexes
- [ ] Test migration up/down
- [ ] Update User model in auth package

#### Profile Viewing (2 endpoints)
- [ ] `GET /v1/users/:userId/profile` - Get user profile (public)
  - [ ] Implement handler
  - [ ] Add stats calculation (posts, followers, following counts)
  - [ ] Add mutual connections count
  - [ ] Respect privacy settings
  - [ ] Test with various users
- [ ] `GET /v1/users/me/profile` - Get own profile (private)
  - [ ] Implement handler
  - [ ] Return all private fields
  - [ ] Include privacy settings
  - [ ] Test response

#### Profile Management (3 endpoints)
- [ ] `PUT /v1/users/me/profile` - Update profile
  - [ ] Implement handler
  - [ ] Add validation (username, bio length, etc.)
  - [ ] Check username uniqueness
  - [ ] Test with valid/invalid data
- [ ] `POST /v1/users/me/avatar` - Upload avatar
  - [ ] Integrate image processing library
  - [ ] Add file validation (type, size)
  - [ ] Resize to 400x400
  - [ ] Store in S3 or local storage
  - [ ] Delete old avatar
  - [ ] Test with various files
- [ ] `DELETE /v1/users/me/avatar` - Delete avatar
  - [ ] Implement handler
  - [ ] Delete file from storage
  - [ ] Update database
  - [ ] Test deletion

#### Privacy Settings (2 endpoints)
- [ ] `GET /v1/users/me/privacy` - Get privacy settings
  - [ ] Implement handler
  - [ ] Test response
- [ ] `PUT /v1/users/me/privacy` - Update privacy settings
  - [ ] Implement handler
  - [ ] Add validation
  - [ ] Test all privacy combinations

#### Mutual Connections (2 endpoints)
- [ ] `GET /v1/users/:userId/mutual-connections` - Get mutual connections
  - [ ] Write optimized SQL query
  - [ ] Implement pagination
  - [ ] Test with various user combinations
- [ ] `GET /v1/users/:userId/mutual-connections/count` - Get count
  - [ ] Implement handler
  - [ ] Add caching (optional)
  - [ ] Test accuracy

#### Privacy Enforcement
- [ ] Add privacy checks to Follow service
  - [ ] Hide follower/following lists based on settings
- [ ] Add privacy checks to Post service
  - [ ] Filter posts based on who_can_see_posts
- [ ] Prepare for Messaging system
  - [ ] Check who_can_message before allowing DM

#### Testing & Documentation
- [ ] Unit tests for profile service
- [ ] Integration tests for all 9 endpoints
- [ ] Privacy enforcement tests
- [ ] Create test script `test_profile.sh`
- [ ] Create Postman collection
- [ ] Write API documentation
- [ ] Update SECTIONS_STATUS.md

**Deliverables:** 9 endpoints, avatar upload, privacy system

---

### Feature 1: Student Verification System 🔨 To Plan
**Time:** 1 week | **Status:** Not Started | **Priority:** P1

#### Planning Phase
- [ ] Create detailed implementation plan
- [ ] Design database schema
- [ ] Define API endpoints
- [ ] Plan verification workflow

#### Implementation Checklist (Create after planning)
- [ ] Database schema (verification_requests table)
- [ ] .edu email auto-verification
- [ ] Student ID upload endpoint
- [ ] Admin review queue
- [ ] Verification badge system
- [ ] Email notifications
- [ ] API endpoints (5-7)
- [ ] Tests
- [ ] Documentation

**Features:**
- .edu email domain verification (auto-approve)
- Student ID upload for non-.edu emails
- Manual review queue for admins
- Verification badge on profiles
- Block unverified users from full access

---

### Feature 2: Institution System 🔨 To Plan
**Time:** 1.5 weeks | **Status:** Not Started | **Priority:** P1

#### Planning Phase
- [ ] Create detailed implementation plan
- [ ] Design database schema
- [ ] Define API endpoints
- [ ] Plan institution discovery

#### Implementation Checklist (Create after planning)
- [ ] Database schema (institutions table)
- [ ] Seed popular institutions
- [ ] Institution CRUD operations
- [ ] Institution search/autocomplete
- [ ] Institution profile pages
- [ ] Follow institution
- [ ] Institution-specific feed
- [ ] API endpoints (8-10)
- [ ] Tests
- [ ] Documentation

**Features:**
- Institution database (colleges, schools)
- Search & select institution during signup
- Institution profile pages
- Follow institution
- Institution-specific feed

---

### Feature 5: Direct Messaging (DM) 🔨 To Plan
**Time:** 2 weeks | **Status:** Planning Complete | **Priority:** P1  
**Plan:** `DIRECT_MESSAGING_SYSTEM_PLAN.md`

#### Planning Review
- [ ] Review existing plan
- [ ] Update if needed based on recent changes

#### Implementation Checklist (Create after planning review)
- [ ] Database schema (4 tables)
- [ ] Conversation management (4 endpoints)
- [ ] Message CRUD (5 endpoints)
- [ ] Block management (4 endpoints)
- [ ] Image upload (1 endpoint)
- [ ] WebSocket server setup
- [ ] Real-time message delivery
- [ ] Polling fallback
- [ ] Unread count system
- [ ] Tests
- [ ] Documentation

**Features:**
- 1-on-1 chat
- Text messages
- Image sharing
- Message history
- Unread count
- Block users
- Real-time delivery (WebSocket or polling)

---

## 🔨 PRIORITY 2: IMPORTANT (Enhance Platform)

### Feature 6: Groups System 🔨 To Plan
**Time:** 2 weeks | **Status:** Not Started | **Priority:** P2

#### Planning Phase
- [ ] Create detailed implementation plan
- [ ] Design database schema
- [ ] Define API endpoints
- [ ] Plan group permissions

#### Features to Plan:
- Create/join/leave groups
- Group feed (posts within group)
- Group chat
- Group admin controls
- Public/private groups
- Search groups

---

### Feature 7: Search & Explore 🔨 To Plan
**Time:** 1.5 weeks | **Status:** Not Started | **Priority:** P2

#### Planning Phase
- [ ] Create detailed implementation plan
- [ ] Design search algorithm
- [ ] Define API endpoints
- [ ] Plan indexing strategy

#### Features to Plan:
- Search users (name, institution, course)
- Search groups (name, category)
- Explore page (trending posts, popular groups)
- Filter by institution/interests
- Autocomplete suggestions

---

### Feature 8: Notifications 🔨 To Plan
**Time:** 1 week | **Status:** Not Started | **Priority:** P2

#### Planning Phase
- [ ] Create detailed implementation plan
- [ ] Design notification system
- [ ] Define API endpoints
- [ ] Plan push notification integration

#### Features to Plan:
- In-app notifications
- Push notifications (mobile)
- Email notifications (optional)
- Notification types (follower, like, comment, message, group invite, mention)

---

## 🌟 PRIORITY 3: NICE TO HAVE (Polish & Growth)

### Feature 9: Enhanced Suggestions 🔨 To Plan
**Time:** 1 week | **Status:** Not Started | **Priority:** P3

#### Planning Phase
- [ ] Create detailed implementation plan
- [ ] Design suggestion algorithm
- [ ] Define API endpoints

#### Features to Plan:
- Suggest users from same institution
- Suggest groups based on interests
- "People you may know" section
- Weekly suggestion emails
- Dismiss suggestions

---

### Feature 10: User Activity & Stats 🔨 To Plan
**Time:** 0.5 week | **Status:** Not Started | **Priority:** P3

#### Features to Plan:
- Profile views counter
- Post impressions
- Engagement stats (likes, comments)
- Activity history

---

## 📊 SUMMARY PROGRESS

### Completion Status

| Feature | Priority | Status | Time | Endpoints | Progress |
|---------|----------|--------|------|-----------|----------|
| Auth System | P1 | ✅ Done | - | 10 | 100% |
| Follow System | P1 | ✅ Done | - | 6 | 100% |
| Posts System | P1 | ✅ Done | - | 6 | 100% |
| Comments System | P1 | ✅ Done | - | 5 | 100% |
| Shares System | P1 | ✅ Done | - | 6 | 100% |
| **User Profiles** | **P1** | **📋 Planned** | **1w** | **9** | **0%** |
| **Student Verification** | **P1** | **🔨 Todo** | **1w** | **5-7** | **0%** |
| **Institution System** | **P1** | **🔨 Todo** | **1.5w** | **8-10** | **0%** |
| **Direct Messaging** | **P1** | **📋 Planned** | **2w** | **17+WS** | **0%** |
| Groups System | P2 | 🔨 Todo | 2w | ~12 | 0% |
| Search & Explore | P2 | 🔨 Todo | 1.5w | ~8 | 0% |
| Notifications | P2 | 🔨 Todo | 1w | ~6 | 0% |
| Enhanced Suggestions | P3 | 🔨 Todo | 1w | ~10 | 0% |
| User Stats | P3 | 🔨 Todo | 0.5w | ~5 | 0% |

### Overall Progress
- **Completed:** 33 endpoints (5 features)
- **Planned:** 26+ endpoints (2 features)
- **To Do:** ~70 endpoints (7 features)

**Total MVP1 Progress: ~40% Complete**

---

## 📅 RECOMMENDED IMPLEMENTATION ORDER

### Week 1 (Current Week)
1. ✅ Plan User Profiles (Done)
2. 🔨 Implement User Profiles System (9 endpoints)
   - Day 1-2: Profile viewing & management
   - Day 3-4: Privacy settings & avatar
   - Day 5: Mutual connections
   - Day 6: Privacy enforcement
   - Day 7: Testing & docs

### Week 2-3
1. 📋 Plan Student Verification
2. 📋 Plan Institution System
3. 🔨 Implement Student Verification (1 week)
4. 🔨 Implement Institution System (1.5 weeks)

### Week 4-5
1. 📋 Review DM Plan
2. 🔨 Implement Direct Messaging (2 weeks)
   - HTTP endpoints
   - WebSocket server
   - Real-time delivery

### Week 6-7
1. 📋 Plan Groups System
2. 🔨 Implement Groups System (2 weeks)

### Week 8
1. 📋 Plan Search & Notifications
2. 🔨 Implement Search & Explore (1.5 weeks)

### Week 9
1. 🔨 Implement Notifications (1 week)

### Week 10 (Optional)
1. 📋 Plan Suggestions & Stats
2. 🔨 Implement nice-to-have features

---

## 🎯 IMMEDIATE NEXT STEPS

### Today (Planning Complete)
- [x] Plan User Profiles & Following ✅

### Tomorrow (Start Implementation)
- [ ] Create migration for profile fields
- [ ] Implement profile viewing endpoints
- [ ] Start profile management

### This Week
- [ ] Complete User Profiles implementation
- [ ] Test all endpoints
- [ ] Create Postman collection
- [ ] Write documentation

### Next Week
- [ ] Plan Student Verification System
- [ ] Plan Institution System
- [ ] Start implementation

---

## 📦 DELIVERABLES CHECKLIST

### Per Feature
- [ ] Database migrations (up/down)
- [ ] Domain models
- [ ] DTOs (request/response)
- [ ] Repository interface
- [ ] Repository implementation
- [ ] Service layer (business logic)
- [ ] HTTP handlers
- [ ] Routes registration
- [ ] Unit tests
- [ ] Integration tests
- [ ] Test scripts
- [ ] Postman collection
- [ ] API documentation
- [ ] Update SECTIONS_STATUS.md

### Before Launch
- [ ] All P1 features complete
- [ ] Performance optimization
- [ ] Security audit
- [ ] Documentation complete
- [ ] Postman collections for all features
- [ ] Deployment scripts
- [ ] Monitoring setup
- [ ] Error logging (Sentry)
- [ ] Analytics integration

---

## 🚀 LAUNCH CRITERIA

### Technical Requirements
- [ ] All P1 features implemented and tested
- [ ] < 500ms API response time (95th percentile)
- [ ] 99.9% uptime
- [ ] Database indexed and optimized
- [ ] SSL/TLS configured
- [ ] Rate limiting implemented
- [ ] Error handling across all endpoints

### Quality Requirements
- [ ] All endpoints tested (unit + integration)
- [ ] No critical bugs
- [ ] Code reviewed
- [ ] Documentation complete
- [ ] Postman collections working

### Business Requirements
- [ ] Student verification working
- [ ] Institution system live
- [ ] Messaging functional
- [ ] Privacy controls working
- [ ] Content moderation ready

---

## 💡 NOTES

### Current Focus
- **Now:** User Profiles & Following implementation
- **Next:** Student Verification & Institution System
- **Then:** Direct Messaging

### Time Estimate
- **P1 Features Remaining:** ~5.5 weeks
- **P2 Features:** ~4.5 weeks
- **P3 Features:** ~1.5 weeks
- **Total Remaining:** ~8-10 weeks for full MVP1

### Dependencies
- User Profiles needed before DM (privacy checks)
- Institution System needed for verification
- Groups need Posts & Profiles complete
- Search needs all content types ready
- Notifications need all features to notify about

---

**Last Updated:** 2024-11-26  
**Status:** 40% Complete (5/13 features done)  
**Next Milestone:** Complete P1 Critical Path (5 more features)  

🎯 **Goal:** Launch-ready MVP1 in 8-10 weeks


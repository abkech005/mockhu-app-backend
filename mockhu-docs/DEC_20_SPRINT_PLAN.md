# 🚀 Mockhu Mini MVP - December 20 Sprint Plan

---

## 🎯 Goal: Launch-Ready Product by Dec 20, 2025

**Available Time:** 25 days (Nov 25 - Dec 20)  
**Strategy:** Ultra-focused MVP with ONLY essential features  
**Target:** Soft launch to 1 college with 100-200 students

---

## ✅ What We Already Have (Week 0 - Done!)

- [x] User authentication (email/phone)
- [x] Email/phone verification  
- [x] User profile creation
- [x] Interest selection & storage
- [x] Contact sync & matching
- [x] User suggestions (interest-based)
- [x] JWT authentication
- [x] Database setup (PostgreSQL)
- [x] API framework (Go + Fiber)

**Estimated Progress:** ~20% of full MVP1 complete

---

## 🎯 December 20 Mini MVP Scope

### **Core Philosophy:**
> "Build the MINIMUM needed for students to:
> 1. Create profile
> 2. Find friends
> 3. Post updates
> 4. Chat with friends"

**Think:** Early Facebook (2004) - Just profiles, friends, and wall posts

---

## 📋 Must-Have Features (Essential)

### **1. Student Verification** ⭐⭐⭐
```
What to Build:
✅ .edu email auto-verification
✅ Manual student ID upload (for non-.edu)
✅ Admin review queue (simple table + API)
✅ Verification badge on profile

Skip for Now:
❌ Automated ID verification
❌ Third-party verification services

Time: 2 days
```

### **2. Institution System (Basic)** ⭐⭐⭐
```
What to Build:
✅ Institution database (pre-populate top 100 colleges)
✅ Search & select during signup
✅ Display institution on profile
✅ Filter users by institution

Skip for Now:
❌ Institution pages
❌ Institution-specific feeds
❌ Follow institutions

Time: 2 days
```

### **3. Posts & Feed** ⭐⭐⭐
```
What to Build:
✅ Create text posts (no images yet)
✅ Simple chronological feed (from people you follow)
✅ Like posts (no comments yet)
✅ Delete own posts

Skip for Now:
❌ Image uploads
❌ Comments
❌ Share functionality
❌ Algorithmic feed
❌ Video posts

Time: 4 days
```

### **4. Follow System** ⭐⭐⭐
```
What to Build:
✅ Follow/unfollow users
✅ View user profiles
✅ Follower/following count
✅ Follower/following lists (simple)

Skip for Now:
❌ Mutual connections display
❌ Friend requests
❌ Privacy settings

Time: 2 days
```

### **5. Direct Messaging (Ultra Simple)** ⭐⭐⭐
```
What to Build:
✅ Send text messages (1-on-1)
✅ Message history
✅ Message list (conversations)
✅ Polling for new messages (5-second interval)

Skip for Now:
❌ WebSocket/real-time
❌ Image sharing
❌ Read receipts
❌ Typing indicators
❌ Group chats

Time: 3 days
```

### **6. Search (Basic)** ⭐⭐⭐
```
What to Build:
✅ Search users by name
✅ Filter by institution

Skip for Now:
❌ Advanced filters
❌ Search groups
❌ Autocomplete
❌ Search posts

Time: 1 day
```

### **7. Notifications (Minimal)** ⭐⭐
```
What to Build:
✅ In-app notification icon with count
✅ Notification list (basic)
✅ Types: new follower, new message, post like

Skip for Now:
❌ Push notifications
❌ Email notifications
❌ Rich notifications
❌ Notification settings

Time: 2 days
```

### **8. Enhanced Suggestions** ⭐⭐
```
What to Build:
✅ "Suggested for you" section on home
✅ Filter by: same institution, shared interests
✅ Dismiss suggestions

Already 70% Done!

Time: 1 day
```

---

## 📊 December 20 Mini MVP Timeline

### **Week 1: Nov 25 - Dec 1 (7 days)**

**Day 1-2 (Nov 25-26): Student Verification**
- [ ] Database migration for verification table
- [ ] .edu email verification logic
- [ ] Student ID upload API
- [ ] Admin review queue
- [ ] Verification badge

**Day 3-4 (Nov 27-28): Institution System**
- [ ] Institution database & migration
- [ ] Seed script (top 100 colleges)
- [ ] Institution search API
- [ ] Update user profile to include institution
- [ ] Institution filter in suggestions

**Day 5-7 (Nov 29-Dec 1): Posts & Feed - Part 1**
- [ ] Post model & migration
- [ ] Create post API (text only)
- [ ] Feed generation (chronological)
- [ ] Like post API
- [ ] Delete post API

---

### **Week 2: Dec 2 - Dec 8 (7 days)**

**Day 8-9 (Dec 2-3): Follow System**
- [ ] Follow/unfollow API
- [ ] User profile API (with posts)
- [ ] Follower/following lists
- [ ] Follow button logic
- [ ] Update feed to show followed users' posts

**Day 10-12 (Dec 4-6): Direct Messaging**
- [ ] Message model & migration
- [ ] Send message API
- [ ] Get messages API
- [ ] Conversation list API
- [ ] Polling endpoint for new messages
- [ ] Message count badge

**Day 13-14 (Dec 7-8): Search & Discover**
- [ ] User search API (by name)
- [ ] Institution filter
- [ ] Enhanced suggestions (same institution)
- [ ] Explore page API

---

### **Week 3: Dec 9 - Dec 15 (7 days)**

**Day 15-16 (Dec 9-10): Notifications**
- [ ] Notification model & migration
- [ ] Create notification on events (follow, like, message)
- [ ] Get notifications API
- [ ] Mark as read API
- [ ] Notification count badge

**Day 17-18 (Dec 11-12): Frontend Development**
- [ ] Set up React app (or Next.js)
- [ ] Authentication pages (login, signup)
- [ ] Profile setup flow
- [ ] Home feed page
- [ ] Create post component

**Day 19-20 (Dec 13-14): Frontend Continued**
- [ ] User profile page
- [ ] Follow/unfollow UI
- [ ] Search page
- [ ] Suggestions section

**Day 21 (Dec 15): Frontend Messaging**
- [ ] Message list page
- [ ] Chat interface
- [ ] Send message UI

---

### **Week 4: Dec 16 - Dec 20 (5 days)**

**Day 22-23 (Dec 16-17): Testing & Bug Fixes**
- [ ] End-to-end testing
- [ ] Fix critical bugs
- [ ] Performance optimization
- [ ] Mobile responsive checks

**Day 24 (Dec 18): Polish & Preparation**
- [ ] Landing page
- [ ] Terms of Service
- [ ] Privacy Policy
- [ ] Help/FAQ page
- [ ] Deploy to production

**Day 25 (Dec 19): Soft Launch Prep**
- [ ] Invite 20 beta testers
- [ ] Create onboarding guide
- [ ] Set up analytics
- [ ] Prepare launch announcement

**Day 26 (Dec 20): 🎉 LAUNCH DAY!**
- [ ] Monitor for issues
- [ ] Respond to feedback
- [ ] Fix urgent bugs
- [ ] Celebrate! 🎊

---

## 🎯 What's IN vs OUT (December 20 Mini MVP)

### **✅ IN (Must Have)**
- Student verification
- Institution selection
- Text posts
- Like posts
- Follow/unfollow
- User profiles
- Direct messages (polling)
- Basic search
- Suggestions (enhanced)
- Basic notifications

### **❌ OUT (Future Updates)**
- Image/video posts
- Comments on posts
- Groups
- Institution pages
- Stories
- Share functionality
- WebSocket real-time
- Push notifications
- Advanced search
- Explore page
- User stats
- Profile views
- Trending
- Email notifications
- Mobile apps (web only)

---

## 💪 Can We Actually Do This?

### **Assuming Full-Time Work (8+ hours/day):**
✅ **YES - Achievable but INTENSE**
- 25 days = 200+ hours of work
- Backend: ~100 hours
- Frontend: ~80 hours
- Testing/Polish: ~20 hours

### **If Working Part-Time (4 hours/day):**
⚠️ **MAYBE - Very tight, cut more features**
- Need to cut: Groups, DMs, or Notifications
- Focus on: Posts + Follow + Suggestions only

### **If Working Weekends Only:**
❌ **NO - Not enough time**
- Need at least 4+ weeks

---

## 🚨 Critical Success Factors

### **1. Stay Focused**
- Don't add ANY features not on the list
- Resist feature creep
- "Good enough" is good enough for launch

### **2. Use Existing Solutions**
- UI Components: Use Tailwind + DaisyUI or MUI
- Auth: You already have it ✅
- Image Uploads: Skip for now or use Cloudinary
- Real-time: Use polling, not WebSocket

### **3. Pre-Built Assets**
- Institution list: Download from Wikipedia
- Profile avatars: Use Dicebear (auto-generated)
- Default images: Use Unsplash API

### **4. Parallel Development**
- Backend and frontend in parallel
- Use mock data for frontend initially
- API-first development

### **5. Cut Scope If Needed**
**If running behind, cut in this order:**
1. Notifications (can add in v1.1)
2. Direct Messages (can add in v1.1)
3. Search (suggestions are enough initially)

**Never Cut:**
- Posts & Feed
- Follow system
- User profiles

---

## 📱 Frontend Strategy

### **Option 1: React SPA (Recommended)**
```
Pros:
✅ Fastest development
✅ You know React
✅ Single codebase
✅ Mobile PWA possible

Cons:
❌ No mobile apps yet
❌ SEO limitations

Time: Fits in 25 days
```

### **Option 2: Next.js (SSR)**
```
Pros:
✅ Better SEO
✅ Fast performance
✅ Modern

Cons:
❌ Slightly longer dev time
❌ More complex

Time: Tight but doable
```

### **Option 3: Mobile-First (React Native)**
```
Pros:
✅ Native apps
✅ Better UX

Cons:
❌ Takes longer (40+ days)
❌ App store approval delays

Time: Won't fit in 25 days ❌
```

**Recommendation:** React SPA with mobile-responsive design (PWA)

---

## 🎨 UI/UX Shortcuts

### **Use Component Libraries**
- **Tailwind CSS**: Fast styling
- **DaisyUI** or **Headless UI**: Pre-built components
- **React Icons**: Icons
- **React Hook Form**: Forms
- **React Query**: API state management

### **Design Inspiration**
- Copy UI from: Instagram (mobile), Twitter (feed), Facebook (profile)
- Use Figma templates (free student templates)
- Keep it SIMPLE - don't over-design

### **Templates to Use**
- Landing page: Tailwind UI
- Dashboard: Tailwind templates
- Chat: Existing React chat components

---

## 🚀 Launch Strategy (Dec 20)

### **Soft Launch Plan**
```
Target: Your college only
Users: 100-200 students
Duration: 2 weeks (Dec 20 - Jan 3)

Week 1 (Dec 20-26):
- Invite 20 close friends
- Get feedback
- Fix critical bugs
- Add small improvements

Week 2 (Dec 27-Jan 3):
- Post in college groups
- Grow to 100-200 users
- Monitor engagement
- Plan v1.1 features
```

### **Success Metrics**
- 100+ signups by Dec 31
- 50+ DAU (daily active users)
- 20+ posts per day
- 80% of users follow 5+ people
- <5 critical bugs

---

## 🛠️ Technology Stack (Final)

### **Backend (Current)**
```
✅ Go + Fiber
✅ PostgreSQL
✅ JWT Auth
✅ RESTful APIs

To Add:
🔲 AWS S3 (future - for images)
🔲 Redis (future - for caching)
🔲 Docker (for deployment)
```

### **Frontend (Choose One)**
```
Option 1: React + Vite + Tailwind ✅ (Recommended)
Option 2: Next.js + Tailwind
Option 3: React Native ❌ (Too slow)
```

### **Infrastructure**
```
Backend: DigitalOcean Droplet ($12/month)
Database: Managed PostgreSQL ($15/month)
Domain: Namecheap ($10/year)
CDN: CloudFlare (Free)

Total: ~$30/month
```

---

## 📊 Risk Assessment

### **High Risk ⚠️**
- **Frontend Development Time**: Never underestimate UI work
  - Mitigation: Use component libraries, keep design simple
  
- **Unexpected Bugs**: Always happen
  - Mitigation: Build in 3 buffer days

### **Medium Risk ⚠️**
- **Feature Creep**: "Just one more thing..."
  - Mitigation: Stick to the list, write "V1.1" for new ideas
  
- **Perfectionism**: "It's not ready yet"
  - Mitigation: Done is better than perfect

### **Low Risk ✅**
- **Backend Complexity**: You're experienced with Go
- **Database Design**: Straightforward models
- **Deployment**: Simple setup

---

## ✅ Daily Checklist Template

```markdown
## Day X - [Date]

### Goals:
- [ ] Feature 1
- [ ] Feature 2
- [ ] Feature 3

### Actual Progress:
- [x] What I completed
- [ ] What's blocked
- [ ] What's deferred

### Tomorrow:
- [ ] Next priority 1
- [ ] Next priority 2

### Blockers:
- None / [List blockers]

### Time Spent: X hours
```

---

## 🎯 Final Recommendation

### **Realistic Plan for Dec 20:**

**If Working Full-Time (8+ hours/day):**
```
✅ Build December 20 Mini MVP
✅ Include: Posts, Follow, DM, Suggestions
✅ Launch to 1 college
✅ 100-200 users by Dec 31
```

**If Working Part-Time (4 hours/day):**
```
⚠️ Cut DMs and Notifications
✅ Focus on: Posts + Follow + Suggestions
✅ Add DMs in v1.1 (early January)
```

**If You Have a Team:**
```
✅ Full Mini MVP possible
✅ Split: 1 backend + 1 frontend
✅ Could add more features
```

---

## 💡 My Honest Advice

**Option A: Aggressive Dec 20 Launch** ⚡
- Go all-in, work full-time
- Build Mini MVP (no images, no groups)
- Launch by Dec 20
- Get real users & feedback
- **Pro:** Real validation, momentum
- **Con:** Intense work, might burn out

**Option B: Conservative Jan 15 Launch** 🎯
- Work part-time (4-6 hours/day)
- Build fuller MVP (with images, groups)
- Launch by Jan 15 (4 more weeks)
- **Pro:** Less stress, better product
- **Con:** 4 weeks later, more time to overthink

**Option C: Hybrid - Dec 20 Private Beta** 🚀
- Build Mini MVP by Dec 20
- Launch to just 20 friends (private)
- Gather feedback over holidays
- Public launch Jan 15 with improvements
- **Pro:** Best of both worlds
- **Con:** Two launches instead of one

---

## 🎓 What Would I Do?

**If I were you:**

1. **Target Dec 20 for Private Beta** (20-30 users)
2. **Build Mini MVP without DMs** (add in v1.1)
3. **Focus on: Posts + Follow + Suggestions**
4. **Use holidays for feedback & improvements**
5. **Public launch Jan 15** with better product

**Why?**
- Less pressure, better quality
- Holiday break = more coding time
- Real user feedback before public launch
- Avoids burning out before launch

---

## 📅 Next Steps (Choose Your Path)

### **Path 1: Aggressive (Dec 20 Public)**
- Start coding TODAY
- Follow 25-day timeline strictly
- Skip holidays
- Launch publicly Dec 20

### **Path 2: Balanced (Dec 20 Private)**
- Start this week
- Build Mini MVP by Dec 20
- Private beta with friends
- Public launch Jan 15

### **Path 3: Conservative (Jan 15)**
- Start this week
- Work part-time
- Build fuller MVP
- Launch Jan 15

---

## 🚀 My Recommendation: Path 2 (Balanced)

**Target: Dec 20 Private Beta → Jan 15 Public Launch**

**Why This Is Best:**
- ✅ Achievable without burning out
- ✅ Real user feedback
- ✅ Better product at public launch
- ✅ Use holidays for improvement
- ✅ Less stress, more learning

**What do you think?** Want to go for Path 1, 2, or 3? 🎯


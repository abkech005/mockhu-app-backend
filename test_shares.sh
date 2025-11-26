#!/bin/bash

# Test script for Share System APIs (Section 4)
# This script tests all share-related endpoints

BASE_URL="http://localhost:8085"
TOKEN=""
POST_ID=""
SHARE_ID=""
USER_ID=""

echo "=========================================="
echo "Testing Share System APIs (Section 4)"
echo "=========================================="
echo ""

# Step 1: Signup a test user
echo "1. Signing up test user..."
SIGNUP_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "email",
    "email": "sharetest@example.com",
    "password": "testpass123"
  }')

echo "Signup Response: $SIGNUP_RESPONSE"
USER_ID=$(echo $SIGNUP_RESPONSE | grep -o '"user_id":"[^"]*' | cut -d'"' -f4)
echo "User ID: $USER_ID"
echo ""

# Step 2: Login to get token
echo "2. Logging in to get access token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "sharetest@example.com",
    "password": "testpass123"
  }')

echo "Login Response: $LOGIN_RESPONSE"
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
echo "Token: ${TOKEN:0:50}..."
echo ""

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token. Exiting."
  exit 1
fi

# Step 3: Get a post to share
echo "3. Getting a post to share..."
POSTS_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/posts?page=1&limit=1" \
  -H "Authorization: Bearer $TOKEN")
  
echo "Posts Response: $POSTS_RESPONSE"
POST_ID=$(echo $POSTS_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo "Post ID: $POST_ID"
echo ""

if [ -z "$POST_ID" ]; then
  echo "⚠️  No posts found. Creating a test post first..."
  # Create a test post
  CREATE_POST_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "content": "Test post for sharing",
      "images": []
    }')
  
  echo "Create Post Response: $CREATE_POST_RESPONSE"
  POST_ID=$(echo $CREATE_POST_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
  echo "Created Post ID: $POST_ID"
  echo ""
fi

if [ -z "$POST_ID" ]; then
  echo "❌ Cannot proceed without a post ID. Exiting."
  exit 1
fi

# Step 4: Test Public Endpoints (No Auth Required)
echo "=========================================="
echo "Testing PUBLIC Endpoints (No Auth)"
echo "=========================================="
echo ""

# 4.1: Get share count for a post (public)
echo "4.1. GET /v1/posts/:postId/shares/count (Public)"
SHARE_COUNT=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/shares/count")
echo "Response: $SHARE_COUNT"
echo ""

# 4.2: Get shares for a post (public)
echo "4.2. GET /v1/posts/:postId/shares (Public)"
POST_SHARES=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/shares?page=1&limit=10")
echo "Response: $POST_SHARES"
echo ""

# 4.3: Get user shares (public)
if [ ! -z "$USER_ID" ]; then
  echo "4.3. GET /v1/users/:userId/shares (Public)"
  USER_SHARES=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID/shares?page=1&limit=10")
  echo "Response: $USER_SHARES"
  echo ""
fi

# Step 5: Test Protected Endpoints (Auth Required)
echo "=========================================="
echo "Testing PROTECTED Endpoints (Auth Required)"
echo "=========================================="
echo ""

# 5.1: Create a share (protected)
echo "5.1. POST /v1/posts/:postId/shares (Protected) - Create Share"
CREATE_SHARE_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/shares" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shared_to_type": "timeline"
  }')

echo "Response: $CREATE_SHARE_RESPONSE"
SHARE_ID=$(echo $CREATE_SHARE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Share ID: $SHARE_ID"
echo ""

# 5.2: Get specific share (public, but test with auth)
if [ ! -z "$SHARE_ID" ]; then
  echo "5.2. GET /v1/shares/:shareId (Public)"
  GET_SHARE=$(curl -s -X GET "$BASE_URL/v1/shares/$SHARE_ID")
  echo "Response: $GET_SHARE"
  echo ""
fi

# 5.3: Try to create duplicate share (should fail)
echo "5.3. POST /v1/posts/:postId/shares (Protected) - Try Duplicate Share"
DUPLICATE_SHARE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/shares" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shared_to_type": "timeline"
  }')
echo "Response (should be error): $DUPLICATE_SHARE"
echo ""

# 5.4: Get updated share count
echo "5.4. GET /v1/posts/:postId/shares/count (Public) - After creating share"
UPDATED_COUNT=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/shares/count")
echo "Response: $UPDATED_COUNT"
echo ""

# 5.5: Get updated shares list
echo "5.5. GET /v1/posts/:postId/shares (Public) - After creating share"
UPDATED_SHARES=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/shares?page=1&limit=10")
echo "Response: $UPDATED_SHARES"
echo ""

# 5.6: Delete share (protected)
if [ ! -z "$SHARE_ID" ]; then
  echo "5.6. DELETE /v1/shares/:shareId (Protected) - Delete Share"
  DELETE_SHARE=$(curl -s -X DELETE "$BASE_URL/v1/shares/$SHARE_ID" \
    -H "Authorization: Bearer $TOKEN")
  echo "Response: $DELETE_SHARE"
  echo ""
fi

# 5.7: Test invalid share type
echo "5.7. POST /v1/posts/:postId/shares (Protected) - Invalid Share Type"
INVALID_TYPE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/shares" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shared_to_type": "invalid_type"
  }')
echo "Response (should be error): $INVALID_TYPE"
echo ""

# 5.8: Test without token (should fail)
echo "5.8. POST /v1/posts/:postId/shares (Protected) - Without Token"
NO_TOKEN=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/shares" \
  -H "Content-Type: application/json" \
  -d '{
    "shared_to_type": "timeline"
  }')
echo "Response (should be error): $NO_TOKEN"
echo ""

echo "=========================================="
echo "Share System API Testing Complete!"
echo "=========================================="


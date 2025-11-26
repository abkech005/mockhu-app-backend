#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8085"
ACCESS_TOKEN=""
USER_ID=""
USER_ID_2=""
POST_ID=""

echo -e "${YELLOW}=== Mockhu API Test Suite ===${NC}\n"

# Test 1: Signup User 1
echo -e "${GREEN}[1] Signup User 1${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "email",
    "email": "user1@test.com",
    "password": "password123"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
USER_ID=$(echo "$RESPONSE" | jq -r '.user_id' 2>/dev/null)
VERIFICATION_CODE=$(echo "$RESPONSE" | jq -r '.verification_code' 2>/dev/null)
echo "User ID: $USER_ID"
echo "Verification Code: $VERIFICATION_CODE"
echo ""

# Test 2: Verify User 1
echo -e "${GREEN}[2] Verify User 1${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/verify" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER_ID\",
    \"method\": \"email\",
    \"code\": \"$VERIFICATION_CODE\"
  }")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 3: Login User 1
echo -e "${GREEN}[3] Login User 1${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user1@test.com",
    "password": "password123"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token' 2>/dev/null)
echo "Access Token: ${ACCESS_TOKEN:0:50}..."
echo ""

# Test 4: Signup User 2
echo -e "${GREEN}[4] Signup User 2${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "email",
    "email": "user2@test.com",
    "password": "password123"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
USER_ID_2=$(echo "$RESPONSE" | jq -r '.user_id' 2>/dev/null)
VERIFICATION_CODE_2=$(echo "$RESPONSE" | jq -r '.verification_code' 2>/dev/null)
echo "User 2 ID: $USER_ID_2"
echo ""

# Test 5: Verify User 2
echo -e "${GREEN}[5] Verify User 2${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/verify" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER_ID_2\",
    \"method\": \"email\",
    \"code\": \"$VERIFICATION_CODE_2\"
  }")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 6: Login User 2
echo -e "${GREEN}[6] Login User 2${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user2@test.com",
    "password": "password123"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
ACCESS_TOKEN_2=$(echo "$RESPONSE" | jq -r '.access_token' 2>/dev/null)
echo ""

# Test 7: Get All Interests
echo -e "${GREEN}[7] Get All Interests${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/interests/")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 8: Complete Onboarding User 1
echo -e "${GREEN}[8] Complete Onboarding User 1${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/onboarding/complete" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "{
    \"first_name\": \"John\",
    \"last_name\": \"Doe\",
    \"username\": \"johndoe\",
    \"avatar_url\": \"https://example.com/avatar.jpg\",
    \"interests\": [\"technology\", \"programming\"]
  }")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 9: Complete Onboarding User 2
echo -e "${GREEN}[9] Complete Onboarding User 2${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/onboarding/complete" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN_2" \
  -d "{
    \"first_name\": \"Jane\",
    \"last_name\": \"Smith\",
    \"username\": \"janesmith\",
    \"avatar_url\": \"https://example.com/avatar2.jpg\",
    \"interests\": [\"technology\", \"design\"]
  }")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 10: Follow User 2
echo -e "${GREEN}[10] Follow User 2${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/users/$USER_ID_2/follow" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 11: Check Follow Status
echo -e "${GREEN}[11] Check Follow Status${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID_2/is-following" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 12: Get Follow Stats
echo -e "${GREEN}[12] Get Follow Stats (Public)${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID_2/follow-stats")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 13: Create Post
echo -e "${GREEN}[13] Create Post${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "content": "This is my first post! Testing the posts API.",
    "images": ["https://example.com/image1.jpg"],
    "is_anonymous": false
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
POST_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
echo "Post ID: $POST_ID"
echo ""

# Test 14: Get Post
echo -e "${GREEN}[14] Get Post${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 15: Toggle Reaction (Fire)
echo -e "${GREEN}[15] Toggle Reaction (Fire)${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/reactions" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 16: Toggle Reaction (Unfire)
echo -e "${GREEN}[16] Toggle Reaction (Unfire)${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/reactions" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 17: Create Another Post
echo -e "${GREEN}[17] Create Another Post${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN_2" \
  -d '{
    "content": "User 2 post! This is a test post from user 2.",
    "images": [],
    "is_anonymous": false
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
POST_ID_2=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
echo "Post 2 ID: $POST_ID_2"
echo ""

# Test 18: Get User Posts
echo -e "${GREEN}[18] Get User Posts${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID/posts?page=1&limit=10")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 19: Get Feed
echo -e "${GREEN}[19] Get Feed${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/feed?page=1&limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 20: Get Followers
echo -e "${GREEN}[20] Get Followers${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID_2/followers?page=1&limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 21: Get Following
echo -e "${GREEN}[21] Get Following${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID/following?page=1&limit=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 22: Delete Post
echo -e "${GREEN}[22] Delete Post${NC}"
RESPONSE=$(curl -s -X DELETE "$BASE_URL/v1/posts/$POST_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 23: Verify Post Deleted
echo -e "${GREEN}[23] Verify Post Deleted (should return 404)${NC}"
RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X GET "$BASE_URL/v1/posts/$POST_ID")
echo "$RESPONSE"
echo ""

echo -e "${YELLOW}=== Test Suite Complete ===${NC}"
echo "User 1 ID: $USER_ID"
echo "User 2 ID: $USER_ID_2"
echo "Post ID: $POST_ID"
echo "Post 2 ID: $POST_ID_2"


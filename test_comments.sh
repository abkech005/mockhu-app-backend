#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8085"
ACCESS_TOKEN=""
USER_ID=""
POST_ID=""
COMMENT_ID=""
REPLY_ID=""

echo -e "${YELLOW}=== Comments System Test Suite ===${NC}\n"

# Test 1: Login User 1
echo -e "${GREEN}[1] Login User 1${NC}"
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

# Test 2: Get or Create a Post
echo -e "${GREEN}[2] Get Existing Post or Create New${NC}"
# Try to get posts for user
POSTS_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/users/$USER_ID/posts?limit=1" \
  -H "Authorization: Bearer $ACCESS_TOKEN" 2>/dev/null)

POST_ID=$(echo "$POSTS_RESPONSE" | jq -r '.posts[0].id' 2>/dev/null)

if [ -z "$POST_ID" ] || [ "$POST_ID" = "null" ]; then
    echo "No existing post found. Creating a new post..."
    POST_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $ACCESS_TOKEN" \
      -d '{
        "content": "This is a test post for comments!",
        "images": [],
        "is_anonymous": false
      }')
    echo "$POST_RESPONSE" | jq '.' 2>/dev/null || echo "$POST_RESPONSE"
    POST_ID=$(echo "$POST_RESPONSE" | jq -r '.id' 2>/dev/null)
fi

echo "Post ID: $POST_ID"
echo ""

# Test 3: Create Comment
echo -e "${GREEN}[3] Create Comment on Post${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "content": "This is a great post! Thanks for sharing.",
    "is_anonymous": false
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
COMMENT_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
echo "Comment ID: $COMMENT_ID"
echo ""

# Test 4: Get Comment
echo -e "${GREEN}[4] Get Single Comment${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/comments/$COMMENT_ID")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 5: Get All Comments for Post
echo -e "${GREEN}[5] Get All Comments for Post${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/comments?page=1&limit=10")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 6: Create Reply to Comment
echo -e "${GREEN}[6] Create Reply to Comment${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "{
    \"content\": \"I totally agree with this comment!\",
    \"parent_comment_id\": \"$COMMENT_ID\",
    \"is_anonymous\": false
  }")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
REPLY_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
echo "Reply ID: $REPLY_ID"
echo ""

# Test 7: Create Another Comment
echo -e "${GREEN}[7] Create Another Comment${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "content": "This is my second comment on this post.",
    "is_anonymous": false
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
COMMENT_ID_2=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
echo "Comment 2 ID: $COMMENT_ID_2"
echo ""

# Test 8: Get All Comments Again (should show replies)
echo -e "${GREEN}[8] Get All Comments Again (with replies)${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/comments?page=1&limit=10")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 9: Update Comment
echo -e "${GREEN}[9] Update Comment${NC}"
RESPONSE=$(curl -s -X PUT "$BASE_URL/v1/comments/$COMMENT_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "content": "Updated: This is a great post! Thanks for sharing. (Edited)"
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 10: Verify Update
echo -e "${GREEN}[10] Verify Comment Updated${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/comments/$COMMENT_ID")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 11: Create Anonymous Comment
echo -e "${GREEN}[11] Create Anonymous Comment${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "content": "This is an anonymous comment.",
    "is_anonymous": true
  }')
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
ANON_COMMENT_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
echo "Anonymous Comment ID: $ANON_COMMENT_ID"
echo ""

# Test 12: Test Invalid Content (too long)
echo -e "${RED}[12] Test Invalid Content (too long - should fail)${NC}"
RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d "{
    \"content\": \"$(python3 -c 'print(\"x\" * 2001)')\",
    \"is_anonymous\": false
  }")
echo "$RESPONSE"
echo ""

# Test 13: Test Invalid Parent Comment
echo -e "${RED}[13] Test Invalid Parent Comment (should fail)${NC}"
RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "content": "Trying to reply to non-existent comment",
    "parent_comment_id": "00000000-0000-0000-0000-000000000000",
    "is_anonymous": false
  }')
echo "$RESPONSE"
echo ""

# Test 14: Delete Comment
echo -e "${GREEN}[14] Delete Comment${NC}"
RESPONSE=$(curl -s -X DELETE "$BASE_URL/v1/comments/$COMMENT_ID_2" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 15: Verify Comment Deleted
echo -e "${GREEN}[15] Verify Comment Deleted (should return 404)${NC}"
RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X GET "$BASE_URL/v1/comments/$COMMENT_ID_2")
echo "$RESPONSE"
echo ""

# Test 16: Get Comments After Delete
echo -e "${GREEN}[16] Get Comments After Delete${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/comments?page=1&limit=10")
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Test 17: Test Unauthorized Update (should fail)
echo -e "${RED}[17] Test Unauthorized Update (should fail)${NC}"
# This would require a different user's token, so we'll skip for now
echo "Skipped - would need second user token"
echo ""

echo -e "${YELLOW}=== Test Suite Complete ===${NC}"
echo "Post ID: $POST_ID"
echo "Comment ID: $COMMENT_ID"
echo "Reply ID: $REPLY_ID"
echo "Comment 2 ID: $COMMENT_ID_2"
echo "Anonymous Comment ID: $ANON_COMMENT_ID"


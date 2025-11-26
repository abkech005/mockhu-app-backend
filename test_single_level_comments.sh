#!/bin/bash

# Test Single-Level Comment System

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

BASE_URL="http://localhost:8085"

echo -e "${YELLOW}=== Testing Single-Level Comment System ===${NC}\n"

# Login
echo -e "${GREEN}[1] Login${NC}"
TOKEN=$(curl -s -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"identifier":"user1@test.com","password":"password123"}' | jq -r '.access_token')
echo "Token obtained: ${TOKEN:0:30}..."
echo ""

# Get or create post
echo -e "${GREEN}[2] Get/Create Post${NC}"
POST_ID=$(curl -s -X GET "$BASE_URL/v1/users/26b32915-7228-46c9-b798-d9b8e6ba4601/posts?limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.posts[0].id')

if [ -z "$POST_ID" ] || [ "$POST_ID" = "null" ]; then
    POST_ID=$(curl -s -X POST "$BASE_URL/v1/posts" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"content":"Test post for single-level comments","images":[],"is_anonymous":false}' | jq -r '.id')
fi
echo "Post ID: $POST_ID"
echo ""

# Create top-level comment
echo -e "${GREEN}[3] Create Top-Level Comment${NC}"
COMMENT_ID=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"This is a top-level comment","is_anonymous":false}' | jq -r '.id')
echo "Comment ID: $COMMENT_ID"
echo ""

# Create reply to comment (should work)
echo -e "${GREEN}[4] Create Reply to Comment (should work)${NC}"
REPLY_RESPONSE=$(curl -s -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"content\":\"This is a reply to the comment\",\"parent_comment_id\":\"$COMMENT_ID\",\"is_anonymous\":false}")
REPLY_ID=$(echo "$REPLY_RESPONSE" | jq -r '.id')
if [ "$REPLY_ID" != "null" ] && [ -n "$REPLY_ID" ]; then
    echo -e "${GREEN}✅ SUCCESS: Reply created${NC}"
    echo "$REPLY_RESPONSE" | jq '{id, parent_comment_id, content}'
else
    echo -e "${RED}❌ FAILED: Could not create reply${NC}"
    echo "$REPLY_RESPONSE"
fi
echo ""

# Try to reply to reply (should fail)
echo -e "${RED}[5] Try to Reply to Reply (should FAIL)${NC}"
NESTED_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/v1/posts/$POST_ID/comments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"content\":\"Trying to reply to a reply\",\"parent_comment_id\":\"$REPLY_ID\",\"is_anonymous\":false}")
echo "$NESTED_RESPONSE"
if echo "$NESTED_RESPONSE" | grep -q "cannot reply to a reply"; then
    echo -e "${GREEN}✅ CORRECT: Reply to reply was rejected${NC}"
else
    echo -e "${RED}❌ ERROR: Should have rejected reply to reply${NC}"
fi
echo ""

# Get comments to verify structure
echo -e "${GREEN}[6] Get All Comments (verify structure)${NC}"
COMMENTS_RESPONSE=$(curl -s -X GET "$BASE_URL/v1/posts/$POST_ID/comments")
echo "$COMMENTS_RESPONSE" | jq '{
  total_comments: .comments | length,
  first_comment: .comments[0] | {id, content, parent_comment_id, reply_count, replies_count: (.replies | length)}
}'
echo ""

echo -e "${YELLOW}=== Test Complete ===${NC}"
echo "Post ID: $POST_ID"
echo "Comment ID: $COMMENT_ID"
echo "Reply ID: $REPLY_ID"


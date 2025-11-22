#!/bin/bash
set -e

echo "1. Checking Health..."
curl -s http://localhost:8080/health | jq .
echo ""

echo "2. Creating Event..."
CREATE_RESP=$(curl -s -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"type":"walk", "note":"test walk"}')
echo $CREATE_RESP | jq .
ID=$(echo $CREATE_RESP | jq .id)
echo "Created ID: $ID"
echo ""

echo "3. Listing Events..."
curl -s "http://localhost:8080/api/v1/events?limit=5" | jq .
echo ""

echo "4. Getting Event $ID..."
curl -s http://localhost:8080/api/v1/events/$ID | jq .
echo ""

echo "5. Deleting Event $ID..."
curl -s -X DELETE http://localhost:8080/api/v1/events/$ID -v 2>&1 | grep "< HTTP"
echo ""

echo "6. Verifying Deletion (should be 404)..."
curl -s http://localhost:8080/api/v1/events/$ID -v 2>&1 | grep "< HTTP"
echo "Done."

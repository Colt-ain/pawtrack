#!/bin/bash
set -e

echo "1. Creating User..."
CREATE_RESP=$(curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Smith", "email":"alice@example.com", "password":"securepass123"}')
echo $CREATE_RESP | jq .
ID=$(echo $CREATE_RESP | jq .id)
echo "Created User ID: $ID"
echo ""

echo "2. Listing Users..."
curl -s "http://localhost:8080/api/v1/users" | jq .
echo ""

echo "3. Updating User $ID..."
UPDATE_RESP=$(curl -s -X PUT http://localhost:8080/api/v1/users/$ID \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson", "email":"alice.johnson@example.com"}')
echo $UPDATE_RESP | jq .
echo ""

echo "4. Getting User $ID..."
curl -s http://localhost:8080/api/v1/users/$ID | jq .
echo ""

echo "5. Trying to create duplicate email (should fail)..."
curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob", "email":"alice.johnson@example.com", "password":"pass456"}' | jq .
echo ""

echo "6. Deleting User $ID..."
curl -s -X DELETE http://localhost:8080/api/v1/users/$ID -v 2>&1 | grep "< HTTP"
echo ""

echo "7. Verifying Deletion (should be 404)..."
curl -s http://localhost:8080/api/v1/users/$ID -v 2>&1 | grep "< HTTP"
echo "Done."

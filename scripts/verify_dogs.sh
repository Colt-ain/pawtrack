#!/bin/bash
set -e

echo "1. Creating Dog..."
CREATE_RESP=$(curl -s -X POST http://localhost:8080/api/v1/dogs \
  -H "Content-Type: application/json" \
  -d '{"name":"Buddy", "breed":"Golden Retriever", "birth_date":"2022-01-01T00:00:00Z"}')
echo $CREATE_RESP | jq .
ID=$(echo $CREATE_RESP | jq .id)
echo "Created Dog ID: $ID"
echo ""

echo "2. Listing Dogs..."
curl -s "http://localhost:8080/api/v1/dogs" | jq .
echo ""

echo "3. Updating Dog $ID..."
UPDATE_RESP=$(curl -s -X PUT http://localhost:8080/api/v1/dogs/$ID \
  -H "Content-Type: application/json" \
  -d '{"name":"Buddy Jr.", "breed":"Golden Retriever", "birth_date":"2022-01-01T00:00:00Z"}')
echo $UPDATE_RESP | jq .
echo ""

echo "4. Getting Dog $ID..."
curl -s http://localhost:8080/api/v1/dogs/$ID | jq .
echo ""

echo "5. Deleting Dog $ID..."
curl -s -X DELETE http://localhost:8080/api/v1/dogs/$ID -v 2>&1 | grep "< HTTP"
echo ""

echo "6. Verifying Deletion (should be 404)..."
curl -s http://localhost:8080/api/v1/dogs/$ID -v 2>&1 | grep "< HTTP"
echo "Done."

#!/bin/bash

# Configuration
API_URL="http://localhost:8080"
ENDPOINT="/leads/bulk"

echo "Checking if server is running at $API_URL..."
if ! curl -s "$API_URL/health" | grep "healthy" > /dev/null; then
  echo "Error: Server is not running. Please start the server first with 'go run cmd/server/main.go'"
  exit 1
fi

echo "Seeding data..."

curl -X POST "$API_URL$ENDPOINT" \
-H "Content-Type: application/json" \
-d '{
  "leads": [
    {
      "name": "Bhavya Jain",
      "email": "bhavya@example.com",
      "phone": "+91-9999999999",
      "source": "website"
    },
    {
      "name": "Alice Johnson",
      "email": "alice@example.com",
      "phone": "+1-555-123456",
      "source": "linkedin"
    },
    {
      "name": "Bob Smith",
      "email": "bob@example.com",
      "phone": "+44-777777777",
      "source": "referral"
    },
    {
      "name": "Charlie Brown",
      "email": "charlie@example.com",
      "phone": "+91-8888888888",
      "source": "ads"
    },
    {
      "name": "David Wilson",
      "email": "david@example.com",
      "phone": "+91-7777777777",
      "source": "cold-email"
    }
  ]
}'

echo -e "\n\nSeeding request sent. Please check the response above for details."

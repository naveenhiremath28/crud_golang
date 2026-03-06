#!/bin/sh

# Employee Service API - cURL Test Commands
# Usage: sh test-api.sh
# Prereq: Keycloak and the app must be running

BASE_URL="http://localhost:3001/api"

# ─────────────────────────────────────────────
# 1. Login (get access token)
# ─────────────────────────────────────────────
echo "=== LOGIN ==="
curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "api.request",
    "ver": "v1",
    "ts": "2026-03-06T00:00:00Z",
    "params": { "msgid": "test-login" },
    "request": {
      "username": "admin",
      "password": "admin"
    }
  }' | tee /tmp/login_response.json

echo "\n"

# Extract access token and refresh token from response
ACCESS_TOKEN=$(cat /tmp/login_response.json | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
REFRESH_TOKEN=$(cat /tmp/login_response.json | sed -n 's/.*"refresh_token":"\([^"]*\)".*/\1/p')

if [ -z "$ACCESS_TOKEN" ]; then
  echo "ERROR: Failed to get access token. Check if Keycloak and the app are running."
  exit 1
fi

echo "Token obtained successfully.\n"

# ─────────────────────────────────────────────
# 2. Refresh Token
# ─────────────────────────────────────────────
echo "=== REFRESH TOKEN ==="
curl -s -X POST "$BASE_URL/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"api.request\",
    \"ver\": \"v1\",
    \"ts\": \"2026-03-06T00:00:00Z\",
    \"params\": { \"msgid\": \"test-refresh\" },
    \"request\": {
      \"refresh_token\": \"$REFRESH_TOKEN\"
    }
  }"

echo "\n"

# ─────────────────────────────────────────────
# 3. Server Status
# ─────────────────────────────────────────────
echo "=== SERVER STATUS ==="
curl -s -X GET "$BASE_URL/v1/" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

echo "\n"

# ─────────────────────────────────────────────
# 4. Add Employee (admin only)
# ─────────────────────────────────────────────
echo "=== ADD EMPLOYEE ==="
curl -s -X POST "$BASE_URL/v1/addEmployee" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "id": "api.request",
    "ver": "v1",
    "ts": "2026-03-06T00:00:00Z",
    "params": { "msgid": "test-add" },
    "request": {
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "salary": 75000.00,
      "mobile": "9876543210"
    }
  }' | tee /tmp/add_employee_response.json

echo "\n"

# Extract employee ID from response
EMPLOYEE_ID=$(cat /tmp/add_employee_response.json | sed -n 's/.*"id":"\([0-9a-f-]*\)".*/\1/p' | head -1)

if [ -z "$EMPLOYEE_ID" ]; then
  echo "WARNING: Could not extract employee ID. Using placeholder for next requests.\n"
  EMPLOYEE_ID="REPLACE_WITH_EMPLOYEE_ID"
fi

# ─────────────────────────────────────────────
# 5. List Employees (user, manager, admin)
# ─────────────────────────────────────────────
echo "=== LIST EMPLOYEES ==="
curl -s -X GET "$BASE_URL/v1/listEmployees" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

echo "\n"

# ─────────────────────────────────────────────
# 6. Get Employee by ID (manager, admin)
# ─────────────────────────────────────────────
echo "=== GET EMPLOYEE ==="
curl -s -X GET "$BASE_URL/v1/getEmployee/$EMPLOYEE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

echo "\n"

# ─────────────────────────────────────────────
# 7. Update Employee (manager, admin)
# ─────────────────────────────────────────────
echo "=== UPDATE EMPLOYEE ==="
curl -s -X PATCH "$BASE_URL/v1/updateEmployee/$EMPLOYEE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "id": "api.request",
    "ver": "v1",
    "ts": "2026-03-06T00:00:00Z",
    "params": { "msgid": "test-update" },
    "request": {
      "first_name": "John",
      "last_name": "Doe Updated",
      "email": "john.updated@example.com",
      "salary": 85000.00,
      "mobile": "9876543210"
    }
  }'

echo "\n"

# ─────────────────────────────────────────────
# 8. Delete Employee (admin only)
# ─────────────────────────────────────────────
echo "=== DELETE EMPLOYEE ==="
curl -s -X DELETE "$BASE_URL/v1/deleteEmployee/$EMPLOYEE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

echo "\n"

# Cleanup
rm -f /tmp/login_response.json /tmp/add_employee_response.json

echo "=== DONE ==="

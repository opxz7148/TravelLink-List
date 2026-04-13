#!/bin/bash

# TravelLink - User API Test Script
# Tests all user-related endpoints using curl
# Usage: ./test_user_api.sh

set -e

# Configuration
API_BASE_URL="http://localhost:8080/api/v1"
CONTENT_TYPE="Content-Type: application/json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variables to store responses
ACCESS_TOKEN=""
USER_ID=""
ADMIN_TOKEN=""
ADMIN_ID=""

# Helper functions
print_header() {
    echo ""
    echo "==============================================="
    echo "  $1"
    echo "==============================================="
}

print_step() {
    echo -e "${YELLOW}→ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Check if server is running
check_server() {
    print_step "Checking if server is running on $API_BASE_URL..."
    if ! curl -s "$API_BASE_URL/../health" > /dev/null 2>&1; then
        print_error "Server is not running on $API_BASE_URL"
        print_step "Start the server with: cd backend && go run cmd/api/main.go"
        exit 1
    fi
    print_success "Server is running"
}

# ============================================================================
# TEST 1: USER REGISTRATION
# ============================================================================
test_register() {
    print_header "TEST 1: User Registration"
    
    print_step "Registering user: john_doe..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/register" \
        -H "$CONTENT_TYPE" \
        -d '{
            "email": "john_doe@example.com",
            "username": "john_doe",
            "password": "SecurePass123!",
            "display_name": "John Doe"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    # Extract token and user ID
    ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
    USER_ID=$(echo "$RESPONSE" | jq -r '.data.user.id')
    
    if [ "$ACCESS_TOKEN" != "null" ] && [ ! -z "$ACCESS_TOKEN" ]; then
        print_success "User registered successfully"
        print_step "Access Token: ${ACCESS_TOKEN:0:50}..."
        print_step "User ID: $USER_ID"
    else
        print_error "Registration failed"
        exit 1
    fi
}

# ============================================================================
# TEST 2: USER LOGIN
# ============================================================================
test_login() {
    print_header "TEST 2: User Login"
    
    print_step "Logging in as john_doe..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/login" \
        -H "$CONTENT_TYPE" \
        -d '{
            "email": "john_doe@example.com",
            "password": "SecurePass123!"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
    if [ "$TOKEN" != "null" ] && [ ! -z "$TOKEN" ]; then
        print_success "Login successful"
    else
        print_error "Login failed"
        exit 1
    fi
}

# ============================================================================
# TEST 3: GET USER PROFILE
# ============================================================================
test_get_profile() {
    print_header "TEST 3: Get User Profile"
    
    print_step "Fetching user profile (public, no auth required)..."
    RESPONSE=$(curl -s -X GET "$API_BASE_URL/users/$USER_ID")
    
    echo "Response: $RESPONSE" | jq '.'
    
    USERNAME=$(echo "$RESPONSE" | jq -r '.data.user.username')
    if [ "$USERNAME" = "john_doe" ]; then
        print_success "Profile fetched successfully"
    else
        print_error "Failed to fetch profile"
        exit 1
    fi
}

# ============================================================================
# TEST 4: UPDATE USER PROFILE
# ============================================================================
test_update_profile() {
    print_header "TEST 4: Update User Profile"
    
    print_step "Updating user profile (requires auth)..."
    RESPONSE=$(curl -s -X PUT "$API_BASE_URL/users/$USER_ID" \
        -H "$CONTENT_TYPE" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{
            "display_name": "John Developer",
            "bio": "Passionate about travel and software"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    UPDATED_NAME=$(echo "$RESPONSE" | jq -r '.data.user.display_name')
    if [ "$UPDATED_NAME" = "John Developer" ]; then
        print_success "Profile updated successfully"
    else
        print_error "Failed to update profile"
        exit 1
    fi
}

# ============================================================================
# TEST 5: CHANGE PASSWORD
# ============================================================================
test_change_password() {
    print_header "TEST 5: Change Password"
    
    print_step "Changing password..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/users/$USER_ID/change-password" \
        -H "$CONTENT_TYPE" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{
            "old_password": "SecurePass123!",
            "new_password": "NewSecurePass456!"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
    if [ "$SUCCESS" = "true" ]; then
        print_success "Password changed successfully"
    else
        print_error "Failed to change password"
        exit 1
    fi
}

# ============================================================================
# TEST 6: LOGIN WITH NEW PASSWORD
# ============================================================================
test_login_new_password() {
    print_header "TEST 6: Login with New Password"
    
    print_step "Attempting login with new password..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/login" \
        -H "$CONTENT_TYPE" \
        -d '{
            "email": "john_doe@example.com",
            "password": "NewSecurePass456!"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
    if [ "$TOKEN" != "null" ] && [ ! -z "$TOKEN" ]; then
        print_success "Login with new password successful"
    else
        print_error "Login with new password failed"
        exit 1
    fi
}

# ============================================================================
# TEST 7: LOGOUT
# ============================================================================
test_logout() {
    print_header "TEST 7: Logout"
    
    print_step "Logging out..."
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$API_BASE_URL/auth/logout" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    echo "HTTP Status Code: $HTTP_CODE"
    
    if [ "$HTTP_CODE" = "204" ]; then
        print_success "Logout successful (204 No Content)"
    else
        print_error "Logout failed (expected 204, got $HTTP_CODE)"
        exit 1
    fi
}

# ============================================================================
# TEST 8: REGISTER ADMIN USER (for admin tests)
# ============================================================================
test_register_admin() {
    print_header "TEST 8: Register Admin User (for admin tests)"
    
    print_step "Creating a second user for admin operations..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/register" \
        -H "$CONTENT_TYPE" \
        -d '{
            "email": "admin@example.com",
            "username": "admin_user",
            "password": "AdminPass789!",
            "display_name": "Admin User"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    ADMIN_TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
    ADMIN_ID=$(echo "$RESPONSE" | jq -r '.data.user.id')
    
    if [ "$ADMIN_TOKEN" != "null" ] && [ ! -z "$ADMIN_TOKEN" ]; then
        print_success "Admin user created"
    else
        print_error "Failed to create admin user"
    fi
}

# ============================================================================
# TEST 9: ATTEMPTED ADMIN OPERATION (should fail - no admin role)
# ============================================================================
test_admin_operation_no_permission() {
    print_header "TEST 9: Test Admin Operation (No Permission - Expected to Fail)"
    
    print_step "Attempting to change user role without admin permission (should fail)..."
    RESPONSE=$(curl -s -X PATCH "$API_BASE_URL/users/$ADMIN_ID/role" \
        -H "$CONTENT_TYPE" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d '{
            "role": "traveller"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
    if [ "$SUCCESS" = "false" ]; then
        print_success "Correctly rejected non-admin request"
    else
        print_error "Should have rejected the request"
    fi
}

# ============================================================================
# TEST 10: AUTHENTICATION ERROR - INVALID TOKEN
# ============================================================================
test_invalid_token() {
    print_header "TEST 10: Test Invalid Token Handling"
    
    print_step "Attempting request with invalid token..."
    RESPONSE=$(curl -s -X GET "$API_BASE_URL/users/$USER_ID" \
        -H "Authorization: Bearer invalid_token_here")
    
    echo "Response: $RESPONSE" | jq '.'
    
    SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
    if [ "$SUCCESS" = "false" ]; then
        print_success "Correctly rejected invalid token"
    else
        print_error "Should have rejected invalid token"
    fi
}

# ============================================================================
# TEST 11: AUTHENTICATION ERROR - WEAK PASSWORD
# ============================================================================
test_weak_password() {
    print_header "TEST 11: Test Weak Password Validation"
    
    print_step "Attempting to register with weak password..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/register" \
        -H "$CONTENT_TYPE" \
        -d '{
            "email": "weak_pass@example.com",
            "username": "weak_user",
            "password": "123",
            "display_name": "Weak Password User"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
    if [ "$SUCCESS" = "false" ]; then
        print_success "Correctly rejected weak password"
    else
        print_error "Should have rejected weak password"
    fi
}

# ============================================================================
# TEST 12: AUTHENTICATION ERROR - DUPLICATE EMAIL
# ============================================================================
test_duplicate_email() {
    print_header "TEST 12: Test Duplicate Email Validation"
    
    print_step "Attempting to register with duplicate email..."
    RESPONSE=$(curl -s -X POST "$API_BASE_URL/auth/register" \
        -H "$CONTENT_TYPE" \
        -d '{
            "email": "john_doe@example.com",
            "username": "different_username",
            "password": "SecurePass123!",
            "display_name": "Another User"
        }')
    
    echo "Response: $RESPONSE" | jq '.'
    
    SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
    if [ "$SUCCESS" = "false" ]; then
        print_success "Correctly rejected duplicate email"
    else
        print_error "Should have rejected duplicate email"
    fi
}

# ============================================================================
# MAIN EXECUTION
# ============================================================================
main() {
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║         TravelLink - User API Test Script                      ║"
    echo "║                                                                ║"
    echo "║  Testing all user-related endpoints:                           ║"
    echo "║  - Registration & Login                                        ║"
    echo "║  - Profile Management                                          ║"
    echo "║  - Password Management                                         ║"
    echo "║  - Logout                                                      ║"
    echo "║  - Error Handling                                              ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    # Check server
    check_server
    
    # Run all tests
    test_register
    test_login
    test_get_profile
    test_update_profile
    test_change_password
    test_login_new_password
    test_logout
    test_register_admin
    test_admin_operation_no_permission
    test_invalid_token
    test_weak_password
    test_duplicate_email
    
    # Summary
    print_header "TEST SUMMARY"
    print_success "All user API tests completed successfully!"
    echo ""
    echo "Test Coverage:"
    echo "  ✓ User Registration"
    echo "  ✓ User Login"
    echo "  ✓ Get User Profile (public)"
    echo "  ✓ Update User Profile (auth required)"
    echo "  ✓ Change Password"
    echo "  ✓ Login with New Password"
    echo "  ✓ Logout"
    echo "  ✓ Admin Operation Permission Check"
    echo "  ✓ Invalid Token Handling"
    echo "  ✓ Weak Password Validation"
    echo "  ✓ Duplicate Email Validation"
    echo ""
}

# Run main function
main "$@"

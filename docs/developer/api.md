# CMDB Lite API Documentation

This document provides detailed information about the CMDB Lite REST API, including endpoints, request/response formats, authentication, and examples.

## Table of Contents

- [API Overview](#api-overview)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [API Endpoints](#api-endpoints)
  - [Authentication Endpoints](#authentication-endpoints)
  - [Configuration Item Endpoints](#configuration-item-endpoints)
  - [Relationship Endpoints](#relationship-endpoints)
  - [Audit Log Endpoints](#audit-log-endpoints)
  - [User Endpoints](#user-endpoints)
- [Data Models](#data-models)
- [API Examples](#api-examples)
- [SDKs and Libraries](#sdks-and-libraries)

## API Overview

The CMDB Lite API is a RESTful API that allows programmatic access to all features of the CMDB Lite application. The API uses JSON for request and response bodies and follows standard HTTP conventions.

### Base URL

The base URL for the API depends on your deployment:

- **Development**: `http://localhost:8080/api`
- **Production**: `https://your-domain.com/api`

### API Versioning

The current API version is `v1`. All endpoints include the version in the path:

```
https://your-domain.com/api/v1/endpoint
```

### Content Types

The API uses the following content types:

- **Request Content-Type**: `application/json`
- **Response Content-Type**: `application/json`

## Authentication

All API endpoints (except authentication endpoints) require authentication using JWT (JSON Web Token) tokens.

### Obtaining a JWT Token

To obtain a JWT token, send a POST request to the `/auth/login` endpoint with your username and password:

```bash
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "your-username", "password": "your-password"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-12-31T23:59:59Z",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "your-username",
    "email": "your-email@example.com",
    "role": "admin"
  }
}
```

### Using the JWT Token

Include the JWT token in the Authorization header of your requests:

```bash
curl -X GET https://your-domain.com/api/v1/cis \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token Expiration

JWT tokens expire after 24 hours by default. When a token expires, you'll receive a 401 Unauthorized response. You need to obtain a new token by logging in again.

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests. In addition, error responses include a JSON object with details about the error.

### Error Response Format

```json
{
  "error": "Error message",
  "details": "Additional error details",
  "code": "ERROR_CODE"
}
```

### Common HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 OK | The request was successful. |
| 201 Created | The resource was successfully created. |
| 400 Bad Request | The request was malformed or contained invalid data. |
| 401 Unauthorized | Authentication is required or has failed. |
| 403 Forbidden | The authenticated user does not have permission to access the resource. |
| 404 Not Found | The requested resource does not exist. |
| 422 Unprocessable Entity | The request was well-formed but contains semantic errors. |
| 429 Too Many Requests | The client has exceeded the rate limit. |
| 500 Internal Server Error | An error occurred on the server. |
| 503 Service Unavailable | The service is temporarily unavailable. |

### Error Codes

| Error Code | Description |
|-------------|-------------|
| `INVALID_REQUEST` | The request was malformed or contained invalid data. |
| `UNAUTHORIZED` | Authentication is required or has failed. |
| `FORBIDDEN` | The authenticated user does not have permission. |
| `NOT_FOUND` | The requested resource does not exist. |
| `VALIDATION_ERROR` | The request data failed validation. |
| `DUPLICATE_ENTRY` | A resource with the same unique identifier already exists. |
| `FOREIGN_KEY_CONSTRAINT` | The operation would violate a foreign key constraint. |
| `RATE_LIMIT_EXCEEDED` | The client has exceeded the rate limit. |
| `INTERNAL_ERROR` | An unexpected error occurred on the server. |

## Rate Limiting

To ensure fair usage and prevent abuse, the API implements rate limiting:

- **Default Limit**: 100 requests per minute per authenticated user
- **Burst Limit**: 200 requests per minute per authenticated user
- **Anonymous Requests**: 10 requests per minute per IP address

When the rate limit is exceeded, the API responds with a 429 Too Many Requests status code:

```json
{
  "error": "Rate limit exceeded",
  "details": "You have exceeded the rate limit of 100 requests per minute.",
  "code": "RATE_LIMIT_EXCEEDED",
  "retry_after": 60
}
```

The `retry_after` field indicates the number of seconds you should wait before making another request.

## API Endpoints

### Authentication Endpoints

#### Login

Authenticate a user and receive a JWT token.

- **Endpoint**: `POST /api/v1/auth/login`
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "token": "string",
    "expires_at": "string",
    "user": {
      "id": "string",
      "username": "string",
      "email": "string",
      "role": "string"
    }
  }
  ```
- **Error Responses**:
  - 400 Bad Request: Invalid request body
  - 401 Unauthorized: Invalid username or password

#### Logout

Invalidate the current JWT token.

- **Endpoint**: `POST /api/v1/auth/logout`
- **Authentication**: Required (JWT token)
- **Request Body**: None
- **Response** (200 OK):
  ```json
  {
    "message": "Successfully logged out"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token

#### Refresh Token

Refresh an expired JWT token.

- **Endpoint**: `POST /api/v1/auth/refresh`
- **Authentication**: Required (JWT token)
- **Request Body**: None
- **Response** (200 OK):
  ```json
  {
    "token": "string",
    "expires_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token

### Configuration Item Endpoints

#### Get All CIs

Retrieve a paginated list of configuration items.

- **Endpoint**: `GET /api/v1/cis`
- **Authentication**: Required (JWT token)
- **Query Parameters**:
  - `page` (integer, optional): Page number (default: 1)
  - `limit` (integer, optional): Number of items per page (default: 20, max: 100)
  - `type` (string, optional): Filter by CI type
  - `sort` (string, optional): Sort field (default: created_at)
  - `order` (string, optional): Sort order (asc or desc, default: desc)
  - `search` (string, optional): Search term for name and attributes
- **Response** (200 OK):
  ```json
  {
    "data": [
      {
        "id": "string",
        "name": "string",
        "type": "string",
        "attributes": {},
        "tags": ["string"],
        "created_at": "string",
        "updated_at": "string"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5
    }
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 422 Unprocessable Entity: Invalid query parameters

#### Get CI by ID

Retrieve a specific configuration item by ID.

- **Endpoint**: `GET /api/v1/cis/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): CI ID
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "name": "string",
    "type": "string",
    "attributes": {},
    "tags": ["string"],
    "created_at": "string",
    "updated_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 404 Not Found: CI not found

#### Create CI

Create a new configuration item.

- **Endpoint**: `POST /api/v1/cis`
- **Authentication**: Required (JWT token)
- **Request Body**:
  ```json
  {
    "name": "string",
    "type": "string",
    "attributes": {},
    "tags": ["string"]
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": "string",
    "name": "string",
    "type": "string",
    "attributes": {},
    "tags": ["string"],
    "created_at": "string",
    "updated_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 422 Unprocessable Entity: Validation error

#### Update CI

Update an existing configuration item.

- **Endpoint**: `PUT /api/v1/cis/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): CI ID
- **Request Body**:
  ```json
  {
    "name": "string",
    "type": "string",
    "attributes": {},
    "tags": ["string"]
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "name": "string",
    "type": "string",
    "attributes": {},
    "tags": ["string"],
    "created_at": "string",
    "updated_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 404 Not Found: CI not found
  - 422 Unprocessable Entity: Validation error

#### Delete CI

Delete a configuration item.

- **Endpoint**: `DELETE /api/v1/cis/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): CI ID
- **Response** (200 OK):
  ```json
  {
    "message": "Configuration item deleted successfully"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 404 Not Found: CI not found

#### Get CI Relationships

Retrieve all relationships for a specific configuration item.

- **Endpoint**: `GET /api/v1/cis/{id}/relationships`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): CI ID
- **Query Parameters**:
  - `direction` (string, optional): Relationship direction (incoming, outgoing, or all, default: all)
- **Response** (200 OK):
  ```json
  {
    "incoming": [
      {
        "id": "string",
        "source_id": "string",
        "target_id": "string",
        "type": "string",
        "created_at": "string"
      }
    ],
    "outgoing": [
      {
        "id": "string",
        "source_id": "string",
        "target_id": "string",
        "type": "string",
        "created_at": "string"
      }
    ]
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 404 Not Found: CI not found
  - 422 Unprocessable Entity: Invalid query parameters

### Relationship Endpoints

#### Get All Relationships

Retrieve a paginated list of relationships.

- **Endpoint**: `GET /api/v1/relationships`
- **Authentication**: Required (JWT token)
- **Query Parameters**:
  - `page` (integer, optional): Page number (default: 1)
  - `limit` (integer, optional): Number of items per page (default: 20, max: 100)
  - `type` (string, optional): Filter by relationship type
  - `source_id` (string, optional): Filter by source CI ID
  - `target_id` (string, optional): Filter by target CI ID
- **Response** (200 OK):
  ```json
  {
    "data": [
      {
        "id": "string",
        "source_id": "string",
        "target_id": "string",
        "type": "string",
        "created_at": "string"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5
    }
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 422 Unprocessable Entity: Invalid query parameters

#### Get Relationship by ID

Retrieve a specific relationship by ID.

- **Endpoint**: `GET /api/v1/relationships/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): Relationship ID
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "source_id": "string",
    "target_id": "string",
    "type": "string",
    "created_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 404 Not Found: Relationship not found

#### Create Relationship

Create a new relationship.

- **Endpoint**: `POST /api/v1/relationships`
- **Authentication**: Required (JWT token)
- **Request Body**:
  ```json
  {
    "source_id": "string",
    "target_id": "string",
    "type": "string"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": "string",
    "source_id": "string",
    "target_id": "string",
    "type": "string",
    "created_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 422 Unprocessable Entity: Validation error

#### Update Relationship

Update an existing relationship.

- **Endpoint**: `PUT /api/v1/relationships/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): Relationship ID
- **Request Body**:
  ```json
  {
    "type": "string"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "source_id": "string",
    "target_id": "string",
    "type": "string",
    "created_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 404 Not Found: Relationship not found
  - 422 Unprocessable Entity: Validation error

#### Delete Relationship

Delete a relationship.

- **Endpoint**: `DELETE /api/v1/relationships/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): Relationship ID
- **Response** (200 OK):
  ```json
  {
    "message": "Relationship deleted successfully"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 404 Not Found: Relationship not found

### Audit Log Endpoints

#### Get All Audit Logs

Retrieve a paginated list of audit logs.

- **Endpoint**: `GET /api/v1/audit-logs`
- **Authentication**: Required (JWT token)
- **Query Parameters**:
  - `page` (integer, optional): Page number (default: 1)
  - `limit` (integer, optional): Number of items per page (default: 20, max: 100)
  - `entity_type` (string, optional): Filter by entity type (ci, relationship, user)
  - `action` (string, optional): Filter by action (create, update, delete)
  - `user_id` (string, optional): Filter by user ID
  - `from_date` (string, optional): Filter by start date (ISO 8601 format)
  - `to_date` (string, optional): Filter by end date (ISO 8601 format)
- **Response** (200 OK):
  ```json
  {
    "data": [
      {
        "id": "string",
        "entity_type": "string",
        "entity_id": "string",
        "action": "string",
        "user_id": "string",
        "username": "string",
        "changed_at": "string",
        "details": {}
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5
    }
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 422 Unprocessable Entity: Invalid query parameters

#### Get Audit Log by ID

Retrieve a specific audit log entry by ID.

- **Endpoint**: `GET /api/v1/audit-logs/{id}`
- **Authentication**: Required (JWT token)
- **Path Parameters**:
  - `id` (string, required): Audit log ID
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "entity_type": "string",
    "entity_id": "string",
    "action": "string",
    "user_id": "string",
    "username": "string",
    "changed_at": "string",
    "details": {}
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 404 Not Found: Audit log not found

### User Endpoints

#### Get All Users

Retrieve a paginated list of users.

- **Endpoint**: `GET /api/v1/users`
- **Authentication**: Required (JWT token, admin role)
- **Query Parameters**:
  - `page` (integer, optional): Page number (default: 1)
  - `limit` (integer, optional): Number of items per page (default: 20, max: 100)
  - `role` (string, optional): Filter by user role
- **Response** (200 OK):
  ```json
  {
    "data": [
      {
        "id": "string",
        "username": "string",
        "email": "string",
        "role": "string",
        "created_at": "string",
        "updated_at": "string"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5
    }
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions (admin role required)
  - 422 Unprocessable Entity: Invalid query parameters

#### Get User by ID

Retrieve a specific user by ID.

- **Endpoint**: `GET /api/v1/users/{id}`
- **Authentication**: Required (JWT token, admin role or own user)
- **Path Parameters**:
  - `id` (string, required): User ID
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "string",
    "created_at": "string",
    "updated_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 404 Not Found: User not found

#### Create User

Create a new user.

- **Endpoint**: `POST /api/v1/users`
- **Authentication**: Required (JWT token, admin role)
- **Request Body**:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string",
    "role": "string"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "string",
    "created_at": "string",
    "updated_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions (admin role required)
  - 422 Unprocessable Entity: Validation error

#### Update User

Update an existing user.

- **Endpoint**: `PUT /api/v1/users/{id}`
- **Authentication**: Required (JWT token, admin role or own user)
- **Path Parameters**:
  - `id` (string, required): User ID
- **Request Body**:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string",
    "role": "string"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "id": "string",
    "username": "string",
    "email": "string",
    "role": "string",
    "created_at": "string",
    "updated_at": "string"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions
  - 404 Not Found: User not found
  - 422 Unprocessable Entity: Validation error

#### Delete User

Delete a user.

- **Endpoint**: `DELETE /api/v1/users/{id}`
- **Authentication**: Required (JWT token, admin role)
- **Path Parameters**:
  - `id` (string, required): User ID
- **Response** (200 OK):
  ```json
  {
    "message": "User deleted successfully"
  }
  ```
- **Error Responses**:
  - 401 Unauthorized: Invalid or expired token
  - 403 Forbidden: Insufficient permissions (admin role required)
  - 404 Not Found: User not found

## Data Models

### Configuration Item

```json
{
  "id": "string",
  "name": "string",
  "type": "string",
  "attributes": {},
  "tags": ["string"],
  "created_at": "string",
  "updated_at": "string"
}
```

### Relationship

```json
{
  "id": "string",
  "source_id": "string",
  "target_id": "string",
  "type": "string",
  "created_at": "string"
}
```

### Audit Log

```json
{
  "id": "string",
  "entity_type": "string",
  "entity_id": "string",
  "action": "string",
  "user_id": "string",
  "username": "string",
  "changed_at": "string",
  "details": {}
}
```

### User

```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "role": "string",
  "created_at": "string",
  "updated_at": "string"
}
```

## API Examples

### Authentication Example

```bash
# Login
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-12-31T23:59:59Z",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}

# Store the token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Configuration Item Examples

```bash
# Get all CIs
curl -X GET https://your-domain.com/api/v1/cis \
  -H "Authorization: Bearer $TOKEN"

# Get a specific CI
curl -X GET https://your-domain.com/api/v1/cis/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer $TOKEN"

# Create a new CI
curl -X POST https://your-domain.com/api/v1/cis \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Web Server",
    "type": "Server",
    "attributes": {
      "ip_address": "192.168.1.100",
      "os": "Ubuntu 20.04",
      "cpu_cores": 4,
      "memory_gb": 16
    },
    "tags": ["production", "web-tier"]
  }'

# Update a CI
curl -X PUT https://your-domain.com/api/v1/cis/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Web Server (Updated)",
    "attributes": {
      "ip_address": "192.168.1.101",
      "os": "Ubuntu 22.04",
      "cpu_cores": 4,
      "memory_gb": 16
    }
  }'

# Delete a CI
curl -X DELETE https://your-domain.com/api/v1/cis/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer $TOKEN"
```

### Relationship Examples

```bash
# Get all relationships
curl -X GET https://your-domain.com/api/v1/relationships \
  -H "Authorization: Bearer $TOKEN"

# Create a new relationship
curl -X POST https://your-domain.com/api/v1/relationships \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "source_id": "123e4567-e89b-12d3-a456-426614174000",
    "target_id": "223e4567-e89b-12d3-a456-426614174001",
    "type": "hosts"
  }'

# Update a relationship
curl -X PUT https://your-domain.com/api/v1/relationships/323e4567-e89b-12d3-a456-426614174002 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "depends_on"
  }'

# Delete a relationship
curl -X DELETE https://your-domain.com/api/v1/relationships/323e4567-e89b-12d3-a456-426614174002 \
  -H "Authorization: Bearer $TOKEN"
```

### Audit Log Examples

```bash
# Get all audit logs
curl -X GET https://your-domain.com/api/v1/audit-logs \
  -H "Authorization: Bearer $TOKEN"

# Get audit logs for a specific entity type
curl -X GET "https://your-domain.com/api/v1/audit-logs?entity_type=ci" \
  -H "Authorization: Bearer $TOKEN"

# Get audit logs for a specific date range
curl -X GET "https://your-domain.com/api/v1/audit-logs?from_date=2023-01-01T00:00:00Z&to_date=2023-12-31T23:59:59Z" \
  -H "Authorization: Bearer $TOKEN"
```

## SDKs and Libraries

Currently, CMDB Lite does not provide official SDKs or libraries. However, you can easily integrate with the API using standard HTTP client libraries in your preferred programming language.

### JavaScript/TypeScript Example

```javascript
// Using fetch API
const API_BASE_URL = 'https://your-domain.com/api/v1';
let authToken = '';

// Login
async function login(username, password) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });
  
  const data = await response.json();
  authToken = data.token;
  return data;
}

// Get all CIs
async function getCIs() {
  const response = await fetch(`${API_BASE_URL}/cis`, {
    headers: {
      'Authorization': `Bearer ${authToken}`,
    },
  });
  
  return response.json();
}

// Create a new CI
async function createCI(ciData) {
  const response = await fetch(`${API_BASE_URL}/cis`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(ciData),
  });
  
  return response.json();
}
```

### Python Example

```python
import requests
import json

API_BASE_URL = 'https://your-domain.com/api/v1'
auth_token = ''

# Login
def login(username, password):
    response = requests.post(
        f'{API_BASE_URL}/auth/login',
        json={'username': username, 'password': password}
    )
    data = response.json()
    global auth_token
    auth_token = data['token']
    return data

# Get all CIs
def get_cis():
    response = requests.get(
        f'{API_BASE_URL}/cis',
        headers={'Authorization': f'Bearer {auth_token}'}
    )
    return response.json()

# Create a new CI
def create_ci(ci_data):
    response = requests.post(
        f'{API_BASE_URL}/cis',
        headers={
            'Authorization': f'Bearer {auth_token}',
            'Content-Type': 'application/json'
        },
        json=ci_data
    )
    return response.json()
```

### Go Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const APIBaseURL = "https://your-domain.com/api/v1"

var authToken string

// Login
func login(username, password string) (map[string]interface{}, error) {
    data := map[string]string{
        "username": username,
        "password": password,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.Post(
        fmt.Sprintf("%s/auth/login", APIBaseURL),
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }
    
    authToken = result["token"].(string)
    return result, nil
}

// Get all CIs
func getCIs() (map[string]interface{}, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/cis", APIBaseURL), nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}

// Create a new CI
func createCI(ciData map[string]interface{}) (map[string]interface{}, error) {
    jsonData, err := json.Marshal(ciData)
    if err != nil {
        return nil, err
    }
    
    client := &http.Client{}
    req, err := http.NewRequest("POST", fmt.Sprintf("%s/cis", APIBaseURL), bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
    req.Header.Add("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }
    
    return result, nil
}
```

For more information on contributing to the API, see the [Developer Guide](README.md) and the [Contribution Guide](../project/contributing.md).
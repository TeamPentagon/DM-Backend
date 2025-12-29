# API Documentation

This document provides detailed API documentation for the DM-Backend service.

## Base URL

```
http://localhost:8085
```

## API Version

Current version: **v1**

All versioned endpoints are prefixed with `/api/v1/`.

## Authentication

Currently, the API supports basic authentication via the `Authorization` header.

```http
Authorization: Bearer <token>
```

> **Note:** Authentication middleware is available but JWT validation is pending implementation.

## Response Format

All API responses follow this structure:

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

## Endpoints

### Health Check

Check if the service is running and healthy.

#### Request
```http
GET /health
```

#### Response
```json
{
  "success": true,
  "data": {
    "status": "healthy"
  }
}
```

#### Status Codes
| Code | Description |
|------|-------------|
| 200 | Service is healthy |
| 503 | Service unavailable |

---

### Chat

Send a message to the AI and get a response.

#### Request
```http
POST /api/v1/chat?response=<message>
```

#### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `response` | string | Yes | The user's message or prompt to send to the AI |

#### Example Request
```bash
curl -X POST "http://localhost:8085/api/v1/chat?response=Hello%2C%20how%20are%20you%3F"
```

#### Success Response
```json
{
  "success": true,
  "data": {
    "results": [
      {
        "text": "Hello! I'm doing well, thank you for asking. How can I assist you today?"
      }
    ]
  }
}
```

#### Error Responses

**Missing Parameter (400)**
```json
{
  "success": false,
  "error": "Missing 'response' query parameter"
}
```

**AI Service Unavailable (503)**
```json
{
  "success": false,
  "error": "AI service unavailable"
}
```

**Internal Error (500)**
```json
{
  "success": false,
  "error": "Failed to process request"
}
```

#### Status Codes
| Code | Description |
|------|-------------|
| 200 | Successful response from AI |
| 400 | Missing or invalid parameters |
| 500 | Internal server error |
| 503 | AI service unavailable |

---

### Legacy Chat Endpoint

For backward compatibility with older clients.

#### Request
```http
POST /chat?response=<message>
```

> **Note:** This endpoint has the same behavior as `/api/v1/chat`.

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Requests per second:** Configurable (default: 10)
- **Burst size:** Configurable (default: 20)

When rate limited, you'll receive:

```http
HTTP/1.1 429 Too Many Requests
```

```json
{
  "success": false,
  "error": "Rate limit exceeded. Please try again later."
}
```

## CORS

Cross-Origin Resource Sharing is enabled with the following defaults:

- **Allowed Origins:** `*` (all origins in development)
- **Allowed Methods:** GET, POST, PUT, DELETE, OPTIONS
- **Allowed Headers:** Origin, Content-Type, Authorization, Accept
- **Credentials:** Allowed
- **Max Age:** 12 hours

## Request ID Tracking

Every request is assigned a unique request ID for debugging:

- If you provide `X-Request-ID` header, that ID will be used
- Otherwise, a new ID is generated

The request ID is returned in the response headers:

```http
X-Request-ID: 20241229103245.123456
```

## Error Handling

### Common Error Codes

| HTTP Code | Error Type | Description |
|-----------|------------|-------------|
| 400 | Bad Request | Invalid or missing parameters |
| 401 | Unauthorized | Missing or invalid authentication |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server-side error |
| 503 | Service Unavailable | External service (AI) unavailable |

### Error Response Structure

```json
{
  "success": false,
  "error": "Human-readable error message"
}
```

## Data Models

### User Model

```json
{
  "user_name": "string",
  "person_id": "string (unique identifier)",
  "profile": "string",
  "password": "string (hashed)",
  "email": "string",
  "profile_pic_url": "string (URL)",
  "account_time": "integer (Unix timestamp)",
  "birth_date": "integer (Unix timestamp)",
  "gender": "string",
  "last_edit": "integer (Unix timestamp)",
  "phone_number": "string"
}
```

### Message Model

```json
{
  "ai_id": "string",
  "user_id": "string",
  "content": "string",
  "timestamp": "integer (Unix timestamp in milliseconds)",
  "conversation_id": "string",
  "msg_id": "string (unique identifier)"
}
```

### Chat History Model

```json
{
  "conversation_id": "string",
  "messages": [
    {
      "ai_id": "string",
      "user_id": "string",
      "content": "string",
      "timestamp": "integer",
      "conversation_id": "string",
      "msg_id": "string"
    }
  ]
}
```

## Examples

### cURL Examples

**Health Check:**
```bash
curl http://localhost:8085/health
```

**Send Chat Message:**
```bash
curl -X POST "http://localhost:8085/api/v1/chat?response=What%20is%20the%20weather%20like%3F"
```

**With Authentication:**
```bash
curl -X POST \
  -H "Authorization: Bearer your-token-here" \
  "http://localhost:8085/api/v1/chat?response=Hello"
```

### Go Client Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   string      `json:"error"`
}

func main() {
    message := "Hello, how are you?"
    endpoint := fmt.Sprintf("http://localhost:8085/api/v1/chat?response=%s", 
        url.QueryEscape(message))
    
    resp, err := http.Post(endpoint, "application/json", nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var result Response
    json.Unmarshal(body, &result)
    
    if result.Success {
        fmt.Printf("AI Response: %v\n", result.Data)
    } else {
        fmt.Printf("Error: %s\n", result.Error)
    }
}
```

### JavaScript/Fetch Example

```javascript
async function sendMessage(message) {
    const response = await fetch(
        `http://localhost:8085/api/v1/chat?response=${encodeURIComponent(message)}`,
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        }
    );
    
    const data = await response.json();
    
    if (data.success) {
        console.log('AI Response:', data.data);
    } else {
        console.error('Error:', data.error);
    }
}

sendMessage('Hello, how are you?');
```

## SDK Support

Official SDKs are planned for:
- [ ] Go
- [ ] JavaScript/TypeScript
- [ ] Python

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for API changes and updates.

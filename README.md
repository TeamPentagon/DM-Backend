# DM-Backend

A scalable Go backend service for AI-powered chat applications with support for medical/healthcare domain recommendations.

[![Go Version](https://img.shields.io/badge/Go-1.22.2-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [API Reference](#api-reference)
- [Project Structure](#project-structure)
- [Database](#database)
- [Testing](#testing)
- [Development](#development)
- [Datasets](#datasets)
- [Models](#models)
- [Contributing](#contributing)
- [License](#license)

## Overview

DM-Backend is a robust backend service designed to power AI chat applications with a focus on healthcare and medical recommendations. It provides a REST API for chat interactions, user management, and conversation history storage.

## Features

- 🚀 **High Performance**: Built with Go and Gin framework for optimal performance
- 💾 **Flexible Storage**: Dual database support with LevelDB (NoSQL) and SQLite (SQL)
- 🔐 **Security**: Built-in middleware for authentication, rate limiting, and CORS
- 📊 **Scalable**: Sharding support for horizontal scaling
- 🧪 **Well Tested**: Comprehensive test coverage for all components
- 📝 **API Versioning**: RESTful API with versioned endpoints

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Applications                      │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Gin HTTP Server                          │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                     Middleware Stack                         ││
│  │  ┌─────────┐ ┌───────┐ ┌──────────┐ ┌────────┐ ┌──────────┐ ││
│  │  │  CORS   │ │Logger │ │RateLimiter│ │  Auth  │ │RequestID │ ││
│  │  └─────────┘ └───────┘ └──────────┘ └────────┘ └──────────┘ ││
│  └─────────────────────────────────────────────────────────────┘│
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                         API Routes                               │
│  ┌─────────┐  ┌─────────────┐  ┌───────────────┐                │
│  │ /health │  │ /api/v1/chat│  │ /api/v1/users │                │
│  └─────────┘  └─────────────┘  └───────────────┘                │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Business Logic                             │
│  ┌─────────────┐  ┌───────────────┐  ┌─────────────────────┐    │
│  │ User Model  │  │  Chat Model   │  │ Fragmentation Logic │    │
│  └─────────────┘  └───────────────┘  └─────────────────────┘    │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Database Layer                              │
│  ┌────────────────────────┐  ┌────────────────────────────────┐ │
│  │       LevelDB          │  │           SQLite               │ │
│  │  (NoSQL Key-Value)     │  │      (Relational Data)         │ │
│  └────────────────────────┘  └────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Getting Started

### Prerequisites

- Go 1.22.2 or higher
- Protocol Buffers compiler (protoc) - for generating model files
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/TeamPentagon/DM-Backend.git
   cd DM-Backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Generate Protocol Buffer files:
   ```bash
   # Install protoc-gen-go if not already installed
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   
   # Generate Go files from proto files
   cd internal/model
   protoc --go_out=. --go_opt=paths=source_relative user.proto
   protoc --go_out=. --go_opt=paths=source_relative chat.proto
   cd ../..
   ```

4. Build the application:
   ```bash
   go build -o dm-backend .
   ```

5. Run the application:
   ```bash
   ./dm-backend
   ```

### Configuration

The application can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8085` |
| `AI_ENDPOINT` | AI service endpoint URL | `https://herbal-pmc-allowance-cognitive.trycloudflare.com/api/v1/generate` |

Example:
```bash
export PORT=8080
export AI_ENDPOINT="https://your-ai-service.com/api/v1/generate"
./dm-backend
```

## API Reference

### Health Check

**Endpoint:** `GET /health`

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy"
  }
}
```

### Chat

**Endpoint:** `POST /api/v1/chat`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `response` | string | Yes | The user's message/prompt |

**Response:**
```json
{
  "success": true,
  "data": {
    "results": [
      {
        "text": "AI response here..."
      }
    ]
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error description"
}
```

### Legacy Endpoint

For backward compatibility, the chat endpoint is also available at:

**Endpoint:** `POST /chat`

## Project Structure

```
DM-Backend/
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── README.md                   # This file
├── LICENSE                     # MIT License
├── internal/                   # Internal packages
│   ├── database/               # Database utilities
│   │   ├── database.go         # DB connection management
│   │   ├── database_test.go    # DB connection tests
│   │   ├── fragmentation.go    # Shard management
│   │   └── fragmentation_test.go
│   ├── middleware/             # HTTP middleware
│   │   ├── middleware.go       # CORS, Auth, RateLimit, etc.
│   │   └── middleware_test.go
│   └── model/                  # Data models
│       ├── user.go             # User CRUD operations
│       ├── user.proto          # User Protocol Buffer schema
│       ├── user.pb.go          # Generated User types
│       ├── chat.go             # Chat operations
│       ├── chat.proto          # Chat Protocol Buffer schema
│       └── chat.pb.go          # Generated Chat types
├── test/                       # Integration tests
│   ├── user_test.go            # User integration tests
│   └── chat_test.go            # Chat integration tests
└── Database/                   # Database files (created at runtime)
    ├── Common/                 # Common data shards
    ├── GlobalSchema/           # Fragmentation schema
    └── NOSQL/                  # NoSQL message storage
```

## Database

### LevelDB (Key-Value Store)

Used for:
- User data storage (serialized with Protocol Buffers)
- Chat messages and history
- Fragmentation schema (shard mapping)

### SQLite (Relational)

Available for:
- Complex queries
- Relational data that requires joins
- Reporting and analytics

### Sharding

The application supports horizontal sharding through the fragmentation module:

```go
// Add a key to shard mapping
database.FragmentationAdd(shardNumber, "user_key_1", "user_key_2")

// Get shard for a key
shard, err := database.FragmentationGet("user_key_1")

// Update shard mapping
database.FragmentationUpdate("user_key_1", newShardNumber)

// Remove mapping
database.FragmentationRemove("user_key_1")
```

## Testing

Run all tests:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Run specific package tests:
```bash
# Database tests
go test ./internal/database/... -v

# Middleware tests
go test ./internal/middleware/... -v

# Integration tests
go test ./test/... -v
```

## Development

### Adding a New Model

1. Create the Protocol Buffer schema in `internal/model/`:
   ```protobuf
   syntax = "proto3";
   package model;
   option go_package = ".;model";
   
   message YourModel {
       string id = 1;
       string name = 2;
   }
   ```

2. Generate Go code:
   ```bash
   cd internal/model
   protoc --go_out=. --go_opt=paths=source_relative yourmodel.proto
   ```

3. Implement CRUD operations in a new Go file

### Adding Middleware

1. Add your middleware function in `internal/middleware/middleware.go`
2. Follow the Gin middleware pattern:
   ```go
   func YourMiddleware() gin.HandlerFunc {
       return func(c *gin.Context) {
           // Pre-processing
           c.Next()
           // Post-processing
       }
   }
   ```
3. Add tests in `middleware_test.go`

## Datasets

The project references the following datasets for training/fine-tuning:

1. [Doctor Review Dataset](https://www.kaggle.com/datasets/avasaralasaipavan/doctor-review-dataset-has-reviews-on-doctors)
2. [Healthcare Practitioner Dataset](https://www.kaggle.com/datasets/kirillshchitaev/moscow-healthcare-practitioners-dataset)
3. [Patient Reviews (German)](https://www.kaggle.com/datasets/thedevastator/german-2021-patient-reviews-and-ratings-of-docto)
4. [Healthcare NLP: LLMs, Transformers](https://www.kaggle.com/datasets/jpmiller/layoutlm)
5. [Medical Recommendation Dataset](https://www.kaggle.com/datasets/joymarhew/medical-reccomadation-dataset)
6. [WebMD Drug Reviews Dataset](https://www.kaggle.com/datasets/rohanharode07/webmd-drug-reviews-dataset)

## Models

Supported AI models for backend integration:

| Model | Provider | Link |
|-------|----------|------|
| Mixtral | Mistral AI | [Kaggle](https://www.kaggle.com/models/mistral-ai/mixtral) |
| Flan-T5 | Google | [Kaggle](https://www.kaggle.com/models/google/flan-t5) |
| Mistral | Mistral AI | [Kaggle](https://www.kaggle.com/models/mistral-ai/mistral) |
| Llama 2 | Meta | [Kaggle](https://www.kaggle.com/models/metaresearch/llama-2) |
| Llama 3 | Meta | [Kaggle](https://www.kaggle.com/models/metaresearch/llama-3/pyTorch/8b) |
| Phi | Microsoft | [Kaggle](https://www.kaggle.com/models/Microsoft/phi) |
| Gemma | Google | [Kaggle](https://www.kaggle.com/models/google/gemma) |

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -am 'Add your feature'`
4. Push to the branch: `git push origin feature/your-feature`
5. Submit a pull request

### Code Style

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` before committing
- Run `go vet` to check for common issues
- Add tests for new functionality

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Team Pentagon** - Building intelligent healthcare solutions

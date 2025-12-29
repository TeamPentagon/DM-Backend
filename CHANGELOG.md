# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive database package with LevelDB and SQLite support
- Middleware package with CORS, rate limiting, authentication, logging, and recovery middleware
- Full CRUD operations for User and Chat models
- Fragmentation support for horizontal sharding
- Health check endpoint at `/health`
- API versioning with `/api/v1/` prefix
- Structured response format with success/error handling
- Comprehensive test suite with 80%+ coverage
- Complete documentation including README, CONTRIBUTING, and API reference
- Makefile for common development tasks
- Environment variable configuration support

### Changed
- Refactored main.go with proper error handling (no more log.Fatal in handlers)
- Updated User model with proper UpdateUserData and DeleteUserData implementations
- Renamed fragmentation functions to follow Go naming conventions (no underscores)
- Improved database connection handling with path validation

### Fixed
- **Critical**: `UpdatedUserData` function was not updating data, only reading
- **Critical**: `DeleteUserData` function was not deleting data, only reading  
- Potential panic in `buildPath` when called with fewer than 2 parameters
- `defer ldb.Close()` placed after return statement in fragmentation functions
- Missing input validation on all database operations
- Missing error wrapping for better debugging

### Security
- Added rate limiting middleware to prevent abuse
- Added CORS middleware with configurable origins
- Added authentication middleware (placeholder for JWT)
- Added request ID tracking for debugging

## [0.1.0] - 2024-XX-XX

### Added
- Initial project structure
- Basic chat endpoint with AI integration
- User model with Protocol Buffer serialization
- Chat model with message history
- LevelDB integration for data storage

---

## Version History

### Versioning Strategy

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality
- **PATCH** version for backwards-compatible bug fixes

### Migration Notes

#### From v0.1.0 to v0.2.0 (Upcoming)

**Breaking Changes:**
1. Fragmentation function names changed:
   - `Fragmentation_Add` → `FragmentationAdd`
   - `Fragmentation_Remove` → `FragmentationRemove`
   - `Fragmentation_Get` → `FragmentationGet`

2. User model function renamed:
   - `UpdatedUserData` → `UpdateUserData`

**Migration Steps:**
1. Update all calls to fragmentation functions to use the new names
2. Update calls to `UpdatedUserData` to use `UpdateUserData`
3. Run `go build ./...` to verify no compilation errors

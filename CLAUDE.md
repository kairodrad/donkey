# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Setup

- Go 1.24.3 or later required
- Run `go mod download` to install dependencies
- Set DATABASE_URL for PostgreSQL or omit for SQLite

## Development Commands

**Run the server:**
```bash
go run ./cmd/server
```

**IMPORTANT:** Do not start or stop servers automatically. Always ask the user to start/stop servers when needed.

**Run tests:**
```bash
go test ./...
```

**Generate Swagger documentation:**
```bash
swag init -g cmd/server/main.go
```

## Project Architecture

This is a Go web application for the "Donkey" card game that supports multiplayer gameplay through HTTP APIs and real-time updates.

### Core Components

**Backend Structure:**
- `cmd/server/main.go` - Application entry point with Swagger configuration
- `internal/server/server.go` - Gin HTTP server setup with all route definitions
- `internal/api/` - HTTP handlers for game operations (register, start, join, finalize, etc.)
- `internal/game/` - Core game logic including card dealing and state management
- `internal/model/` - GORM data models (User, Game, GamePlayer, GameCard, GameState, GameSessionLog)
- `internal/db/db.go` - Database initialization supporting both PostgreSQL and SQLite

**Database Design:**
- Uses GORM ORM with automatic migrations
- PostgreSQL for production (via DATABASE_URL env var), SQLite in-memory for development/testing
- Many-to-many relationship between Users and Games through GamePlayer junction table
- GameCard tracks individual card assignments to players

**Frontend:**
- `web/index.html` - Single-page React application using CDN imports
- `web/ui/` - JavaScript modules for different game screens (registration, game flow, settings)
- `web/assets/` - Static card images and game assets
- Uses Tailwind CSS for styling

**Game Flow:**
1. User registration creates a User with random hex ID
2. Game creation by requester, others join via game ID
3. Manual finalization by requester or auto-start at 8 players
4. Card dealing distributes full deck among players starting after requester
5. Real-time state updates via SSE streams and polling

**Key Features:**
- Server-Sent Events for real-time game updates (`/api/game/:gameId/stream/:userId`)
- Chat system with game session logging
- Admin endpoints for debugging game state
- Swagger API documentation at `/api/swagger/index.html`
- Responsive design with light/dark mode support

**Testing:**
- Unit tests in `internal/api/game_flow_test.go`
- Integration tests verify complete game flow scenarios
- Test coverage goals: Aim for 80%+ coverage on core game logic
- Test database setup: Uses SQLite in-memory for fast test execution
- Run integration tests with `go test ./... -tags=integration`

## Code Quality

- Run `go fmt` to format code according to Go standards
- Run `go vet` to check for common errors and suspicious constructs
- Follow Go best practices and idiomatic patterns
- Code review expectations:
  - All PRs require review before merging
  - Focus on readability, maintainability, and performance
  - Ensure adequate test coverage for new features
  - Check for proper error handling and edge cases
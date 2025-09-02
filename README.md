# Donkey Card Game
A multiplayer card game that brings family together - now with a modern Vue.js interface and mobile-first design.

## Features

- 🎮 **2-8 player multiplayer** support
- 📱 **Mobile-first design** with touch-optimized card selection
- 🎨 **Modern UI** using Vue.js with pastel design system
- 🌙 **Light/dark theme** support with system detection
- 🔄 **Real-time updates** via Server-Sent Events
- 💬 **In-game chat** and session logging
- 🃏 **Fan card arrangements** for intuitive gameplay

## Quick Start

### Prerequisites
- **Node.js v16.14.1** (managed via NVM) - **Required for Vite 4 compatibility**
- **Go 1.24.3+** for backend server

### Setup
```bash
# Use correct Node.js version
nvm use

# Install frontend dependencies
npm install

# Start development servers
npm run dev     # Frontend (localhost:3000)
npm run server  # Backend (localhost:8080)
```

For detailed setup instructions, see [DEVELOPMENT.md](./DEVELOPMENT.md).

## Production

```bash
# Build and start production server
npm start

# Or manually
npm run build
go run ./cmd/server
```

## API Documentation

Swagger UI available at: [http://localhost:8080/api/swagger/index.html](http://localhost:8080/api/swagger/index.html)

## Testing

```bash
go test ./...
```

## Technology Stack

- **Frontend**: Vue.js 3, Vite, Tailwind CSS
- **Backend**: Go, Gin, GORM
- **Database**: SQLite (dev) / PostgreSQL (prod)
- **Real-time**: Server-Sent Events

## Architecture

The game features a circular seating arrangement for up to 8 players, with real-time state synchronization and mobile-optimized card interactions using fan layouts for intuitive touch gameplay.

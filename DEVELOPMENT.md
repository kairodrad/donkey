# Development Setup

This document explains how to set up the development environment for the Donkey card game.

## Prerequisites

### Node.js Version Management

This project uses **Node.js v16.14.1** and **npm v7.18.1**. These specific versions are required for Vite 4 compatibility and to avoid crypto errors. We use NVM (Node Version Manager) to ensure everyone uses the exact same versions.

#### 1. Install NVM (if not already installed)

**macOS/Linux:**
```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
```

**Windows:**
Download and install [nvm-windows](https://github.com/coreybutler/nvm-windows/releases)

#### 2. Use the Project's Node Version

```bash
# Navigate to project directory
cd /path/to/donkey

# Install and use the specified Node.js version
nvm install
nvm use

# Verify versions
node --version  # Should output: v16.14.1
npm --version   # Should output: 7.18.1
```

The `.nvmrc` file in the project root contains the exact Node.js version (v16.14.1) required for this project.

### Go Environment

- **Go 1.24.3** or later required for the backend server
- Set `DATABASE_URL` environment variable for PostgreSQL (optional, defaults to SQLite)

## Quick Start

### 1. Install Dependencies
```bash
npm install
```

### 2. Development Mode
```bash
# Terminal 1: Start Vue.js development server (with hot reload)
npm run dev
# Serves at http://localhost:3000 with API proxy to :8080

# Terminal 2: Start Go backend server
npm run server
# or
go run ./cmd/server
# Serves at http://localhost:8080
```

### 3. Production Build
```bash
# Build Vue.js app and start production server
npm start
# or manually:
npm run build
go run ./cmd/server
```

## Project Structure

```
donkey/
├── .nvmrc                 # Node.js version lock file
├── package.json           # npm dependencies and scripts
├── vite.config.js         # Vite configuration
├── tailwind.config.js     # Tailwind CSS configuration
├── src/                   # Vue.js source code
│   ├── main.js           # App entry point
│   ├── App.vue           # Main Vue component
│   ├── components/       # Reusable Vue components
│   ├── composables/      # Vue composition functions
│   ├── utils/            # JavaScript utilities
│   └── styles/           # CSS styles
├── dist/                 # Built Vue.js app (created by npm run build)
├── web/assets/           # Static assets (card images, icons)
├── cmd/server/           # Go application entry point
├── internal/             # Go backend code
│   ├── api/             # HTTP handlers
│   ├── db/              # Database layer
│   ├── game/            # Game logic
│   ├── model/           # Data models
│   └── server/          # HTTP server setup
└── ../design-system/    # External design system dependency
    └── src/card_assets/ # Card image assets
```

## Available Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server with hot reload |
| `npm run build` | Build production Vue.js app |
| `npm run server` | Start Go backend server |
| `npm start` | Build and start production server |
| `npm run nvm` | Use project's Node.js version |
| `npm run lint` | Run ESLint code linting |
| `npm run format` | Format code with Prettier |

## Troubleshooting

### Node.js Version Issues
**Important:** This project requires exactly Node.js v16.14.1 and npm v7.18.1. These are the tested versions that work with Vite 4 without crypto errors.

If you encounter crypto errors or other compatibility issues:
```bash
# Ensure you're using the correct Node.js version
nvm use
node --version  # Should be v16.14.1

# Clean install if needed
rm -rf node_modules package-lock.json
npm install
```

### Port Conflicts
- Vue.js dev server: http://localhost:3000
- Go backend server: http://localhost:8080
- Make sure these ports are available

### Asset Loading Issues
The app references card assets from the design system at `../design-system/src/card_assets/`. Ensure this directory exists and contains the card image files.

## Technology Stack

### Frontend
- **Vue.js 3** - Progressive JavaScript framework
- **Vite 4** - Fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **Design System** - Custom Vue component library with pastel theme

### Backend
- **Go** - HTTP server and game logic
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **SQLite/PostgreSQL** - Database options

### Development
- **NVM** - Node.js version management
- **ESLint** - Code linting
- **Prettier** - Code formatting

## Design System Integration

This project uses the kairodrad/design-system for UI components and styling:
- Card components with fan arrangements
- Mobile-first responsive design
- Light/dark theme support
- Pastel color palette
- Advanced CSS custom properties for theming
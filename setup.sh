#!/bin/bash

# Donkey Card Game - Development Setup Script

echo "🃏 Setting up Donkey Card Game development environment..."
echo

# Check if NVM is installed
if ! command -v nvm &> /dev/null; then
    echo "❌ NVM not found. Installing NVM..."
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    echo "✅ NVM installed successfully"
else
    echo "✅ NVM is already installed"
fi

# Use correct Node.js version
if [ -f ".nvmrc" ]; then
    echo "📦 Installing and using Node.js version from .nvmrc..."
    nvm install
    nvm use
    echo "✅ Using Node.js $(node --version) and npm $(npm --version)"
    if [[ "$(node --version)" != "v16.14.1" ]]; then
        echo "⚠️  Warning: Expected Node.js v16.14.1 for Vite 4 compatibility"
    fi
else
    echo "❌ .nvmrc file not found"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go not found. Please install Go 1.24.3+ from https://golang.org/dl/"
    exit 1
else
    echo "✅ Go $(go version | cut -d' ' -f3) found"
fi

# Install Node.js dependencies
echo "📦 Installing npm dependencies..."
if npm install; then
    echo "✅ npm dependencies installed successfully"
else
    echo "❌ Failed to install npm dependencies"
    exit 1
fi

# Test build
echo "🔨 Testing production build..."
if npm run build; then
    echo "✅ Production build successful"
else
    echo "❌ Production build failed"
    exit 1
fi

echo
echo "🎉 Setup complete! You can now run:"
echo
echo "  Development servers:"
echo "    npm run dev     # Vue.js frontend (localhost:3000)"
echo "    npm run server  # Go backend (localhost:8080)"
echo
echo "  Production:"
echo "    npm start       # Build and start production server"
echo
echo "  Documentation:"
echo "    See DEVELOPMENT.md for detailed instructions"
echo
echo "🃏 Happy coding!"
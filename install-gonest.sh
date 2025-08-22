#!/bin/bash

# GoNest CLI One-Line Installer
# Usage: curl -sSL https://raw.githubusercontent.com/ulims/GoNest/main/install-gonest.sh | bash

set -e

echo "🚀 Installing GoNest CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION+"
    exit 1
fi

echo "✅ Go $GO_VERSION detected"

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "📥 Cloning GoNest repository..."
git clone --quiet https://github.com/ulims/GoNest.git
cd GoNest

echo "🔨 Building and installing CLI tool..."
go install ./cmd/gonest

# Clean up
cd /
rm -rf "$TEMP_DIR"

echo "✅ GoNest CLI installed successfully!"
echo ""
echo "🎯 Usage:"
echo "   gonest --help                    # Show help"
echo "   gonest new my-app               # Create new project"
echo "   gonest new my-api --template=api --strict  # With template and strict mode"
echo ""
echo "🚀 Happy coding with GoNest!"

#!/bin/bash

# GoNest CLI Installer using Go Workspaces (Modern Approach)
# Usage: curl -sSL https://raw.githubusercontent.com/ulims/GoNest/main/install-gonest-workspace.sh | bash

set -e

echo "🚀 Installing GoNest CLI using Go Workspaces..."

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

# Create workspace directory
WORKSPACE_DIR="$HOME/.gonest-workspace"
mkdir -p "$WORKSPACE_DIR"
cd "$WORKSPACE_DIR"

# Initialize Go workspace
if [ ! -f "go.work" ]; then
    echo "🔧 Initializing Go workspace..."
    go work init
fi

# Add GoNest module
echo "📥 Adding GoNest module to workspace..."
go work use github.com/ulims/GoNest

# Install CLI tool
echo "🔨 Installing CLI tool..."
go install github.com/ulims/GoNest/cmd/gonest@latest

echo "✅ GoNest CLI installed successfully!"
echo ""
echo "🎯 Usage:"
echo "   gonest --help                    # Show help"
echo "   gonest new my-app               # Create new project"
echo "   gonest new my-api --template=api --strict  # With template and strict mode"
echo ""
echo "🚀 Happy coding with GoNest!"
echo ""
echo "💡 To update later:"
echo "   cd ~/.gonest-workspace && go work sync && go install github.com/ulims/GoNest/cmd/gonest@latest"

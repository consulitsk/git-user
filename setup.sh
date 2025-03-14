#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Define variables
VERSION="v1.0.0"
REPO_URL="https://github.com/consulitsk/git-user/releases/download"
BINARY_NAME="git-user-linux"
INSTALL_DIR="$HOME/.local/bin"
INSTALL_PATH="$INSTALL_DIR/git-user"

# Detect system architecture (currently supports amd64 only)
ARCH=$(uname -m)
if [ "$ARCH" != "x86_64" ]; then
    echo "❌ Unsupported architecture: $ARCH. This installer supports only x86_64 (amd64)."
    exit 1
fi

# Create installation directory if it does not exist
mkdir -p "$INSTALL_DIR"

# Download the binary (silent mode, no progress output)
echo "🔽 Downloading git-user..."
if ! curl -s -L -o "$INSTALL_PATH" "$REPO_URL/$VERSION/$BINARY_NAME"; then
    echo "❌ Failed to download git-user. Check your internet connection and try again."
    exit 1
fi

# Make the binary executable
chmod +x "$INSTALL_PATH"

# Ensure ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "⚠️  ~/.local/bin is not in your PATH. Adding it..."
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc"
    export PATH="$HOME/.local/bin:$PATH"
fi

# Verify installation
if command -v git-user >/dev/null 2>&1; then
    echo "✅ git-user installed successfully!"
    git-user --help  # Show help output to confirm installation
else
    echo "❌ Installation failed!"
    exit 1
fi

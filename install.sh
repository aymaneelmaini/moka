#!/bin/bash

set -e

REPO="aymaneelmaini/moka"
INSTALL_DIR="$HOME/.local/bin"
DATA_DIR="$HOME/.moka"
SERVICE_NAME="moka"

echo "🚀 Installing Moka Finance Tracker..."
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "📦 Detected: $OS-$ARCH"
echo ""

# Get latest release version
echo "🔍 Fetching latest release..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "❌ Could not fetch latest version"
    exit 1
fi

echo "📥 Downloading Moka $LATEST_VERSION..."
BINARY_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/moka-$OS-$ARCH"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download binary
curl -L -o "$INSTALL_DIR/moka" "$BINARY_URL"
chmod +x "$INSTALL_DIR/moka"

echo "✅ Binary installed to $INSTALL_DIR/moka"
echo ""

# Create data directory
echo "📁 Creating data directory..."
mkdir -p "$DATA_DIR"
echo "✅ Data directory created at $DATA_DIR"
echo ""

# Add to PATH if not already there
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "📝 Adding $INSTALL_DIR to PATH..."

    # Detect shell
    if [ -n "$BASH_VERSION" ]; then
        SHELL_RC="$HOME/.bashrc"
    elif [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="$HOME/.zshrc"
    else
        SHELL_RC="$HOME/.profile"
    fi

    echo "" >> "$SHELL_RC"
    echo "# Moka Finance Tracker" >> "$SHELL_RC"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"

    echo "✅ Added to $SHELL_RC"
    echo "⚠️  Run 'source $SHELL_RC' or restart your terminal"
fi

echo ""

# Download service file
echo "🔧 Setting up systemd service..."
curl -L -o /tmp/moka.service "https://raw.githubusercontent.com/$REPO/$LATEST_VERSION/moka.service"

# Replace placeholders in service file
sed -i "s|%u|$USER|g" /tmp/moka.service
sed -i "s|%h|$HOME|g" /tmp/moka.service

sudo cp /tmp/moka.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
rm /tmp/moka.service

echo "✅ Systemd service installed"
echo ""

# Add moka.local to /etc/hosts if not exists
if ! grep -q "moka.local" /etc/hosts; then
    echo "📝 Adding moka.local to /etc/hosts..."
    echo "127.0.0.1 moka.local" | sudo tee -a /etc/hosts > /dev/null
    echo "✅ Domain added"
fi

echo ""

# Start the service
echo "▶️  Starting Moka service..."
sudo systemctl start $SERVICE_NAME

echo ""
echo "✨ Installation complete!"
echo ""
echo "🌐 Access your app at: http://moka.local:9876"
echo ""
echo "Useful commands:"
echo "  - Check status:  sudo systemctl status moka"
echo "  - View logs:     sudo journalctl -u moka -f"
echo "  - Restart:       sudo systemctl restart moka"
echo "  - Stop:          sudo systemctl stop moka"
echo ""
echo "Data is stored in: $DATA_DIR"
echo ""

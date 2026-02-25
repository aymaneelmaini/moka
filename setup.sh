#!/bin/bash

echo "🚀 Setting up Moka Finance Tracker..."
echo ""

# 1. Add moka.local to /etc/hosts
echo "📝 Adding moka.local to /etc/hosts..."
echo "127.0.0.1 moka.local" | sudo tee -a /etc/hosts
echo "✅ Domain added!"
echo ""

# 2. Copy systemd service file
echo "📦 Installing systemd service..."
sudo cp moka.service /etc/systemd/system/
echo "✅ Service file copied!"
echo ""

# 3. Reload systemd
echo "🔄 Reloading systemd..."
sudo systemctl daemon-reload
echo "✅ Systemd reloaded!"
echo ""

# 4. Enable service (start on boot)
echo "⚙️  Enabling moka service..."
sudo systemctl enable moka
echo "✅ Service enabled (will start on boot)!"
echo ""

# 5. Start the service
echo "▶️  Starting moka service..."
sudo systemctl start moka
echo "✅ Service started!"
echo ""

# 6. Wait a moment for service to start
sleep 2

# 7. Check status
echo "📊 Service status:"
sudo systemctl status moka --no-pager
echo ""

echo "✨ Setup complete!"
echo ""
echo "🌐 Access your app at: http://moka.local:9876"
echo ""
echo "Useful commands:"
echo "  - Check status:  sudo systemctl status moka"
echo "  - View logs:     sudo journalctl -u moka -f"
echo "  - Restart:       sudo systemctl restart moka"
echo "  - Stop:          sudo systemctl stop moka"

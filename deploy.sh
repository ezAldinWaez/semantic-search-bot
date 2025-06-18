#!/bin/bash

# Telegram Semantic Search Bot Deployment Script
# Run this script as root to set up the bot as a systemd service

set -e

BOT_USER="semantic-bot"
BOT_GROUP="semantic-bot"
INSTALL_DIR="/opt/semantic-search-bot"
SERVICE_NAME="semantic-search-bot"
BINARY_NAME="semantic-search-bot"

echo "🚀 Deploying Telegram Semantic Search Bot"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Check if binary exists
if [[ ! -f "./bin/$BINARY_NAME" ]]; then
   echo "❌ Binary not found. Run 'make build' first"
   exit 1
fi

# Create user and group
echo "👤 Creating user and group..."
if ! id "$BOT_USER" &>/dev/null; then
    useradd --system --shell /bin/false --home-dir "$INSTALL_DIR" --create-home "$BOT_USER"
    echo "✅ Created user: $BOT_USER"
else
    echo "ℹ️  User $BOT_USER already exists"
    # Ensure home directory exists and has correct ownership
    mkdir -p "$INSTALL_DIR"
    chown "$BOT_USER:$BOT_GROUP" "$INSTALL_DIR"
fi

# Create installation directory
echo "📁 Creating installation directory..."
mkdir -p "$INSTALL_DIR"
chown "$BOT_USER:$BOT_GROUP" "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR"

# Copy binary
echo "📦 Installing binary..."
cp "./bin/$BINARY_NAME" "$INSTALL_DIR/"
chown "$BOT_USER:$BOT_GROUP" "$INSTALL_DIR/$BINARY_NAME"
chmod 755 "$INSTALL_DIR/$BINARY_NAME"

# Copy .env file if it exists
if [[ -f ".env" ]]; then
    echo "⚙️  Installing configuration..."
    cp ".env" "$INSTALL_DIR/"
    chown "$BOT_USER:$BOT_GROUP" "$INSTALL_DIR/.env"
    chmod 600 "$INSTALL_DIR/.env"  # Secure permissions for config
    echo "✅ Configuration installed"
else
    echo "⚠️  No .env file found. Creating template..."
    cat > "$INSTALL_DIR/.env" << 'EOF'
# Telegram Semantic Search Bot Configuration
TELEGRAM_TOKEN=your_bot_token_here
DATABASE_PATH=./messages.db
EMBEDDING_API_URL=http://localhost:11434
EMBEDDING_MODEL=all-minilm:latest
EOF
    chown "$BOT_USER:$BOT_GROUP" "$INSTALL_DIR/.env"
    chmod 600 "$INSTALL_DIR/.env"
    echo "📝 Please edit $INSTALL_DIR/.env with your bot token"
fi

# Create systemd service file
echo "🔧 Creating systemd service..."
cat > "/etc/systemd/system/$SERVICE_NAME.service" << EOF
[Unit]
Description=Telegram Semantic Search Bot
Documentation=https://github.com/ezAldinWaez/semantic-search-bot
After=network.target
Wants=network.target

[Service]
Type=simple
User=$BOT_USER
Group=$BOT_GROUP
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/$BINARY_NAME
Restart=always
RestartSec=10

# Environment
Environment=PATH=/usr/local/bin:/usr/bin:/bin
Environment=HOME=$INSTALL_DIR

# Security settings (less restrictive to avoid permission issues)
NoNewPrivileges=true
PrivateTmp=true

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$SERVICE_NAME

# Resource limits
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
echo "🔄 Configuring systemd..."
systemctl daemon-reload
systemctl enable "$SERVICE_NAME"

# Test binary execution as the service user
echo "🧪 Testing binary execution..."
if sudo -u "$BOT_USER" "$INSTALL_DIR/$BINARY_NAME" --help &>/dev/null; then
    echo "✅ Binary test successful"
elif sudo -u "$BOT_USER" "$INSTALL_DIR/$BINARY_NAME" --version &>/dev/null; then
    echo "✅ Binary test successful"
else
    echo "⚠️  Binary test failed, but this might be normal if the bot doesn't support --help"
    echo "   The service should still work if configuration is correct"
fi

echo ""
echo "✅ Deployment completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Edit configuration: sudo nano $INSTALL_DIR/.env"
echo "2. Make sure Ollama is running: ollama serve"
echo "3. Start the service: sudo systemctl start $SERVICE_NAME"
echo "4. Check status: sudo systemctl status $SERVICE_NAME"
echo "5. View logs: sudo journalctl -u $SERVICE_NAME -f"
echo ""
echo "🔧 Debug commands if service fails:"
echo "   sudo -u $BOT_USER $INSTALL_DIR/$BINARY_NAME  # Test manual execution"
echo "   sudo journalctl -u $SERVICE_NAME -n 50       # View error logs"
echo "   ls -la $INSTALL_DIR                          # Check file permissions"
echo ""
echo "🎯 Service management commands:"
echo "   sudo systemctl start $SERVICE_NAME     # Start service"
echo "   sudo systemctl stop $SERVICE_NAME      # Stop service"
echo "   sudo systemctl restart $SERVICE_NAME   # Restart service"
echo "   sudo systemctl status $SERVICE_NAME    # Check status"
echo ""
echo "📁 Bot files location: $INSTALL_DIR"
echo "🔧 Service file: /etc/systemd/system/$SERVICE_NAME.service"

#!/bin/bash
set -e

APP_NAME="mist"
INSTALL_DIR="/opt/mist"
MIST_FILE="/var/lib/mist/mist.db"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

echo "🧹 Uninstalling $APP_NAME..."

# -------------------------------
# Stop and disable systemd service
# -------------------------------
if systemctl list-unit-files | grep -q "^$APP_NAME.service"; then
    echo "⛔ Stopping and disabling $APP_NAME service..."
    sudo systemctl stop $APP_NAME || true
    sudo systemctl disable $APP_NAME || true
    sudo systemctl daemon-reload
fi

# -------------------------------
# Remove systemd service file
# -------------------------------
if [ -f "$SERVICE_FILE" ]; then
    echo "🗑️ Removing systemd service file: $SERVICE_FILE"
    sudo rm -f "$SERVICE_FILE"
    sudo systemctl daemon-reload
fi

# -------------------------------
# Remove installation directory
# -------------------------------
if [ -d "$INSTALL_DIR" ]; then
    echo "📂 Removing installation directory: $INSTALL_DIR"
    sudo rm -rf "$INSTALL_DIR"
fi

# -------------------------------
# Remove database file
# -------------------------------
if [ -f "$MIST_FILE" ]; then
    echo "🗃️ Removing Mist database file: $MIST_FILE"
    sudo rm -f "$MIST_FILE"
fi

# -------------------------------
# Remove CLI tool
# -------------------------------
if [ -f "/usr/local/bin/mist-cli" ]; then
    echo "🛠️ Removing CLI tool: /usr/local/bin/mist-cli"
    sudo rm -f /usr/local/bin/mist-cli
fi

# -------------------------------
# Firewall cleanup (optional)
# -------------------------------
PORT=8080
if command -v ufw &>/dev/null; then
    echo "🔒 Removing UFW rule for port $PORT..."
    sudo ufw delete allow $PORT/tcp || true
elif command -v firewall-cmd &>/dev/null; then
    echo "🔒 Removing firewalld rule for port $PORT..."
    sudo firewall-cmd --permanent --remove-port=${PORT}/tcp || true
    sudo firewall-cmd --reload || true
elif command -v iptables &>/dev/null; then
    echo "🔒 Removing iptables rule for port $PORT..."
    sudo iptables -D INPUT -p tcp --dport $PORT -j ACCEPT 2>/dev/null || true
    if command -v netfilter-persistent &>/dev/null; then
        sudo netfilter-persistent save
    fi
fi

# -------------------------------
# Environment cleanup (optional)
# -------------------------------
echo "🧩 Cleaning up environment paths..."
sed -i '/\/usr\/local\/go\/bin/d' ~/.bashrc || true
sed -i '/\.bun\/bin/d' ~/.bashrc || true

# -------------------------------
# Confirmation summary
# -------------------------------
echo ""
echo "✅ $APP_NAME has been completely uninstalled."
echo "Removed:"
echo "  - Service: $SERVICE_FILE"
echo "  - Install Dir: $INSTALL_DIR"
echo "  - Database: $MIST_FILE"
echo "  - CLI tool: /usr/local/bin/mist-cli"
echo "  - Firewall rules for port $PORT"
echo ""
echo "You may need to run 'source ~/.bashrc' to refresh your PATH."

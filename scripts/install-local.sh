#!/bin/bash
set -Eeo pipefail

# Local installation script for Mist.
# This copies the local repository to /opt/mist instead of cloning from GitHub.
# Useful for local development and testing.

LOG_FILE="/tmp/mist-install-local.log"
sudo rm -f "$LOG_FILE" 2>/dev/null || true
: > "$LOG_FILE"

REAL_USER="${SUDO_USER:-$USER}"
REAL_HOME="$(getent passwd "$REAL_USER" | cut -d: -f6)"

APP_NAME="mist"
INSTALL_DIR="/opt/mist"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SERVER_DIR="apps/server"
WEB_DIR="apps/web"
TRAEFIK_COMPOSE_PATH="deploy/traefik-compose.yml"
GO_BINARY_NAME="mist"
PORT=8080
MIST_FILE="/var/lib/mist/mist.db"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

SPINNER_PID=""
SUDO_KEEPALIVE_PID=""

export PATH="/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:$PATH"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[INFO]${NC} $1" | tee -a "$LOG_FILE"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1" | tee -a "$LOG_FILE"; }
error() { echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"; }
debug() { echo -e "${BLUE}[DEBUG]${NC} $1" | tee -a "$LOG_FILE"; }

spinner() {
    local i=0 chars='|/-\'
    while :; do
        i=$(( (i + 1) % 4 ))
        printf "\r⏳ %c" "${chars:$i:1}"
        sleep 0.1
    done
}

run_step() {
    printf "\n▶ $1\n"
    spinner & SPINNER_PID=$!
    if bash -c "$2" >>"$LOG_FILE" 2>&1; then
        kill "$SPINNER_PID" 2>/dev/null || true
        wait "$SPINNER_PID" 2>/dev/null || true
        printf "\r\033[K✔ Done\n"
        return 0
    else
        kill "$SPINNER_PID" 2>/dev/null || true
        wait "$SPINNER_PID" 2>/dev/null || true
        printf "\r\033[K✘ Failed\n"
        return 1
    fi
}

cleanup() {
    kill "$SPINNER_PID" 2>/dev/null || true
    kill "$SUDO_KEEPALIVE_PID" 2>/dev/null || true
}
trap cleanup EXIT

rollback() {
    error "Installation failed! Cleaning up..."
    if [ -f "$SERVICE_FILE" ]; then
        sudo systemctl stop "$APP_NAME" 2>>"$LOG_FILE" || true
        sudo systemctl disable "$APP_NAME" 2>>"$LOG_FILE" || true
        sudo rm -f "$SERVICE_FILE" 2>>"$LOG_FILE" || true
    fi
    error "Check logs: $LOG_FILE"
    tail -30 "$LOG_FILE" 2>/dev/null
    cleanup
    exit 1
}

trap - ERR
echo
echo "╔════════════════════════════════════════════╗"
echo "║ 🚀 Mist Local Installation                 ║"
echo "║ (Development Mode)                         ║"
echo "╚════════════════════════════════════════════╝"
echo
log "Starting Mist local installation from: $REPO_ROOT"

if [ "$EUID" -ne 0 ] && [ -z "${SUDO_USER:-}" ]; then
    error "This script requires sudo. Run: sudo bash scripts/install-local.sh"
    exit 1
fi

# Verify we're in a valid Mist directory
if [ ! -f "$REPO_ROOT/deploy/traefik-static.yml" ] || [ ! -d "$REPO_ROOT/$SERVER_DIR" ]; then
    error "This doesn't look like a Mist directory!"
    error "Expected to find deploy/traefik-static.yml and $SERVER_DIR/"
    exit 1
fi

echo "🔐 Verifying sudo access..."
sudo -v || exit 1

(while true; do sleep 60; sudo -n true || exit; done) 2>/dev/null &
SUDO_KEEPALIVE_PID=$!

log "Checking disk space..."
AVAILABLE=$(df /opt 2>/dev/null | tail -1 | awk '{print $4}' || echo "0")
if [ "$AVAILABLE" -lt 2000000 ]; then
    error "Need at least 2GB free in /opt"
    exit 1
fi

trap rollback ERR

log "Installing dependencies..."
if command -v apt >/dev/null 2>&1; then
    run_step "Installing packages (apt)" "sudo DEBIAN_FRONTEND=noninteractive apt update && sudo DEBIAN_FRONTEND=noninteractive apt install -y git curl build-essential wget unzip rsync" || exit 1
elif command -v dnf >/dev/null 2>&1; then
    run_step "Installing packages (dnf)" "sudo dnf install -y git curl gcc make wget unzip rsync" || exit 1
elif command -v yum >/dev/null 2>&1; then
    run_step "Installing packages (yum)" "sudo yum install -y git curl gcc make wget unzip rsync" || exit 1
elif command -v pacman >/dev/null 2>&1; then
    run_step "Installing packages (pacman)" "sudo pacman -Sy --noconfirm git curl base-devel wget unzip rsync" || exit 1
else
    error "Unsupported package manager"
    exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
    error "Docker not found. Install from: https://docs.docker.com/engine/install/"
    exit 1
fi
if ! docker ps >/dev/null 2>&1; then
    error "Docker not running or no permissions"
    exit 1
fi
log "Docker verified"

if ! command -v go >/dev/null 2>&1; then
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64) GO_ARCH="amd64";;
        aarch64|arm64) GO_ARCH="arm64";;
        armv7l) GO_ARCH="armv6l";;
        *) error "Unsupported arch: $ARCH"; exit 1;;
    esac
    
    run_step "Downloading Go" "wget -q https://go.dev/dl/go1.22.11.linux-${GO_ARCH}.tar.gz -O /tmp/go.tar.gz" || exit 1
    run_step "Installing Go" "sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf /tmp/go.tar.gz && rm -f /tmp/go.tar.gz" || exit 1
    
    grep -q '/usr/local/go/bin' "$REAL_HOME/.bashrc" 2>/dev/null || echo 'export PATH=$PATH:/usr/local/go/bin' >>"$REAL_HOME/.bashrc"
    export PATH="$PATH:/usr/local/go/bin"
fi
go version >>"$LOG_FILE" 2>&1 || { error "Go not working"; exit 1; }
log "Go ready"

# ---------------- Local Copy ----------------

log "Copying local files to $INSTALL_DIR..."
debug "Source: $REPO_ROOT"
debug "Destination: $INSTALL_DIR"

# Stop existing service if running
if sudo systemctl is-active --quiet "$APP_NAME" 2>/dev/null; then
    log "Stopping existing service..."
    sudo systemctl stop "$APP_NAME" || true
fi

# Stop Traefik if running
if docker ps --format '{{.Names}}' | grep -q "^traefik$"; then
    log "Stopping Traefik container..."
    docker compose -f "$INSTALL_DIR/$TRAEFIK_COMPOSE_PATH" down 2>/dev/null || docker stop traefik 2>/dev/null || true
fi

# Create install directory
run_step "Creating install directory" "sudo mkdir -p '$INSTALL_DIR'" || exit 1

# Copy files using rsync (excludes .git, node_modules, build artifacts)
run_step "Copying files" "sudo rsync -av --delete \
    --exclude='.git' \
    --exclude='node_modules' \
    --exclude='dist' \
    --exclude='build' \
    --exclude='*.log' \
    --exclude='mist.db' \
    --exclude='letsencrypt' \
    --exclude='apps/server/mist' \
    --exclude='apps/server/mist-cli' \
    --exclude='.env' \
    '$REPO_ROOT/' '$INSTALL_DIR/'" || exit 1

[ -d "$INSTALL_DIR/$SERVER_DIR" ] || { error "Server directory missing"; exit 1; }

run_step "Setting ownership" "sudo chown -R root:root '$INSTALL_DIR'" || exit 1

log "Local files copied successfully"

run_step "Creating data directories" "sudo mkdir -p /var/lib/mist/{traefik,logs,backups} && sudo touch '$MIST_FILE' && sudo chown -R root:root /var/lib/mist && sudo chmod -R 755 /var/lib/mist" || exit 1
run_step "Creating Traefik dynamic config" "sudo tee /var/lib/mist/traefik/dynamic.yml >/dev/null <<'EOF'
http:
  routers: {}
  services: {}
EOF
" || exit 1

cd "$INSTALL_DIR/$SERVER_DIR"
[ -f "go.mod" ] || { error "go.mod missing"; exit 1; }
run_step "Downloading dependencies" "cd '$INSTALL_DIR/$SERVER_DIR' && go mod download && go mod tidy" || exit 1
run_step "Building backend" "cd '$INSTALL_DIR/$SERVER_DIR' && go build -v -o '$GO_BINARY_NAME' ./cmd/mist" || exit 1
[ -f "$GO_BINARY_NAME" ] || { error "Binary not created"; exit 1; }
chmod +x "$GO_BINARY_NAME"
log "Build complete"

# ---------------- Dashboard Build ----------------

if [ -d "$INSTALL_DIR/$WEB_DIR" ]; then
    if ! command -v node >/dev/null 2>&1; then
        error "Node.js not found. Install from: https://nodejs.org/"
        exit 1
    fi
    if ! command -v bun >/dev/null 2>&1; then
        error "bun not found"
        exit 1
    fi
    
    run_step "Installing dashboard dependencies" "cd '$INSTALL_DIR/$WEB_DIR' && bun install" || exit 1
    run_step "Building dashboard" "cd '$INSTALL_DIR/$WEB_DIR' && bun run build" || exit 1
    
    # Move build output to server/static
    DASH_BUILD_DIR="$INSTALL_DIR/$WEB_DIR/dist"
    STATIC_DIR="$INSTALL_DIR/$SERVER_DIR/static"
    
    if [ -d "$DASH_BUILD_DIR" ]; then
        run_step "Moving dashboard build to static folder" "sudo rm -rf '$STATIC_DIR' && sudo mv '$DASH_BUILD_DIR' '$STATIC_DIR'" || exit 1
        log "Dashboard built and deployed to $SERVER_DIR/static"
    else
        error "Dashboard build output not found"
        exit 1
    fi
else
    warn "Web app directory not found, skipping dashboard build"
fi

# ---------------- CLI Tool ----------------

if [ -d "$INSTALL_DIR/$SERVER_DIR/cmd/cli" ]; then
    if run_step "Building CLI tool" "cd '$INSTALL_DIR/$SERVER_DIR' && go build -o mist-cli ./cmd/cli"; then
        if run_step "Installing CLI tool" "sudo cp '$INSTALL_DIR/$SERVER_DIR/mist-cli' /usr/local/bin/mist-cli && sudo chmod +x /usr/local/bin/mist-cli"; then
            log "CLI tool installed: mist-cli"
        else
            warn "Failed to install CLI tool, but continuing..."
        fi
    else
        warn "Failed to build CLI tool, but continuing..."
    fi
fi

run_step "Creating systemd service" "sudo tee '$SERVICE_FILE' >/dev/null <<'EOF'
[Unit]
Description=Mist Service (Local Dev)
After=network.target docker.service
Requires=docker.service

[Service]
WorkingDirectory=/opt/mist/apps/server
ExecStart=/opt/mist/apps/server/mist
Restart=always
RestartSec=5
User=root
Environment=PORT=8080
Environment=PATH=/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin

[Install]
WantedBy=multi-user.target
EOF
" || exit 1

run_step "Starting service" "sudo systemctl daemon-reload && sudo systemctl enable '$APP_NAME' && sudo systemctl start '$APP_NAME'" || exit 1
sleep 3
sudo systemctl is-active --quiet "$APP_NAME" || { error "Service failed to start"; sudo journalctl -u "$APP_NAME" -n 20; exit 1; }
log "Service running"

[ -f "$INSTALL_DIR/$TRAEFIK_COMPOSE_PATH" ] || { error "deploy/traefik-compose.yml missing"; exit 1; }
run_step "Creating Docker network" "docker network inspect traefik-net >/dev/null 2>&1 || docker network create traefik-net" || warn "Network creation failed"
run_step "Starting Traefik" "docker compose -f '$INSTALL_DIR/$TRAEFIK_COMPOSE_PATH' up -d" || warn "Traefik failed"

if command -v ufw >/dev/null 2>&1; then
    sudo ufw allow $PORT/tcp 2>&1 || true
elif command -v firewall-cmd >/dev/null 2>&1; then
    sudo firewall-cmd --permanent --add-port=${PORT}/tcp 2>&1 || true
    sudo firewall-cmd --reload 2>&1 || true
fi

log "Running health check..."
sleep 5
HTTP_OK=false
for i in {1..10}; do
    if curl -f -s -o /dev/null "http://localhost:$PORT/api/health" 2>>"$LOG_FILE"; then
        HTTP_OK=true
        break
    fi
    sleep 2
done
[ "$HTTP_OK" = true ] && log "HTTP check passed" || warn "HTTP check failed (may still be initializing)"

SERVER_IP=$(hostname -i 2>/dev/null | awk '{print $1}' || ip route get 1 2>/dev/null | awk '{print $7; exit}' || echo "localhost")
CLI_INSTALLED=""
if [ -f "/usr/local/bin/mist-cli" ]; then
    CLI_INSTALLED="║ 💻 CLI Tool: mist-cli --help              ║"
fi

echo
echo "╔════════════════════════════════════════════╗"
echo "║ 🎉 Mist local installation complete        ║"
printf "║ 👉 %-40s║\n" "http://$SERVER_IP:$PORT"
printf "║ 👉 %-40s║\n" "http://localhost:$PORT"
if [ -n "$CLI_INSTALLED" ]; then
    printf "║ %-42s ║\n" "$CLI_INSTALLED"
fi
echo "╚════════════════════════════════════════════╝"
echo
echo "📄 Logs: $LOG_FILE"
echo "🔧 Service: sudo systemctl status $APP_NAME"
echo "📋 Logs: sudo journalctl -u $APP_NAME -f"
if [ -f "/usr/local/bin/mist-cli" ]; then
    echo "💻 CLI Tool: mist-cli --help"
fi
echo
echo "💡 To reinstall after changes:"
echo "   sudo bash scripts/install-local.sh"
echo

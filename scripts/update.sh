#!/bin/bash
set -Eeo pipefail

LOG_FILE="/tmp/mist-update.log"
sudo rm -f "$LOG_FILE" 2>/dev/null || true
: > "$LOG_FILE"

REAL_USER="${SUDO_USER:-$USER}"
REAL_HOME="$(getent passwd "$REAL_USER" | cut -d: -f6)"

REPO="https://github.com/corecollectives/mist"
APP_NAME="mist"
INSTALL_DIR="/opt/mist"
GO_BACKEND_DIR="server"
GO_BINARY_NAME="mist"
DB_FILE="/var/lib/mist/mist.db"
BACKUP_DIR="/var/lib/mist/backups"
LOG_DIR="/var/lib/mist/logs"
LOCK_FILE="/var/lib/mist/update.lock"

if [ -d "$INSTALL_DIR/.git" ]; then
    CURRENT_BRANCH=$(cd "$INSTALL_DIR" && git rev-parse --abbrev-ref HEAD 2>/dev/null)
    if [ -z "$CURRENT_BRANCH" ] || [ "$CURRENT_BRANCH" = "HEAD" ]; then
        CURRENT_BRANCH=$(cd "$INSTALL_DIR" && git for-each-ref --format='%(upstream:short)' "$(git symbolic-ref -q HEAD)" 2>/dev/null | sed 's|^origin/||')
    fi
    
    if [ -z "$CURRENT_BRANCH" ] || [ "$CURRENT_BRANCH" = "HEAD" ]; then
        CURRENT_BRANCH="release"
    fi
    
    BRANCH="$CURRENT_BRANCH"
else
    BRANCH="release"
fi

SPINNER_PID=""
SUDO_KEEPALIVE_PID=""
BACKUP_COMMIT=""
BACKUP_DB=""

export PATH="/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:$PATH"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

spinner() {
    local i=0
    local chars='|/-\'
    while :; do
        i=$(( (i + 1) % 4 ))
        printf "\r⏳ %c" "${chars:$i:1}"
        sleep 0.1
    done
}

run_step() {
    local msg="$1"
    local cmd="$2"

    printf "\n▶ %s\n" "$msg"
    spinner &
    SPINNER_PID=$!

    if bash -c "$cmd" >>"$LOG_FILE" 2>&1; then
        kill "$SPINNER_PID" >/dev/null 2>&1 || true
        wait "$SPINNER_PID" 2>/dev/null || true
        printf "\r\033[K✔ Done\n"
        return 0
    else
        kill "$SPINNER_PID" >/dev/null 2>&1 || true
        wait "$SPINNER_PID" 2>/dev/null || true
        printf "\r\033[K✘ Failed\n"
        return 1
    fi
}

cleanup() {
    kill "$SPINNER_PID" >/dev/null 2>&1 || true
    kill "$SUDO_KEEPALIVE_PID" >/dev/null 2>&1 || true
    
    # Remove lock file
    rm -f "$LOCK_FILE"
}
trap cleanup EXIT

rollback() {
    error "Update failed! Attempting rollback..."
    
    export PATH="/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:$PATH"
    
    git config --global --add safe.directory "$INSTALL_DIR" 2>>"$LOG_FILE" || true
    
    if [ -n "$BACKUP_COMMIT" ]; then
        warn "Rolling back to commit: $BACKUP_COMMIT"
        cd "$INSTALL_DIR"
        git reset --hard "$BACKUP_COMMIT" >>"$LOG_FILE" 2>&1 || true
        
        # Rebuild from rollback
        cd "$INSTALL_DIR/$GO_BACKEND_DIR"
        go mod tidy >>"$LOG_FILE" 2>&1 || true
        go build -o "$GO_BINARY_NAME" >>"$LOG_FILE" 2>&1 || true
        
        sudo systemctl restart "$APP_NAME" >>"$LOG_FILE" 2>&1 || true
    fi
    
    if [ -n "$BACKUP_DB" ] && [ -f "$BACKUP_DB" ]; then
        warn "Rolling back database to backup"
        cp "$BACKUP_DB" "$DB_FILE" || true
    fi
    
    error "Rollback completed. Check logs at: $LOG_FILE"
    cleanup
    exit 1
}
trap rollback ERR


log "Starting Mist update process..."

trap - ERR

if [ "$EUID" -eq 0 ] || [ -n "${SUDO_USER:-}" ]; then
    log "Running with elevated privileges"
else
    error "This script requires sudo privileges"
    cleanup
    exit 1
fi

echo "🔐 Verifying sudo access..."
if ! sudo -v 2>>"$LOG_FILE"; then
    error "Failed to verify sudo access"
    cleanup
    exit 1
fi

(
    while true; do
        sleep 60
        sudo -n true || exit
    done
) 2>/dev/null &
SUDO_KEEPALIVE_PID=$!

if [ -f "$LOCK_FILE" ]; then
    error "Another update is already in progress!"
    error "If this is incorrect, remove: $LOCK_FILE"
    cleanup
    exit 1
fi

echo "$$" > "$LOCK_FILE"
log "Update lock acquired"

# Check if Mist is installed
if [ ! -d "$INSTALL_DIR/.git" ]; then
    error "Mist is not installed at $INSTALL_DIR"
    error "Please run install.sh first"
    cleanup
    exit 1
fi

log "Detected branch: $BRANCH"

if ! sudo systemctl is-active --quiet "$APP_NAME" 2>>"$LOG_FILE"; then
    warn "Mist service is not running!"
    warn "Continuing anyway, but this is unusual..."
fi

AVAILABLE_SPACE=$(df "$INSTALL_DIR" | tail -1 | awk '{print $4}')
if [ "$AVAILABLE_SPACE" -lt 500000 ]; then
    error "Insufficient disk space! Need at least 500MB free"
    cleanup
    exit 1
fi
log "Disk space check passed"

if ! command -v go >/dev/null 2>&1; then
    warn "Go not found in PATH, checking common locations..."
    
    GO_LOCATIONS=(
        "/usr/local/go/bin/go"
        "/usr/bin/go"
        "/opt/go/bin/go"
        "$REAL_HOME/.local/go/bin/go"
    )
    
    GO_FOUND=false
    for go_path in "${GO_LOCATIONS[@]}"; do
        if [ -x "$go_path" ]; then
            export PATH="$(dirname "$go_path"):$PATH"
            log "Found Go at: $go_path"
            GO_FOUND=true
            break
        fi
    done
    
    if [ "$GO_FOUND" = false ]; then
        error "Go is not installed or not found in any common location"
        error "Checked locations: ${GO_LOCATIONS[*]}"
        cleanup
        exit 1
    fi
fi

if ! go version >>"$LOG_FILE" 2>&1; then
    error "Go is installed but not working correctly"
    cleanup
    exit 1
fi
log "Go installation verified: $(go version | awk '{print $3}')"

if ! command -v docker >/dev/null 2>&1; then
    warn "Docker not found in PATH, checking common locations..."
    
    DOCKER_LOCATIONS=(
        "/usr/bin/docker"
        "/usr/local/bin/docker"
        "/opt/docker/bin/docker"
    )
    
    DOCKER_FOUND=false
    for docker_path in "${DOCKER_LOCATIONS[@]}"; do
        if [ -x "$docker_path" ]; then
            export PATH="$(dirname "$docker_path"):$PATH"
            log "Found Docker at: $docker_path"
            DOCKER_FOUND=true
            break
        fi
    done
    
    if [ "$DOCKER_FOUND" = false ]; then
        error "Docker is not installed or not found in any common location"
        error "Checked locations: ${DOCKER_LOCATIONS[*]}"
        cleanup
        exit 1
    fi
fi

if ! docker --version >>"$LOG_FILE" 2>&1; then
    error "Docker is installed but not working correctly"
    cleanup
    exit 1
fi
log "Docker installation verified: $(docker --version | awk '{print $3}' | tr -d ',')"

trap rollback ERR

log "Checking network connectivity..."
if ! curl -s --connect-timeout 10 https://github.com >/dev/null 2>&1; then
    error "No network connectivity to GitHub"
    error "Please check your internet connection"
    exit 1
fi
log "Network connectivity verified"


sudo mkdir -p "$BACKUP_DIR"
sudo mkdir -p "$LOG_DIR"
sudo chown "$REAL_USER:$REAL_USER" "$BACKUP_DIR"
sudo chown "$REAL_USER:$REAL_USER" "$LOG_DIR"
log "Backup directory ready"


TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_DB="$BACKUP_DIR/mist-$TIMESTAMP.db"
PERMANENT_LOG="$LOG_DIR/update-$TIMESTAMP.log"

if [ -f "$DB_FILE" ]; then
    log "Creating database backup..."
    if cp "$DB_FILE" "$BACKUP_DB"; then
        log "Database backed up to: $BACKUP_DB"
    else
        warn "Failed to backup database, but continuing..."
        BACKUP_DB=""
    fi
fi

cd "$INSTALL_DIR"
git config --local advice.detachedHead false >>"$LOG_FILE" 2>&1 || true

git config --global --add safe.directory "$INSTALL_DIR" >>"$LOG_FILE" 2>&1 || true

log "Verifying branch '$BRANCH' exists on remote..."
if ! git ls-remote --heads origin "$BRANCH" 2>>"$LOG_FILE" | grep -q "$BRANCH"; then
    error "Branch '$BRANCH' does not exist on remote"
    error "Available branches:"
    git ls-remote --heads origin 2>>"$LOG_FILE" | sed 's|.*refs/heads/||' | tee -a "$LOG_FILE"
    
    # Try to fall back to 'release' branch if it exists
    if git ls-remote --heads origin "release" 2>>"$LOG_FILE" | grep -q "release"; then
        warn "Falling back to 'release' branch"
        BRANCH="release"
    else
        error "Could not find a suitable branch to update from"
        exit 1
    fi
fi

log "Using branch: $BRANCH"

MAX_RETRIES=3
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if run_step "Fetching latest updates from $BRANCH branch (attempt $((RETRY_COUNT + 1))/$MAX_RETRIES)" "
        cd '$INSTALL_DIR' &&
        git fetch origin '$BRANCH'
    "; then
        break
    fi
    
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
        warn "Git fetch failed, retrying in 5 seconds..."
        sleep 5
    else
        error "Failed to fetch updates after $MAX_RETRIES attempts"
        error "Last fetch error:"
        git fetch origin "$BRANCH" 2>&1 | tail -10 | tee -a "$LOG_FILE"
        exit 1
    fi
done

cd "$INSTALL_DIR"
LOCAL_COMMIT=$(git rev-parse HEAD 2>>"$LOG_FILE")
REMOTE_COMMIT=$(git rev-parse origin/$BRANCH 2>>"$LOG_FILE")
BACKUP_COMMIT="$LOCAL_COMMIT"

if [ "$LOCAL_COMMIT" = "$REMOTE_COMMIT" ]; then
    log "Already up to date"
    log "Current version: ${LOCAL_COMMIT:0:7}"
    
    # Save log even when no update is needed
    if [ -n "$PERMANENT_LOG" ]; then
        cp "$LOG_FILE" "$PERMANENT_LOG" 2>/dev/null || true
        log "Update check log saved: $PERMANENT_LOG"
    fi
    
    exit 0
fi

log "Updates available"
log "Current commit: ${LOCAL_COMMIT:0:7}"
log "Latest commit:  ${REMOTE_COMMIT:0:7}"

log "Changes in this update:"
git log --oneline "$LOCAL_COMMIT..$REMOTE_COMMIT" | head -10 | tee -a "$LOG_FILE"


run_step "Creating backup tag" "
    cd '$INSTALL_DIR' &&
    git tag -f backup-$TIMESTAMP
"


if ! run_step "Updating repository to latest version" "
    cd '$INSTALL_DIR' &&
    git reset --hard origin/'$BRANCH'
"; then
    error "Failed to update repository"
    rollback
fi

NEW_COMMIT=$(git rev-parse HEAD)
if [ "$NEW_COMMIT" != "$REMOTE_COMMIT" ]; then
    error "Repository update verification failed!"
    rollback
fi
log "Repository update verified"


log "Running database migrations (if any)..."
cd "$INSTALL_DIR/$GO_BACKEND_DIR"

if [ -d "db/migrations" ]; then
    log "Migration files detected"
fi

log "Cleaning Go build cache..."
go clean -cache -modcache -i -r >>"$LOG_FILE" 2>&1 || true

if ! run_step "Rebuilding backend binary" "
    cd '$INSTALL_DIR/$GO_BACKEND_DIR' &&
    go mod download &&
    go mod tidy &&
    go build -v -o '$GO_BINARY_NAME'
"; then
    error "Failed to rebuild backend"
    rollback
fi

if [ ! -f "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" ]; then
    error "Backend binary was not created!"
    rollback
fi

if [ ! -x "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" ]; then
    error "Backend binary is not executable!"
    rollback
fi

if ! file "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" | grep -q "executable" >>"$LOG_FILE" 2>&1; then
    error "Backend binary appears to be corrupted!"
    rollback
fi

log "Backend binary verified"


if [ -d "$INSTALL_DIR/cli" ]; then
    if ! run_step "Rebuilding CLI tool" "
        cd '$INSTALL_DIR/cli' &&
        go mod tidy &&
        go build -o mist-cli
    "; then
        warn "Failed to rebuild CLI, but continuing..."
    else
        if ! run_step "Installing CLI tool" "
            sudo cp '$INSTALL_DIR/cli/mist-cli' /usr/local/bin/mist-cli &&
            sudo chmod +x /usr/local/bin/mist-cli
        "; then
            warn "Failed to install CLI, but continuing..."
        fi
    fi
fi


log "Stopping Mist service..."
STOP_RETRIES=0
MAX_STOP_RETRIES=3

while [ $STOP_RETRIES -lt $MAX_STOP_RETRIES ]; do
    if sudo systemctl stop "$APP_NAME" --no-block >>"$LOG_FILE" 2>&1; then
        # Wait for service to stop
        sleep 3
        if ! sudo systemctl is-active --quiet "$APP_NAME"; then
            log "Service stopped successfully"
            break
        fi
    fi
    
    STOP_RETRIES=$((STOP_RETRIES + 1))
    if [ $STOP_RETRIES -lt $MAX_STOP_RETRIES ]; then
        warn "Failed to stop service, retrying..."
        sleep 2
    else
        warn "Failed to stop service gracefully, forcing..."
        sudo systemctl kill "$APP_NAME" >>"$LOG_FILE" 2>&1 || true
        sleep 3
    fi
done

if ! run_step "Reloading systemd and starting service" "
    sudo systemctl daemon-reload &&
    sudo systemctl start '$APP_NAME'
"; then
    error "Failed to start service after update"
    rollback
fi


log "Waiting for service to start..."
RETRY_COUNT=0
MAX_RETRIES=30

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if sudo systemctl is-active --quiet "$APP_NAME" 2>>"$LOG_FILE"; then
        log "Service is active"
        break
    fi
    sleep 1
    RETRY_COUNT=$((RETRY_COUNT + 1))
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    error "Service failed to start within 30 seconds!"
    error "Service status:"
    sudo systemctl status "$APP_NAME" --no-pager >>"$LOG_FILE" 2>&1 || true
    error "Recent logs:"
    sudo journalctl -u "$APP_NAME" -n 50 --no-pager >>"$LOG_FILE" 2>&1 || true
    error "Rolling back..."
    rollback
fi

log "Performing extended health check..."
sleep 5

if ! sudo systemctl is-active --quiet "$APP_NAME" 2>>"$LOG_FILE"; then
    error "Service started but crashed immediately!"
    error "Service status:"
    sudo systemctl status "$APP_NAME" --no-pager >>"$LOG_FILE" 2>&1 || true
    error "Recent logs:"
    sudo journalctl -u "$APP_NAME" -n 50 --no-pager >>"$LOG_FILE" 2>&1 || true
    error "Check logs: sudo journalctl -u $APP_NAME -n 50"
    rollback
fi

log "Service health check passed"

log "Checking HTTP endpoint..."
HTTP_RETRIES=0
MAX_HTTP_RETRIES=10
HTTP_SUCCESS=false

while [ $HTTP_RETRIES -lt $MAX_HTTP_RETRIES ]; do
    if curl -f -s -o /dev/null -w "%{http_code}" --connect-timeout 5 "http://localhost:8080/api/health" 2>>"$LOG_FILE" | grep -q "200"; then
        log "HTTP health check passed"
        HTTP_SUCCESS=true
        break
    fi
    sleep 2
    HTTP_RETRIES=$((HTTP_RETRIES + 1))
done

if [ "$HTTP_SUCCESS" = false ]; then
    warn "HTTP health check failed after $MAX_HTTP_RETRIES attempts"
    warn "Service is running but may still be initializing"
    warn "Check logs if issues persist: sudo journalctl -u $APP_NAME -n 50"
fi


if [ -f "$INSTALL_DIR/traefik-compose.yml" ]; then
    if ! run_step "Updating Traefik" "
        docker compose -f '$INSTALL_DIR/traefik-compose.yml' up -d
    "; then
        warn "Failed to update Traefik, but Mist update succeeded"
    fi
fi


log "Cleaning up old backups (keeping last 5)..."
cd "$BACKUP_DIR"
ls -t mist-*.db 2>/dev/null | tail -n +6 | xargs -r rm -f || true

cd "$LOG_DIR"
ls -t update-*.log 2>/dev/null | tail -n +11 | xargs -r rm -f || true

cd "$INSTALL_DIR"
git tag | grep "^backup-" | sort -r | tail -n +11 | xargs -r git tag -d 2>/dev/null || true


log "Saving update log to: $PERMANENT_LOG"
cp "$LOG_FILE" "$PERMANENT_LOG" 2>/dev/null || true


log ""
log "╔════════════════════════════════════════════╗"
log "║ 🎉 Mist has been updated successfully      ║"
log "║                                            ║"
log "║ From: ${LOCAL_COMMIT:0:7}                             ║"
log "║ To:   ${REMOTE_COMMIT:0:7}                             ║"
log "╚════════════════════════════════════════════╝"
log ""
log "📄 Update log saved: $PERMANENT_LOG"
log "💾 Database backup: $BACKUP_DB"
log "🔄 Service status: sudo systemctl status $APP_NAME"
log "📋 Service logs: sudo journalctl -u $APP_NAME -n 50"
log ""
log "If you experience issues, you can rollback with:"
log "  cd $INSTALL_DIR && git reset --hard backup-$TIMESTAMP"
log "  sudo systemctl restart $APP_NAME"

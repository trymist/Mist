# Configuration

Configure Mist for your environment.

## Configuration System

Mist uses a **database-driven configuration system**. Most settings are stored in the SQLite database at `/var/lib/mist/mist.db` and managed through the web UI or API.

::: tip
There is no `/etc/mist/config.yml` file. System settings are configured through the **Settings → System** page in the web interface.
:::

## Fixed Configuration Values

Some values are hardcoded in the application:

### Server Port

- **Port:** `8080` (hardcoded)
- **Location:** `server/api/main.go:33`
- **Cannot be changed** via environment variables

### Session Configuration

- **JWT Expiry:** 31 days (hardcoded)
- **Location:** `server/api/middleware/auth.go:21`

### File Paths

- **Root Path:** `/var/lib/mist`
- **Logs:** `/var/lib/mist/logs`
- **Avatars:** `/var/lib/mist/uploads/avatar`
- **Traefik Config:** `/var/lib/mist/traefik`
- **Certificates:** `/opt/mist/letsencrypt/acme.json`

## Database Configuration

All system settings and credentials are stored in the SQLite database.

### GitHub App Configuration

GitHub App credentials are stored in the `github_app` table:
- App ID
- Client ID
- Client Secret
- Webhook Secret
- Private Key (RSA)

Configure through: **Settings → Integrations → GitHub App**

See [GitHub App Setup](./github-app) for configuration instructions.

## Docker Configuration

Mist uses the Docker socket at `/var/run/docker.sock`.

Ensure the user running Mist has Docker permissions:

```bash
sudo usermod -aG docker mist-user
```

## Traefik Integration

Mist automatically generates Traefik configuration in `/var/lib/mist/traefik/`:
- `traefik.yml` - Static configuration
- `dynamic.yml` - Dynamic routing rules (auto-generated)

See [Traefik Setup](./traefik) for detailed configuration.

## System Settings

System-wide settings are configured through the Mist UI (requires owner role).

### Wildcard Domain

Configure automatic domain generation for web applications:

1. Navigate to **Settings** → **System**
2. Enter your wildcard domain (e.g., `example.com` or `*.example.com`)
3. Optionally configure the Mist dashboard subdomain name (default: `mist`)

When configured, new web applications automatically receive domains in the format:
```
{project-name}-{app-name}.{wildcard-domain}
```

**Example:**
- Wildcard domain: `apps.example.com`
- Project: `production`, App: `api`
- Auto-generated domain: `production-api.apps.example.com`

**DNS Requirements:**
Configure a wildcard DNS record:
```
Type: A
Name: *
Value: YOUR_SERVER_IP
```

Or for nested subdomains:
```
Type: A
Name: *.apps
Value: YOUR_SERVER_IP
```

[Learn more about wildcard domains →](../guide/domains#wildcard-domain-configuration)

## Coming Soon

The following configuration features are planned:

- **Configurable Server Port** - Change port via UI or configuration
- **Configurable JWT Secret** - Set JWT secret during installation
- **SMTP Configuration** - Email settings for notifications
- **Backup Settings** - Automatic backup configuration
- **Resource Limits** - Global resource limits per project/application
- **Logging Configuration** - Log levels and destinations
- **Custom SSL Certificates** - Upload custom TLS certificates

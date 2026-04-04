# Domains & SSL

Configure custom domains for your applications.

## Wildcard Domain Configuration

Mist supports automatic domain generation for web applications using a wildcard domain configuration.

### What is Wildcard Domain?

When configured, Mist automatically generates a subdomain for every new web application using the pattern:
```
{project-name}-{app-name}.{wildcard-domain}
```

### Setting Up Wildcard Domain

1. Go to **Settings** (requires owner role)
2. Enter your wildcard domain (e.g., `example.com` or `*.example.com`)
3. Optionally configure the Mist dashboard subdomain name (default: `mist`)
4. Click **"Save"**

### DNS Configuration for Wildcard

To use wildcard domains, configure a wildcard DNS record:

```
Type: A
Name: *
Value: YOUR_SERVER_IP
TTL: 3600
```

Or for subdomains:
```
Type: A
Name: *.apps
Value: YOUR_SERVER_IP
TTL: 3600
```

This allows any subdomain (e.g., `production-api.example.com`, `staging-web.example.com`) to route to your server.

### Benefits

- **Automatic domains**: No manual domain configuration needed for new apps
- **Consistent naming**: Standardized domain format across all applications
- **SSL automation**: Auto-generated domains get automatic Let's Encrypt certificates
- **Easy management**: Change the base domain in one place

### Examples

**Wildcard domain**: `example.com`
- Project: `production`, App: `api` → `production-api.example.com`
- Project: `staging`, App: `frontend` → `staging-frontend.example.com`

**Wildcard domain**: `apps.mysite.com`
- Project: `personal`, App: `blog` → `personal-blog.apps.mysite.com`

::: tip
Auto-generated domains only apply to web applications. Service and database applications do not receive automatic domains.
:::

## Adding a Domain

1. Navigate to your application
2. Go to **Domains** tab
3. Click **"Add Domain"**
4. Enter your domain name (e.g., `app.example.com`)
5. Click **"Add"**

## DNS Configuration

Point your domain to your Mist server:

### A Record

```
Type: A
Name: app (or @ for root)
Value: YOUR_SERVER_IP
TTL: 3600
```

### CNAME Record (Subdomain)

```
Type: CNAME
Name: app
Value: yourdomain.com
TTL: 3600
```

## Verify DNS Propagation

Check if DNS has propagated:

```bash
dig app.example.com
nslookup app.example.com
```

Propagation can take 5 minutes to 48 hours.

## SSL/TLS Certificates

Mist automatically provisions and manages SSL/TLS certificates for your domains using Let's Encrypt and Traefik.

### Automatic SSL

When you add a domain to your application:

1. **Domain Verification**: Mist verifies your DNS is correctly configured
2. **Certificate Request**: Traefik automatically requests a certificate from Let's Encrypt
3. **HTTPS Enabled**: Your domain is secured with SSL/TLS
4. **HTTP Redirect**: HTTP traffic is automatically redirected to HTTPS

### Auto-Renewal

Certificates are automatically renewed by Traefik:
- Renewal happens before expiration (typically 30 days before)
- Zero downtime during renewal
- No manual intervention required

### Requirements

For automatic SSL to work:
- Domain DNS must point to your Mist server (A record)
- Ports 80 and 443 must be open and accessible
- Domain must be publicly accessible for Let's Encrypt verification

::: tip Certificate Status
You can view the SSL status of your domains in the Domains tab of your application. Status will show as "pending" during issuance and "active" once the certificate is installed.
:::

### Coming Soon

- **Wildcard Support** - `*.example.com` certificates
- **Custom Certificates** - Upload your own SSL certs

## Multiple Domains

Add multiple domains to the same application:

```
app.example.com
www.app.example.com
app.yourdomain.com
```

All domains will route to the same container.

## Domain Management

### Edit Domain

Update domain configuration in the domains list.

### Delete Domain

Remove a domain from your application:

1. Click the delete icon next to the domain
2. Confirm deletion

## Troubleshooting

### Domain Not Resolving

- Verify DNS records are correct
- Wait for DNS propagation (can take up to 48 hours)
- Check DNS with: `dig app.example.com` or `nslookup app.example.com`
- Clear local DNS cache: `sudo systemd-resolve --flush-caches`

### SSL Certificate Not Issuing

- Verify domain DNS points to your server IP
- Ensure ports 80 and 443 are open in your firewall
- Check that domain is publicly accessible
- View Traefik logs for errors: `docker logs traefik`
- Let's Encrypt has rate limits - wait if you've made many requests

### HTTPS Not Working

- Wait a few minutes for certificate issuance
- Check SSL status in Domains tab (should show "active")
- Verify Traefik container is running: `docker ps | grep traefik`
- Check Traefik configuration in `/etc/traefik/`

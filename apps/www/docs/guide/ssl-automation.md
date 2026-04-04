# SSL Automation

Automatic SSL/TLS certificate provisioning with Let's Encrypt.

Mist provides automatic SSL certificate management through Traefik and Let's Encrypt integration, ensuring your applications are always secured with HTTPS.

## Features

### Automatic Certificate Issuance

- **Let's Encrypt Integration**: Free SSL certificates via ACME protocol
- **Automatic Generation**: Certificates issued when you add a domain
- **HTTP-01 Challenge**: Domain verification through HTTP
- **Zero Configuration**: No manual certificate management needed

### Auto-Renewal

- **Automatic Renewal**: Certificates renewed before expiry (typically 30 days before)
- **Zero Downtime**: Renewals happen seamlessly in the background
- **No Manual Intervention**: Traefik handles the entire renewal process

### Security Features

- **Force HTTPS**: Automatic HTTP to HTTPS redirect enabled by default
- **TLS 1.2+**: Modern encryption standards
- **Secure Headers**: Best practice security headers configured

## How It Works

### Adding a Domain

1. **Add Domain**: Enter your domain name in the application's Domains tab
2. **DNS Verification**: Ensure your domain's A record points to your Mist server
3. **Automatic Issuance**: Traefik detects the new domain and requests a certificate
4. **Certificate Installation**: Let's Encrypt issues the certificate and Traefik installs it
5. **HTTPS Enabled**: Your application is now accessible via HTTPS

The entire process is automatic and typically takes 1-2 minutes after DNS propagation.

### Domain Labels

Mist automatically configures Traefik labels on your containers:

```yaml
traefik.enable=true
traefik.http.routers.app-{id}.rule=Host(`your-domain.com`)
traefik.http.routers.app-{id}.entrypoints=websecure
traefik.http.routers.app-{id}.tls=true
traefik.http.routers.app-{id}.tls.certresolver=le
traefik.http.routers.app-{id}-http.rule=Host(`your-domain.com`)
traefik.http.routers.app-{id}-http.entrypoints=web
traefik.http.routers.app-{id}-http.middlewares=https-redirect
traefik.http.middlewares.https-redirect.redirectscheme.scheme=https
```

### Certificate Storage

Certificates are stored in Traefik's ACME storage:
- Location: `/etc/traefik/acme.json`
- Format: JSON file with encrypted certificates
- Backup: Recommended to backup this file regularly

## Requirements

For SSL automation to work properly:

1. **Valid Domain**: Domain must be registered and active
2. **DNS Configuration**: A record pointing to your server's public IP
3. **Port Access**: Ports 80 (HTTP) and 443 (HTTPS) must be publicly accessible
4. **Public Server**: Server must be reachable from the internet for Let's Encrypt verification
5. **No Rate Limits**: Be aware of Let's Encrypt rate limits (50 certificates per domain per week)

## Monitoring Certificate Status

### In Dashboard

View certificate status in the Domains tab:
- **Pending**: Certificate request in progress
- **Active**: Certificate successfully issued and installed
- **Failed**: Certificate issuance failed (check logs)

### Via Logs

Check Traefik logs for certificate operations:

```bash
docker logs traefik
```

Look for messages like:
- `Obtaining certificate for domain`
- `Certificate obtained successfully`
- `Renewing certificate`

## Troubleshooting

### Certificate Not Issuing

**DNS Not Propagated**
- Wait for DNS propagation (up to 48 hours)
- Verify with: `dig your-domain.com`

**Port Not Accessible**
- Check firewall: `sudo ufw status`
- Ensure ports 80 and 443 are open
- Test accessibility: `curl -I http://your-domain.com`

**Let's Encrypt Rate Limit**
- Let's Encrypt has limits: 50 certificates per registered domain per week
- Wait if you've hit the limit
- Use Let's Encrypt staging environment for testing

**Domain Verification Failed**
- Ensure domain points to correct IP
- Check that no other service is using port 80
- Verify Traefik is running: `docker ps | grep traefik`

### Certificate Not Renewing

**Check ACME Storage**
```bash
ls -la /etc/traefik/acme.json
```

**Check Traefik Configuration**
```bash
cat /etc/traefik/traefik.yml
```

Ensure certificate resolver is configured:
```yaml
certificatesResolvers:
  le:
    acme:
      email: your-email@example.com
      storage: /acme.json
      httpChallenge:
        entryPoint: web
```

### Debugging

Enable debug logging in Traefik:
1. Edit Traefik configuration
2. Set log level to DEBUG
3. Restart Traefik: `docker restart traefik`
4. Monitor logs: `docker logs -f traefik`

## Best Practices

### Email Configuration

Configure a valid email in Traefik for Let's Encrypt notifications:
- Notifies about certificate expiration
- Required for Let's Encrypt terms acceptance
- Used for important updates

### Backup ACME Storage

Regularly backup `/etc/traefik/acme.json`:
```bash
cp /etc/traefik/acme.json /backup/acme.json.$(date +%Y%m%d)
```

### Test Domains

Use staging environment for testing:
- Avoids rate limits
- Tests certificate issuance without counting against production limits
- Configure in Traefik: `caServer: https://acme-staging-v02.api.letsencrypt.org/directory`

### Multiple Domains

You can add multiple domains to a single application:
- Each domain gets its own certificate
- All domains route to the same container
- Certificates managed independently

## Coming Soon

- **Wildcard Certificates**: Support for `*.example.com` domains
- **Custom Certificates**: Upload and manage your own SSL certificates
- **Certificate Dashboard**: View all certificates and expiration dates
- **DNS-01 Challenge**: Support for wildcard certificates via DNS validation
- **Certificate Notifications**: Email alerts before certificate expiration

## Related Documentation

- [Domains Guide](./domains) - Managing application domains
- [Traefik Configuration](../deployment/traefik) - Advanced Traefik setup
- [Applications](./applications) - Application management

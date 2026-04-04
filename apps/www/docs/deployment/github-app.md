# GitHub App Setup

Configure GitHub App integration for automatic deployments from GitHub repositories.

## Overview

Mist uses **GitHub App integration** (not OAuth Apps or Personal Access Tokens) to:
- Access your repositories
- Receive webhook notifications for code pushes
- Trigger automatic deployments when you push code
- Clone private repositories securely

## Automatic GitHub App Creation

Mist provides a streamlined process to create a GitHub App directly from the dashboard.

### Step 1: Access GitHub App Setup

1. Log in to Mist as **Owner** (only Owners can create GitHub Apps)
2. Navigate to **Settings → Integrations → GitHub**
3. Click **"Create GitHub App"** or **"Connect GitHub"**

### Step 2: Create GitHub App

1. Log in to Mist as **Owner** (only Owners can create GitHub Apps)
2. Navigate to **Settings → Integrations → GitHub**
3. Click **"Create GitHub App"** or **"Connect GitHub"**

### Step 3: Automatic GitHub App Creation

When you click "Create GitHub App" in Mist:

1. **Mist generates a manifest** with:
   - App name: `Mist-{random-5-digit-number}` (e.g., `Mist-12345`)
   - Webhook URL: Auto-configured based on your Mist installation
   - Callback URLs for OAuth
   - Required permissions (contents, metadata, webhooks)
   - Event subscriptions (push, pull_request, deployment_status)

2. **You're redirected to GitHub** with pre-filled configuration

3. **GitHub creates the app** and redirects back to Mist

4. **Mist stores credentials** in the database:
   - App ID
   - Client ID
   - Client Secret
   - Private Key (RSA)
   - Webhook Secret

::: tip Fully Automated
You don't need to manually configure permissions, generate keys, or copy credentials. Mist handles everything automatically.
:::

## GitHub App Permissions

The automatically created GitHub App requests these permissions:

**Repository Permissions:**
- **Contents**: Read - Access repository files and commits
- **Metadata**: Read - Basic repository information
- **Pull Requests**: Read - View pull requests
- **Deployments**: Read - Deployment status
- **Administration**: Write - Manage repository settings
- **Webhooks**: Write - Create and manage webhooks

**Subscribed Events:**
- **Push** - Triggered when code is pushed to repository
- **Pull Request** - Triggered on PR events
- **Deployment Status** - Triggered on deployment updates

**Webhook Security:**
- Webhook signatures are validated using HMAC-SHA256
- Only requests with valid signatures are processed
- Webhook secret is automatically generated during app creation

## Install GitHub App on Repositories

After creating the GitHub App, you need to install it on the repositories you want to deploy:

### Install on Your Account/Organization

1. **From Mist Dashboard**:
   - After successful app creation, you're redirected to installation page
   - Or go to GitHub → Settings → Applications → Your Mist App → Install

2. **Select Installation Target**:
   - Personal account
   - Or organization (if you have access)

3. **Choose Repository Access**:
   - **All repositories**: Mist can access all current and future repos
   - **Only select repositories**: Choose specific repos to grant access

4. **Complete Installation**:
   - Click **"Install"**
   - You'll be redirected back to Mist

### Verify Installation

In Mist dashboard:

1. **Create a new application** or edit existing
2. Go to **Git Configuration** section
3. **Repository dropdown** should show your installed repositories
4. Select repository and branch
5. Configure Dockerfile path
6. Deploy!

::: tip Multiple Organizations
If you have multiple GitHub organizations, you can install the same GitHub App on each organization separately. Each installation is independent.
:::

## Configure Webhook URL for Production

If you need to update the GitHub App webhook URL after initial creation:

1. Go to [GitHub Apps Settings](https://github.com/settings/apps)
2. Select your Mist app
3. Update **Webhook URL** to: `https://mist.yourdomain.com/api/github/webhook`
4. Update **Callback URL** to: `https://mist.yourdomain.com/api/github/callback`
5. Update **Setup URL** to: `https://mist.yourdomain.com/api/github/installation/callback`
6. Click **"Save changes"**

## Automatic Deployments

Once configured, Mist automatically deploys when you push code:

### How it Works

1. **You push code** to a branch configured in Mist
2. **GitHub sends webhook** to Mist
3. **Mist receives push event** and starts deployment
4. **Clones repository** using GitHub App authentication
5. **Builds Docker image** from Dockerfile
6. **Deploys container** with configured settings
7. **Updates deployment status** in Mist dashboard

### Monitor Webhook Deliveries

Check if webhooks are being delivered successfully:

1. Go to [GitHub Apps Settings](https://github.com/settings/apps)
2. Select your Mist app
3. Click **"Advanced" tab**
4. View **Recent Deliveries**

Look for:
- ✅ **200 OK**: Webhook received successfully
- ❌ **Failed**: Check error details and Mist logs

## Test Webhook Integration

Verify that automatic deployments work correctly:

### Method 1: Push to Repository

1. **Configure application** in Mist with GitHub repository
2. **Make a change** to your code
3. **Commit and push**:
   ```bash
   git add .
   git commit -m "Test deployment"
   git push origin main
   ```
4. **Check Mist dashboard**:
   - Navigate to application → Deployments
   - You should see a new deployment starting automatically
   - Monitor build logs in real-time

### Method 2: Redeliver Webhook

Test without pushing new code:

1. Go to [GitHub Apps Settings](https://github.com/settings/apps)
2. Click your Mist app → **Advanced** tab
3. Find a recent **push** event in Recent Deliveries
4. Click **"Redeliver"**
5. Check Mist dashboard for new deployment

### Method 3: Check Mist Logs

Monitor webhook reception in real-time:

```bash
# Watch for incoming webhooks
sudo journalctl -u mist -f | grep -i "webhook\|github"

# Check recent webhook activity
sudo journalctl -u mist -n 100 | grep -i webhook
```

Look for log messages like:
- `Received GitHub webhook`
- `Processing push event for repository: owner/repo`
- `Starting deployment for application: app-name`

## Troubleshooting

### Webhooks Not Being Received

**Check webhook deliveries in GitHub:**

1. GitHub Apps → Your App → Advanced tab
2. Look at Recent Deliveries
3. If failed:
   - ❌ **Connection timeout**: Firewall blocking GitHub IPs
   - ❌ **Connection refused**: Mist not running or port 8080 not accessible
   - ❌ **404 Not Found**: Wrong webhook URL
   - ❌ **500 Internal Server Error**: Check Mist logs for errors

**Verify webhook URL is accessible:**

```bash
# Test from external server (or use webhook.site for testing)
curl -X POST https://your-mist-domain.com/api/github/webhook \
  -H "Content-Type: application/json" \
  -d '{"test": true}'

# Should return 200 or similar (not 404/timeout)
```

**Check firewall rules:**

```bash
# Verify port 80/443 are open
sudo ufw status | grep -E "80|443"

# If using Mist API directly (port 8080), verify it's accessible
curl http://your-domain:8080/health  # Should work
```

**Allow GitHub webhook IPs (optional, for security):**

```bash
# Get GitHub's current IP ranges
curl https://api.github.com/meta | jq -r '.hooks[]'

# Example: Allow GitHub IPs for webhook endpoint
# (Implement via nginx/firewall rules)
```

### Repositories Not Showing in Dropdown

**Verify GitHub App installation:**

1. Go to GitHub → Settings → Applications → Installed GitHub Apps
2. Find your Mist app
3. Verify it's installed on the correct account/organization
4. Check repository access (All repositories or Selected ones)

**Reinstall GitHub App:**

If repositories still don't appear:
1. Uninstall the GitHub App from GitHub settings
2. Go back to Mist → Settings → Integrations
3. Reinstall the app
4. Select repositories to grant access

**Check Mist logs:**

```bash
sudo journalctl -u mist -n 100 | grep -i "github\|repository"
```

### Private Key Errors

If you see errors related to GitHub authentication:

**Verify GitHub App exists in database:**

```bash
sqlite3 /var/lib/mist/mist.db "SELECT app_id, client_id FROM github_app;"
```

**Check private key format:**

The private key should be stored as RSA format in the database. If corrupted:
1. Delete the GitHub App from Mist settings
2. Recreate it using the automatic setup process

### Deployment Not Triggering

**Verify branch configuration:**

- Make sure you're pushing to the **branch configured in Mist**
- Default is often `main` but could be `master`, `develop`, etc.

**Check deployment logs:**

```bash
# View recent deployments
sudo journalctl -u mist | grep -i "deployment\|deploy"

# Check for errors
sudo journalctl -u mist -p err
```

**Manually trigger deployment:**

From Mist dashboard:
1. Go to application details
2. Click **"Deploy"** button
3. Monitor build logs

### Webhook URL Changed

If you change your Mist domain or IP:

1. Update webhook URL in GitHub App settings (see "Configure Webhook URL for Production" above)
2. Or delete and recreate the GitHub App through Mist UI

## Security Best Practices

### Webhook Security

Mist validates all incoming webhook requests using HMAC-SHA256 signatures. Only authenticated requests from GitHub are processed.

**Additional Security (Optional):**

If you want an extra layer of security, restrict webhook endpoint access to GitHub's IP ranges:

**Option 1: IP Whitelist (Nginx)**

```nginx
# /etc/nginx/sites-available/mist
location /api/github/webhook {
    # GitHub webhook IP ranges (update periodically from https://api.github.com/meta)
    allow 140.82.112.0/20;
    allow 143.55.64.0/20;
    allow 185.199.108.0/22;
    allow 192.30.252.0/22;
    deny all;
    
    proxy_pass http://localhost:8080;
}
```

**Option 2: Firewall Rules**

Use UFW or iptables to allow only GitHub IP ranges for the webhook endpoint.

## Advanced Configuration

### Multiple GitHub Apps

You can only have **one GitHub App per Mist instance** currently. If you need multiple:

- Use different Mist instances for different organizations
- Or install the same GitHub App on multiple organizations

### GitHub Enterprise Server

Mist uses GitHub.com APIs. For GitHub Enterprise Server support:
- Currently not supported
- Coming soon: Configurable GitHub API base URL

## Further Reading

- [GitHub Apps Documentation](https://docs.github.com/en/apps)
- [GitHub Webhooks](https://docs.github.com/en/webhooks)
- [GitHub API Meta Endpoint](https://docs.github.com/en/rest/meta) - Get webhook IP ranges
<!-- - [Mist Git Integration Guide](../guide/git) -->

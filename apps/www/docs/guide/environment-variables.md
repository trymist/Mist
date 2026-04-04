# Environment Variables

Environment variables allow you to configure your applications without hardcoding values in your code.

## Overview

Environment variables in Mist are:

- **Secure**: Values are masked in the UI
- **Flexible**: Support both build-time and runtime variables
- **Easy to Manage**: Single and bulk import modes

## Adding Environment Variables

### Single Variable

1. Navigate to your application
2. Go to **Environment** tab
3. Click **"Add Variable"**
4. Select **"Single"** mode
5. Enter **Key** and **Value**
6. Click **"Add"**

### Bulk Import <span class="status-badge status-implemented">New</span>

Add multiple variables at once by pasting them in `KEY=VALUE` format:

1. Click **"Add Variable"**
2. Select **"Bulk Paste"** mode
3. Paste your environment variables:

```bash
NODE_ENV=production
DATABASE_URL=postgresql://user:pass@host:5432/db
API_KEY=your-api-key-here
REDIS_URL=redis://localhost:6379
PORT=3000
```

4. Preview detected variables
5. Click **"Add X Variables"**

#### Supported Formats

The parser supports:

```bash
# Standard format
KEY=value

# With quotes
KEY="value with spaces"
KEY='single quoted value'

# Comments (ignored)
# This is a comment
DATABASE_URL=postgres://localhost

# Empty lines (ignored)

NEXT_PUBLIC_API=https://api.example.com
```

## Variable Types

### Runtime Variables

Available when your container runs:

```javascript
// Node.js
const apiKey = process.env.API_KEY;

// Python
import os
api_key = os.environ.get('API_KEY')
```

### Build-time Variables

Available during the build process:

```dockerfile
ARG NODE_ENV
RUN npm run build
```

## Editing Variables

1. Find the variable in the list
2. Click the **Edit** icon (pencil)
3. Modify key or value
4. Click **"Save"**

## Deleting Variables

1. Find the variable in the list
2. Click the **Delete** icon (trash)
3. Confirm deletion

::: warning
Changes to environment variables require redeployment to take effect.
:::

## Common Use Cases

### Database Connection

```bash
DATABASE_URL=postgresql://user:password@host:5432/dbname
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=appuser
DB_PASSWORD=secret123
```

### API Keys

```bash
STRIPE_API_KEY=sk_test_xxxxx
SENDGRID_API_KEY=SG.xxxxx
AWS_ACCESS_KEY_ID=AKIAXXXXX
AWS_SECRET_ACCESS_KEY=xxxxx
```

### Application Configuration

```bash
NODE_ENV=production
PORT=3000
LOG_LEVEL=info
SESSION_SECRET=random-secret-string
```

### Feature Flags

```bash
ENABLE_ANALYTICS=true
ENABLE_BETA_FEATURES=false
MAINTENANCE_MODE=false
```

## Best Practices

### Security

- ✅ Never commit secrets to Git
- ✅ Use environment variables for all sensitive data
- ✅ Rotate API keys regularly
- ✅ Use different values for staging and production

### Naming Conventions

```bash
# Use UPPER_SNAKE_CASE
DATABASE_URL=...

# Prefix public variables (Next.js, Vite, etc.)
NEXT_PUBLIC_API_URL=...
VITE_API_URL=...

# Group related variables
AWS_REGION=us-east-1
AWS_BUCKET=my-bucket
AWS_ACCESS_KEY_ID=...
```

### Organization

```bash
# Database
DATABASE_URL=...
DATABASE_POOL_SIZE=10

# Redis
REDIS_URL=...
REDIS_TTL=3600

# External Services
STRIPE_KEY=...
SENDGRID_KEY=...

# App Config
NODE_ENV=production
PORT=3000
```

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Planned Features</h4>
  <p>Future enhancements for environment variables:</p>
</div>

- **Encryption at Rest** - Encrypted storage for sensitive values
- **Variable Groups** - Organize variables into reusable groups
- **Templates** - Pre-defined variable sets for common stacks
- **Version History** - Track changes to variables over time
- **Import/Export** - Download as `.env` file
- **Secrets Management** - Integration with Vault or similar
- **Variable Validation** - Ensure required variables are set

## Troubleshooting

### Variables Not Available

**Problem**: Environment variables not accessible in application

**Solutions**:
- Redeploy the application after adding variables
- Check variable names (case-sensitive)
- Verify the variable is not overridden elsewhere

### Build Failing

**Problem**: Build fails after adding variables

**Solutions**:
- Check for special characters in values
- Use quotes for values with spaces
- Ensure variables don't conflict with system variables

### Parsing Errors (Bulk Import)

**Problem**: Some variables not detected

**Solutions**:
- Ensure format is `KEY=VALUE`
- Keys must start with letter or underscore
- Keys can only contain letters, numbers, and underscores
- Remove invalid characters

## Examples by Framework

### Node.js / Express

```bash
NODE_ENV=production
PORT=3000
DATABASE_URL=postgresql://localhost/myapp
SESSION_SECRET=your-secret-key
```

### Next.js

```bash
# Server-side
DATABASE_URL=postgresql://localhost/myapp

# Client-side (must be prefixed)
NEXT_PUBLIC_API_URL=https://api.example.com
```

### Python / Django

```bash
DJANGO_SECRET_KEY=your-secret-key
DEBUG=False
DATABASE_URL=postgresql://localhost/myapp
ALLOWED_HOSTS=example.com,www.example.com
```

### Go

```bash
PORT=8080
DATABASE_URL=postgresql://localhost/myapp
JWT_SECRET=your-jwt-secret
```

## API Reference

Manage environment variables programmatically via the API:

- [Create Variable](./environment-variables#create)
- [List Variables](./environment-variables#list)
- [Update Variable](./environment-variables#update)
- [Delete Variable](./environment-variables#delete)

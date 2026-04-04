# Database Services

One-click database provisioning for your applications.

Mist provides managed database services through pre-configured Docker templates, allowing you to quickly deploy databases without manual configuration.

## Available Databases

### PostgreSQL
- Latest PostgreSQL version
- Pre-configured with recommended settings
- Automatic environment variable setup
- Internal network access via `app-{appId}` hostname
- Template-driven CPU and memory limits

### MySQL
- Latest MySQL version
- Optimized configuration
- User and password management via environment variables
- Internal network connectivity

### MariaDB
- MySQL-compatible database
- Drop-in replacement for MySQL
- Same connection interface

### Redis
- In-memory data store and cache
- Latest Redis version
- Password protection via environment variables
- Pub/sub and caching support

### MongoDB
- Document-oriented database
- Latest MongoDB version
- User authentication support
- Internal DNS resolution

## Creating a Database

1. Navigate to your project
2. Click **"New Application"**
3. Select **"Database"** type
4. Choose your database template:
   - PostgreSQL
   - MySQL
   - MariaDB
   - Redis
   - MongoDB
5. Enter instance name and description
6. Configure environment variables (passwords, database names)
7. Click **"Create"**

The database will be provisioned automatically with:
- Pre-configured port from template
- Recommended CPU and memory limits
- Default environment variables
- Connected to `traefik-net` network

## Connecting to Databases

### From Your Applications

Databases are accessible via internal hostname using the format `app-{appId}`:

**Example Connection Strings:**

```bash
# PostgreSQL
DATABASE_URL=postgresql://user:password@app-123:5432/myapp

# MySQL
DATABASE_URL=mysql://root:password@app-124:3306/myapp

# Redis
REDIS_URL=redis://:password@app-125:6379

# MongoDB
MONGODB_URI=mongodb://admin:password@app-126:27017/myapp
```

Replace `app-{appId}` with your actual database application ID shown in the dashboard.

### Configuration via Environment Variables

Add connection details to your web/service application's environment variables:

1. Go to your application's **Environment** tab
2. Add database connection variables
3. Use the database's container name (`app-{appId}`)
4. Deploy your application

## Database Management

### Environment Variables

Configure your database through environment variables:

**PostgreSQL:**
```
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=mydatabase
```

**MySQL:**
```
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=mydatabase
MYSQL_USER=myuser
MYSQL_PASSWORD=mypassword
```

**MongoDB:**
```
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password
```

**Redis:**
```
REDIS_PASSWORD=yourpassword
```

### Resource Limits

Template-based defaults are applied:
- **CPU**: 1.0 cores (adjustable)
- **Memory**: 512MB (adjustable)
- **Restart Policy**: unless-stopped

You can modify these in the application settings.

### Container Management

- **Start/Stop**: Control database lifecycle from dashboard
- **Restart**: Restart database containers
- **Logs**: View database logs in real-time
- **Status**: Monitor container state and uptime

## Networking

All databases are connected to the `traefik-net` network:
- Internal DNS resolution
- Accessible by container name (`app-{appId}`)
- Isolated from public internet
- Secure inter-container communication

## Best Practices

### Security

- **Strong Passwords**: Use complex passwords for database users
- **Environment Variables**: Never hardcode credentials
- **Least Privilege**: Create database users with minimal required permissions
- **Network Isolation**: Databases are not exposed externally by default

### Performance

- **Resource Allocation**: Monitor usage and adjust CPU/memory limits
- **Connection Pooling**: Use connection pools in your applications
- **Indexing**: Create appropriate indexes for query performance
- **Monitoring**: Check logs regularly for errors or warnings

### Backup

::: warning Manual Backups Required
Automatic backup scheduling is not yet available. Manually backup your databases regularly using database-specific tools.
:::

**Backup Methods:**

**PostgreSQL:**
```bash
docker exec app-{appId} pg_dump -U user dbname > backup.sql
```

**MySQL:**
```bash
docker exec app-{appId} mysqldump -u root -p dbname > backup.sql
```

**MongoDB:**
```bash
docker exec app-{appId} mongodump --out /backup
```

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Upcoming Features</h4>
</div>

- **External Access** - Public database access with security controls
- **Management UIs** - Built-in admin interfaces (pgAdmin, phpMyAdmin, Mongo Express, Redis Commander)
- **Automated Backups** - Scheduled backup and restore
- **Connection Pooling** - Optimize database connections
- **Replication** - High availability setups
- **Multiple Versions** - Choose specific database versions
- **Volume Management** - Persistent volume configuration
- **Database Cloning** - Duplicate databases for staging/testing
- **Performance Metrics** - Database-specific monitoring
- **Auto-Configuration** - Automatic environment variable injection to apps

## Troubleshooting

### Connection Refused

- Verify database container is running
- Check container logs for errors
- Ensure correct hostname format (`app-{appId}`)
- Verify port matches database default

### Authentication Failed

- Check environment variables are set correctly
- Verify username and password
- Ensure database user has been created
- Check database logs for auth errors

### Application Can't Connect

- Verify both containers are on `traefik-net` network
- Check firewall rules aren't blocking internal traffic
- Test connection from application container logs
- Verify DNS resolution of container name

## Related Documentation

- [Applications](./applications) - Application types and management
- [Environment Variables](./environment-variables) - Managing configuration
- [Service Templates](./applications#service-templates) - How templates work

# CLI Tool

Command-line interface for managing Mist users and system settings.

The Mist CLI (`mist-cli`) is a powerful command-line tool that allows administrators to manage their Mist installation directly from the terminal. It's automatically installed when you install Mist using the standard installation script.

## Installation

The CLI tool is automatically installed during Mist installation and is available at `/usr/local/bin/mist-cli`.

To verify installation:

```bash
mist-cli version
```

## Available Commands

## Available Commands

The CLI supports the following main commands:

- `user` - Manage users
- `settings` - Manage system settings
- `version` - Show CLI version
- `help` - Show help message

Use `mist-cli <command> --help` for more information about each command.

---

## User Management

Manage user accounts including password changes and listing users.

### Change User Password

Change a user's password with optional interactive or non-interactive mode.

**Interactive mode (recommended):**
```bash
mist-cli user change-password --username admin
```

This will prompt you to enter and confirm the new password securely.

**Non-interactive mode:**
```bash
mist-cli user change-password --username admin --password newpass123
```

**Examples:**
```bash
# Change admin password (interactive)
mist-cli user change-password --username admin

# Change user password directly
mist-cli user change-password --username john --password SecurePass123
```

### List Users

Display all users in the system with their details.

```bash
mist-cli user list
```

**Output:**
```
Users:
----------------------------------------------
ID    Username             Email                          Role      
----------------------------------------------
1     admin                admin@example.com              owner     
2     developer            dev@example.com                member    
----------------------------------------------
Total: 2 users
```

---

## System Settings Management

Manage system-wide settings that control Mist behavior.

### List All Settings

View all current system settings:

```bash
mist-cli settings list
```

**Output:**
```
System Settings:
----------------------------------------------
Setting                        Value
----------------------------------------------
wildcard_domain                example.com
mist_app_name                  mist
production_mode                true
secure_cookies                 true
auto_cleanup_containers        false
auto_cleanup_images            false
----------------------------------------------
```

### Get Specific Setting

Retrieve the value of a specific setting:

```bash
mist-cli settings get --key wildcard_domain
```

**Available setting keys:**
- `wildcard_domain` - Wildcard domain for auto-generated app domains
- `mist_app_name` - Subdomain name for Mist dashboard
- `production_mode` - Enable production mode (true/false)
- `secure_cookies` - Enable secure cookies for HTTPS (true/false)
- `auto_cleanup_containers` - Auto cleanup stopped containers (true/false)
- `auto_cleanup_images` - Auto cleanup dangling images (true/false)

### Set Setting Value

Update a system setting:

```bash
mist-cli settings set --key <key> --value <value>
```

**Examples:**

```bash
# Set wildcard domain
mist-cli settings set --key wildcard_domain --value apps.example.com

# Enable production mode
mist-cli settings set --key production_mode --value true

# Enable secure cookies
mist-cli settings set --key secure_cookies --value true

# Enable auto cleanup
mist-cli settings set --key auto_cleanup_containers --value true
mist-cli settings set --key auto_cleanup_images --value true

# Change Mist dashboard subdomain
mist-cli settings set --key mist_app_name --value dashboard
```

---

## Common Use Cases

### Initial Setup

After installing Mist, set up your wildcard domain:

```bash
# Configure wildcard domain
mist-cli settings set --key wildcard_domain --value apps.example.com

# Set Mist dashboard subdomain
mist-cli settings set --key mist_app_name --value mist

# Enable production mode
mist-cli settings set --key production_mode --value true
```

### Password Recovery

If you've forgotten the admin password:

```bash
# Reset admin password
sudo mist-cli user change-password --username admin
```

### Enable Auto Cleanup

To automatically clean up Docker resources:

```bash
# Enable container cleanup
mist-cli settings set --key auto_cleanup_containers --value true

# Enable image cleanup
mist-cli settings set --key auto_cleanup_images --value true
```

### Check Current Configuration

View all settings and users:

```bash
# List all settings
mist-cli settings list

# List all users
mist-cli user list
```

---

## Permissions

The CLI requires direct access to the Mist database at `/var/lib/mist/mist.db`. You may need to run commands with `sudo`:

```bash
sudo mist-cli user change-password --username admin
```

---

## Troubleshooting

### Database Not Found

**Error:**
```
Error: database file not found at /var/lib/mist/mist.db. Please ensure Mist is installed and running
```

**Solution:** Ensure Mist is properly installed and the service is running:
```bash
sudo systemctl status mist
```

### Permission Denied

**Error:** Permission issues when accessing the database

**Solution:** Run the command with sudo:
```bash
sudo mist-cli settings list
```

### Unknown Setting Key

**Error:**
```
Error: Unknown setting key 'invalid_key'
```

**Solution:** Use `mist-cli settings --help` to see available setting keys.

---

## Help & Support

For more information on any command, use the `--help` flag:

```bash
mist-cli --help
mist-cli user --help
mist-cli settings --help
```

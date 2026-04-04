# Forgot Password Recovery

Step-by-step guide to recover access when you've forgotten your password.

---

## Overview

If you've forgotten your Mist password, you can reset it using the `mist-cli` command-line tool. This tool provides direct database access to change user passwords without requiring dashboard access.

::: warning Server Access Required
You need SSH or direct access to the server where Mist is installed to use this method.
:::

---

## Quick Recovery Steps

### 1. SSH into Your Server

Connect to your server where Mist is installed:

```bash
ssh user@your-server-ip
```

### 2. Reset Password Using CLI

Run the password reset command:

```bash
sudo mist-cli user change-password --username YOUR_USERNAME
```

You'll be prompted to enter and confirm your new password:

```
Enter new password: 
Confirm new password: 
✓ Password changed successfully for user 'YOUR_USERNAME'
```

### 3. Log Back In

Navigate to your Mist dashboard and log in with your new password:

```
https://your-mist-domain.com
```

---

## Detailed Instructions

### For Admin Users

If you're the admin and forgot the admin password:

```bash
# SSH to server
ssh user@your-server-ip

# Reset admin password
sudo mist-cli user change-password --username admin

# Enter and confirm new password when prompted
```

**Example session:**
```bash
$ sudo mist-cli user change-password --username admin
Enter new password: ****************
Confirm new password: ****************
✓ Password changed successfully for user 'admin'
```

### For Other Users

If you're helping another user who forgot their password:

```bash
# First, list all users to verify the username
sudo mist-cli user list

# Output:
# Users:
# ----------------------------------------------
# ID    Username             Email                          Role      
# ----------------------------------------------
# 1     admin                admin@example.com              owner     
# 2     developer            dev@example.com                member    
# ----------------------------------------------

# Reset the specific user's password
sudo mist-cli user change-password --username developer
```

---

## Non-Interactive Mode

If you need to reset a password in a script or automation, you can provide the password directly:

```bash
sudo mist-cli user change-password --username admin --password NewSecurePass123
```

::: warning Security Notice
Using the `--password` flag will expose the password in:
- Shell history
- Process list (visible via `ps` command)
- Log files

Use this method only when necessary and clear your shell history afterward:
```bash
history -c
```
:::

---

## Troubleshooting

### CLI Command Not Found

**Error:**
```bash
mist-cli: command not found
```

**Solution:**

Check if the CLI is installed:
```bash
ls -la /usr/local/bin/mist-cli
```

If not found, rebuild and install the CLI:
```bash
cd /opt/mist/cli
sudo go build -o mist-cli
sudo cp mist-cli /usr/local/bin/
sudo chmod +x /usr/local/bin/mist-cli

# Verify installation
mist-cli version
```

### Permission Denied

**Error:**
```bash
Error: open /var/lib/mist/mist.db: permission denied
```

**Solution:**

Run the command with `sudo`:
```bash
sudo mist-cli user change-password --username admin
```

### Database Not Found

**Error:**
```bash
Error: database file not found at /var/lib/mist/mist.db
```

**Solution:**

Verify Mist is installed and running:
```bash
# Check service status
sudo systemctl status mist

# Check if database file exists
ls -la /var/lib/mist/mist.db

# If service is not running, start it
sudo systemctl start mist
```

### User Not Found

**Error:**
```bash
Error: User 'johndoe' not found
```

**Solution:**

List all users to find the correct username:
```bash
sudo mist-cli user list
```

Usernames are case-sensitive, so ensure you're using the exact username as shown in the list.

### Passwords Don't Match

**Error:**
```bash
Error: Passwords do not match
```

**Solution:**

When entering passwords in interactive mode, make sure you type the same password twice. If you keep getting this error:

1. Copy your desired password to clipboard
2. Paste it when prompted for "Enter new password"
3. Paste it again when prompted for "Confirm new password"

Or use non-interactive mode:
```bash
sudo mist-cli user change-password --username admin --password YourNewPassword123
```

---

## Security Best Practices

### After Password Reset

1. **Clear shell history** if you used non-interactive mode:
```bash
history -c
```

2. **Use a strong password:**
   - At least 12 characters
   - Mix of uppercase, lowercase, numbers, and symbols
   - Avoid common words or patterns

3. **Consider using a password manager** like:
   - Bitwarden
   - 1Password
   - KeePassXC

### Prevent Future Lockouts

1. **Store passwords securely:**
   - Use a password manager
   - Keep encrypted backup of credentials
   - Document recovery procedures

2. **Set up multiple admin accounts:**
```bash
# Create a backup admin account (via dashboard)
# Or add an additional user with owner role
```

3. **Keep SSH access secure:**
   - Use SSH keys instead of passwords
   - Document server access procedures
   - Maintain backup access methods

---

## Alternative Recovery Methods

### If You Don't Have Server Access

If you cannot access the server directly, you'll need to:

1. **Contact your hosting provider** or system administrator
2. **Access the server console** through your cloud provider's dashboard (AWS, DigitalOcean, etc.)
3. **Use server recovery mode** if available

### If CLI is Not Available

If the CLI tool is not installed or not working, you can reset the password directly in the database:

::: danger Advanced Users Only
Direct database manipulation can corrupt your data. Use the CLI method when possible.
:::

```bash
# Connect to database
sqlite3 /var/lib/mist/mist.db

# Generate password hash (use a temporary Go script or online bcrypt generator)
# Then update the database
UPDATE users SET password_hash = 'YOUR_BCRYPT_HASH' WHERE username = 'admin';
.quit
```

---

## Testing the New Password

After resetting your password:

1. **Open your browser** and navigate to Mist dashboard
2. **Clear browser cookies** (recommended):
   - Chrome: Settings → Privacy → Clear browsing data → Cookies
   - Firefox: Settings → Privacy → Clear Data → Cookies
3. **Log in** with your username and new password
4. **Verify access** to all dashboard features

---

## Common Scenarios

### Scenario 1: Locked Out as Only Admin

```bash
# SSH to server
ssh user@server

# Reset your admin password
sudo mist-cli user change-password --username admin

# Log back into dashboard
# https://your-mist-domain.com
```

### Scenario 2: Reset Password for Another User

```bash
# List users first
sudo mist-cli user list

# Reset specific user's password
sudo mist-cli user change-password --username developer

# Notify the user of their password reset
```

### Scenario 3: Bulk Password Reset

```bash
# Reset multiple users (requires bash loop)
for user in developer tester; do
  sudo mist-cli user change-password --username $user --password TempPassword123
  echo "Reset password for $user"
done

# Notify users to change their temporary passwords
```

---

## Getting Help

If you're still having trouble recovering access:

### Check Documentation
- [CLI Tool Guide](../guide/cli) - Complete CLI documentation
- [User Management](../guide/users) - User roles and permissions
- [Authentication](../guide/authentication) - Login and security

### Community Support
- **GitHub Issues:** [github.com/trymist/mist/issues](https://github.com/trymist/mist/issues)
- **Discord:** [discord.gg/kxK8XHR6](https://discord.gg/kxK8XHR6)

### Provide This Information
When asking for help, include:

```bash
# Mist version
curl http://localhost:8080/api/updates/version

# Service status
sudo systemctl status mist

# CLI version
mist-cli version

# User list (remove sensitive info)
sudo mist-cli user list
```

---

## Related Documentation

- [CLI Tool Guide](../guide/cli) - Full CLI documentation
- [User Management](../guide/users) - Managing users and roles
- [Common Issues](./) - Other troubleshooting guides

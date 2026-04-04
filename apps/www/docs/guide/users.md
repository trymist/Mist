# Users & Roles

Manage users and their permissions in Mist.

## System Roles

Mist has a hierarchical role-based access control system with three levels:

### Owner

**Highest privilege level** - Complete system control

- Only ONE owner exists per Mist instance
- Created automatically during first-time setup
- Cannot be duplicated through UI or API
- Full access to all features and settings
- Can view and update system settings
- Can create and delete users (including admins)
- Can view all audit logs
- Can delete themselves (resets system to setup mode)

### Admin

**High-level administrative privileges** - Second to owner

- Can create and manage projects
- Can create users (admin and user roles)
- Can delete users (except owners and other admins)
- Can view all audit logs
- Access to assigned projects only
- **Cannot** view or modify system settings
- **Cannot** delete owners or other admins
- **Cannot** create owner accounts

### User

**Standard role** - Basic project access

- Can access assigned projects
- Can deploy applications
- Can manage applications in assigned projects
- Can view logs and metrics
- **Cannot** create projects
- **Cannot** create other users
- **Cannot** view system settings
- **Cannot** view audit logs

## Role Comparison

| Permission | Owner | Admin | User |
|-----------|-------|-------|------|
| View system settings | ✅ | ❌ | ❌ |
| Update system settings | ✅ | ❌ | ❌ |
| Create projects | ✅ | ✅ | ❌ |
| Create users | ✅ | ✅ | ❌ |
| Delete owners | ✅ | ❌ | ❌ |
| Delete admins | ✅ | ❌ | ❌ |
| Delete users | ✅ | ✅ | ❌ |
| View audit logs | ✅ | ✅ | ❌ |
| Access assigned projects | ✅ | ✅ | ✅ |
| Manage applications | ✅ | ✅ | ✅ |

## Managing Users

### Creating Users

**Required role**: Owner or Admin

1. Navigate to **Settings** → **Users**
2. Click **"Add User"** or **"New User"**
3. Fill in user information:
   - Username
   - Email address
   - Password
4. Select role:
   - **Admin**: For administrative users
   - **User**: For standard users
5. Click **"Create"**

::: tip
You cannot create additional owner accounts. Only one owner exists per Mist instance.
:::

### Viewing Users

**Required role**: Owner or Admin

View all users in the system:
- Username
- Email
- Role (displayed with color-coded badges)
- Account creation date

**Role badge colors:**
- Owner: Purple
- Admin: Blue
- User: Gray

### Updating Users

**Required role**: Owner or Admin

Modify existing user details:
1. Click on the user you want to edit
2. Update information:
   - Username
   - Email
   - Password
   - Role (owner/admin only)
3. Click **"Save"**

All changes are logged in the audit system.

### Deleting Users

**Deletion permissions:**
- **Owner**: Can delete any user (admins and users)
- **Admin**: Can only delete users (not owners or other admins)
- **User**: Cannot delete anyone

To delete a user:
1. Navigate to the Users page
2. Click the delete icon next to the user
3. Confirm the deletion

::: warning Deletion Restrictions
- Admins cannot delete owners or other admins
- Projects and applications created by the deleted user are preserved
- If the owner deletes themselves, the system resets to setup mode
:::

## Project Roles

::: info Coming Soon
Project-level roles (Owner, Admin, Member) are not yet implemented. Currently, all project members have equal access to project resources.

See [Projects documentation](./projects) for more details on project membership.
:::

## Best Practices

### Role Assignment

- **Owner**: Reserve for the primary system administrator
- **Admin**: For trusted team members who need to manage users and projects
- **User**: For team members who only need to work on specific projects

### Security Recommendations

1. **Limit Admins**: Only assign admin role to users who need administrative privileges
2. **Use Users**: Most team members should have the user role
3. **Protect Owner**: The owner account should have a strong password and be well-protected
4. **Regular Audits**: Review user list regularly and remove inactive accounts
5. **Audit Logs**: Monitor audit logs for suspicious activity

### Account Management

- Set strong password requirements for all users
- Remove users who no longer need access
- Review role assignments periodically
- Document who has owner and admin access

## Related Documentation

- [Authentication](./authentication) - Detailed authentication information
- [Audit Logs](./audit-logs) - Track user actions
- [Projects](./projects) - Project membership management

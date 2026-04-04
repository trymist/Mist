# Authentication & Users

Manage users and authentication in your Mist instance.

## User Roles

Mist implements a hierarchical role-based access control (RBAC) system with three distinct roles:

### Owner

The highest privilege level in Mist with complete system control.

**Characteristics:**
- **Unique**: Only ONE owner can exist in the system
- **Created**: Automatically assigned during first-time setup
- **Cannot be duplicated**: No API or UI allows creating additional owners

**Permissions:**
- ✅ Full system access to all features
- ✅ View and update system settings (wildcard domain, etc.)
- ✅ Create and delete projects
- ✅ Create users (admin and user roles only)
- ✅ Delete any user including admins
- ✅ View all audit logs
- ✅ Access all projects and applications
- ✅ Manage Traefik configuration

**Special Considerations:**
- Owner can delete themselves, which resets the system to setup mode
- If the owner is deleted, a new owner can be created via signup

### Admin

High-level administrative privileges, second to owner.

**Permissions:**
- ✅ Create and manage projects
- ✅ Create users (admin and user roles only)
- ✅ Delete users (except owners and other admins)
- ✅ View all audit logs
- ✅ Access assigned projects as member
- ✅ Manage applications in assigned projects
- ❌ Cannot view or update system settings
- ❌ Cannot delete owners
- ❌ Cannot delete other admins
- ❌ Cannot create owner accounts

### User

Standard user role with limited permissions.

**Permissions:**
- ✅ Access assigned projects as member
- ✅ Deploy applications in assigned projects
- ✅ View logs and metrics for assigned projects
- ✅ Manage application settings in assigned projects
- ✅ View own profile
- ❌ Cannot create projects
- ❌ Cannot create users
- ❌ Cannot view system settings
- ❌ Cannot view audit logs

## Authentication

### Login

Access Mist at `http://your-server-ip:8080`:

1. Enter email and password
2. Click **"Sign In"**
3. JWT token stored in HTTP-only cookie

### Session Management

- Sessions expire after 31 days
- Tokens stored in secure HTTP-only cookies
- Role embedded in JWT for fast authorization
- Role verified against database on each request

### First-Time Setup

On first access, create the owner account:

1. Navigate to `http://your-server-ip:8080`
2. You'll see the setup page
3. Enter email and password
4. Click **"Create Admin Account"**
5. You'll be automatically logged in as **owner**

::: tip
The first user is always assigned the owner role. After the first user is created, the signup page is disabled.
:::

## User Management

### Creating Users

**Who can create users**: Owner and Admin roles only

1. Go to **Settings** → **Users**
2. Click **"Add User"** or **"New User"**
3. Enter user details:
   - Username
   - Email
   - Password
4. Select role:
   - **Admin**: Full administrative privileges (cannot create owners)
   - **User**: Standard user access
5. Click **"Create"**

::: warning Cannot Create Owners
The owner role cannot be assigned when creating users. Only one owner exists per Mist instance.
:::

### Updating Users

**Who can update users**: Owner and Admin

Edit user details:
- Update username
- Update email
- Reset password
- Change role (owner/admin only)

All user updates are logged in the audit log.

### Deleting Users

**Who can delete users**:
- **Owner**: Can delete any user including admins
- **Admin**: Can only delete users (not owners or other admins)
- **User**: Cannot delete any users

To delete a user:
1. Navigate to **Settings** → **Users**
2. Click delete icon next to the user
3. Confirm deletion

::: warning
- Admins cannot delete owners or other admins
- Deleting a user removes their access but preserves projects and applications they created
- If the owner deletes themselves, the system resets to setup mode
:::

## Role Hierarchy

The role hierarchy is strictly enforced:

```
Owner (highest)
  ↓
Admin
  ↓
User (lowest)
```

**Deletion Rules:**
- Owner can delete: Admin, User
- Admin can delete: User only
- User cannot delete anyone

## Security Features

### Password Security

- Minimum 8 characters required
- Hashed with bcrypt algorithm
- Salted for additional security
- Never stored in plain text

### JWT Tokens

- Signed with HS256 algorithm
- Includes user ID, email, and role
- 31-day expiration
- Stored in HTTP-only cookies (prevents XSS attacks)

### Session Security

- HTTP-only cookies prevent JavaScript access
- Secure flag enabled in production (HTTPS)
- Role verified on every API request
- Invalid tokens result in automatic logout

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Authentication Enhancements</h4>
</div>

- **Password Reset** - Email-based password recovery
- **Email Verification** - Verify email addresses
- **Two-Factor Authentication** - TOTP/Authenticator app support
- **SSO Integration** - OAuth2, SAML support
- **API Tokens** - Generate tokens for CLI/API access
- **Session Management UI** - View and revoke active sessions
- **User Invitations** - Invite users via email
- **Fine-grained Permissions** - Custom permission sets

# Projects

Projects help you organize and group related applications together.

## Overview

A project is a container for multiple applications, allowing you to:

- Group related applications (e.g., frontend, backend, worker)
- Manage team access and permissions
- Organize by client, environment, or purpose
- Track resources and deployments collectively

## Creating a Project

1. Click **"New Project"** from the dashboard
2. Enter project details:
   - **Name**: Descriptive project name
   - **Description**: Optional description
   - **Tags**: Comma-separated tags for filtering

3. Click **"Create"**

## Project Members

Add team members to collaborate on projects:

1. Open your project
2. Go to **Members** tab
3. Click **"Add Member"**
4. Select user(s) from list
5. Click **"Save"**

### Member Access

All project members have equal access to:
- View all applications in the project
- Deploy applications
- View logs and metrics
- Manage application settings

::: tip Coming Soon
Role-based permissions (Owner, Admin, Member) and ownership transfer are upcoming features. Currently, all members have full access to the project.
:::

## Managing Applications

All applications within a project are listed on the project page:

- Create new applications
- View deployment status
- Access application details
- Monitor resource usage

## Project Settings

### General

- Update project name and description
- Modify tags for organization
- View project statistics
- Manage project members

### Danger Zone

- Delete project (requires confirmation)
  - Only the project owner can delete the project
  - Deletes all applications and their data

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Upcoming Features</h4>
</div>

- **Role-Based Permissions** - Owner, Admin, and Member roles with different access levels
- **Project Ownership Transfer** - Transfer project ownership to another user
- **Project-level Environment Variables** - Share variables across apps
- **Project Templates** - Quick-start templates for common stacks
- **Resource Quotas** - Limit CPU, memory, and storage per project
- **Billing Integration** - Track costs per project
- **Project Archives** - Temporarily disable all apps

# Audit Logs

Track all user actions in your Mist instance.

## Overview

Audit logs record:
- User authentication events
- Application deployments
- Configuration changes
- User management actions
- Project modifications

## Viewing Audit Logs

<div class="coming-soon-banner">
  <h4>🚧 Audit Log UI Coming Soon</h4>
  <p>Audit logs are currently stored in the database but UI for viewing is under development.</p>
</div>

## What's Logged

### Authentication
- User login attempts
- Session creation
- Logout events

### Applications
- Application creation
- Application updates
- Application deletion
- Deployment triggers

### Environment Variables
- Variable creation
- Variable updates
- Variable deletion

### Projects
- Project creation
- Member addition/removal
- Project settings changes

### Domains
- Domain addition
- Domain removal

### Users (Admin)
- User creation
- Role changes
- User deletion

## Log Structure

Each audit log entry includes:
- **Timestamp** - When action occurred
- **User** - Who performed the action
- **Action** - What was done
- **Resource** - What was affected
- **Details** - Additional context
- **IP Address** - Source IP

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Audit Log Features</h4>
</div>

- **Audit Log Viewer UI** - Browse logs in dashboard
- **Search & Filter** - Find specific events
- **Export Logs** - Download as CSV/JSON
- **Retention Policies** - Configure log retention
- **Compliance Reports** - Generate audit reports
- **Real-time Alerts** - Notify on suspicious activity

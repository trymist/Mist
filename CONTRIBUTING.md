# Contributing to Mist

Thank you for your interest in contributing to Mist! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Project Structure](#project-structure)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Code of Conduct

This project and everyone participating in it is governed by our commitment to:

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Accept constructive criticism gracefully
- Focus on what's best for the community and users
- Show empathy towards others

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Linux** - any linux distro or wsl
- **Go 1.22+** - [Download](https://go.dev/dl/)
- **Docker** - [Installation Guide](https://docs.docker.com/get-docker/)
- **Bun** (for frontend development) - [Install](https://bun.sh/)
- **Git** - [Download](https://git-scm.com/downloads)

Optional but recommended:
- **fyrer** (for hot reloading during development):
  ```bash
  cargo install fyrer
  ```

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/trymist/mist
   cd mist
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/trymist/mist
   ```

## Development Environment

### Using Fyrer (Recommended)

The easiest way to start development:

```bash
docker compose -f "traefik-compose.yml" up -d
fyrer
```

This will automatically:
- Watch for file changes
- Rebuild the server when Go files change
- Rebuild the frontend when TypeScript/React files change
- Handle dependencies

### Manual Setup

If you prefer manual control:

**Backend (Go):**
```bash
cd server
go mod tidy 
go run .
```

**Frontend (React/TypeScript):**
```bash
cd dash
bun install
bun run dev
```

**Traefik**
```bash
docker compose -f "traefik-compose.yml" up -d
```

**CLI:**
```bash
cd cli
go mod download
go build .
```

## Project Structure

```
mist/
├── server/           
│   ├── api/          # API server 
│   ├── compose/      # Docker Compose operations
│   ├── db/           # DB configurations/migrations
│   ├── git/          # Git helper functions
│   ├── models/       # Database models (GORM)
│   ├── queue/        # Deployment queue 
│   ├── store/        # Database store
│   ├── utils/        # Utility functions
│   └── websockets/   # WebSocket handlers
├── dash/             # React frontend (Vite + TypeScript)
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── features/    # Feature modules
│   │   ├── hooks/      # React hooks
│   │   ├── pages/      # Page components
│   │   ├── services/   # API services
│   │   └── types/      # TypeScript types
│   └── package.json
├── cli/              # Go CLI tool
├── test/             # Go test suite
├── www/docs/         # VitePress documentation
└── traefik-compose.yml  # Traefik configuration
```

### Technology Stack

- **Backend:** Go 1.22+, Gin framework, GORM, SQLite, Docker API
- **Frontend:** React 19, TypeScript, Vite, Tailwind CSS, Radix UI
- **CLI:** Go with cobra
- **Documentation:** VitePress
- **Database:** SQLite with GORM

## How to Contribute

### Reporting Bugs

Before creating a bug report:

1. Check the [existing issues](https://github.com/trymist/mist/issues) to avoid duplicates
2. Update to the latest version to see if the issue persists

When reporting bugs, include:
- Clear, descriptive title
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Docker version, Mist version)
- Relevant logs or screenshots

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating one:

- Use a clear, descriptive title
- Provide detailed description
- Explain why this enhancement would be useful
- Include mockups or examples if applicable

### Contributing Code

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   
   Branch naming conventions:
   - `feature/` - New features
   - `fix/` - Bug fixes
   - `docs/` - Documentation updates
   - `refactor/` - Code refactoring
   - `test/` - Test additions/improvements

2. **Make your changes** following our coding standards

3. **Test your changes** (see [Testing](#testing))

4. **Update documentation** if needed

5. **Commit** with a clear message (see [Commit Messages](#commit-messages))

6. **Push** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** on GitHub

## Coding Standards

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `go vet` to catch common issues
- Keep functions focused and small
- Write meaningful comments for exported functions
- Use meaningful variable names

Example:
```go
// DeployApplication deploys an application from a git repository
func DeployApplication(app *models.Application, deployment *models.Deployment) error {
    // Implementation
}
```

### TypeScript/React Code

- Use TypeScript for all new code
- Follow the existing component structure
- Use functional components with hooks
- Prefer named exports over default exports
- Use Tailwind CSS for styling
- Follow the project's component patterns (see existing components)

Example:
```typescript
interface DeploymentCardProps {
  deployment: Deployment;
  onDelete: (id: string) => void;
}

export function DeploymentCard({ deployment, onDelete }: DeploymentCardProps) {
  return (
    <Card>
      <CardHeader>{deployment.name}</CardHeader>
    </Card>
  );
}
```

### File Organization

- Keep related files close together
- Use index.ts files for clean exports
- Place shared utilities in appropriate `utils/` or `lib/` folders
- Co-locate tests with source files or in `__tests__/` directories

## Testing

### Backend Tests

Located in the `/test` directory:

```bash
cd test
go test -v ./...
```

When adding features:
- Write unit tests for new functions
- Add integration tests for API endpoints
- Test edge cases and error conditions
- Aim for good coverage of critical paths

### Frontend Tests

Currently, the project uses manual testing. When adding features:
- Test in multiple browsers
- Verify responsive design on mobile/desktop
- Check accessibility with keyboard navigation
- Test error states and loading states

### Manual Testing Checklist

Before submitting:
- [ ] Application builds without errors
- [ ] Server starts and responds to requests
- [ ] Frontend loads and renders correctly
- [ ] Docker operations work (if modified)
- [ ] Git operations work (if modified)
- [ ] Database migrations work (if modified)
- [ ] WebSocket connections work (if modified)
- [ ] No console errors in browser

## Commit Messages

Use clear, descriptive commit messages that explain the "why" not just the "what":

```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `style:` - Formatting, missing semicolons, etc.
- `refactor:` - Code restructuring
- `test:` - Adding tests
- `chore:` - Build process, dependencies, etc.

Examples:
```
feat(deployments): add ability to stop ongoing deployments

fix(ssl): correct certificate renewal timing

docs(api): add examples for environment variable endpoints
```

## Pull Request Process

1. **Update your branch** with the latest `main`:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Ensure all tests pass**

3. **Update CHANGELOG.md** if your change is notable:
   - Add under the "Unreleased" section
   - Follow the existing format

4. **Fill out the PR template** (if provided) or include:
   - Description of changes
   - Motivation/context
   - Testing performed
   - Screenshots (for UI changes)
   - Breaking changes (if any)

5. **Request review** from maintainers

6. **Address feedback** promptly and professionally

7. **Once approved**, a maintainer will merge your PR

### PR Requirements

- [ ] Branch is up to date with `main`
- [ ] All CI checks pass
- [ ] Code follows style guidelines
- [ ] Tests added/updated for new functionality
- [ ] Documentation updated (if needed)
- [ ] CHANGELOG.md updated (if notable change)
- [ ] No merge conflicts

## Release Process

Mist follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality (backward compatible)
- **PATCH** version for bug fixes (backward compatible)

The CHANGELOG.md is maintained following [Keep a Changelog](https://keepachangelog.com/) format.

## Getting Help

- **Discord:** [Join our community](https://discord.gg/hr6TCQDDkj)
- **GitHub Issues:** [Report bugs or request features](https://github.com/trymist/mist/issues)
- **Documentation:** [trymist.cloud](https://trymist.cloud/guide/getting-started.html)

## Recognition

Contributors will be:
- Listed in the project's README (for significant contributions)
- Mentioned in release notes
- Acknowledged in our Discord community

Thank you for contributing to Mist! 

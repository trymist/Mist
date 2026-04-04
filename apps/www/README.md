# Mist Documentation

This directory contains the official documentation for Mist PaaS.

## Overview

The documentation is built using [VitePress](https://vitepress.dev/), a Vite & Vue powered static site generator.

## Development

### Prerequisites

- Node.js 18+
- npm or yarn

### Install Dependencies

```bash
npm install
```

### Development Server

Start the development server with hot-reload:

```bash
npm run docs:dev
```

Visit `http://localhost:5173` to view the documentation.

### Build

Build the static documentation site:

```bash
npm run docs:build
```

Output will be in `docs/.vitepress/dist/`.

### Preview Build

Preview the built documentation:

```bash
npm run docs:preview
```

## Structure

```
www/
├── docs/
│   ├── .vitepress/
│   │   ├── config.js          # VitePress configuration
│   │   └── theme/             # Custom theme
│   │       ├── index.js       # Theme entry
│   │       └── custom.css     # Custom styles (matches dash theme)
│   ├── guide/                 # User guides
│   │   ├── what-is-mist.md
│   │   ├── getting-started.md
│   │   ├── applications.md
│   │   ├── deployments.md
│   │   ├── environment-variables.md
│   │   ├── domains.md
│   │   ├── git-integration.md
│   │   ├── logs.md
│   │   ├── metrics.md
│   │   ├── authentication.md
│   │   ├── users.md
│   │   ├── audit-logs.md
│   │   ├── databases.md       # Coming Soon
│   │   ├── ssl-automation.md  # Coming Soon
│   │   ├── rollback.md        # Coming Soon
│   │   ├── notifications.md   # Coming Soon
│   │   └── cli.md             # Coming Soon
│   ├── api/                   # API documentation
│   │   ├── overview.md
│   │   ├── authentication.md
│   │   ├── projects.md
│   │   ├── applications.md
│   │   ├── deployments.md
│   │   ├── environment-variables.md
│   │   ├── domains.md
│   │   ├── users.md
│   │   ├── github.md
│   │   └── websockets.md
│   ├── deployment/            # Deployment guides
│   │   ├── installation.md
│   │   ├── configuration.md
│   │   ├── traefik.md
│   │   ├── github-app.md
│   │   ├── upgrading.md
│   │   ├── requirements.md
│   │   ├── security.md
│   │   └── backup.md
│   ├── public/                # Static assets
│   │   └── mist.png          # Logo
│   └── index.md               # Home page
├── package.json
└── README.md                  # This file
```

## Theme

The documentation uses a custom theme that matches the Mist dashboard's dark purple color scheme using OKLCH colors. See `docs/.vitepress/theme/custom.css` for the theme implementation.

### Color Scheme

- Primary: `oklch(0.488 0.243 264.376)` - Purple
- Background: `oklch(0.141 0.005 285.823)` - Very dark purple
- Card: `oklch(0.21 0.006 285.885)` - Dark purple

## Contributing

To add or update documentation:

1. Create or edit markdown files in the appropriate directory
2. Update `docs/.vitepress/config.js` navigation if adding new pages
3. Test locally with `npm run docs:dev`
4. Build with `npm run docs:build` to verify no errors
5. Submit a pull request

## Coming Soon Sections

Pages marked with "Coming Soon" badges represent planned features. These sections include:

- Database Services
- SSL Automation
- Deployment Rollback
- Notifications
- CLI Tool

As these features are implemented, update the corresponding documentation pages and remove the "Coming Soon" banners.

## Deployment

The documentation can be deployed to any static hosting service:

- GitHub Pages
- Netlify
- Vercel
- Cloudflare Pages
- Your own web server

Simply deploy the contents of `docs/.vitepress/dist/` after building.

## License

MIT

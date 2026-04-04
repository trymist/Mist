import { defineConfig } from 'vitepress'

export default defineConfig({

  title: 'Mist',
  description: 'Self-hostable Platform-as-a-Service',
  base: '/',
  head: [
    ['link', { rel: 'icon', type: 'image/png', href: '/mist.png' }]
  ],
  sitemap: {
    hostname: 'https://trymist.cloud'
  },

  themeConfig: {
    logo: '/mist.png',

    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/getting-started' },
      // { text: 'API', link: '/api/overview' },
      { text: 'Deployment', link: '/deployment/installation' },
      { text: 'Troubleshooting', link: '/troubleshooting' },
      { text: 'Changelog', link: '/changelog' }
      // { text: 'GitHub', link: 'https://github.com/trymist/mist' }
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Introduction',
          items: [
            { text: 'What is Mist?', link: '/guide/what-is-mist' },
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'Architecture', link: '/guide/architecture' }
          ]
        },
        {
          text: 'Core Features',
          items: [
            { text: 'Projects', link: '/guide/projects' },
            { text: 'Applications', link: '/guide/applications' },
            { text: 'Deployments', link: '/guide/deployments' },
            { text: 'Environment Variables', link: '/guide/environment-variables' },
            { text: 'Domains & SSL', link: '/guide/domains' },
            { text: 'Git Integration', link: '/guide/git-integration' },
            { text: 'Docker Compose', link: '/guide/docker-compose' },
            { text: 'Database Services', link: '/guide/databases' },
            { text: 'SSL Automation', link: '/guide/ssl-automation' },
          ]
        },
        {
          text: 'Monitoring',
          items: [
            { text: 'Logs', link: '/guide/logs' },
            { text: 'System Metrics', link: '/guide/metrics' },
            { text: 'Audit Logs', link: '/guide/audit-logs' }
          ]
        },
        {
          text: 'User Management',
          items: [
            { text: 'Authentication', link: '/guide/authentication' },
            { text: 'Users & Roles', link: '/guide/users' }
          ]
        },
        {
          text: 'Administration',
          items: [
            { text: 'CLI Tool', link: '/guide/cli' }
          ]
        },
        {
          text: 'Coming Soon',
          items: [
            { text: 'Rollback Deployments', link: '/guide/rollback' },
            { text: 'Notifications', link: '/guide/notifications' }
          ]
        }
      ],
      '/api/': [
        {
          text: 'API Reference',
          items: [
            { text: 'Overview', link: '/api/overview' },
            { text: 'Authentication', link: '/api/authentication' },
            { text: 'Projects', link: '/api/projects' },
            { text: 'Applications', link: '/api/applications' },
            { text: 'Deployments', link: '/api/deployments' },
            { text: 'Environment Variables', link: '/api/environment-variables' },
            { text: 'Domains', link: '/api/domains' },
            { text: 'Users', link: '/api/users' },
            { text: 'GitHub Integration', link: '/api/github' },
            { text: 'WebSockets', link: '/api/websockets' }
          ]
        }
      ],
      '/deployment/': [
        {
          text: 'Deployment',
          items: [
            { text: 'Installation', link: '/deployment/installation' },
            { text: 'Configuration', link: '/deployment/configuration' },
            { text: 'Traefik Setup', link: '/deployment/traefik' },
            { text: 'GitHub App Setup', link: '/deployment/github-app' },
            { text: 'Upgrading', link: '/deployment/upgrading' }
          ]
        },
        {
          text: 'Production',
          items: [
            { text: 'System Requirements', link: '/deployment/requirements' },
            { text: 'Security Best Practices', link: '/deployment/security' },
            { text: 'Backup & Recovery', link: '/deployment/backup' }
          ]
        }
      ],
      '/troubleshooting/': [
        {
          text: 'Troubleshooting',
          items: [
            { text: 'Common Issues', link: '/troubleshooting/' },
            { text: 'Forgot Password', link: '/troubleshooting/forgot-password' },
            { text: 'Manual Update', link: '/troubleshooting/manual-update' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/trymist/mist' },
      { icon: 'discord', link: 'https://discord.gg/hr6TCQDDkj' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025 Mist PaaS'
    },

    search: {
      provider: 'local'
    }
  }
})

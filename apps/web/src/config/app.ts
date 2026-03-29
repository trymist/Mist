
interface AppConfig {
  api: {
    baseURL: string;
    timeout: number;
  };
  app: {
    name: string;
    version: string;
    description: string;
  };
  features: {
    enableDevTools: boolean;
    enableErrorReporting: boolean;
  };
  ui: {
    theme: {
      default: 'light' | 'dark';
    };
    pagination: {
      defaultPageSize: number;
    };
  };
}

const config: AppConfig = {
  api: {
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 30000,
  },
  app: {
    name: 'Mist Dashboard',
    version: import.meta.env.VITE_APP_VERSION || '1.0.0',
    description: 'Modern project management dashboard',
  },
  features: {
    enableDevTools: import.meta.env.DEV || false,
    enableErrorReporting: import.meta.env.PROD || false,
  },
  ui: {
    theme: {
      default: 'light',
    },
    pagination: {
      defaultPageSize: 10,
    },
  },
};

export default config;

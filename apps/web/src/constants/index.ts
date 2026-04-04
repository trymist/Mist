
export const APP_CONSTANTS = {
  APP_NAME: 'Mist Dashboard',
  APP_VERSION: '1.0.0',
  
  API_ENDPOINTS: {
    AUTH: '/auth',
    PROJECTS: '/projects',
    USERS: '/users',
    GITHUB: '/github',
  },

  UI: {
    SIDEBAR_WIDTH: '16rem',
    SIDEBAR_WIDTH_MOBILE: '18rem',
    SIDEBAR_WIDTH_ICON: '3rem',
    DEFAULT_PAGE_SIZE: 10,
    MAX_PROJECTS_PER_PAGE: 50,
  },

  VALIDATION: {
    MIN_PASSWORD_LENGTH: 8,
    MAX_PROJECT_NAME_LENGTH: 100,
    MAX_DESCRIPTION_LENGTH: 500,
    MAX_TAGS_COUNT: 10,
  },

  STORAGE_KEYS: {
    THEME: 'mist-theme',
    SIDEBAR_STATE: 'mist-sidebar-state',
    USER_PREFERENCES: 'mist-user-preferences',
  },

  USER_ROLES: {
    OWNER: 'owner',
    ADMIN: 'admin',
    USER: 'user',
  } as const,

  PROJECT_STATUS: {
    ACTIVE: 'active',
    INACTIVE: 'inactive',
    ARCHIVED: 'archived',
  } as const,
} as const;

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  SETUP: '/setup',
  PROJECTS: '/projects',
  PROJECT_DETAIL: '/projects/:id',
  USERS: '/users',
  GIT: '/git',
  SETTINGS: '/settings',
  LOGS: '/logs',
  DEPLOYMENTS: '/deployments',
  DATABASES: '/databases',
  CALLBACK: '/callback',
} as const;

export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  FORBIDDEN: 'Access denied.',
  NOT_FOUND: 'The requested resource was not found.',
  SERVER_ERROR: 'Internal server error. Please try again later.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  GENERIC_ERROR: 'Something went wrong. Please try again.',
} as const;

export const SUCCESS_MESSAGES = {
  PROJECT_CREATED: 'Project created successfully',
  PROJECT_UPDATED: 'Project updated successfully',
  PROJECT_DELETED: 'Project deleted successfully',
  USER_CREATED: 'User created successfully',
  USER_UPDATED: 'User updated successfully',
  LOGIN_SUCCESS: 'Logged in successfully',
  LOGOUT_SUCCESS: 'Logged out successfully',
  SETUP_COMPLETE: 'Setup completed successfully',
} as const;

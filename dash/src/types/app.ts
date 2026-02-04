export type AppType = 'web' | 'service' | 'database' | 'compose';

export type RestartPolicy = 'no' | 'always' | 'on-failure' | 'unless-stopped';

export type App = {
  id: number;
  projectId: number;
  createdBy: number;
  name: string;
  description: string | null;
  appType: AppType;
  templateName: string | null;
  gitProviderId: number | null;
  gitRepository: string | null;
  gitBranch: string;
  gitCloneUrl: string | null;
  deploymentStrategy: string;
  port: number | null;
  shouldExpose: boolean | null;
  exposePort: number | null;
  rootDirectory: string;
  buildCommand: string | null;
  startCommand: string | null;
  dockerfilePath: string | null;
  cpuLimit: number | null;
  memoryLimit: number | null;
  restartPolicy: RestartPolicy;
  healthcheckPath: string | null;
  healthcheckInterval: number;
  status: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateAppRequest = {
  projectId: number;
  name: string;
  description?: string;
  appType: AppType;
  templateName?: string;
  gitRepository?: string;
  gitBranch?: string;
  port?: number;
  shouldExpose?: boolean;
  exposePort?: number;
  rootDirectory?: string;
  buildCommand?: string;
  startCommand?: string;
  cpuLimit?: number;
  memoryLimit?: number;
  restartPolicy?: RestartPolicy;
  envVars?: Record<string, string>;
};

export type UpdateAppRequest = Partial<Omit<App, 'id' | 'createdAt' | 'updatedAt'>>;

export type EnvVariable = {
  id: number;
  appId: number;
  key: string;
  value: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateEnvVariableRequest = {
  appId: number;
  key: string;
  value: string;
};

export type UpdateEnvVariableRequest = {
  id: number;
  key: string;
  value: string;
};

export type Domain = {
  id: number;
  appId: number;
  domain: string;
  sslStatus: string;
  dnsConfigured: boolean;
  dnsVerifiedAt?: string;
  lastDnsCheck?: string;
  dnsCheckError?: string;
  createdAt: string;
  updatedAt: string;
};

export type CreateDomainRequest = {
  appId: number;
  domain: string;
};

export type UpdateDomainRequest = {
  id: number;
  domain: string;
};

export type DNSRecord = {
  type: string;
  name: string;
  value: string;
};

export type DNSInstructions = {
  domain: string;
  serverIP: string;
  records: DNSRecord[];
};

export type DNSVerificationResponse = {
  domain: Domain;
  valid: boolean;
  serverIP: string;
  error?: string;
};

export type ContainerStatus = {
  name: string;
  status: string;
  state: string;
  uptime: string;
  healthy: boolean;
};

export type ServiceTemplateCategory = 'database' | 'cache' | 'queue' | 'storage' | 'other';

export type ServiceTemplate = {
  id: number;
  name: string;
  displayName: string;
  category: ServiceTemplateCategory;
  description: string | null;
  iconUrl: string | null;
  dockerImage: string;
  dockerImageVersion: string | null;
  defaultPort: number;
  defaultEnvVars: string | null; // JSON string
  requiredEnvVars: string | null; // JSON string
  defaultVolumePath: string | null;
  volumeRequired: boolean;
  recommendedCpu: number | null;
  recommendedMemory: number | null;
  minMemory: number | null;
  healthcheckCommand: string | null;
  healthcheckInterval: number;
  adminUiImage: string | null;
  adminUiPort: number | null;
  setupInstructions: string | null;
  isActive: boolean;
  isFeatured: boolean;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
};

export type Volume = {
  id: number;
  appId: number;
  name: string;
  hostPath: string;
  containerPath: string;
  readOnly: boolean;
  createdAt: string;
};

export type CreateVolumeRequest = {
  appId: number;
  name: string;
  hostPath: string;
  containerPath: string;
  readOnly?: boolean;
};

export type UpdateVolumeRequest = {
  id: number;
  name: string;
  hostPath: string;
  containerPath: string;
  readOnly: boolean;
};

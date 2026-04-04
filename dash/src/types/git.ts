
export interface GitHubApp {
  id: number;
  name: string;
  appId: number;
  clientId: string;
  slug: string;
  createdAt: string;
}

export interface Repository {
  id: string;
  name: string;
  fullName: string;
  provider: string;
  url: string;
  isPrivate: boolean;
  description?: string;
  defaultBranch: string;
  createdAt?: string;
  updatedAt?: string;
}

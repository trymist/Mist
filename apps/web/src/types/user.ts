
export interface User {
  id: number;
  username: string;
  email: string;
  role: 'owner' | 'admin' | 'user';
  avatarUrl?: string | null;
  isAdmin?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignupCredentials {
  username: string;
  email: string;
  password: string;
}

export interface UserCreateInput {
  username: string;
  email: string;
  password: string;
  role: 'admin' | 'user';
}

export type CreateUserData = UserCreateInput;

export interface UpdateProfileData {
  username?: string;
  email?: string;
}

export interface UpdatePasswordData {
  userId: number;
  currentPassword: string;
  newPassword: string;
}

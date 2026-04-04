import type { User, CreateUserData } from '@/types';

const API_BASE = '/api';

export const usersService = {
  async getAll(): Promise<User[]> {
    const response = await fetch(`${API_BASE}/users/getAll`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to fetch users');
    }

    const users: User[] = (data.data || []).map((u: User) => ({
      ...u,
      isAdmin: u.role === 'admin' || u.role === 'owner',
    }));

    return users;
  },

  async create(userData: CreateUserData): Promise<User> {
    const response = await fetch(`${API_BASE}/users/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(userData),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to create user');
    }

    return data.data;
  },

  async update(userId: number, updates: Partial<User>): Promise<User> {
    const response = await fetch(`${API_BASE}/users/update`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ id: userId, ...updates }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.message || 'Failed to update user');
    }

    return data.data;
  },

  /**
   * Update password
   */
  async updatePassword(userId: number, currentPassword: string, newPassword: string): Promise<void> {
    const response = await fetch(`${API_BASE}/users/password`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        userId,
        currentPassword,
        newPassword,
      }),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.message || 'Failed to update password');
    }
  },

  /**
   * Upload avatar
   */
  async uploadAvatar(file: File): Promise<{ avatarUrl: string; user: User }> {
    const formData = new FormData();
    formData.append('avatar', file);

    const response = await fetch(`${API_BASE}/users/avatar`, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to upload avatar');
    }

    return data.data;
  },

  /**
   * Delete avatar
   */
  async deleteAvatar(): Promise<User> {
    const response = await fetch(`${API_BASE}/users/avatar`, {
      method: 'DELETE',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete avatar');
    }

    return data.data;
  },

  /**
   * Delete user
   */
  async delete(userId: number): Promise<void> {
    const response = await fetch(`${API_BASE}/users/delete?id=${userId}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to delete user');
    }
  },
};

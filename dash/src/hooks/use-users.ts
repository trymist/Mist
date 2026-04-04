import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { usersService } from '@/services';
import type { User, CreateUserData } from '@/types';

interface UseUsersOptions {
  autoFetch?: boolean;
}

interface UseUsersReturn {
  users: User[];
  loading: boolean;
  error: string | null;
  fetchUsers: () => Promise<void>;
  createUser: (data: CreateUserData) => Promise<User | null>;
  updateUser: (id: number, data: Partial<User>) => Promise<User | null>;
  deleteUser: (id: number) => Promise<boolean>;
  refreshUsers: () => Promise<void>;
}

export const useUsers = (options: UseUsersOptions = {}): UseUsersReturn => {
  const { autoFetch = true } = options;
  
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await usersService.getAll();
      setUsers(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch users';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const createUser = useCallback(async (data: CreateUserData): Promise<User | null> => {
    try {
      const user = await usersService.create(data);
      setUsers(prev => [...prev, user]);
      toast.success('User created successfully');
      return user;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create user';
      toast.error(errorMessage);
      return null;
    }
  }, []);

  const updateUser = useCallback(async (id: number, data: Partial<User>): Promise<User | null> => {
    try {
      const updatedUser = await usersService.update(id, data);
      setUsers(prev => prev.map(u => u.id === id ? updatedUser : u));
      toast.success('User updated successfully');
      return updatedUser;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update user';
      toast.error(errorMessage);
      return null;
    }
  }, []);

  const deleteUser = useCallback(async (id: number): Promise<boolean> => {
    try {
      await usersService.delete(id);
      setUsers(prev => prev.filter(u => u.id !== id));
      toast.success('User deleted successfully');
      return true;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete user';
      toast.error(errorMessage);
      return false;
    }
  }, []);

  const refreshUsers = useCallback(() => fetchUsers(), [fetchUsers]);

  useEffect(() => {
    if (autoFetch) {
      fetchUsers();
    }
  }, [autoFetch, fetchUsers]);

  return {
    users,
    loading,
    error,
    fetchUsers,
    createUser,
    updateUser,
    deleteUser,
    refreshUsers,
  };
};

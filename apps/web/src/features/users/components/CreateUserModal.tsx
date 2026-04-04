import { FormModal } from '@/components/common/form-modal';
import { getAvailableRoles } from '../utils';
import type { User, CreateUserData } from '@/types';

interface CreateUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (userData: CreateUserData) => Promise<void>;
  currentUser: User | null;
}

export function CreateUserModal({ 
  isOpen, 
  onClose, 
  onSubmit, 
  currentUser 
}: CreateUserModalProps) {
  const handleSubmit = async (formData: Record<string, unknown>) => {
    await onSubmit(formData as unknown as CreateUserData);
  };

  return (
    <FormModal
      isOpen={isOpen}
      onClose={onClose}
      onSubmit={handleSubmit}
      title="Create New User"
      fields={[
        { 
          name: 'username', 
          label: 'Username', 
          type: 'text', 
          required: true
        },
        { 
          name: 'email', 
          label: 'Email', 
          type: 'email', 
          required: true
        },
        { 
          name: 'password', 
          label: 'Password', 
          type: 'password', 
          required: true
        },
        {
          name: 'role',
          label: 'Role',
          type: 'select',
          options: getAvailableRoles(currentUser),
          required: true,
          defaultValue: 'user'
        },
      ]}
    />
  );
}

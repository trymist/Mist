import { useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';

interface CreateAppModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
}

export function CreateAppModal({ isOpen, onClose, onConfirm }: CreateAppModalProps) {
  const [isCreating, setIsCreating] = useState(false);

  const handleConfirm = async () => {
    setIsCreating(true);
    try {
      await onConfirm();
      onClose();
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create GitHub App</DialogTitle>
          <DialogDescription>
            This will create a new GitHub App in your account with permissions for:
            <ul className="list-disc list-inside mt-2 text-muted-foreground">
              <li>Accessing your repositories</li>
              <li>Receiving push & deployment events</li>
              <li>Managing webhooks for automation</li>
              <li>Other users will be able to use this app for deployments</li>
            </ul>
            <p className="mt-2">
              You'll be redirected to GitHub to complete the process.
            </p>
          </DialogDescription>
        </DialogHeader>

        <DialogFooter className="flex justify-end space-x-2">
          <Button 
            variant="outline" 
            onClick={onClose}
            disabled={isCreating}
          >
            Cancel
          </Button>
          <Button 
            onClick={handleConfirm}
            disabled={isCreating}
          >
            {isCreating ? 'Creating...' : 'Continue'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Users, Crown, AlertCircle } from "lucide-react";
import { toast } from "sonner";
import { usersService, projectsService } from "@/services";
import type { User, Project } from "@/types";

interface ManageMembersModalProps {
  isOpen: boolean;
  onClose: () => void;
  project: Project;
  onSuccess: () => void;
}

export const ManageMembersModal = ({
  isOpen,
  onClose,
  project,
  onSuccess,
}: ManageMembersModalProps) => {
  const [allUsers, setAllUsers] = useState<User[]>([]);
  const [selectedUserIds, setSelectedUserIds] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (isOpen) {
      loadUsers();
    }
  }, [isOpen, project]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      setError("");
      const users = await usersService.getAll();
      setAllUsers(users);

      // Initialize selected users with current project members
      const memberIds = new Set(
        project.projectMembers.map((member) => Number(member.id))
      );
      setSelectedUserIds(memberIds);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to load users";
      setError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleUser = (userId: number) => {
    // Prevent unchecking the owner
    if (userId === Number(project.ownerId)) {
      toast.error("Cannot remove the project owner from members");
      return;
    }

    setSelectedUserIds((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(userId)) {
        newSet.delete(userId);
      } else {
        newSet.add(userId);
      }
      return newSet;
    });
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      setError("");

      // Ensure owner is always included
      const userIdsToSave = Array.from(selectedUserIds);
      const ownerId = Number(project.ownerId);
      if (!userIdsToSave.includes(ownerId)) {
        userIdsToSave.push(ownerId);
      }

      await projectsService.updateMembers(project.id, userIdsToSave);
      toast.success("Project members updated successfully");
      onSuccess();
      onClose();
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to update members";
      setError(message);
      toast.error(message);
    } finally {
      setSaving(false);
    }
  };

  const isOwner = (userId: number) => Number(project.ownerId) === userId;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-md max-h-[80vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Users className="h-5 w-5 text-primary" />
            Manage Project Members
          </DialogTitle>
          <DialogDescription>
            Select users who can access this project. The owner cannot be removed.
          </DialogDescription>
        </DialogHeader>

        <div className="flex-1 overflow-y-auto py-4">
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {loading ? (
            <div className="text-center py-8 text-muted-foreground">
              Loading users...
            </div>
          ) : allUsers.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No users found
            </div>
          ) : (
            <div className="space-y-3">
              {allUsers.map((user) => {
                const isUserOwner = isOwner(Number(user.id));
                const isChecked = selectedUserIds.has(Number(user.id));

                return (
                  <div
                    key={user.id}
                    className="flex items-center space-x-3 p-3 rounded-lg border border-border hover:bg-muted/50 transition-colors"
                  >
                    <Checkbox
                      id={`user-${user.id}`}
                      checked={isChecked}
                      onCheckedChange={() => handleToggleUser(Number(user.id))}
                      disabled={isUserOwner || saving}
                    />
                    <Label
                      htmlFor={`user-${user.id}`}
                      className="flex-1 flex items-center gap-2 cursor-pointer"
                    >
                      <div className="flex-1">
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{user.username}</span>
                          {isUserOwner && (
                            <Badge variant="secondary" className="flex items-center gap-1">
                              <Crown className="h-3 w-3" />
                              Owner
                            </Badge>
                          )}
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {user.email}
                        </span>
                      </div>
                      <Badge variant="outline">{user.role}</Badge>
                    </Label>
                  </div>
                );
              })}
            </div>
          )}
        </div>

        <DialogFooter className="flex-shrink-0">
          <Button variant="outline" onClick={onClose} disabled={saving}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={loading || saving}>
            {saving ? "Saving..." : "Save Changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

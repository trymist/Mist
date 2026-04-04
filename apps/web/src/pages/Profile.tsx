import { useState, useRef } from 'react';
import { useAuth } from '@/providers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { User, Mail, Lock, Shield, Calendar, AlertCircle, CheckCircle2, Camera, Trash2, Upload } from 'lucide-react';
import { usersService } from '@/services';
import { toast } from 'sonner';

export const ProfilePage = () => {
  const { user, setUser } = useAuth();
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false);
  const [isUpdatingPassword, setIsUpdatingPassword] = useState(false);
  const [isUploadingAvatar, setIsUploadingAvatar] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [username, setUsername] = useState(user?.username || '');
  const [email, setEmail] = useState(user?.email || '');
  const [profileError, setProfileError] = useState('');

  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordError, setPasswordError] = useState('');

  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    setProfileError('');

    if (!username.trim() || !email.trim()) {
      setProfileError('Username and email are required');
      return;
    }

    if (!user) return;

    setIsUpdatingProfile(true);

    try {
      const updatedUser = await usersService.update(user.id, {
        username: username.trim(),
        email: email.trim(),
      });

      setUser(updatedUser);
      toast.success('Profile updated successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update profile';
      setProfileError(message);
      toast.error(message);
    } finally {
      setIsUpdatingProfile(false);
    }
  };

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setPasswordError('');

    if (!currentPassword || !newPassword || !confirmPassword) {
      setPasswordError('All password fields are required');
      return;
    }

    if (newPassword.length < 6) {
      setPasswordError('New password must be at least 6 characters long');
      return;
    }

    if (newPassword !== confirmPassword) {
      setPasswordError('New passwords do not match');
      return;
    }

    if (!user) return;

    setIsUpdatingPassword(true);

    try {
      await usersService.updatePassword(user.id, currentPassword, newPassword);
      
      // Clear form
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
      
      toast.success('Password updated successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update password';
      setPasswordError(message);
      toast.error(message);
    } finally {
      setIsUpdatingPassword(false);
    }
  };

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      toast.error('Invalid file type. Please upload JPG, PNG, GIF, or WebP');
      return;
    }

    // Validate file size (5MB)
    if (file.size > 5 * 1024 * 1024) {
      toast.error('File size must be less than 5MB');
      return;
    }

    // Show preview
    const reader = new FileReader();
    reader.onloadend = () => {
      setAvatarPreview(reader.result as string);
    };
    reader.readAsDataURL(file);

    // Upload
    setIsUploadingAvatar(true);
    try {
      const result = await usersService.uploadAvatar(file);
      console.log('Upload result:', result);
      console.log('Avatar URL:', result.avatarUrl);
      console.log('Updated user:', result.user);
      setUser(result.user);
      toast.success('Avatar updated successfully');
      setAvatarPreview(null);
    } catch (error) {
      console.error('Avatar upload error:', error);
      const message = error instanceof Error ? error.message : 'Failed to upload avatar';
      toast.error(message);
      setAvatarPreview(null);
    } finally {
      setIsUploadingAvatar(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleDeleteAvatar = async () => {
    if (!user?.avatarUrl) return;

    setIsUploadingAvatar(true);
    try {
      const updatedUser = await usersService.deleteAvatar();
      setUser(updatedUser);
      toast.success('Avatar removed successfully');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete avatar';
      toast.error(message);
    } finally {
      setIsUploadingAvatar(false);
    }
  };

  const triggerFileInput = () => {
    fileInputRef.current?.click();
  };

  const getAvatarDisplay = () => {
    if (avatarPreview) return avatarPreview;
    if (user?.avatarUrl) {
      console.log('User avatar URL:', user.avatarUrl);
      return user.avatarUrl;
    }
    console.log('No avatar URL found for user');
    return null;
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const getRoleBadge = (role: string) => {
    const roleConfig = {
      owner: { label: 'Owner', className: 'bg-purple-500 text-white' },
      admin: { label: 'Admin', className: 'bg-blue-500 text-white' },
      user: { label: 'User', className: 'bg-gray-500 text-white' },
    };

    const config = roleConfig[role as keyof typeof roleConfig] || roleConfig.user;
    return <Badge className={config.className}>{config.label}</Badge>;
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background ">
      <div className="py-6 border-b border-border">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Profile
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage your profile and account settings
          </p>
        </div>
      </div>

      {/* Content */}
      <div className="py-6 max-w-4xl">
        {/* Profile Picture Section */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Camera className="h-5 w-5 text-primary" />
              Profile Picture
            </CardTitle>
            <CardDescription>
              Upload a profile picture to personalize your account
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row items-center gap-6">
              {/* Avatar Display */}
              <div className="relative">
                <div className="h-32 w-32 rounded-full overflow-hidden border-4 border-border bg-muted flex items-center justify-center">
                  {getAvatarDisplay() ? (
                    <img
                      src={getAvatarDisplay()!}
                      alt="Profile"
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div className="text-4xl font-bold text-muted-foreground">
                      {getInitials(user.username)}
                    </div>
                  )}
                </div>
                {isUploadingAvatar && (
                  <div className="absolute inset-0 flex items-center justify-center bg-black/50 rounded-full">
                    <span className="text-white text-2xl animate-spin">⏳</span>
                  </div>
                )}
              </div>

              {/* Upload Controls */}
              <div className="flex-1 space-y-3">
                <div className="flex flex-wrap gap-2">
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
                    onChange={handleAvatarChange}
                    className="hidden"
                  />
                  <Button
                    onClick={triggerFileInput}
                    disabled={isUploadingAvatar}
                    variant="outline"
                  >
                    <Upload className="h-4 w-4 mr-2" />
                    Upload Photo
                  </Button>
                  {user.avatarUrl && (
                    <Button
                      onClick={handleDeleteAvatar}
                      disabled={isUploadingAvatar}
                      variant="outline"
                      className="text-destructive hover:text-destructive"
                    >
                      <Trash2 className="h-4 w-4 mr-2" />
                      Remove
                    </Button>
                  )}
                </div>
                <p className="text-sm text-muted-foreground">
                  Accepted formats: JPG, PNG, GIF, WebP. Max size: 5MB
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Account Overview */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-primary" />
              Account Information
            </CardTitle>
            <CardDescription>
              Your account details and role information
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="flex items-center gap-3">
                <User className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Username</p>
                  <p className="font-medium">{user.username}</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <Mail className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Email</p>
                  <p className="font-medium">{user.email}</p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Role</p>
                  <div className="mt-1">{getRoleBadge(user.role)}</div>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <Calendar className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Member Since</p>
                  <p className="font-medium">{formatDate(user.createdAt)}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Update Profile */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5 text-primary" />
              Update Profile
            </CardTitle>
            <CardDescription>
              Update your username and email address
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleUpdateProfile} className="space-y-4">
              {profileError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{profileError}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-2">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="Enter username"
                  disabled={isUpdatingProfile}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="Enter email"
                  disabled={isUpdatingProfile}
                />
              </div>

              <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2">
                <Button
                  type="submit"
                  disabled={isUpdatingProfile || (username === user.username && email === user.email)}
                  className="w-full sm:w-auto"
                >
                  {isUpdatingProfile ? (
                    <>
                      <span className="animate-spin mr-2">⏳</span>
                      Updating...
                    </>
                  ) : (
                    <>
                      <CheckCircle2 className="h-4 w-4 mr-2" />
                      Update Profile
                    </>
                  )}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    setUsername(user.username);
                    setEmail(user.email);
                    setProfileError('');
                  }}
                  disabled={isUpdatingProfile}
                  className="w-full sm:w-auto"
                >
                  Reset
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* Change Password */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Lock className="h-5 w-5 text-primary" />
              Change Password
            </CardTitle>
            <CardDescription>
              Update your password to keep your account secure
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleUpdatePassword} className="space-y-4">
              {passwordError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{passwordError}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-2">
                <Label htmlFor="currentPassword">Current Password</Label>
                <Input
                  id="currentPassword"
                  type="password"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  placeholder="Enter current password"
                  disabled={isUpdatingPassword}
                />
              </div>

              <Separator />

              <div className="space-y-2">
                <Label htmlFor="newPassword">New Password</Label>
                <Input
                  id="newPassword"
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="Enter new password"
                  disabled={isUpdatingPassword}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirmPassword">Confirm New Password</Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="Confirm new password"
                  disabled={isUpdatingPassword}
                />
              </div>

              <div className="flex flex-col sm:flex-row items-start sm:items-center gap-2">
                <Button
                  type="submit"
                  disabled={isUpdatingPassword || !currentPassword || !newPassword || !confirmPassword}
                  className="w-full sm:w-auto"
                >
                  {isUpdatingPassword ? (
                    <>
                      <span className="animate-spin mr-2">⏳</span>
                      Updating...
                    </>
                  ) : (
                    <>
                      <Lock className="h-4 w-4 mr-2" />
                      Update Password
                    </>
                  )}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    setCurrentPassword('');
                    setNewPassword('');
                    setConfirmPassword('');
                    setPasswordError('');
                  }}
                  disabled={isUpdatingPassword}
                  className="w-full sm:w-auto"
                >
                  Clear
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

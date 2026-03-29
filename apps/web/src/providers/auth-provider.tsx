import React, { createContext, useContext, useEffect, useState } from "react";
import type { User } from "../types";
import { toast } from "sonner";
import { authApi } from "@/api/endpoints/auth";
type AuthContextType = {
  setupRequired: boolean | null;
  setSetupRequired: React.Dispatch<React.SetStateAction<boolean | null>>;
  user: User | null;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [setupRequired, setSetupRequired] = React.useState<boolean | null>(null);
  const [user, setUser] = useState<User | null>(null)

  const logout = async () => {
    const response = await authApi.logout();
    if (response.success) {
      setUser(null);
      toast.success("Logged out successfully");
    }
    else {
      toast.error(response.message || "An error occurred during logout.");
    }

  }

  useEffect(() => {
    const fetchAuth = async () => {
      const response = await authApi.getMe();
      if (response.success) {
        if (response.data.setupRequired === true) {
          setSetupRequired(true);
        }
        else {
          setSetupRequired(false);
          if (response.data.user) {
            setUser({
              ...response
                .data.user, isAdmin: response.data.user.role === "owner" || response.data.user.role === "admin"
            });
          }
        }
      } else {
        toast.error(response.message || "An error occurred while fetching authentication status.");
      }
    }
    fetchAuth();
  }, []);

  return (
    <AuthContext.Provider value={{ setupRequired, setSetupRequired, user, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );

}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

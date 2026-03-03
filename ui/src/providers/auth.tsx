import { type ReactNode, useEffect } from "react";
import { useUser } from "@/hooks/useUser";
import { useAuthStore } from "@/stores/auth";

export function AuthProvider({ children }: { children: ReactNode }) {
  const initializeAuth = useAuthStore((state) => state.initializeAuth);

  // Initialize auth on mount
  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  // Use the useUser hook to enable automatic refetching on window focus
  // This ensures user data (including tier) is always fresh when user returns
  // from external sites like Polar checkout
  useUser();

  return <>{children}</>;
}

export function useAuth() {
  return useAuthStore();
}

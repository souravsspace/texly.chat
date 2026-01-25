import { type ReactNode, useEffect } from "react";
import { useAuthStore } from "@/stores/auth";

export function AuthProvider({ children }: { children: ReactNode }) {
  const initializeAuth = useAuthStore((state) => state.initializeAuth);

  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  return <>{children}</>;
}

export function useAuth() {
  return useAuthStore();
}

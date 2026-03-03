import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { userQueryOptions } from "@/api/queries";
import { useAuthStore } from "@/stores/auth";

/**
 * Hook to get current user with automatic refetching
 *
 * This hook uses TanStack Query for data fetching with:
 * - Automatic refetch on window focus (e.g., returning from Polar checkout)
 * - Automatic refetch on mount
 * - 30 second stale time
 *
 * It also syncs the user data with the Zustand auth store for compatibility
 * with existing code.
 */
export function useUser() {
  const token = useAuthStore((state) => state.token);
  const setUser = useAuthStore((state) => state.setUser);

  const query = useQuery({
    ...userQueryOptions,
    enabled: !!token, // Only fetch if user is authenticated
  });

  // Sync query data with auth store
  useEffect(() => {
    if (query.data && !query.isError) {
      setUser(query.data);
    }
  }, [query.data, query.isError, setUser]);

  return query;
}

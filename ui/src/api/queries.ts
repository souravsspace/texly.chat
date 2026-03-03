import { queryOptions } from "@tanstack/react-query";
import { api } from "./index";

/**
 * User query with automatic refetching on window focus
 * This ensures user data (including tier) is always fresh
 */
export const userQueryOptions = queryOptions({
  queryKey: ["user", "me"],
  queryFn: () => api.users.getMe(),
  staleTime: 30_000, // 30 seconds - consider data fresh for 30s
  refetchOnWindowFocus: true, // Refetch when user returns to tab
  refetchOnMount: true, // Always refetch when component mounts
  retry: 1, // Only retry once on failure
});

/**
 * Billing usage query with aggressive refetching
 */
export const billingUsageQueryOptions = queryOptions({
  queryKey: ["billing-usage"],
  queryFn: () => api.billing.usage(),
  staleTime: 0, // Always consider stale
  refetchOnWindowFocus: true,
  refetchOnMount: "always",
});

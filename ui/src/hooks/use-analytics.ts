import { useQuery } from "@tanstack/react-start";
import { api } from "@/api";

export function useBotAnalytics(botId: string) {
  return useQuery({
    queryKey: ["analytics", "bot", botId],
    queryFn: () => api.analytics.getBotAnalytics(botId),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useBotDailyStats(botId: string, days = 30) {
  return useQuery({
    queryKey: ["analytics", "bot", botId, "daily", days],
    queryFn: () => api.analytics.getBotDailyStats(botId, days),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useUserAnalytics() {
  return useQuery({
    queryKey: ["analytics", "user"],
    queryFn: () => api.analytics.getUserAnalytics(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useSessionMessages(sessionId: string) {
  return useQuery({
    queryKey: ["analytics", "session", sessionId, "messages"],
    queryFn: () => api.analytics.getSessionMessages(sessionId),
    enabled: !!sessionId,
  });
}

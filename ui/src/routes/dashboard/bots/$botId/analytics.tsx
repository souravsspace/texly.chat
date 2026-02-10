import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { AlertCircle, CalendarDays, MessageSquare } from "lucide-react";
import { useState } from "react";
import { api } from "@/api";
import { Skeleton } from "@/components/ui/skeleton";
import { AnalyticsMetrics } from "./-components/analytics-metrics";
import { MessageChart } from "./-components/message-chart";
import { TokenUsageChart } from "./-components/token-usage-chart";

export const Route = createFileRoute("/dashboard/bots/$botId/analytics")({
  component: AnalyticsPage,
});

function AnalyticsPage() {
  const { botId } = Route.useParams();
  const [days, setDays] = useState(30);

  const {
    data: bot,
    isLoading: isBotLoading,
    error: botError,
  } = useQuery({
    queryKey: ["bot", botId],
    queryFn: () => api.bots.get(botId),
  });

  const {
    data: analytics,
    isLoading: isAnalyticsLoading,
    error: analyticsError,
  } = useQuery({
    queryKey: ["analytics", "bot", botId],
    queryFn: () => api.analytics.getBotAnalytics(botId),
  });

  const {
    data: dailyStats,
    isLoading: isDailyStatsLoading,
    error: dailyStatsError,
  } = useQuery({
    queryKey: ["analytics", "bot", botId, "daily", days],
    queryFn: () => api.analytics.getBotDailyStats(botId, days),
  });

  const isLoading = isBotLoading || isAnalyticsLoading || isDailyStatsLoading;
  const error = botError || analyticsError || dailyStatsError;

  // Calculate metrics from dailyStats based on selected time range
  const calculatePeriodMetrics = () => {
    if (!dailyStats || dailyStats.length === 0) {
      return null;
    }

    const totalMessages = dailyStats.reduce(
      (sum, stat) => sum + stat.message_count,
      0
    );
    const totalTokens = dailyStats.reduce(
      (sum, stat) => sum + stat.total_tokens,
      0
    );
    const totalSessions = dailyStats.reduce(
      (sum, stat) => sum + stat.unique_sessions,
      0
    );

    // Find last activity date from dailyStats
    const sortedStats = [...dailyStats].sort(
      (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
    );
    const lastActivity = sortedStats[0]?.date || null;

    return {
      total_messages: totalMessages,
      total_tokens: totalTokens,
      total_sessions: totalSessions,
      avg_messages_per_day: totalMessages / days,
      avg_tokens_per_day: totalTokens / days,
      avg_messages_per_session:
        totalSessions > 0 ? totalMessages / totalSessions : 0,
      last_message_at: lastActivity,
    };
  };

  const periodMetrics = calculatePeriodMetrics();

  if (error) {
    return (
      <div className="space-y-8">
        <div className="flex min-h-[400px] flex-col items-center justify-center gap-4 border border-destructive border-dashed p-8">
          <AlertCircle className="h-12 w-12 text-destructive" />
          <div className="text-center">
            <h3 className="font-semibold text-lg">Error loading analytics</h3>
            <p className="text-muted-foreground text-sm">
              {error instanceof Error
                ? error.message
                : "An unknown error occurred"}
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (isLoading || !bot || !analytics || !dailyStats) {
    return (
      <div className="space-y-8">
        {/* Metrics Cards Skeleton */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[1, 2, 3, 4].map((i) => (
            <div className="border bg-card p-6 shadow-sm" key={i}>
              <Skeleton className="mb-3 h-4 w-32" />
              <Skeleton className="h-8 w-20" />
              <Skeleton className="mt-2 h-3 w-24" />
            </div>
          ))}
        </div>

        {/* Charts Skeleton */}
        <div className="space-y-6">
          <Skeleton className="h-[400px] w-full" />
          <Skeleton className="h-[400px] w-full" />
        </div>
      </div>
    );
  }

  const hasData = analytics.total_messages > 0;

  return (
    <div className="space-y-6">
      {hasData && periodMetrics ? (
        <>
          {/* Summary Metrics */}
          <AnalyticsMetrics analytics={periodMetrics} />

          {/* Cost Estimation - Prominent */}
          <div className="border bg-gradient-to-r from-primary/10 to-primary/5 p-8 shadow-md">
            <div className="mb-6 flex items-center justify-between border-primary/20 border-b pb-4">
              <div>
                <h2 className="font-bold text-2xl">Cost Analysis</h2>
                <p className="mt-1 text-muted-foreground text-sm">
                  Last {days} days • gpt-4o-mini model
                </p>
              </div>
            </div>

            <div className="grid gap-6 md:grid-cols-3">
              {/* Total Cost */}
              <div className="space-y-2">
                <p className="font-medium text-muted-foreground text-sm uppercase tracking-wide">
                  Total Cost
                </p>
                <div className="font-bold text-5xl text-primary">
                  ${((periodMetrics.total_tokens / 1000) * 0.0001).toFixed(4)}
                </div>
              </div>

              {/* Token Usage */}
              <div className="space-y-2">
                <p className="font-medium text-muted-foreground text-sm uppercase tracking-wide">
                  Tokens Used
                </p>
                <div className="font-bold text-3xl">
                  {periodMetrics.total_tokens.toLocaleString()}
                </div>
                <p className="text-muted-foreground text-xs">
                  ~
                  {Math.round(
                    periodMetrics.total_tokens / days
                  ).toLocaleString()}{" "}
                  per day
                </p>
              </div>

              {/* Rate Info */}
              <div className="space-y-2">
                <p className="font-medium text-muted-foreground text-sm uppercase tracking-wide">
                  Rate
                </p>
                <div className="font-bold text-3xl">$0.0001</div>
                <p className="text-muted-foreground text-xs">per 1K tokens</p>
              </div>
            </div>

            {/* Projected Monthly Cost */}
            <div className="mt-6 border-primary/20 border-t pt-4">
              <p className="text-muted-foreground text-sm">
                Projected monthly cost:{" "}
                <span className="font-semibold text-foreground">
                  $
                  {(
                    ((periodMetrics.total_tokens / days) * 30 * 0.0001) /
                    1000
                  ).toFixed(4)}
                </span>{" "}
                (based on current usage)
              </p>
            </div>
          </div>

          {/* Time Range Selector */}
          <div className="flex items-center gap-2">
            <CalendarDays className="h-4 w-4 text-muted-foreground" />
            <span className="text-muted-foreground text-sm">Show last:</span>
            <div className="flex gap-2">
              {[7, 14, 30, 60, 90].map((value) => (
                <button
                  className={`px-3 py-1 text-sm transition-colors ${
                    days === value
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted text-muted-foreground hover:bg-muted/80"
                  }`}
                  key={value}
                  onClick={() => setDays(value)}
                  type="button"
                >
                  {value} days
                </button>
              ))}
            </div>
          </div>

          {/* Charts */}
          <div className="grid gap-6 lg:grid-cols-1">
            <MessageChart data={dailyStats} />
            <TokenUsageChart data={dailyStats} />
          </div>
        </>
      ) : (
        <div className="flex min-h-[400px] flex-col items-center justify-center gap-6 border border-dashed p-12">
          <div className="bg-muted p-6">
            <MessageSquare className="h-12 w-12 text-muted-foreground" />
          </div>
          <div className="space-y-2 text-center">
            <h3 className="font-semibold text-xl">No analytics data yet</h3>
            <p className="max-w-md text-muted-foreground text-sm">
              Your bot hasn't received any messages yet. Once users start
              chatting with your bot, you'll see detailed analytics here
              including message counts, token usage, and conversation patterns.
            </p>
          </div>
          <div className="flex gap-2">
            <Link
              className="bg-primary px-4 py-2 font-medium text-primary-foreground text-sm transition-colors hover:bg-primary/90"
              to={`/dashboard/bots/${botId}/chat`}
            >
              Test Your Bot
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}

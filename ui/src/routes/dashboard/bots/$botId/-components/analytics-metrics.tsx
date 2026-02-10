import { format } from "date-fns";
import { Activity, Hash, MessageSquare, TrendingUp } from "lucide-react";
import type { BotAnalytics } from "@/api/index.types";

interface AnalyticsMetricsProps {
  analytics: BotAnalytics;
}

function formatLastActivity(dateString: string | null): {
  value: string;
  description: string;
} {
  if (!dateString) {
    return { value: "No activity", description: "" };
  }

  const date = new Date(dateString);
  const now = new Date();
  const diffInHours = Math.floor(
    (now.getTime() - date.getTime()) / (1000 * 60 * 60)
  );

  if (diffInHours < 1) {
    return {
      value: "Just now",
      description: format(date, "h:mm a"),
    };
  }

  if (diffInHours < 24) {
    return {
      value: `${diffInHours}h ago`,
      description: format(date, "h:mm a"),
    };
  }

  if (diffInHours < 48) {
    return {
      value: "Yesterday",
      description: format(date, "h:mm a"),
    };
  }

  // Format as "4pm Feb 10"
  return {
    value: format(date, "ha MMM d"),
    description: format(date, "yyyy"),
  };
}

export function AnalyticsMetrics({ analytics }: AnalyticsMetricsProps) {
  const lastActivityFormat = formatLastActivity(analytics.last_message_at);

  const metrics = [
    {
      label: "Total Messages",
      value: analytics.total_messages.toLocaleString(),
      icon: MessageSquare,
      description: `${analytics.avg_messages_per_day.toFixed(1)} per day`,
    },
    {
      label: "Total Tokens",
      value: analytics.total_tokens.toLocaleString(),
      icon: Hash,
      description: `${analytics.avg_tokens_per_day.toLocaleString()} per day`,
    },
    {
      label: "Total Sessions",
      value: analytics.total_sessions.toLocaleString(),
      icon: Activity,
      description: `${analytics.avg_messages_per_session.toFixed(1)} msgs/session`,
    },
    {
      label: "Last Activity",
      value: lastActivityFormat.value,
      icon: TrendingUp,
      description: lastActivityFormat.description,
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {metrics.map((metric) => {
        const Icon = metric.icon;
        return (
          <div
            className="border bg-card p-6 text-card-foreground shadow-sm"
            key={metric.label}
          >
            <div className="flex items-center justify-between">
              <p className="font-medium text-muted-foreground text-sm">
                {metric.label}
              </p>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </div>
            <div className="mt-3">
              <p className="font-bold text-2xl">{metric.value}</p>
              {metric.description && (
                <p className="mt-1 text-muted-foreground text-xs">
                  {metric.description}
                </p>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}

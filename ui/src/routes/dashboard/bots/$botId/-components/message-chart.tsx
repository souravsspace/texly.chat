import { format } from "date-fns";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { MessageStats } from "@/api/index.types";

interface MessageChartProps {
  data: MessageStats[];
}

interface TooltipPayload {
  name: string;
  value: number;
  color: string;
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: TooltipPayload[];
  label?: string;
}

function CustomTooltip({ active, payload, label }: CustomTooltipProps) {
  const shouldShowTooltip = active && payload && payload.length > 0;

  if (!shouldShowTooltip) {
    return null;
  }

  return (
    <div className="border bg-card p-3 shadow-lg">
      <p className="mb-2 font-medium text-card-foreground text-sm">{label}</p>
      {payload.map((entry) => (
        <p className="text-xs" key={entry.name} style={{ color: entry.color }}>
          {entry.name}: {entry.value.toLocaleString()}
        </p>
      ))}
    </div>
  );
}

export function MessageChart({ data }: MessageChartProps) {
  const chartData = data.map((stat) => ({
    date: format(new Date(stat.date), "MMM dd"),
    "Total Messages": stat.message_count,
    "User Messages": stat.user_messages,
    "Bot Messages": stat.bot_messages,
  }));

  return (
    <div className="border bg-card p-6 shadow-sm">
      <h3 className="mb-4 font-semibold text-lg">Daily Message Activity</h3>
      <ResponsiveContainer height={300} width="100%">
        <LineChart
          data={chartData}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid stroke="rgba(255,255,255,0.1)" strokeDasharray="3 3" />
          <XAxis
            dataKey="date"
            stroke="rgba(255,255,255,0.3)"
            tick={{ fill: "rgba(255,255,255,0.9)", fontSize: 12 }}
          />
          <YAxis
            stroke="rgba(255,255,255,0.3)"
            tick={{ fill: "rgba(255,255,255,0.9)", fontSize: 12 }}
          />
          <Tooltip content={<CustomTooltip />} />
          <Legend
            wrapperStyle={{
              fontSize: "14px",
              color: "rgba(255,255,255,0.9)",
            }}
          />
          <Line
            activeDot={{ r: 8 }}
            dataKey="Total Messages"
            stroke="#e2e8f0"
            strokeWidth={2}
            type="monotone"
          />
          <Line
            dataKey="User Messages"
            stroke="#94a3b8"
            strokeWidth={2}
            type="monotone"
          />
          <Line
            dataKey="Bot Messages"
            stroke="#475569"
            strokeWidth={2}
            type="monotone"
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

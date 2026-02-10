import { format } from "date-fns";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { MessageStats } from "@/api/index.types";

interface TokenUsageChartProps {
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

export function TokenUsageChart({ data }: TokenUsageChartProps) {
  const chartData = data.map((stat) => ({
    date: format(new Date(stat.date), "MMM dd"),
    "Total Tokens": stat.total_tokens,
    Sessions: stat.unique_sessions,
  }));

  return (
    <div className="border bg-card p-6 shadow-sm">
      <h3 className="mb-4 font-semibold text-lg">Token Usage & Sessions</h3>
      <ResponsiveContainer height={300} width="100%">
        <BarChart
          data={chartData}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
          style={{ cursor: "default" }}
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
            yAxisId="left"
          />
          <YAxis
            orientation="right"
            stroke="rgba(255,255,255,0.3)"
            tick={{ fill: "rgba(255,255,255,0.9)", fontSize: 12 }}
            yAxisId="right"
          />
          <Tooltip
            content={<CustomTooltip />}
            cursor={{ fill: "rgba(255,255,255,0.05)" }}
          />
          <Legend
            wrapperStyle={{
              fontSize: "14px",
              color: "rgba(255,255,255,0.9)",
            }}
          />
          <Bar
            activeBar={{ fill: "#94a3b8", stroke: "none", fillOpacity: 1 }}
            dataKey="Total Tokens"
            fill="#94a3b8"
            fillOpacity={1}
            yAxisId="left"
          />
          <Bar
            activeBar={{ fill: "#cbd5e1", stroke: "none", fillOpacity: 1 }}
            dataKey="Sessions"
            fill="#cbd5e1"
            fillOpacity={1}
            yAxisId="right"
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}

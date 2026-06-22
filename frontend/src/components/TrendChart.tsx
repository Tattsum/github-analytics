import {
  Area,
  AreaChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { TimeSeriesPoint } from "../lib/timeSeries";

// additions/deletions are excluded: their magnitude dwarfs the others and would
// flatten the activity lines. Palette matches YearlyTrendChart.
const SERIES: ReadonlyArray<{ key: string; name: string; color: string }> = [
  { key: "commitCount", name: "コミット", color: "#2563eb" },
  { key: "prCreated", name: "PR作成", color: "#16a34a" },
  { key: "prMerged", name: "PRマージ", color: "#15803d" },
  { key: "reviewCount", name: "レビュー", color: "#d97706" },
  { key: "issueCount", name: "Issue", color: "#9333ea" },
];

export interface TrendChartProps {
  /** Pre-filtered, pre-bucketed series sorted ascending by date. */
  data: ReadonlyArray<TimeSeriesPoint>;
  height?: number;
  emptyMessage?: string;
}

export function TrendChart({ data, height = 320, emptyMessage = "対象期間の活動データはありません。" }: TrendChartProps) {
  if (data.length === 0) {
    return <p>{emptyMessage}</p>;
  }

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={[...data]}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis dataKey="date" tick={{ fontSize: 12 }} />
        <YAxis tick={{ fontSize: 12 }} />
        <Tooltip />
        <Legend />
        {SERIES.map((s) => (
          <Area
            key={s.key}
            type="monotone"
            dataKey={s.key}
            name={s.name}
            stroke={s.color}
            fill={s.color}
            fillOpacity={0.15}
            dot={false}
            isAnimationActive={false}
          />
        ))}
      </AreaChart>
    </ResponsiveContainer>
  );
}

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
import type { YearlyStatistics } from "../../gql/graphql";

export interface YearlyTrendChartProps {
  /** Per-year statistics for the member, oldest-to-newest is not assumed. */
  yearlyStats: ReadonlyArray<YearlyStatistics>;
  height?: number;
}

// Series rendered as stacked-but-overlaid areas. Activity metrics share a Y
// axis; additions/deletions are intentionally excluded here because their
// magnitude dwarfs the others and would flatten the activity lines.
const SERIES: ReadonlyArray<{ key: keyof YearlyStatistics; name: string; color: string }> = [
  { key: "commitCount", name: "コミット", color: "#2563eb" },
  { key: "prCreated", name: "PR作成", color: "#16a34a" },
  { key: "prMerged", name: "PRマージ", color: "#15803d" },
  { key: "reviewCount", name: "レビュー", color: "#d97706" },
  { key: "issueCount", name: "Issue", color: "#9333ea" },
];

// YearlyTrendChart plots the member's activity metrics over time. The data is
// sorted ascending by year so the X axis reads left-to-right chronologically
// regardless of the order the API returns rows in.
export function YearlyTrendChart({ yearlyStats, height = 320 }: YearlyTrendChartProps) {
  if (yearlyStats.length === 0) {
    return <p>年別の活動データはありません。</p>;
  }

  const data = [...yearlyStats].sort((a, b) => a.year - b.year);

  return (
    <ResponsiveContainer width="100%" height={height}>
      <AreaChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis dataKey="year" tick={{ fontSize: 12 }} />
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

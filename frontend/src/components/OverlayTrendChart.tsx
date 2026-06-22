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
import type { TimeSeriesPoint } from "../lib/timeSeries";

// One overlaid line: a stable data key (matching a column in the points) plus
// its display name and color.
export interface OverlaySeries {
  key: string;
  name: string;
  color: string;
}

export interface OverlayTrendChartProps {
  /** Pre-filtered, pre-bucketed points sorted ascending by date, with one
   *  numeric column per series key. */
  data: ReadonlyArray<TimeSeriesPoint>;
  series: ReadonlyArray<OverlaySeries>;
  height?: number;
  emptyMessage?: string;
}

// OverlayTrendChart overlays one line per entity (repository or member) for a
// single chosen metric, so trends can be compared on a shared time axis. Lines
// (not areas) are used so many overlapping series stay legible. Missing points
// leave gaps rather than being drawn as zero.
export function OverlayTrendChart({
  data,
  series,
  height = 320,
  emptyMessage = "対象期間の活動データはありません。",
}: OverlayTrendChartProps) {
  if (data.length === 0 || series.length === 0) {
    return <p>{emptyMessage}</p>;
  }

  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={[...data]}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis dataKey="date" tick={{ fontSize: 12 }} />
        <YAxis tick={{ fontSize: 12 }} />
        <Tooltip />
        <Legend />
        {series.map((s) => (
          <Line
            key={s.key}
            type="monotone"
            dataKey={s.key}
            name={s.name}
            stroke={s.color}
            dot={false}
            connectNulls={false}
            isAnimationActive={false}
          />
        ))}
      </LineChart>
    </ResponsiveContainer>
  );
}

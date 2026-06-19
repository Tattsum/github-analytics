import {
  Bar,
  BarChart as RechartsBarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export interface BarChartSeries<T> {
  /** dataKey on each datum to plot as a bar. */
  dataKey: keyof T & string;
  /** Optional bar color; falls back to the chart default. */
  color?: string;
  /** Optional legend/tooltip label; defaults to dataKey. */
  name?: string;
}

export interface BarChartProps<T> {
  data: readonly T[];
  /** dataKey used for the category (X) axis labels. */
  categoryKey: keyof T & string;
  series: ReadonlyArray<BarChartSeries<T>>;
  height?: number;
}

const DEFAULT_COLOR = "#2563eb";

// BarChart is a thin, typed wrapper around Recharts so pages declare data +
// series instead of repeating axis/container/tooltip boilerplate. Phase 1 pages
// (rankings, comparisons, repository metrics) render bars through this.
export function BarChart<T>({ data, categoryKey, series, height = 320 }: BarChartProps<T>) {
  return (
    <ResponsiveContainer width="100%" height={height}>
      <RechartsBarChart data={data as T[]}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis dataKey={categoryKey} tick={{ fontSize: 12 }} />
        <YAxis tick={{ fontSize: 12 }} />
        <Tooltip />
        {series.map((s) => (
          <Bar
            key={s.dataKey}
            dataKey={s.dataKey}
            name={s.name ?? s.dataKey}
            fill={s.color ?? DEFAULT_COLOR}
            isAnimationActive={false}
          />
        ))}
      </RechartsBarChart>
    </ResponsiveContainer>
  );
}

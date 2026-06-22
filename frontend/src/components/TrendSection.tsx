import { useMemo, useState } from "react";
import { bounds, prepareSeries, type TimeSeriesPoint } from "../lib/timeSeries";
import { DateRangePicker, type DateRange } from "./DateRangePicker";
import { TrendChart } from "./TrendChart";

export interface TrendSectionProps {
  /** Raw daily points from the API (one per active day, ascending or not). */
  dailyStats: ReadonlyArray<TimeSeriesPoint>;
  emptyMessage?: string;
}

// Defaults to the full available range bucketed by month until the user narrows it.
export function TrendSection({ dailyStats, emptyMessage }: TrendSectionProps) {
  const seriesBounds = useMemo(() => bounds(dailyStats), [dailyStats]);
  const [range, setRange] = useState<DateRange | null>(null);

  const effectiveRange: DateRange = range ?? {
    from: seriesBounds?.min ?? "",
    to: seriesBounds?.max ?? "",
    granularity: "month",
  };

  const data = useMemo(
    () =>
      prepareSeries(
        dailyStats,
        effectiveRange.granularity,
        effectiveRange.from || undefined,
        effectiveRange.to || undefined,
      ),
    [dailyStats, effectiveRange.from, effectiveRange.to, effectiveRange.granularity],
  );

  if (!seriesBounds) {
    return <p>{emptyMessage ?? "活動データはありません。"}</p>;
  }

  return (
    <div css={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
      <DateRangePicker value={effectiveRange} onChange={setRange} min={seriesBounds.min} max={seriesBounds.max} />
      <TrendChart data={data} emptyMessage={emptyMessage} />
    </div>
  );
}

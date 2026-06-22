// Pure helpers for the day-level activity time series. The backend returns a
// flat list of daily points (ISO "YYYY-MM-DD" dates, UTC-normalized); arbitrary
// date-range filtering and week/month bucketing are computed here on the
// frontend, consistent with how ranking/comparison are computed client-side.

// Granularity controls how daily points are grouped for the trend chart.
export type Granularity = "day" | "week" | "month";

// TimeSeriesPoint is an ISO-dated point with arbitrary numeric metrics. Non
// numeric fields (e.g. a GraphQL __typename) are ignored when summing buckets.
export interface TimeSeriesPoint {
  date: string;
  [metric: string]: string | number;
}

// filterByDateRange keeps points whose date is within [from, to], inclusive on
// both ends. An undefined bound is treated as open. ISO dates compare correctly
// as plain strings, so no Date parsing is needed.
export function filterByDateRange<T extends TimeSeriesPoint>(
  points: readonly T[],
  from?: string,
  to?: string,
): T[] {
  return points.filter((p) => (from === undefined || p.date >= from) && (to === undefined || p.date <= to));
}

export function bounds(points: readonly TimeSeriesPoint[]): { min: string; max: string } | undefined {
  if (points.length === 0) {
    return undefined;
  }
  let min = points[0]!.date;
  let max = points[0]!.date;
  for (const p of points) {
    if (p.date < min) min = p.date;
    if (p.date > max) max = p.date;
  }
  return { min, max };
}

// weekStart returns the Monday (ISO week start) of the given date, as an ISO
// date string. Computed in UTC so it never shifts across local timezones.
function weekStart(dateStr: string): string {
  const d = new Date(`${dateStr}T00:00:00Z`);
  const dayOfWeek = d.getUTCDay(); // 0=Sun .. 6=Sat
  const shiftToMonday = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
  d.setUTCDate(d.getUTCDate() + shiftToMonday);
  return d.toISOString().slice(0, 10);
}

// bucketKey maps a daily date to the label of the bucket it belongs to.
function bucketKey(dateStr: string, granularity: Granularity): string {
  switch (granularity) {
    case "day":
      return dateStr;
    case "week":
      return weekStart(dateStr);
    case "month":
      return `${dateStr.slice(0, 7)}-01`;
  }
}

// Each bucket is labelled by its start date so the chart's X axis reads chronologically.
export function bucketBy(points: readonly TimeSeriesPoint[], granularity: Granularity): TimeSeriesPoint[] {
  const buckets = new Map<string, TimeSeriesPoint>();

  for (const point of points) {
    const key = bucketKey(point.date, granularity);
    let bucket = buckets.get(key);
    if (bucket === undefined) {
      bucket = { date: key };
      buckets.set(key, bucket);
    }

    for (const [metric, value] of Object.entries(point)) {
      if (metric === "date" || typeof value !== "number") {
        continue;
      }
      bucket[metric] = ((bucket[metric] as number | undefined) ?? 0) + value;
    }
  }

  return [...buckets.values()].sort((a, b) => a.date.localeCompare(b.date));
}

export function prepareSeries(
  points: readonly TimeSeriesPoint[],
  granularity: Granularity,
  from?: string,
  to?: string,
): TimeSeriesPoint[] {
  return bucketBy(filterByDateRange(points, from, to), granularity);
}

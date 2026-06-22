// Pure helpers for the time-series comparison (multi-line overlay). Given several
// entities (repositories, or members within a repository), each with its own
// day-level series, these build the merged multi-column points a line chart needs
// (one numeric column per entity holding the chosen metric), pick a sensible
// default selection, and derive the owner filter options. All date filtering and
// week/month bucketing reuse the timeSeries helpers, consistent with the rest of
// the frontend computing ranking/comparison client-side.

import { prepareSeries, type Granularity, type TimeSeriesPoint } from "./timeSeries";

// DailyMetricPoint is one entity's metrics for a single day, mirroring the
// GraphQL DailyStatistics fields the comparison queries select.
export interface DailyMetricPoint {
  date: string;
  commitCount: number;
  prCreated: number;
  prMerged: number;
  reviewCount: number;
  issueCount: number;
  totalAdditions: number;
  totalDeletions: number;
}

// DailyMetricKey is the set of overlay-able numeric metrics (date excluded).
export type DailyMetricKey = Exclude<keyof DailyMetricPoint, "date">;

// ComparableSeries is one entity to overlay: a stable key (repository
// nameWithOwner, or member login) and its day-level series.
export interface ComparableSeries {
  key: string;
  daily: readonly DailyMetricPoint[];
}

// buildComparisonSeries projects each entity's daily series onto a single metric,
// labelling that entity's column by its key, then merges everything by date and
// applies the date-range filter and week/month bucketing. The result feeds a
// chart with one line per entity key.
export function buildComparisonSeries(
  entities: readonly ComparableSeries[],
  metric: DailyMetricKey,
  granularity: Granularity,
  from?: string,
  to?: string,
): TimeSeriesPoint[] {
  const projected: TimeSeriesPoint[] = [];
  for (const entity of entities) {
    for (const point of entity.daily) {
      projected.push({ date: point.date, [entity.key]: point[metric] });
    }
  }

  return prepareSeries(projected, granularity, from, to);
}

// totalForMetric sums one entity's metric across its whole series.
export function totalForMetric(entity: ComparableSeries, metric: DailyMetricKey): number {
  return entity.daily.reduce((sum, point) => sum + point[metric], 0);
}

// topEntityKeysByMetric returns the keys of the top-n entities by total of the
// chosen metric (descending). Ties break by key ascending so the default
// selection is deterministic.
export function topEntityKeysByMetric(
  entities: readonly ComparableSeries[],
  metric: DailyMetricKey,
  n: number,
): string[] {
  return [...entities]
    .map((entity) => ({ key: entity.key, total: totalForMetric(entity, metric) }))
    .sort((a, b) => (b.total - a.total) || a.key.localeCompare(b.key))
    .slice(0, Math.max(0, n))
    .map((entry) => entry.key);
}

// OwnerOption is one distinct repository owner for the org-internal filter.
export interface OwnerOption {
  owner: string;
  ownerType: string;
  count: number;
}

// RepoOwnerLike is the minimal shape needed to derive owner filter options.
export interface RepoOwnerLike {
  owner: string;
  ownerType: string;
}

// distinctOwners returns the unique owners across the repositories, sorted by
// owner ascending, each with the count of repositories it owns. Used to let the
// user scope a cross-repository comparison to a single org/owner.
export function distinctOwners(repos: readonly RepoOwnerLike[]): OwnerOption[] {
  const byOwner = new Map<string, OwnerOption>();
  for (const repo of repos) {
    const existing = byOwner.get(repo.owner);
    if (existing === undefined) {
      byOwner.set(repo.owner, { owner: repo.owner, ownerType: repo.ownerType, count: 1 });
      continue;
    }
    existing.count += 1;
    // Prefer a known owner type if the first row lacked it.
    if (existing.ownerType === "" && repo.ownerType !== "") {
      existing.ownerType = repo.ownerType;
    }
  }

  return [...byOwner.values()].sort((a, b) => a.owner.localeCompare(b.owner));
}

// SERIES_PALETTE is the line-color palette for the overlay. Length also defines
// the soft cap on simultaneously plotted series (kept readable; colors recycle
// only if a caller exceeds it).
export const SERIES_PALETTE: readonly string[] = [
  "#2563eb",
  "#dc2626",
  "#16a34a",
  "#d97706",
  "#9333ea",
  "#0891b2",
  "#db2777",
  "#65a30d",
];

// Maximum number of series shown at once, for chart readability.
export const MAX_SERIES = SERIES_PALETTE.length;

// colorForIndex returns a stable palette color for a series position.
export function colorForIndex(index: number): string {
  return SERIES_PALETTE[index % SERIES_PALETTE.length]!;
}

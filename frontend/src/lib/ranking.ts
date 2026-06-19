// Client-side ranking/sort helpers. Per the architecture, the GraphQL API
// returns flat, comparable metric lists and the frontend computes all sorting,
// ranking and comparison. These pure helpers are the single place that logic
// lives so member-axis and repository-axis pages stay consistent.

export type SortDirection = "asc" | "desc";

/** A function that extracts the numeric metric to rank/sort an item by. */
export type MetricSelector<T> = (item: T) => number;

/** Resolve a metric selector from either a key or a function. */
function toSelector<T>(metric: (keyof T & string) | MetricSelector<T>): MetricSelector<T> {
  if (typeof metric === "function") {
    return metric;
  }
  return (item: T) => Number(item[metric]);
}

/**
 * sortBy returns a new array sorted by the given numeric metric. It never
 * mutates the input. Direction defaults to descending (largest first), which
 * is the common case for "top contributors" style rankings.
 */
export function sortBy<T>(
  items: readonly T[],
  metric: (keyof T & string) | MetricSelector<T>,
  direction: SortDirection = "desc"
): T[] {
  const select = toSelector(metric);
  const factor = direction === "desc" ? -1 : 1;
  return [...items].sort((a, b) => (select(a) - select(b)) * factor);
}

export interface Ranked<T> {
  item: T;
  /** 1-based position in the ranking. */
  rank: number;
}

/**
 * rank sorts the items by the metric and assigns 1-based positions. Items with
 * an equal metric value share the same rank ("standard competition ranking",
 * e.g. 1, 2, 2, 4), so ties are not arbitrarily ordered by position.
 */
export function rank<T>(
  items: readonly T[],
  metric: (keyof T & string) | MetricSelector<T>,
  direction: SortDirection = "desc"
): Array<Ranked<T>> {
  const select = toSelector(metric);
  const sorted = sortBy(items, metric, direction);

  let lastValue: number | undefined;
  let lastRank = 0;
  return sorted.map((item, index) => {
    const value = select(item);
    if (lastValue === undefined || value !== lastValue) {
      lastRank = index + 1;
      lastValue = value;
    }
    return { item, rank: lastRank };
  });
}

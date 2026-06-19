import type { MetricSelector } from "../../lib/ranking";
import type { MemberStats } from "../../gql/graphql";

// The cross-member metrics a user can rank/compare on. The selector is the
// single source of truth for "how to derive this metric from a member", so the
// ranking table, the bar chart and any future widget all agree. `additions +
// deletions` and `pr/review ratio` are derived rather than plain fields, which
// is why metrics live here (a selector) instead of being a bare field key.

export type MetricKey =
  | "totalCommits"
  | "totalPRCreated"
  | "totalPRMerged"
  | "totalReviews"
  | "totalChangedLines"
  | "prToReviewRatio";

export interface MetricDefinition {
  key: MetricKey;
  /** Human-readable label for the picker, chart and column header. */
  label: string;
  /** Extracts the comparable numeric value from a member. */
  select: MetricSelector<MemberStats>;
  /** How to render a value as text (ratios get one decimal, counts are integers). */
  format: (value: number) => string;
}

const integerFormatter = new Intl.NumberFormat();

function formatInteger(value: number): string {
  return integerFormatter.format(value);
}

function formatRatio(value: number): string {
  return value.toFixed(1);
}

// changedLines combines additions and deletions; per the architecture these are
// PR-derived (commit line counts are unavailable from the API).
export function changedLines(member: MemberStats): number {
  return member.totalAdditions + member.totalDeletions;
}

export const METRICS: readonly MetricDefinition[] = [
  { key: "totalCommits", label: "コミット", select: (m) => m.totalCommits, format: formatInteger },
  { key: "totalPRCreated", label: "PR作成", select: (m) => m.totalPRCreated, format: formatInteger },
  { key: "totalPRMerged", label: "PRマージ", select: (m) => m.totalPRMerged, format: formatInteger },
  { key: "totalReviews", label: "レビュー", select: (m) => m.totalReviews, format: formatInteger },
  { key: "totalChangedLines", label: "追加行 + 削除行", select: changedLines, format: formatInteger },
  { key: "prToReviewRatio", label: "PR/レビュー比", select: (m) => m.prToReviewRatio, format: formatRatio },
];

const METRICS_BY_KEY: ReadonlyMap<MetricKey, MetricDefinition> = new Map(
  METRICS.map((metric) => [metric.key, metric])
);

/** Look up a metric definition by key, falling back to the first metric. */
export function metricByKey(key: MetricKey): MetricDefinition {
  return METRICS_BY_KEY.get(key) ?? METRICS[0];
}

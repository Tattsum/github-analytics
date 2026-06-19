// Metric definitions for the repository-axis pages. The GraphQL API returns
// flat, comparable metrics and the frontend owns all sorting/ranking, so these
// pure helpers are the single source of truth for "which numbers can we rank
// repositories / contributors by" and how to read them off the typed objects.

// Shape of a single repository's totals as returned by the `repositories` and
// `repository` queries. Kept structural (not importing generated types) so the
// helpers stay usable from tests with plain fixtures.
export interface RepoTotals {
  commits: number;
  prCreated: number;
  prMerged: number;
  issues: number;
  reviews: number;
  additions: number;
  deletions: number;
}

export interface RepoStatsLike {
  nameWithOwner: string;
  contributorCount: number;
  total: RepoTotals;
}

export interface ContributorLike {
  login: string;
  commitCount: number;
  prCreated: number;
  reviewCount: number;
  additions: number;
  deletions: number;
}

// A metric the user can pick from the metric selector. value extracts the
// number to rank/compare by; the keys are stable so they can drive a <select>.
export interface MetricOption<T> {
  key: string;
  label: string;
  value: (item: T) => number;
}

// Repository-axis metrics: one bar per repository, ranked across repositories.
export const repositoryMetrics: ReadonlyArray<MetricOption<RepoStatsLike>> = [
  { key: "commits", label: "コミット", value: (r) => r.total.commits },
  { key: "prCreated", label: "PR作成", value: (r) => r.total.prCreated },
  { key: "prMerged", label: "PRマージ", value: (r) => r.total.prMerged },
  { key: "issues", label: "Issue", value: (r) => r.total.issues },
  { key: "reviews", label: "レビュー", value: (r) => r.total.reviews },
  { key: "additions", label: "追加行", value: (r) => r.total.additions },
  { key: "deletions", label: "削除行", value: (r) => r.total.deletions },
  { key: "contributors", label: "コントリビューター数", value: (r) => r.contributorCount },
];

// Contributor-axis metrics: one bar per contributor within a single repository.
export const contributorMetrics: ReadonlyArray<MetricOption<ContributorLike>> = [
  { key: "commitCount", label: "コミット", value: (c) => c.commitCount },
  { key: "prCreated", label: "PR作成", value: (c) => c.prCreated },
  { key: "reviewCount", label: "レビュー", value: (c) => c.reviewCount },
  { key: "additions", label: "追加行", value: (c) => c.additions },
  { key: "deletions", label: "削除行", value: (c) => c.deletions },
];

// findMetric returns the option matching key, falling back to the first option
// when the key is unknown (e.g. a stale selection). Options must be non-empty.
export function findMetric<T>(
  options: ReadonlyArray<MetricOption<T>>,
  key: string
): MetricOption<T> {
  return options.find((o) => o.key === key) ?? options[0];
}

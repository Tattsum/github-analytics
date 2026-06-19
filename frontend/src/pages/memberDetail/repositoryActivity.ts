import type { RepositoryActivity } from "../../gql/graphql";

// Span in (fractional) years between a repository's first and last recorded
// activity. Used to surface long-term involvement. Both timestamps are ISO-8601
// strings from the API; an unparseable value yields 0 so it sorts last.
export function activitySpanYears(activity: Pick<RepositoryActivity, "firstActivity" | "lastActivity">): number {
  const first = Date.parse(activity.firstActivity);
  const last = Date.parse(activity.lastActivity);
  if (Number.isNaN(first) || Number.isNaN(last) || last < first) {
    return 0;
  }
  const millisPerYear = 365.25 * 24 * 60 * 60 * 1000;
  return (last - first) / millisPerYear;
}

// byCommitsDesc orders repositories by commit count, highest first, breaking
// ties on repository name for a stable, deterministic order.
export function byCommitsDesc(
  a: Pick<RepositoryActivity, "commitCount" | "repository">,
  b: Pick<RepositoryActivity, "commitCount" | "repository">,
): number {
  if (b.commitCount !== a.commitCount) {
    return b.commitCount - a.commitCount;
  }
  return a.repository.localeCompare(b.repository);
}

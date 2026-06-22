import { useMemo, useState } from "react";
import { useQuery } from "urql";
import { graphql } from "../../gql";
import { EntityTrendOverlay } from "./EntityTrendOverlay";
import { distinctOwners, type ComparableSeries } from "../../lib/comparison";

// Per-repository daily series with owner metadata, summed across members on the
// server. Backs the cross-repository trend overlay and the org-internal filter.
const RepositoryTrendComparisonQuery = graphql(`
  query RepositoryTrendComparison {
    repositoryDailyStats {
      nameWithOwner
      owner
      ownerType
      dailyStats {
        date
        commitCount
        prCreated
        prMerged
        reviewCount
        issueCount
        totalAdditions
        totalDeletions
      }
    }
  }
`);

const ALL_OWNERS = "";

// RepoTrendComparison overlays several repositories' day-level activity for one
// metric, optionally scoped to a single owner (e.g. an organization), to compare
// trends across repositories. It owns the owner filter; the metric/date/series
// selection and chart are delegated to the shared EntityTrendOverlay.
export function RepoTrendComparison() {
  const [{ data, fetching, error }] = useQuery({ query: RepositoryTrendComparisonQuery });
  const [ownerFilter, setOwnerFilter] = useState<string>(ALL_OWNERS);

  const repos = useMemo(() => data?.repositoryDailyStats ?? [], [data]);
  const owners = useMemo(() => distinctOwners(repos), [repos]);

  const entities: ComparableSeries[] = useMemo(
    () =>
      repos
        .filter((r) => ownerFilter === ALL_OWNERS || r.owner === ownerFilter)
        .map((r) => ({ key: r.nameWithOwner, daily: r.dailyStats })),
    [repos, ownerFilter],
  );

  if (fetching) {
    return <p>比較データを読み込み中…</p>;
  }
  if (error) {
    return <p>比較データを読み込めませんでした: {error.message}</p>;
  }
  if (repos.length === 0) {
    return <p>比較できるリポジトリがありません。</p>;
  }

  return (
    <div css={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
      <label css={{ display: "inline-flex", alignItems: "center", gap: "0.5rem" }}>
        <span css={{ fontSize: "0.875rem", color: "#374151" }}>オーナー</span>
        <select
          value={ownerFilter}
          onChange={(e) => setOwnerFilter(e.target.value)}
          css={{ padding: "0.375rem 0.5rem", borderRadius: "0.375rem", border: "1px solid #d1d5db", fontSize: "0.875rem" }}
        >
          <option value={ALL_OWNERS}>すべてのオーナー</option>
          {owners.map((o) => (
            <option key={o.owner} value={o.owner}>
              {o.owner}（{o.ownerType || "不明"}・{o.count}）
            </option>
          ))}
        </select>
      </label>

      <EntityTrendOverlay entities={entities} entityLabel="リポジトリ" />
    </div>
  );
}

import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { useQuery } from "urql";
import { graphql } from "../gql";
import { BarChart } from "../components/BarChart";
import { sortBy } from "../lib/ranking";
import { MetricPicker } from "./repositories/MetricPicker";
import { RankingTable, type RankingColumn } from "./repositories/RankingTable";
import { RepoTrendComparison } from "./repositories/RepoTrendComparison";
import { findMetric, repositoryMetrics, type RepoStatsLike } from "./repositories/metrics";

// Repository-axis cross aggregation. Fetches the flat `repositories` list and
// lets the user pick the metric to rank/compare by; ranking, sorting and the
// chart are all computed client-side per the architecture.
const RepositoriesQuery = graphql(`
  query Repositories {
    repositories {
      nameWithOwner
      contributorCount
      total {
        commits
        prCreated
        prMerged
        issues
        reviews
        additions
        deletions
      }
    }
  }
`);

// How many repositories to show in the comparison chart. Tables show all rows;
// the chart is capped so a large team's repo list stays readable.
const CHART_LIMIT = 15;

const columns: ReadonlyArray<RankingColumn<RepoStatsLike>> = [
  {
    key: "repository",
    header: "リポジトリ",
    render: (r) => (
      <Link to={`/repositories/${encodeURIComponent(r.nameWithOwner)}`}>{r.nameWithOwner}</Link>
    ),
  },
  { key: "commits", header: "コミット", numeric: true, render: (r) => r.total.commits },
  { key: "prCreated", header: "PR作成", numeric: true, render: (r) => r.total.prCreated },
  { key: "prMerged", header: "PRマージ", numeric: true, render: (r) => r.total.prMerged },
  { key: "issues", header: "Issue", numeric: true, render: (r) => r.total.issues },
  { key: "reviews", header: "レビュー", numeric: true, render: (r) => r.total.reviews },
  { key: "additions", header: "追加行", numeric: true, render: (r) => r.total.additions },
  { key: "deletions", header: "削除行", numeric: true, render: (r) => r.total.deletions },
  { key: "contributors", header: "コントリビューター数", numeric: true, render: (r) => r.contributorCount },
];

export function Repositories() {
  const [{ data, fetching, error }] = useQuery({ query: RepositoriesQuery });
  const [metricKey, setMetricKey] = useState(repositoryMetrics[0].key);

  const repositories: RepoStatsLike[] = useMemo(() => data?.repositories ?? [], [data]);
  const metric = findMetric(repositoryMetrics, metricKey);

  // Chart data: the top CHART_LIMIT repositories by the selected metric,
  // flattened to { label, value } so the generic BarChart can plot one series.
  const chartData = useMemo(
    () =>
      sortBy(repositories, metric.value)
        .slice(0, CHART_LIMIT)
        .map((r) => ({ label: r.nameWithOwner, value: metric.value(r) })),
    [repositories, metric]
  );

  return (
    <section>
      <h1>リポジトリ</h1>
      {fetching && <p>読み込み中…</p>}
      {error && <p>リポジトリを読み込めませんでした: {error.message}</p>}
      {data && (
        <>
          <div style={{ margin: "1rem 0" }}>
            <MetricPicker
              label="並び替え"
              options={repositoryMetrics}
              value={metricKey}
              onChange={setMetricKey}
            />
          </div>

          {chartData.length > 0 && (
            <BarChart
              data={chartData}
              categoryKey="label"
              series={[{ dataKey: "value", name: metric.label }]}
            />
          )}

          <div style={{ marginTop: "1.5rem", overflowX: "auto" }}>
            <RankingTable items={repositories} metric={metric} columns={columns} />
          </div>
        </>
      )}

      <section style={{ marginTop: "2.5rem" }}>
        <h2>活動推移の比較</h2>
        <p style={{ color: "#6b7280", marginTop: 0 }}>
          複数リポジトリの日別活動を重ね合わせて比較します。オーナーで絞り込めば組織内リポジトリ横断の推移を確認できます。
        </p>
        <RepoTrendComparison />
      </section>
    </section>
  );
}

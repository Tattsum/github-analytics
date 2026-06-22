import { useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { useQuery } from "urql";
import { graphql } from "../gql";
import { BarChart } from "../components/BarChart";
import { sortBy } from "../lib/ranking";
import { MetricPicker } from "./repositories/MetricPicker";
import { RankingTable, type RankingColumn } from "./repositories/RankingTable";
import { MemberTrendComparison } from "./repositories/MemberTrendComparison";
import { contributorMetrics, findMetric, type ContributorLike } from "./repositories/metrics";

// Single repository drill-down. Fetches one repository's totals plus its flat
// contributor list, then ranks contributors client-side by the chosen metric.
const RepositoryQuery = graphql(`
  query Repository($nameWithOwner: String!) {
    repository(nameWithOwner: $nameWithOwner) {
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
      contributors {
        login
        commitCount
        prCreated
        reviewCount
        additions
        deletions
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
  }
`);

// How many contributors to show in the in-repo comparison chart. The table
// lists everyone; the chart is capped to stay readable.
const CHART_LIMIT = 15;

const columns: ReadonlyArray<RankingColumn<ContributorLike>> = [
  { key: "login", header: "コントリビューター", render: (c) => c.login },
  { key: "commitCount", header: "コミット", numeric: true, render: (c) => c.commitCount },
  { key: "prCreated", header: "PR作成", numeric: true, render: (c) => c.prCreated },
  { key: "reviewCount", header: "レビュー", numeric: true, render: (c) => c.reviewCount },
  { key: "additions", header: "追加行", numeric: true, render: (c) => c.additions },
  { key: "deletions", header: "削除行", numeric: true, render: (c) => c.deletions },
];

interface TotalRow {
  label: string;
  value: number;
}

export function RepositoryDetail() {
  const { name } = useParams<{ name: string }>();
  const nameWithOwner = name ?? "";

  const [{ data, fetching, error }] = useQuery({
    query: RepositoryQuery,
    variables: { nameWithOwner },
    pause: nameWithOwner === "",
  });
  const [metricKey, setMetricKey] = useState(contributorMetrics[0].key);

  const repository = data?.repository ?? null;
  const contributors: ContributorLike[] = useMemo(
    () => repository?.contributors ?? [],
    [repository]
  );
  const metric = findMetric(contributorMetrics, metricKey);

  const totalRows: TotalRow[] = useMemo(() => {
    if (!repository) {
      return [];
    }
    const t = repository.total;
    return [
      { label: "コミット", value: t.commits },
      { label: "PR作成", value: t.prCreated },
      { label: "PRマージ", value: t.prMerged },
      { label: "Issue", value: t.issues },
      { label: "レビュー", value: t.reviews },
      { label: "追加行", value: t.additions },
      { label: "削除行", value: t.deletions },
    ];
  }, [repository]);

  // Chart data: top CHART_LIMIT contributors by the selected metric, flattened
  // to { label, value } for the generic BarChart.
  const chartData = useMemo(
    () =>
      sortBy(contributors, metric.value)
        .slice(0, CHART_LIMIT)
        .map((c) => ({ label: c.login, value: metric.value(c) })),
    [contributors, metric]
  );

  return (
    <section>
      <p>
        <Link to="/repositories">&larr; 一覧へ戻る</Link>
      </p>
      <h1>リポジトリ: {nameWithOwner}</h1>
      {fetching && <p>読み込み中…</p>}
      {error && <p>リポジトリを読み込めませんでした: {error.message}</p>}
      {data && !repository && <p>最新スナップショットにこのリポジトリが見つかりません。</p>}
      {repository && (
        <>
          <p css={{ color: "#6b7280" }}>コントリビューター {repository.contributorCount} 人</p>

          <h2>累計</h2>
          <ul>
            {totalRows.map((row) => (
              <li key={row.label}>
                {row.label}: {row.value}
              </li>
            ))}
          </ul>

          <h2>コントリビューターランキング</h2>
          <div css={{ margin: "1rem 0" }}>
            <MetricPicker
              label="並び替え"
              options={contributorMetrics}
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

          <div css={{ marginTop: "1.5rem", overflowX: "auto" }}>
            <RankingTable items={contributors} metric={metric} columns={columns} />
          </div>

          <section css={{ marginTop: "2.5rem" }}>
            <h2>メンバーの活動推移の比較</h2>
            <p css={{ color: "#6b7280", marginTop: 0 }}>
              このリポジトリ内での各メンバーの日別活動を重ね合わせて比較します。
            </p>
            <MemberTrendComparison contributors={repository.contributors} />
          </section>
        </>
      )}
    </section>
  );
}

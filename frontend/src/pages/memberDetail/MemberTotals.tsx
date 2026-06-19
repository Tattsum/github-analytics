import type { UserStatistics } from "../../gql/graphql";

export interface MemberTotalsProps {
  stats: Pick<
    UserStatistics,
    | "totalCommits"
    | "totalPRCreated"
    | "totalPRMerged"
    | "totalReviews"
    | "totalIssues"
    | "totalAdditions"
    | "totalDeletions"
    | "prToReviewRatio"
    | "firstActivityYear"
    | "peakActivityYear"
    | "peakActivityCommits"
  >;
}

interface Metric {
  label: string;
  value: number;
  /** When true, the value is a ratio and is rendered with two decimals. */
  ratio?: boolean;
  /** When true, the value is a calendar year and is rendered without grouping separators. */
  year?: boolean;
}

// MemberTotals renders the member's lifetime scalar metrics as a definition
// grid. Ranking/comparison against other members lives on other pages; this is
// a read-only single-member summary.
export function MemberTotals({ stats }: MemberTotalsProps) {
  const metrics: ReadonlyArray<Metric> = [
    { label: "コミット", value: stats.totalCommits },
    { label: "PR作成", value: stats.totalPRCreated },
    { label: "PRマージ", value: stats.totalPRMerged },
    { label: "レビュー", value: stats.totalReviews },
    { label: "Issue", value: stats.totalIssues },
    { label: "追加行", value: stats.totalAdditions },
    { label: "削除行", value: stats.totalDeletions },
    { label: "PR/レビュー比", value: stats.prToReviewRatio, ratio: true },
    { label: "活動開始年", value: stats.firstActivityYear, year: true },
    { label: "ピーク年", value: stats.peakActivityYear, year: true },
    { label: "ピーク時コミット", value: stats.peakActivityCommits },
  ];

  return (
    <dl
      style={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))",
        gap: "0.75rem",
        margin: 0,
      }}
    >
      {metrics.map((m) => (
        <div key={m.label}>
          <dt style={{ fontSize: 12, color: "#6b7280" }}>{m.label}</dt>
          <dd style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>
            {m.ratio ? m.value.toFixed(2) : m.year ? String(m.value) : m.value.toLocaleString()}
          </dd>
        </div>
      ))}
    </dl>
  );
}

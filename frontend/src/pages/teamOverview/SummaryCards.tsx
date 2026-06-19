import type { TeamSummary } from "../../gql/graphql";

const numberFormatter = new Intl.NumberFormat();

interface Card {
  label: string;
  value: number;
}

// SummaryCards renders the team-wide totals from `teamSummary` as a row of
// metric cards. It is presentational only; the page owns data fetching.
export function SummaryCards({ summary }: { summary: TeamSummary }) {
  const cards: Card[] = [
    { label: "メンバー数", value: summary.memberCount },
    { label: "リポジトリ", value: summary.repositoryCount },
    { label: "コミット", value: summary.totalCommits },
    { label: "PR作成", value: summary.totalPRCreated },
    { label: "PRマージ", value: summary.totalPRMerged },
    { label: "レビュー", value: summary.totalReviews },
    { label: "Issue", value: summary.totalIssues },
    { label: "追加行", value: summary.totalAdditions },
    { label: "削除行", value: summary.totalDeletions },
  ];

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fill, minmax(150px, 1fr))",
        gap: "0.75rem",
      }}
    >
      {cards.map((card) => (
        <div
          key={card.label}
          style={{
            border: "1px solid #e5e7eb",
            borderRadius: "0.5rem",
            padding: "0.75rem 1rem",
            backgroundColor: "#ffffff",
          }}
        >
          <div style={{ fontSize: "0.8rem", color: "#6b7280" }}>{card.label}</div>
          <div style={{ fontSize: "1.5rem", fontWeight: 600 }}>
            {numberFormatter.format(card.value)}
          </div>
        </div>
      ))}
    </div>
  );
}

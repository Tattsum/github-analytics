import type { RepositoryActivity } from "../../gql/graphql";
import { activitySpanYears, byCommitsDesc } from "./repositoryActivity";

export interface RepositoryActivityListProps {
  repositories: ReadonlyArray<RepositoryActivity>;
  /** When true, show the first..last activity span instead of commit-derived columns ordering. */
  showSpan?: boolean;
  emptyMessage?: string;
}

// RepositoryActivityList renders a member's repository involvement as a table.
// It is shared by both the "top repositories" and "long-term repositories"
// sections; `showSpan` adds first/last activity columns for the long-term view.
export function RepositoryActivityList({
  repositories,
  showSpan = false,
  emptyMessage = "リポジトリのデータはありません。",
}: RepositoryActivityListProps) {
  if (repositories.length === 0) {
    return <p>{emptyMessage}</p>;
  }

  const rows = [...repositories].sort(byCommitsDesc);

  return (
    <table style={{ width: "100%", borderCollapse: "collapse", fontSize: 14 }}>
      <thead>
        <tr style={{ textAlign: "left", borderBottom: "1px solid #e5e7eb" }}>
          <th style={cellStyle}>リポジトリ</th>
          <th style={cellStyle}>コミット</th>
          <th style={cellStyle}>PR</th>
          <th style={cellStyle}>レビュー</th>
          <th style={cellStyle}>Issue</th>
          <th style={cellStyle}>+/-</th>
          {showSpan && <th style={cellStyle}>期間（年）</th>}
        </tr>
      </thead>
      <tbody>
        {rows.map((r) => (
          <tr key={r.repository} style={{ borderBottom: "1px solid #f3f4f6" }}>
            <td style={cellStyle}>{r.repository}</td>
            <td style={cellStyle}>{r.commitCount.toLocaleString()}</td>
            <td style={cellStyle}>{r.prCount.toLocaleString()}</td>
            <td style={cellStyle}>{r.reviewCount.toLocaleString()}</td>
            <td style={cellStyle}>{r.issueCount.toLocaleString()}</td>
            <td style={cellStyle}>
              +{r.totalAdditions.toLocaleString()} / -{r.totalDeletions.toLocaleString()}
            </td>
            {showSpan && <td style={cellStyle}>{activitySpanYears(r).toFixed(1)}</td>}
          </tr>
        ))}
      </tbody>
    </table>
  );
}

const cellStyle: React.CSSProperties = { padding: "0.4rem 0.6rem" };

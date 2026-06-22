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
    <div css={{ overflowX: "auto" }}>
      <table css={{ width: "100%", minWidth: 560, borderCollapse: "collapse", fontSize: 14 }}>
        <thead>
          <tr css={{ textAlign: "left", borderBottom: "1px solid #e5e7eb" }}>
            <th css={cellStyle}>リポジトリ</th>
            <th css={cellStyle}>コミット</th>
            <th css={cellStyle}>PR</th>
            <th css={cellStyle}>レビュー</th>
            <th css={cellStyle}>Issue</th>
            <th css={cellStyle}>+/-</th>
            {showSpan && <th css={cellStyle}>期間（年）</th>}
          </tr>
        </thead>
        <tbody>
          {rows.map((r) => (
            <tr key={r.repository} css={{ borderBottom: "1px solid #f3f4f6" }}>
              <td css={cellStyle}>{r.repository}</td>
              <td css={cellStyle}>{r.commitCount.toLocaleString()}</td>
              <td css={cellStyle}>{r.prCount.toLocaleString()}</td>
              <td css={cellStyle}>{r.reviewCount.toLocaleString()}</td>
              <td css={cellStyle}>{r.issueCount.toLocaleString()}</td>
              <td css={cellStyle}>
                +{r.totalAdditions.toLocaleString()} / -{r.totalDeletions.toLocaleString()}
              </td>
              {showSpan && <td css={cellStyle}>{activitySpanYears(r).toFixed(1)}</td>}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

const cellStyle = { padding: "0.4rem 0.6rem" } as const;

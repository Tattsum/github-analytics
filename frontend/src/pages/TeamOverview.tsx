import { useQuery } from "urql";
import { graphql } from "../gql";
import { SummaryCards } from "./teamOverview/SummaryCards";
import { RankingBoard } from "./teamOverview/RankingBoard";
import { TrendSection } from "../components/TrendSection";

// Team-wide aggregates for the summary cards. Reads the latest snapshot.
const TeamSummaryQuery = graphql(`
  query TeamOverviewSummary {
    teamSummary {
      memberCount
      repositoryCount
      totalCommits
      totalPRCreated
      totalPRMerged
      totalIssues
      totalReviews
      totalAdditions
      totalDeletions
    }
  }
`);

// Team-wide daily totals as an ascending time series. Date-range filtering and
// week/month bucketing are computed on the frontend (TrendSection).
const TeamDailyStatsQuery = graphql(`
  query TeamOverviewDailyStats {
    teamDailyStats {
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
`);

// Cross-member comparable scalars. Ranking/sorting/comparison is computed on
// the frontend (RankingBoard), so this just pulls the flat list.
const MembersQuery = graphql(`
  query TeamOverviewMembers {
    members {
      login
      name
      totalCommits
      totalPRCreated
      totalPRMerged
      totalIssues
      totalReviews
      totalAdditions
      totalDeletions
      prToReviewRatio
    }
  }
`);

// TeamOverview is the main page (route "/"): team summary cards plus a
// cross-member ranking board with a metric picker and bar-chart comparison.
export function TeamOverview() {
  const [summaryResult] = useQuery({ query: TeamSummaryQuery });
  const [membersResult] = useQuery({ query: MembersQuery });
  const [dailyResult] = useQuery({ query: TeamDailyStatsQuery });

  const fetching = summaryResult.fetching || membersResult.fetching || dailyResult.fetching;
  const error = summaryResult.error ?? membersResult.error ?? dailyResult.error;

  return (
    <section css={{ display: "flex", flexDirection: "column", gap: "2rem" }}>
      <h1>チーム概要</h1>

      {fetching && <p>読み込み中…</p>}
      {error && <p css={{ color: "#b91c1c" }}>概要を読み込めませんでした: {error.message}</p>}

      {summaryResult.data && <SummaryCards summary={summaryResult.data.teamSummary} />}

      {dailyResult.data && (
        <div>
          <h2>活動推移（期間指定）</h2>
          <TrendSection
            dailyStats={dailyResult.data.teamDailyStats}
            emptyMessage="最新スナップショットに日次の活動データがまだありません。"
          />
        </div>
      )}

      {membersResult.data && (
        <div>
          <h2>メンバーランキング・比較</h2>
          {membersResult.data.members.length === 0 ? (
            <p>最新スナップショットにメンバーデータがまだありません。</p>
          ) : (
            <RankingBoard members={membersResult.data.members} />
          )}
        </div>
      )}
    </section>
  );
}

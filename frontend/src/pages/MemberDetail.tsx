import { useParams } from "react-router-dom";
import { useQuery } from "urql";
import { graphql } from "../gql";
import { MemberTotals } from "./memberDetail/MemberTotals";
import { RepositoryActivityList } from "./memberDetail/RepositoryActivityList";
import { RoleTransition } from "./memberDetail/RoleTransition";
import { YearlyTrendChart } from "./memberDetail/YearlyTrendChart";
import { TrendSection } from "../components/TrendSection";

// Drill-down query for a single member. Reads the latest snapshot server-side;
// all ranking/sorting/comparison is handled on other pages, so this is purely
// the per-member detail (totals, yearly trend, repositories, role transition).
const MemberDetailQuery = graphql(`
  query MemberDetail($login: String!) {
    member(login: $login) {
      login
      name
      totalCommits
      totalPRCreated
      totalPRMerged
      totalReviews
      totalIssues
      totalAdditions
      totalDeletions
      prToReviewRatio
      firstActivityYear
      peakActivityYear
      peakActivityCommits
      yearlyStats {
        year
        commitCount
        prCreated
        prMerged
        reviewCount
        issueCount
        totalAdditions
        totalDeletions
      }
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
      topRepositories {
        repository
        commitCount
        prCount
        reviewCount
        issueCount
        totalAdditions
        totalDeletions
        firstActivity
        lastActivity
      }
      longTermRepositories {
        repository
        commitCount
        prCount
        reviewCount
        issueCount
        totalAdditions
        totalDeletions
        firstActivity
        lastActivity
      }
      roleTransition {
        year
        description
        prCreated
        reviewCount
        ratio
      }
    }
  }
`);

const sectionStyle: React.CSSProperties = { marginTop: "2rem" };

// MemberDetail is the read-only per-member dashboard at /members/:login.
export function MemberDetail() {
  const { login } = useParams<{ login: string }>();
  const [{ data, fetching, error }] = useQuery({
    query: MemberDetailQuery,
    variables: { login: login ?? "" },
    pause: !login,
  });

  if (!login) {
    return <p>メンバーが指定されていません。</p>;
  }
  if (fetching) {
    return <p>{login} を読み込み中…</p>;
  }
  if (error) {
    return <p>メンバー {login} を読み込めませんでした: {error.message}</p>;
  }
  const member = data?.member;
  if (!member) {
    return <p>メンバー {login} が見つかりません。</p>;
  }

  return (
    <section>
      <h1>
        {member.name} <span style={{ color: "#6b7280", fontWeight: 400 }}>@{member.login}</span>
      </h1>

      <section style={sectionStyle}>
        <h2>累計</h2>
        <MemberTotals stats={member} />
      </section>

      <section style={sectionStyle}>
        <h2>活動推移（期間指定）</h2>
        <TrendSection dailyStats={member.dailyStats} emptyMessage="対象期間の活動データはありません。" />
      </section>

      <section style={sectionStyle}>
        <h2>年別推移</h2>
        <YearlyTrendChart yearlyStats={member.yearlyStats} />
      </section>

      <section style={sectionStyle}>
        <h2>主なリポジトリ</h2>
        <RepositoryActivityList
          repositories={member.topRepositories}
          emptyMessage="リポジトリの活動データはありません。"
        />
      </section>

      <section style={sectionStyle}>
        <h2>長期間関与リポジトリ</h2>
        <RepositoryActivityList
          repositories={member.longTermRepositories}
          showSpan
          emptyMessage="長期間関与リポジトリはありません。"
        />
      </section>

      <section style={sectionStyle}>
        <h2>役割の変化</h2>
        <RoleTransition points={member.roleTransition} />
      </section>
    </section>
  );
}

import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { RoleTransitionPoint } from "../../gql/graphql";

export interface RoleTransitionProps {
  points: ReadonlyArray<RoleTransitionPoint>;
  height?: number;
}

// RoleTransition visualises how a member's contribution shifts between authoring
// (PRs created) and reviewing over the years, plus the review/PR ratio that
// captures the author-to-reviewer transition. The per-year `description`
// (e.g. "author", "reviewer") is surfaced in the table below the chart.
export function RoleTransition({ points, height = 280 }: RoleTransitionProps) {
  if (points.length === 0) {
    return <p>役割変化のデータはありません。</p>;
  }

  const data = [...points].sort((a, b) => a.year - b.year);

  return (
    <>
      <ResponsiveContainer width="100%" height={height}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
          <XAxis dataKey="year" tick={{ fontSize: 12 }} />
          <YAxis tick={{ fontSize: 12 }} />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="prCreated" name="PR作成" stroke="#16a34a" dot={false} isAnimationActive={false} />
          <Line type="monotone" dataKey="reviewCount" name="レビュー" stroke="#d97706" dot={false} isAnimationActive={false} />
        </LineChart>
      </ResponsiveContainer>
      <div css={{ overflowX: "auto", marginTop: "0.75rem" }}>
        <table css={{ width: "100%", minWidth: 480, borderCollapse: "collapse", fontSize: 14 }}>
          <thead>
            <tr css={{ textAlign: "left", borderBottom: "1px solid #e5e7eb" }}>
              <th css={cellStyle}>年</th>
              <th css={cellStyle}>役割</th>
              <th css={cellStyle}>PR作成</th>
              <th css={cellStyle}>レビュー</th>
              <th css={cellStyle}>レビュー/PR比</th>
            </tr>
          </thead>
          <tbody>
            {data.map((p) => (
              <tr key={p.year} css={{ borderBottom: "1px solid #f3f4f6" }}>
                <td css={cellStyle}>{p.year}</td>
                <td css={cellStyle}>{p.description}</td>
                <td css={cellStyle}>{p.prCreated.toLocaleString()}</td>
                <td css={cellStyle}>{p.reviewCount.toLocaleString()}</td>
                <td css={cellStyle}>{p.ratio.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}

const cellStyle = { padding: "0.4rem 0.6rem" } as const;

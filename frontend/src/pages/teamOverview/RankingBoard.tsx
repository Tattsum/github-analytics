import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import type { MemberStats } from "../../gql/graphql";
import { rank } from "../../lib/ranking";
import { BarChart } from "../../components/BarChart";
import { METRICS, metricByKey, type MetricKey } from "./metrics";

// How many top members to plot in the comparison chart. The full ranking table
// always lists everyone; the chart stays readable by showing the leaders only.
const CHART_LIMIT = 15;

interface ChartDatum {
  login: string;
  value: number;
}

// RankingBoard lets the user pick a metric and shows the members ranked by it
// (client-side, via the shared ranking util) plus a bar-chart comparison of the
// same metric. Ranking/sorting is intentionally done here, not in GraphQL.
export function RankingBoard({ members }: { members: readonly MemberStats[] }) {
  const [metricKey, setMetricKey] = useState<MetricKey>("totalCommits");
  const metric = metricByKey(metricKey);

  const ranked = useMemo(
    () => rank(members, metric.select, "desc"),
    [members, metric]
  );

  const chartData = useMemo<ChartDatum[]>(
    () =>
      ranked.slice(0, CHART_LIMIT).map(({ item }) => ({
        login: item.login,
        value: metric.select(item),
      })),
    [ranked, metric]
  );

  return (
    <div>
      <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", marginBottom: "1rem" }}>
        <label htmlFor="ranking-metric" style={{ fontWeight: 600 }}>
          指標
        </label>
        <select
          id="ranking-metric"
          value={metricKey}
          onChange={(event) => setMetricKey(event.target.value as MetricKey)}
          style={{ padding: "0.35rem 0.5rem", borderRadius: "0.375rem", border: "1px solid #d1d5db" }}
        >
          {METRICS.map((m) => (
            <option key={m.key} value={m.key}>
              {m.label}
            </option>
          ))}
        </select>
      </div>

      {chartData.length > 0 && (
        <BarChart
          data={chartData}
          categoryKey="login"
          series={[{ dataKey: "value", name: metric.label }]}
        />
      )}

      <table style={{ width: "100%", borderCollapse: "collapse", marginTop: "1rem" }}>
        <thead>
          <tr style={{ textAlign: "left", borderBottom: "2px solid #e5e7eb" }}>
            <th style={{ padding: "0.5rem", width: "3rem" }}>#</th>
            <th style={{ padding: "0.5rem" }}>メンバー</th>
            <th style={{ padding: "0.5rem", textAlign: "right" }}>{metric.label}</th>
          </tr>
        </thead>
        <tbody>
          {ranked.map(({ item, rank: position }) => (
            <tr key={item.login} style={{ borderBottom: "1px solid #f3f4f6" }}>
              <td style={{ padding: "0.5rem", color: "#6b7280" }}>{position}</td>
              <td style={{ padding: "0.5rem" }}>
                <Link to={`/members/${encodeURIComponent(item.login)}`}>
                  {item.name || item.login}
                </Link>
              </td>
              <td style={{ padding: "0.5rem", textAlign: "right", fontVariantNumeric: "tabular-nums" }}>
                {metric.format(metric.select(item))}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

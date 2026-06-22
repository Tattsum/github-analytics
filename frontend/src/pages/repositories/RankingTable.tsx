import type { ReactNode } from "react";
import { rank, type SortDirection } from "../../lib/ranking";
import type { MetricOption } from "./metrics";

export interface RankingColumn<T> {
  key: string;
  header: string;
  /** Cell content; receives the item so it can render links/numbers. */
  render: (item: T) => ReactNode;
  /** Right-align numeric columns. */
  numeric?: boolean;
}

export interface RankingTableProps<T> {
  items: readonly T[];
  /** The metric the rows are ranked by (client-side). */
  metric: MetricOption<T>;
  columns: ReadonlyArray<RankingColumn<T>>;
  direction?: SortDirection;
}

const cellStyle = (numeric: boolean | undefined) =>
  ({
    padding: "0.5rem 0.75rem",
    textAlign: numeric ? "right" : "left",
    borderBottom: "1px solid #f3f4f6",
    fontSize: "0.875rem",
  }) as const;

const headStyle = (numeric: boolean | undefined) =>
  ({
    ...cellStyle(numeric),
    borderBottom: "2px solid #e5e7eb",
    color: "#6b7280",
    fontWeight: 600,
  }) as const;

// RankingTable renders a client-side ranking: rows ordered by the chosen
// metric, a 1-based rank column (ties share a rank), then caller-supplied
// columns. The same component serves the repository list and the in-repo
// contributor ranking so both stay visually consistent.
export function RankingTable<T>({ items, metric, columns, direction = "desc" }: RankingTableProps<T>) {
  const ranked = rank(items, metric.value, direction);

  if (ranked.length === 0) {
    return <p>表示するデータがありません。</p>;
  }

  return (
    <table css={{ width: "100%", minWidth: 640, borderCollapse: "collapse" }}>
      <thead>
        <tr>
          <th css={headStyle(true)}>#</th>
          {columns.map((c) => (
            <th key={c.key} css={headStyle(c.numeric)}>
              {c.header}
            </th>
          ))}
        </tr>
      </thead>
      <tbody>
        {ranked.map(({ item, rank: position }, index) => (
          <tr key={index}>
            <td css={cellStyle(true)}>{position}</td>
            {columns.map((c) => (
              <td key={c.key} css={cellStyle(c.numeric)}>
                {c.render(item)}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}

import { useMemo } from "react";
import { EntityTrendOverlay } from "./EntityTrendOverlay";
import type { ComparableSeries, DailyMetricPoint } from "../../lib/comparison";

// One repository contributor with their day-level series within that repository,
// as returned by the `repository` query.
export interface ContributorDailySeries {
  login: string;
  dailyStats: readonly DailyMetricPoint[];
}

export interface MemberTrendComparisonProps {
  contributors: readonly ContributorDailySeries[];
}

// MemberTrendComparison overlays several members' day-level activity within a
// single repository, for one metric, so contributors' trends can be compared in
// that repository's context. The selection/metric/date controls and the chart
// are the shared EntityTrendOverlay; here the entities are members (keyed by
// login) rather than repositories.
export function MemberTrendComparison({ contributors }: MemberTrendComparisonProps) {
  const entities: ComparableSeries[] = useMemo(
    () => contributors.map((c) => ({ key: c.login, daily: c.dailyStats })),
    [contributors],
  );

  return (
    <EntityTrendOverlay
      entities={entities}
      entityLabel="メンバー"
      emptyMessage="このリポジトリに比較できるメンバーの活動データがありません。"
    />
  );
}

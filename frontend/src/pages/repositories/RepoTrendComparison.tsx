import { useMemo, useState } from "react";
import { useQuery } from "urql";
import { graphql } from "../../gql";
import { DateRangePicker, type DateRange } from "../../components/DateRangePicker";
import { OverlayTrendChart, type OverlaySeries } from "../../components/OverlayTrendChart";
import { MetricPicker } from "./MetricPicker";
import type { MetricOption } from "./metrics";
import {
  buildComparisonSeries,
  colorForIndex,
  distinctOwners,
  MAX_SERIES,
  topEntityKeysByMetric,
  type ComparableSeries,
  type DailyMetricKey,
  type DailyMetricPoint,
} from "../../lib/comparison";
import { bounds } from "../../lib/timeSeries";

// Per-repository daily series with owner metadata, summed across members on the
// server. Backs the cross-repository trend overlay and the org-internal filter.
const RepositoryTrendComparisonQuery = graphql(`
  query RepositoryTrendComparison {
    repositoryDailyStats {
      nameWithOwner
      owner
      ownerType
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
    }
  }
`);

// additions/deletions are excluded: their magnitude dwarfs the activity counts
// and would flatten every other line in the overlay.
const trendMetrics: ReadonlyArray<MetricOption<DailyMetricPoint>> = [
  { key: "commitCount", label: "コミット", value: (d) => d.commitCount },
  { key: "prCreated", label: "PR作成", value: (d) => d.prCreated },
  { key: "prMerged", label: "PRマージ", value: (d) => d.prMerged },
  { key: "reviewCount", label: "レビュー", value: (d) => d.reviewCount },
  { key: "issueCount", label: "Issue", value: (d) => d.issueCount },
];

const ALL_OWNERS = "";

// RepoTrendComparison lets the user overlay several repositories' day-level
// activity for one metric, optionally scoped to a single owner (e.g. an
// organization), to compare trends across repositories.
export function RepoTrendComparison() {
  const [{ data, fetching, error }] = useQuery({ query: RepositoryTrendComparisonQuery });

  const [ownerFilter, setOwnerFilter] = useState<string>(ALL_OWNERS);
  const [metricKey, setMetricKey] = useState<DailyMetricKey>("commitCount");
  // null = follow the automatic "top-N by metric" default; a Set = explicit pick.
  const [manualSelection, setManualSelection] = useState<Set<string> | null>(null);
  const [range, setRange] = useState<DateRange | null>(null);

  const repos = useMemo(() => data?.repositoryDailyStats ?? [], [data]);
  const owners = useMemo(() => distinctOwners(repos), [repos]);

  const entities: ComparableSeries[] = useMemo(
    () =>
      repos
        .filter((r) => ownerFilter === ALL_OWNERS || r.owner === ownerFilter)
        .map((r) => ({ key: r.nameWithOwner, daily: r.dailyStats })),
    [repos, ownerFilter],
  );

  // Effective selection: the explicit pick (kept to currently-visible repos) or,
  // when none, the automatic top-N by the chosen metric.
  const effectiveKeys = useMemo(() => {
    if (manualSelection) {
      return entities.map((e) => e.key).filter((k) => manualSelection.has(k));
    }
    return topEntityKeysByMetric(entities, metricKey, MAX_SERIES);
  }, [manualSelection, entities, metricKey]);

  // Only the first MAX_SERIES are plotted so the overlay stays readable.
  const plottedKeys = useMemo(() => effectiveKeys.slice(0, MAX_SERIES), [effectiveKeys]);

  const seriesBounds = useMemo(
    () => bounds(entities.flatMap((e) => e.daily.map((d) => ({ date: d.date })))),
    [entities],
  );

  const effectiveRange: DateRange = range ?? {
    from: seriesBounds?.min ?? "",
    to: seriesBounds?.max ?? "",
    granularity: "month",
  };

  const chartData = useMemo(() => {
    const selected = entities.filter((e) => plottedKeys.includes(e.key));
    return buildComparisonSeries(
      selected,
      metricKey,
      effectiveRange.granularity,
      effectiveRange.from || undefined,
      effectiveRange.to || undefined,
    );
  }, [entities, plottedKeys, metricKey, effectiveRange.from, effectiveRange.to, effectiveRange.granularity]);

  const chartSeries: OverlaySeries[] = useMemo(
    () => plottedKeys.map((key, i) => ({ key, name: key, color: colorForIndex(i) })),
    [plottedKeys],
  );

  function toggleRepo(key: string) {
    setManualSelection((prev) => {
      const base = new Set(prev ?? effectiveKeys);
      if (base.has(key)) {
        base.delete(key);
      } else {
        base.add(key);
      }
      return base;
    });
  }

  if (fetching) {
    return <p>比較データを読み込み中…</p>;
  }
  if (error) {
    return <p>比較データを読み込めませんでした: {error.message}</p>;
  }
  if (repos.length === 0) {
    return <p>比較できるリポジトリがありません。</p>;
  }

  const selectedSet = new Set(effectiveKeys);
  const hiddenCount = effectiveKeys.length - plottedKeys.length;

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
      <div style={{ display: "flex", gap: "1rem", alignItems: "center", flexWrap: "wrap" }}>
        <label style={{ display: "inline-flex", alignItems: "center", gap: "0.5rem" }}>
          <span style={{ fontSize: "0.875rem", color: "#374151" }}>オーナー</span>
          <select
            value={ownerFilter}
            onChange={(e) => {
              setOwnerFilter(e.target.value);
              setManualSelection(null);
            }}
            style={{ padding: "0.375rem 0.5rem", borderRadius: "0.375rem", border: "1px solid #d1d5db", fontSize: "0.875rem" }}
          >
            <option value={ALL_OWNERS}>すべてのオーナー</option>
            {owners.map((o) => (
              <option key={o.owner} value={o.owner}>
                {o.owner}（{o.ownerType || "不明"}・{o.count}）
              </option>
            ))}
          </select>
        </label>
        <MetricPicker
          label="指標"
          options={trendMetrics}
          value={metricKey}
          onChange={(key) => setMetricKey(key as DailyMetricKey)}
        />
      </div>

      <DateRangePicker
        value={effectiveRange}
        onChange={setRange}
        min={seriesBounds?.min}
        max={seriesBounds?.max}
      />

      <fieldset style={{ border: "1px solid #e5e7eb", borderRadius: "0.5rem", padding: "0.75rem", margin: 0 }}>
        <legend style={{ fontSize: "0.875rem", color: "#374151", padding: "0 0.5rem" }}>
          比較するリポジトリ（最大 {MAX_SERIES} 件まで表示）
        </legend>
        <div style={{ display: "flex", flexWrap: "wrap", gap: "0.5rem 1.25rem", maxHeight: "10rem", overflowY: "auto" }}>
          {entities.map((e) => (
            <label key={e.key} style={{ display: "inline-flex", alignItems: "center", gap: "0.375rem", fontSize: "0.85rem" }}>
              <input type="checkbox" checked={selectedSet.has(e.key)} onChange={() => toggleRepo(e.key)} />
              {e.key}
            </label>
          ))}
        </div>
        {hiddenCount > 0 && (
          <p style={{ fontSize: "0.8rem", color: "#b45309", margin: "0.5rem 0 0" }}>
            選択 {effectiveKeys.length} 件のうち上位 {MAX_SERIES} 件のみ表示しています（残り {hiddenCount} 件は非表示）。
          </p>
        )}
      </fieldset>

      <OverlayTrendChart data={chartData} series={chartSeries} emptyMessage="選択したリポジトリの対象期間に活動データはありません。" />
    </div>
  );
}

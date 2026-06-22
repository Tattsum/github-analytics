import { useMemo, useState } from "react";
import { DateRangePicker, type DateRange } from "../../components/DateRangePicker";
import { OverlayTrendChart, type OverlaySeries } from "../../components/OverlayTrendChart";
import { MetricPicker } from "./MetricPicker";
import { trendMetrics } from "./metrics";
import {
  buildComparisonSeries,
  colorForIndex,
  MAX_SERIES,
  topEntityKeysByMetric,
  type ComparableSeries,
  type DailyMetricKey,
} from "../../lib/comparison";
import { bounds } from "../../lib/timeSeries";

export interface EntityTrendOverlayProps {
  /** Entities to overlay; key is the series label (repository nameWithOwner or
   *  member login), daily is its day-level series. */
  entities: readonly ComparableSeries[];
  /** Noun for the selectable entity, e.g. "リポジトリ" / "メンバー". */
  entityLabel: string;
  emptyMessage?: string;
}

// EntityTrendOverlay is the shared control surface for the time-series
// comparison: pick a metric, a date range/granularity and which entities to
// overlay, then render one line per entity. It is axis-agnostic — the
// repository-axis and member-within-repository views both drive it with their
// own ComparableSeries, so the selection/top-N/soft-cap logic lives once here.
export function EntityTrendOverlay({ entities, entityLabel, emptyMessage }: EntityTrendOverlayProps) {
  const [metricKey, setMetricKey] = useState<DailyMetricKey>("commitCount");
  // null = follow the automatic "top-N by metric" default; a Set = explicit pick.
  const [manualSelection, setManualSelection] = useState<Set<string> | null>(null);
  const [range, setRange] = useState<DateRange | null>(null);

  // Effective selection: the explicit pick (kept to currently-visible entities)
  // or, when none, the automatic top-N by the chosen metric.
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

  function toggleEntity(key: string) {
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

  if (entities.length === 0) {
    return <p>{emptyMessage ?? `比較できる${entityLabel}がありません。`}</p>;
  }

  const selectedSet = new Set(effectiveKeys);
  const hiddenCount = effectiveKeys.length - plottedKeys.length;

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
      <MetricPicker
        label="指標"
        options={trendMetrics}
        value={metricKey}
        onChange={(key) => setMetricKey(key as DailyMetricKey)}
      />

      <DateRangePicker
        value={effectiveRange}
        onChange={setRange}
        min={seriesBounds?.min}
        max={seriesBounds?.max}
      />

      <fieldset style={{ border: "1px solid #e5e7eb", borderRadius: "0.5rem", padding: "0.75rem", margin: 0 }}>
        <legend style={{ fontSize: "0.875rem", color: "#374151", padding: "0 0.5rem" }}>
          比較する{entityLabel}（最大 {MAX_SERIES} 件まで表示）
        </legend>
        <div style={{ display: "flex", flexWrap: "wrap", gap: "0.5rem 1.25rem", maxHeight: "10rem", overflowY: "auto" }}>
          {entities.map((e) => (
            <label key={e.key} style={{ display: "inline-flex", alignItems: "center", gap: "0.375rem", fontSize: "0.85rem" }}>
              <input type="checkbox" checked={selectedSet.has(e.key)} onChange={() => toggleEntity(e.key)} />
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

      <OverlayTrendChart
        data={chartData}
        series={chartSeries}
        emptyMessage={`選択した${entityLabel}の対象期間に活動データはありません。`}
      />
    </div>
  );
}

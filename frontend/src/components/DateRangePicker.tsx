import type { Granularity } from "../lib/timeSeries";

export interface DateRange {
  from: string;
  to: string;
  granularity: Granularity;
}

export interface DateRangePickerProps {
  value: DateRange;
  onChange: (next: DateRange) => void;
  /** Earliest/latest selectable dates (typically the series bounds). */
  min?: string;
  max?: string;
}

const GRANULARITIES: ReadonlyArray<{ value: Granularity; label: string }> = [
  { value: "day", label: "日" },
  { value: "week", label: "週" },
  { value: "month", label: "月" },
];

const fieldStyle = { display: "flex", flexDirection: "column", gap: "0.25rem" } as const;
const labelStyle = { fontSize: "0.95rem", color: "#6b7280" } as const;

export function DateRangePicker({ value, onChange, min, max }: DateRangePickerProps) {
  return (
    <div css={{ display: "flex", gap: "1rem", alignItems: "flex-end", flexWrap: "wrap" }}>
      <label css={fieldStyle}>
        <span css={labelStyle}>開始日</span>
        <input
          type="date"
          value={value.from}
          min={min}
          max={value.to || max}
          onChange={(e) => onChange({ ...value, from: e.target.value })}
        />
      </label>
      <label css={fieldStyle}>
        <span css={labelStyle}>終了日</span>
        <input
          type="date"
          value={value.to}
          min={value.from || min}
          max={max}
          onChange={(e) => onChange({ ...value, to: e.target.value })}
        />
      </label>
      <label css={fieldStyle}>
        <span css={labelStyle}>粒度</span>
        <select
          value={value.granularity}
          onChange={(e) => onChange({ ...value, granularity: e.target.value as Granularity })}
        >
          {GRANULARITIES.map((g) => (
            <option key={g.value} value={g.value}>
              {g.label}
            </option>
          ))}
        </select>
      </label>
    </div>
  );
}

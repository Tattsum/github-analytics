import type { MetricOption } from "./metrics";

export interface MetricPickerProps<T> {
  label: string;
  options: ReadonlyArray<MetricOption<T>>;
  value: string;
  onChange: (key: string) => void;
}

// MetricPicker is a labelled <select> for choosing which metric drives the
// client-side ranking. Generic over the ranked item type so it serves both the
// repository list and the in-repo contributor ranking.
export function MetricPicker<T>({ label, options, value, onChange }: MetricPickerProps<T>) {
  return (
    <label css={{ display: "inline-flex", alignItems: "center", gap: "0.5rem" }}>
      <span css={{ fontSize: "0.875rem", color: "#374151" }}>{label}</span>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        css={{
          padding: "0.375rem 0.5rem",
          borderRadius: "0.375rem",
          border: "1px solid #d1d5db",
          fontSize: "0.875rem",
        }}
      >
        {options.map((o) => (
          <option key={o.key} value={o.key}>
            {o.label}
          </option>
        ))}
      </select>
    </label>
  );
}

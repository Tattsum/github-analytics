import { describe, expect, it } from "vitest";
import { bounds, bucketBy, filterByDateRange, prepareSeries, type TimeSeriesPoint } from "./timeSeries";

// Spread across a Sun/Mon week boundary and two months so bucketing edges are
// exercised. Concrete non-zero values so a dropped/zeroed metric is detectable.
const series: TimeSeriesPoint[] = [
  { date: "2024-01-07", commitCount: 3, reviewCount: 1 }, // Sunday
  { date: "2024-01-08", commitCount: 5, reviewCount: 2 }, // Monday (new ISO week)
  { date: "2024-01-09", commitCount: 4, reviewCount: 0 }, // Tuesday
  { date: "2024-02-01", commitCount: 7, reviewCount: 6 },
];

describe("filterByDateRange", () => {
  it("includes both endpoints (inclusive range)", () => {
    const got = filterByDateRange(series, "2024-01-08", "2024-02-01");
    expect(got.map((p) => p.date)).toEqual(["2024-01-08", "2024-01-09", "2024-02-01"]);
  });

  it("treats undefined bounds as open", () => {
    expect(filterByDateRange(series, undefined, "2024-01-08").map((p) => p.date)).toEqual([
      "2024-01-07",
      "2024-01-08",
    ]);
    expect(filterByDateRange(series, "2024-02-01", undefined).map((p) => p.date)).toEqual(["2024-02-01"]);
  });

  it("returns empty when no point falls in range", () => {
    expect(filterByDateRange(series, "2025-01-01", "2025-12-31")).toEqual([]);
  });

  it("does not mutate the input", () => {
    const input = [...series];
    filterByDateRange(input, "2024-01-08", "2024-01-09");
    expect(input).toHaveLength(4);
  });
});

describe("bucketBy", () => {
  it("returns daily points summed and sorted for day granularity", () => {
    const got = bucketBy(
      [
        { date: "2024-01-09", commitCount: 1 },
        { date: "2024-01-08", commitCount: 5 },
        { date: "2024-01-08", commitCount: 2 },
      ],
      "day",
    );
    expect(got).toEqual([
      { date: "2024-01-08", commitCount: 7 },
      { date: "2024-01-09", commitCount: 1 },
    ]);
  });

  it("groups into ISO weeks starting Monday", () => {
    const got = bucketBy(series, "week");
    // 2024-01-07 (Sun) is its own week; Mon 01-08 + Tue 01-09 share a week.
    expect(got).toEqual([
      { date: "2024-01-01", commitCount: 3, reviewCount: 1 },
      { date: "2024-01-08", commitCount: 9, reviewCount: 2 },
      { date: "2024-01-29", commitCount: 7, reviewCount: 6 },
    ]);
  });

  it("groups into calendar months labelled by the first of the month", () => {
    const got = bucketBy(series, "month");
    expect(got).toEqual([
      { date: "2024-01-01", commitCount: 12, reviewCount: 3 },
      { date: "2024-02-01", commitCount: 7, reviewCount: 6 },
    ]);
  });

  it("ignores non-numeric fields such as a GraphQL __typename", () => {
    const got = bucketBy([{ date: "2024-01-08", __typename: "DailyStatistics", commitCount: 5 }], "month");
    expect(got).toEqual([{ date: "2024-01-01", commitCount: 5 }]);
  });
});

describe("bounds", () => {
  it("returns the min and max dates", () => {
    expect(bounds(series)).toEqual({ min: "2024-01-07", max: "2024-02-01" });
  });

  it("returns undefined for an empty series", () => {
    expect(bounds([])).toBeUndefined();
  });
});

describe("prepareSeries", () => {
  it("filters to the range then buckets by month", () => {
    const got = prepareSeries(series, "month", "2024-01-08", "2024-02-01");
    expect(got).toEqual([
      { date: "2024-01-01", commitCount: 9, reviewCount: 2 },
      { date: "2024-02-01", commitCount: 7, reviewCount: 6 },
    ]);
  });
});

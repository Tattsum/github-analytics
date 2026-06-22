import { describe, expect, it } from "vitest";
import {
  buildComparisonSeries,
  colorForIndex,
  distinctOwners,
  MAX_SERIES,
  SERIES_PALETTE,
  topEntityKeysByMetric,
  totalForMetric,
  type ComparableSeries,
} from "./comparison";

// Two repositories with an overlapping date and an extra day, so merging by date
// and per-entity columns are both exercised. Non-zero distinct values so a
// dropped/zeroed metric would be detectable.
const foo: ComparableSeries = {
  key: "Tattsum/foo",
  daily: [
    { date: "2024-01-08", commitCount: 5, prCreated: 2, prMerged: 1, reviewCount: 3, issueCount: 1, totalAdditions: 100, totalDeletions: 40 },
    { date: "2024-01-09", commitCount: 4, prCreated: 1, prMerged: 1, reviewCount: 2, issueCount: 0, totalAdditions: 80, totalDeletions: 10 },
  ],
};
const bar: ComparableSeries = {
  key: "Tattsum/bar",
  daily: [
    { date: "2024-01-08", commitCount: 9, prCreated: 4, prMerged: 3, reviewCount: 1, issueCount: 2, totalAdditions: 300, totalDeletions: 90 },
  ],
};

describe("buildComparisonSeries", () => {
  it("merges entities into one column per key keyed by date", () => {
    const got = buildComparisonSeries([foo, bar], "commitCount", "day");
    expect(got).toEqual([
      { date: "2024-01-08", "Tattsum/foo": 5, "Tattsum/bar": 9 },
      { date: "2024-01-09", "Tattsum/foo": 4 },
    ]);
  });

  it("projects the chosen metric, not commitCount", () => {
    const got = buildComparisonSeries([bar], "reviewCount", "day");
    expect(got).toEqual([{ date: "2024-01-08", "Tattsum/bar": 1 }]);
  });

  it("buckets by month, summing each entity column within the bucket", () => {
    const got = buildComparisonSeries([foo], "commitCount", "month");
    expect(got).toEqual([{ date: "2024-01-01", "Tattsum/foo": 9 }]);
  });

  it("applies the inclusive date-range filter before bucketing", () => {
    const got = buildComparisonSeries([foo], "commitCount", "day", "2024-01-09", "2024-01-09");
    expect(got).toEqual([{ date: "2024-01-09", "Tattsum/foo": 4 }]);
  });
});

describe("totalForMetric", () => {
  it("sums the metric across the whole series", () => {
    expect(totalForMetric(foo, "commitCount")).toBe(9);
    expect(totalForMetric(foo, "totalAdditions")).toBe(180);
  });
});

describe("topEntityKeysByMetric", () => {
  it("returns the top-n keys by metric total, descending", () => {
    expect(topEntityKeysByMetric([foo, bar], "commitCount", 1)).toEqual(["Tattsum/bar"]);
  });

  it("breaks ties by key ascending for a deterministic default", () => {
    const a: ComparableSeries = { key: "z/repo", daily: [{ date: "2024-01-01", commitCount: 5, prCreated: 0, prMerged: 0, reviewCount: 0, issueCount: 0, totalAdditions: 0, totalDeletions: 0 }] };
    const b: ComparableSeries = { key: "a/repo", daily: [{ date: "2024-01-01", commitCount: 5, prCreated: 0, prMerged: 0, reviewCount: 0, issueCount: 0, totalAdditions: 0, totalDeletions: 0 }] };
    expect(topEntityKeysByMetric([a, b], "commitCount", 2)).toEqual(["a/repo", "z/repo"]);
  });

  it("caps at n", () => {
    expect(topEntityKeysByMetric([foo, bar], "commitCount", 1)).toHaveLength(1);
  });
});

describe("distinctOwners", () => {
  it("dedupes owners, counts repos, fills owner type, sorts by owner", () => {
    const got = distinctOwners([
      { owner: "myorg", ownerType: "" },
      { owner: "myorg", ownerType: "Organization" },
      { owner: "alice", ownerType: "User" },
    ]);
    expect(got).toEqual([
      { owner: "alice", ownerType: "User", count: 1 },
      { owner: "myorg", ownerType: "Organization", count: 2 },
    ]);
  });
});

describe("color palette", () => {
  it("recycles colors past the palette length", () => {
    expect(colorForIndex(0)).toBe(SERIES_PALETTE[0]);
    expect(colorForIndex(MAX_SERIES)).toBe(SERIES_PALETTE[0]);
  });
});

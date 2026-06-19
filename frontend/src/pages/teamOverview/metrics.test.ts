import { describe, expect, it } from "vitest";
import type { MemberStats } from "../../gql/graphql";
import { METRICS, changedLines, metricByKey, type MetricKey } from "./metrics";

function member(overrides: Partial<MemberStats>): MemberStats {
  return {
    __typename: "MemberStats",
    login: "octocat",
    name: "Octo Cat",
    totalCommits: 42,
    totalPRCreated: 7,
    totalPRMerged: 5,
    totalIssues: 3,
    totalReviews: 11,
    totalAdditions: 100,
    totalDeletions: 40,
    prToReviewRatio: 0.63,
    ...overrides,
  };
}

describe("changedLines", () => {
  it("sums additions and deletions", () => {
    expect(changedLines(member({ totalAdditions: 100, totalDeletions: 40 }))).toBe(140);
  });
});

describe("metric selectors", () => {
  const cases: Array<{ key: MetricKey; input: Partial<MemberStats>; want: number }> = [
    { key: "totalCommits", input: { totalCommits: 42 }, want: 42 },
    { key: "totalPRCreated", input: { totalPRCreated: 7 }, want: 7 },
    { key: "totalPRMerged", input: { totalPRMerged: 5 }, want: 5 },
    { key: "totalReviews", input: { totalReviews: 11 }, want: 11 },
    { key: "totalChangedLines", input: { totalAdditions: 100, totalDeletions: 40 }, want: 140 },
    { key: "prToReviewRatio", input: { prToReviewRatio: 0.63 }, want: 0.63 },
  ];

  for (const tc of cases) {
    it(`selects ${tc.key}`, () => {
      expect(metricByKey(tc.key).select(member(tc.input))).toBe(tc.want);
    });
  }
});

describe("metric formatting", () => {
  it("formats counts as integers", () => {
    expect(metricByKey("totalCommits").format(1500)).toBe(new Intl.NumberFormat().format(1500));
  });

  it("formats the ratio with one decimal", () => {
    expect(metricByKey("prToReviewRatio").format(0.6349)).toBe("0.6");
  });
});

describe("metricByKey", () => {
  it("returns the matching definition", () => {
    expect(metricByKey("totalReviews").key).toBe("totalReviews");
  });

  it("falls back to the first metric for an unknown key", () => {
    expect(metricByKey("nope" as MetricKey).key).toBe(METRICS[0].key);
  });
});

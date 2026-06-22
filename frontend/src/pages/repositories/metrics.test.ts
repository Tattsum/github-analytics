import { describe, expect, it } from "vitest";
import {
  contributorMetrics,
  findMetric,
  repositoryMetrics,
  trendMetrics,
  type ContributorLike,
  type RepoStatsLike,
} from "./metrics";
import type { DailyMetricPoint } from "../../lib/comparison";

const repo: RepoStatsLike = {
  nameWithOwner: "octo/example",
  contributorCount: 7,
  total: {
    commits: 42,
    prCreated: 31,
    prMerged: 28,
    issues: 13,
    reviews: 25,
    additions: 900,
    deletions: 120,
  },
};

const contributor: ContributorLike = {
  login: "alice",
  commitCount: 17,
  prCreated: 9,
  reviewCount: 23,
  additions: 410,
  deletions: 55,
};

describe("repositoryMetrics value selectors", () => {
  const cases: Array<{ key: string; want: number }> = [
    { key: "commits", want: 42 },
    { key: "prCreated", want: 31 },
    { key: "prMerged", want: 28 },
    { key: "issues", want: 13 },
    { key: "reviews", want: 25 },
    { key: "additions", want: 900 },
    { key: "deletions", want: 120 },
    { key: "contributors", want: 7 },
  ];

  for (const tc of cases) {
    it(`reads ${tc.key} from total/contributorCount`, () => {
      const metric = findMetric(repositoryMetrics, tc.key);
      expect(metric.key).toBe(tc.key);
      expect(metric.value(repo)).toBe(tc.want);
    });
  }
});

describe("contributorMetrics value selectors", () => {
  const cases: Array<{ key: string; want: number }> = [
    { key: "commitCount", want: 17 },
    { key: "prCreated", want: 9 },
    { key: "reviewCount", want: 23 },
    { key: "additions", want: 410 },
    { key: "deletions", want: 55 },
  ];

  for (const tc of cases) {
    it(`reads ${tc.key} from contributor`, () => {
      const metric = findMetric(contributorMetrics, tc.key);
      expect(metric.key).toBe(tc.key);
      expect(metric.value(contributor)).toBe(tc.want);
    });
  }
});

describe("trendMetrics", () => {
  const point: DailyMetricPoint = {
    date: "2024-01-08",
    commitCount: 5,
    prCreated: 2,
    prMerged: 1,
    reviewCount: 3,
    issueCount: 4,
    totalAdditions: 100,
    totalDeletions: 40,
  };

  it("excludes additions/deletions from the overlay metrics", () => {
    expect(trendMetrics.map((m) => m.key)).toEqual([
      "commitCount",
      "prCreated",
      "prMerged",
      "reviewCount",
      "issueCount",
    ]);
  });

  it("reads each metric off a daily point", () => {
    expect(findMetric(trendMetrics, "prMerged").value(point)).toBe(1);
    expect(findMetric(trendMetrics, "issueCount").value(point)).toBe(4);
  });
});

describe("findMetric", () => {
  it("returns the matching option by key", () => {
    expect(findMetric(repositoryMetrics, "reviews").label).toBe("レビュー");
  });

  it("falls back to the first option for an unknown key", () => {
    expect(findMetric(repositoryMetrics, "does-not-exist")).toBe(repositoryMetrics[0]);
  });
});

import { describe, expect, it } from "vitest";
import {
  contributorMetrics,
  findMetric,
  repositoryMetrics,
  type ContributorLike,
  type RepoStatsLike,
} from "./metrics";

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

describe("findMetric", () => {
  it("returns the matching option by key", () => {
    expect(findMetric(repositoryMetrics, "reviews").label).toBe("レビュー");
  });

  it("falls back to the first option for an unknown key", () => {
    expect(findMetric(repositoryMetrics, "does-not-exist")).toBe(repositoryMetrics[0]);
  });
});

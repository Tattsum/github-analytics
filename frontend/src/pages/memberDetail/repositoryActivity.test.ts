import { describe, expect, it } from "vitest";
import { activitySpanYears, byCommitsDesc } from "./repositoryActivity";

describe("activitySpanYears", () => {
  const cases: Array<{
    name: string;
    first: string;
    last: string;
    wantApprox: number;
  }> = [
    { name: "one full year", first: "2022-01-01T00:00:00Z", last: "2023-01-01T00:00:00Z", wantApprox: 1 },
    { name: "two years", first: "2021-01-01T00:00:00Z", last: "2023-01-01T00:00:00Z", wantApprox: 2 },
    { name: "same instant is zero", first: "2023-01-01T00:00:00Z", last: "2023-01-01T00:00:00Z", wantApprox: 0 },
    { name: "reversed range clamps to zero", first: "2023-01-01T00:00:00Z", last: "2021-01-01T00:00:00Z", wantApprox: 0 },
    { name: "unparseable dates yield zero", first: "not-a-date", last: "2023-01-01T00:00:00Z", wantApprox: 0 },
  ];

  for (const tc of cases) {
    it(tc.name, () => {
      const got = activitySpanYears({ firstActivity: tc.first, lastActivity: tc.last });
      expect(got).toBeCloseTo(tc.wantApprox, 1);
    });
  }
});

describe("byCommitsDesc", () => {
  it("orders by commit count descending", () => {
    const repos = [
      { repository: "a", commitCount: 10 },
      { repository: "b", commitCount: 90 },
      { repository: "c", commitCount: 50 },
    ];
    const sorted = [...repos].sort(byCommitsDesc);
    expect(sorted.map((r) => r.repository)).toEqual(["b", "c", "a"]);
  });

  it("breaks ties on repository name for a stable order", () => {
    const repos = [
      { repository: "zeta", commitCount: 42 },
      { repository: "alpha", commitCount: 42 },
    ];
    const sorted = [...repos].sort(byCommitsDesc);
    expect(sorted.map((r) => r.repository)).toEqual(["alpha", "zeta"]);
  });
});

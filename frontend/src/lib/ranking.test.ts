import { describe, expect, it } from "vitest";
import { rank, sortBy } from "./ranking";

interface Member {
  login: string;
  commits: number;
}

const members: Member[] = [
  { login: "alice", commits: 50 },
  { login: "bob", commits: 90 },
  { login: "carol", commits: 90 },
  { login: "dave", commits: 10 },
];

describe("sortBy", () => {
  const cases: Array<{
    name: string;
    direction: "asc" | "desc" | undefined;
    wantLogins: string[];
  }> = [
    { name: "descending by default", direction: undefined, wantLogins: ["bob", "carol", "alice", "dave"] },
    { name: "ascending when asked", direction: "asc", wantLogins: ["dave", "alice", "bob", "carol"] },
    { name: "explicit descending", direction: "desc", wantLogins: ["bob", "carol", "alice", "dave"] },
  ];

  for (const tc of cases) {
    it(tc.name, () => {
      const sorted =
        tc.direction === undefined
          ? sortBy(members, "commits")
          : sortBy(members, "commits", tc.direction);
      expect(sorted.map((m) => m.login)).toEqual(tc.wantLogins);
    });
  }

  it("does not mutate the input", () => {
    const input = [...members];
    sortBy(input, "commits");
    expect(input.map((m) => m.login)).toEqual(["alice", "bob", "carol", "dave"]);
  });

  it("accepts a selector function", () => {
    const sorted = sortBy(members, (m) => m.commits, "asc");
    expect(sorted[0]?.login).toBe("dave");
  });
});

describe("rank", () => {
  it("assigns 1-based ranks and shares rank on ties (competition ranking)", () => {
    const ranked = rank(members, "commits");
    expect(ranked).toEqual([
      { item: { login: "bob", commits: 90 }, rank: 1 },
      { item: { login: "carol", commits: 90 }, rank: 1 },
      { item: { login: "alice", commits: 50 }, rank: 3 },
      { item: { login: "dave", commits: 10 }, rank: 4 },
    ]);
  });

  it("ranks ascending when asked", () => {
    const ranked = rank(members, "commits", "asc");
    expect(ranked.map((r) => [r.item.login, r.rank])).toEqual([
      ["dave", 1],
      ["alice", 2],
      ["bob", 3],
      ["carol", 3],
    ]);
  });

  it("handles an empty list", () => {
    expect(rank<Member>([], "commits")).toEqual([]);
  });
});

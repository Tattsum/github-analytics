// Fixed GraphQL responses keyed by operation name. The Go backend / Postgres
// are never started for visual tests; page.route intercepts /query and returns
// these so snapshots are deterministic. Values are arbitrary but stable.

const dailyStats = [
  { date: "2024-01-15", commitCount: 12, prCreated: 4, prMerged: 3, reviewCount: 6, issueCount: 2, totalAdditions: 540, totalDeletions: 120 },
  { date: "2024-02-15", commitCount: 18, prCreated: 6, prMerged: 5, reviewCount: 9, issueCount: 3, totalAdditions: 880, totalDeletions: 210 },
  { date: "2024-03-15", commitCount: 9, prCreated: 3, prMerged: 2, reviewCount: 4, issueCount: 1, totalAdditions: 320, totalDeletions: 90 },
  { date: "2024-04-15", commitCount: 22, prCreated: 8, prMerged: 7, reviewCount: 12, issueCount: 4, totalAdditions: 1200, totalDeletions: 340 },
  { date: "2024-05-15", commitCount: 15, prCreated: 5, prMerged: 4, reviewCount: 7, issueCount: 2, totalAdditions: 610, totalDeletions: 150 },
  { date: "2024-06-15", commitCount: 27, prCreated: 9, prMerged: 8, reviewCount: 14, issueCount: 5, totalAdditions: 1480, totalDeletions: 420 },
];

// Derive a distinct-but-stable series per entity so the overlay chart shows
// separated lines. Factors stay >= 0.6 so no metric rounds down to a zero value
// (which would be indistinguishable from a missing field).
const scaleDaily = (factor: number) =>
  dailyStats.map((d) => ({
    date: d.date,
    commitCount: Math.round(d.commitCount * factor),
    prCreated: Math.round(d.prCreated * factor),
    prMerged: Math.round(d.prMerged * factor),
    reviewCount: Math.round(d.reviewCount * factor),
    issueCount: Math.round(d.issueCount * factor),
    totalAdditions: Math.round(d.totalAdditions * factor),
    totalDeletions: Math.round(d.totalDeletions * factor),
  }));

const member = {
  login: "octocat",
  name: "Octo Cat",
  totalCommits: 1024,
  totalPRCreated: 312,
  totalPRMerged: 287,
  totalIssues: 96,
  totalReviews: 540,
  totalAdditions: 84210,
  totalDeletions: 31044,
  prToReviewRatio: 0.58,
};

const yearlyStats = [
  { year: 2021, commitCount: 120, prCreated: 40, prMerged: 36, reviewCount: 55, issueCount: 12, totalAdditions: 9800, totalDeletions: 3400 },
  { year: 2022, commitCount: 240, prCreated: 78, prMerged: 70, reviewCount: 130, issueCount: 24, totalAdditions: 21000, totalDeletions: 7600 },
  { year: 2023, commitCount: 360, prCreated: 110, prMerged: 102, reviewCount: 200, issueCount: 33, totalAdditions: 31000, totalDeletions: 11000 },
  { year: 2024, commitCount: 304, prCreated: 84, prMerged: 79, reviewCount: 155, issueCount: 27, totalAdditions: 22410, totalDeletions: 9044 },
];

const repositoryActivity = (repository: string, commitCount: number) => ({
  repository,
  commitCount,
  prCount: Math.round(commitCount / 3),
  reviewCount: Math.round(commitCount / 2),
  issueCount: Math.round(commitCount / 8),
  totalAdditions: commitCount * 40,
  totalDeletions: commitCount * 14,
  firstActivity: "2021-03-01",
  lastActivity: "2024-06-30",
});

export const graphqlResponses: Record<string, unknown> = {
  TeamOverviewSummary: {
    teamSummary: {
      memberCount: 12,
      repositoryCount: 34,
      totalCommits: 18420,
      totalPRCreated: 4210,
      totalPRMerged: 3980,
      totalIssues: 1240,
      totalReviews: 6710,
      totalAdditions: 1284200,
      totalDeletions: 498300,
    },
  },
  TeamOverviewDailyStats: { teamDailyStats: dailyStats },
  TeamOverviewMembers: {
    members: [
      { login: "octocat", name: "Octo Cat", totalCommits: 1024, totalPRCreated: 312, totalPRMerged: 287, totalIssues: 96, totalReviews: 540, totalAdditions: 84210, totalDeletions: 31044, prToReviewRatio: 0.58 },
      { login: "hubot", name: "Hu Bot", totalCommits: 880, totalPRCreated: 260, totalPRMerged: 240, totalIssues: 70, totalReviews: 410, totalAdditions: 64000, totalDeletions: 22000, prToReviewRatio: 0.63 },
      { login: "monalisa", name: "Mona Lisa", totalCommits: 640, totalPRCreated: 190, totalPRMerged: 178, totalIssues: 54, totalReviews: 300, totalAdditions: 48000, totalDeletions: 17000, prToReviewRatio: 0.63 },
    ],
  },
  Repositories: {
    repositories: [
      { nameWithOwner: "acme/web", contributorCount: 18, total: { commits: 5200, prCreated: 1400, prMerged: 1320, issues: 420, reviews: 2100, additions: 420000, deletions: 160000 } },
      { nameWithOwner: "acme/api", contributorCount: 12, total: { commits: 3800, prCreated: 980, prMerged: 920, issues: 280, reviews: 1500, additions: 310000, deletions: 120000 } },
      { nameWithOwner: "acme/infra", contributorCount: 7, total: { commits: 1900, prCreated: 520, prMerged: 480, issues: 140, reviews: 760, additions: 150000, deletions: 58000 } },
    ],
  },
  Repository: {
    repository: {
      nameWithOwner: "acme/web",
      contributorCount: 18,
      total: { commits: 5200, prCreated: 1400, prMerged: 1320, issues: 420, reviews: 2100, additions: 420000, deletions: 160000 },
      contributors: [
        { login: "octocat", commitCount: 1200, prCreated: 340, reviewCount: 560, additions: 98000, deletions: 36000, dailyStats: scaleDaily(1.0) },
        { login: "hubot", commitCount: 900, prCreated: 260, reviewCount: 410, additions: 72000, deletions: 24000, dailyStats: scaleDaily(0.8) },
        { login: "monalisa", commitCount: 640, prCreated: 180, reviewCount: 300, additions: 51000, deletions: 18000, dailyStats: scaleDaily(0.6) },
      ],
    },
  },
  MemberDetail: {
    member: {
      ...member,
      firstActivityYear: 2021,
      peakActivityYear: 2023,
      peakActivityCommits: 360,
      yearlyStats,
      dailyStats,
      topRepositories: [
        repositoryActivity("acme/web", 480),
        repositoryActivity("acme/api", 300),
        repositoryActivity("acme/infra", 120),
      ],
      longTermRepositories: [
        repositoryActivity("acme/web", 480),
        repositoryActivity("acme/api", 300),
      ],
      roleTransition: [
        { year: 2021, description: "author", prCreated: 40, reviewCount: 20, ratio: 0.5 },
        { year: 2022, description: "author", prCreated: 78, reviewCount: 60, ratio: 0.77 },
        { year: 2023, description: "reviewer", prCreated: 110, reviewCount: 200, ratio: 1.82 },
        { year: 2024, description: "reviewer", prCreated: 84, reviewCount: 155, ratio: 1.85 },
      ],
    },
  },
  RepositoryTrendComparison: {
    repositoryDailyStats: [
      { nameWithOwner: "acme/web", owner: "acme", ownerType: "Organization", dailyStats: scaleDaily(1.0) },
      { nameWithOwner: "acme/api", owner: "acme", ownerType: "Organization", dailyStats: scaleDaily(0.8) },
      { nameWithOwner: "acme/infra", owner: "acme", ownerType: "Organization", dailyStats: scaleDaily(0.6) },
      { nameWithOwner: "octocat/sandbox", owner: "octocat", ownerType: "User", dailyStats: scaleDaily(0.7) },
    ],
  },
};

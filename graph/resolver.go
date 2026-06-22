// Package graph contains the gqlgen GraphQL resolvers and the mapping helpers
// that translate application/domain values into the generated GraphQL models.
package graph

import (
	"sort"
	"time"

	"github.com/Tattsum/github-analytics/application"
	"github.com/Tattsum/github-analytics/domain"
	"github.com/Tattsum/github-analytics/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

// Resolver is the GraphQL resolver root. It holds the application-level
// SnapshotReader and maps its results to generated GraphQL models. It performs
// no business logic and has no knowledge of the persistence layer.
type Resolver struct {
	reader application.SnapshotReader
}

// NewResolver constructs a Resolver backed by the given SnapshotReader.
// The concrete reader (ent/Postgres) is wired by the composition root.
func NewResolver(reader application.SnapshotReader) *Resolver {
	return &Resolver{reader: reader}
}

// toMemberStats maps an application.MemberStats to its GraphQL model.
func toMemberStats(m *application.MemberStats) *model.MemberStats {
	return &model.MemberStats{
		Login:           m.Login,
		Name:            m.Name,
		TotalCommits:    m.TotalCommits,
		TotalPRCreated:  m.TotalPRCreated,
		TotalPRMerged:   m.TotalPRMerged,
		TotalIssues:     m.TotalIssues,
		TotalReviews:    m.TotalReviews,
		TotalAdditions:  m.TotalAdditions,
		TotalDeletions:  m.TotalDeletions,
		PrToReviewRatio: m.PRToReviewRatio,
	}
}

// toTeamSummary maps an application.TeamSummary to its GraphQL model.
func toTeamSummary(s *application.TeamSummary) *model.TeamSummary {
	return &model.TeamSummary{
		MemberCount:     s.MemberCount,
		RepositoryCount: s.RepositoryCount,
		TotalCommits:    s.TotalCommits,
		TotalPRCreated:  s.TotalPRCreated,
		TotalPRMerged:   s.TotalPRMerged,
		TotalIssues:     s.TotalIssues,
		TotalReviews:    s.TotalReviews,
		TotalAdditions:  s.TotalAdditions,
		TotalDeletions:  s.TotalDeletions,
	}
}

// toRepositoryStats maps an application.RepositoryStats to its GraphQL model.
func toRepositoryStats(r *application.RepositoryStats) *model.RepositoryStats {
	contributors := make([]*model.RepositoryContributor, 0, len(r.Contributors))
	for _, c := range r.Contributors {
		contributors = append(contributors, &model.RepositoryContributor{
			Login:       c.Login,
			CommitCount: c.CommitCount,
			PrCreated:   c.PRCreated,
			ReviewCount: c.ReviewCount,
			Additions:   c.Additions,
			Deletions:   c.Deletions,
		})
	}
	return &model.RepositoryStats{
		NameWithOwner: r.NameWithOwner,
		Total: &model.RepositoryTotals{
			Commits:   r.TotalCommits,
			PrCreated: r.TotalPRCreated,
			PrMerged:  r.TotalPRMerged,
			Issues:    r.TotalIssues,
			Reviews:   r.TotalReviews,
			Additions: r.TotalAdditions,
			Deletions: r.TotalDeletions,
		},
		ContributorCount: r.ContributorCount,
		Contributors:     contributors,
	}
}

// toUserStatistics maps a domain.UserStatistics to its GraphQL model.
// The yearly stats map is flattened into a slice sorted by year ascending so
// the frontend receives a stable, chronological trend.
func toUserStatistics(s *domain.UserStatistics) *model.UserStatistics {
	out := &model.UserStatistics{
		TotalCommits:         s.TotalCommits,
		TotalPRCreated:       s.TotalPRCreated,
		TotalPRMerged:        s.TotalPRMerged,
		TotalIssues:          s.TotalIssues,
		TotalReviews:         s.TotalReviews,
		TotalAdditions:       s.TotalAdditions,
		TotalDeletions:       s.TotalDeletions,
		PrToReviewRatio:      s.PRToReviewRatio,
		FirstActivityYear:    s.FirstActivityYear,
		PeakActivityYear:     s.PeakActivityYear,
		PeakActivityCommits:  s.PeakActivityCommits,
		YearlyStats:          toYearlyStatistics(s.YearlyStats),
		DailyStats:           toDailyStatistics(s.DailyStats),
		TopRepositories:      toRepositoryActivities(s.TopRepositories),
		LongTermRepositories: toRepositoryActivities(s.LongTermRepositories),
		RoleTransition:       toRoleTransitions(s.RoleTransition),
	}
	if s.User != nil {
		out.Login = s.User.Login
		out.Name = s.User.Name
	}
	return out
}

// toYearlyStatistics flattens a year-keyed map into a slice sorted by year ascending.
func toYearlyStatistics(stats map[int]*domain.YearlyStatistics) []*model.YearlyStatistics {
	years := make([]int, 0, len(stats))
	for year := range stats {
		years = append(years, year)
	}
	sort.Ints(years)

	out := make([]*model.YearlyStatistics, 0, len(years))
	for _, year := range years {
		y := stats[year]
		out = append(out, &model.YearlyStatistics{
			Year:           y.Year,
			CommitCount:    y.CommitCount,
			PrCreated:      y.PRCreated,
			PrMerged:       y.PRMerged,
			IssueCount:     y.IssueCount,
			ReviewCount:    y.ReviewCount,
			TotalAdditions: y.TotalAdditions,
			TotalDeletions: y.TotalDeletions,
		})
	}
	return out
}

// toDailyStatistics maps a date-keyed map into a slice sorted by date ascending
// so the frontend receives a stable, chronological time series.
func toDailyStatistics(stats map[string]*domain.DailyStatistics) []*model.DailyStatistics {
	dates := make([]string, 0, len(stats))
	for date := range stats {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	out := make([]*model.DailyStatistics, 0, len(dates))
	for _, date := range dates {
		out = append(out, toDailyStatistic(stats[date]))
	}
	return out
}

// toDailyStatistic maps a single domain.DailyStatistics to its GraphQL model.
func toDailyStatistic(d *domain.DailyStatistics) *model.DailyStatistics {
	return &model.DailyStatistics{
		Date:           d.Date,
		CommitCount:    d.CommitCount,
		PrCreated:      d.PRCreated,
		PrMerged:       d.PRMerged,
		IssueCount:     d.IssueCount,
		ReviewCount:    d.ReviewCount,
		TotalAdditions: d.TotalAdditions,
		TotalDeletions: d.TotalDeletions,
	}
}

// toRepositoryActivities maps domain repository activities to their GraphQL models.
func toRepositoryActivities(activities []*domain.RepositoryActivity) []*model.RepositoryActivity {
	out := make([]*model.RepositoryActivity, 0, len(activities))
	for _, a := range activities {
		out = append(out, &model.RepositoryActivity{
			Repository:     a.Repository,
			CommitCount:    a.CommitCount,
			PrCount:        a.PRCount,
			IssueCount:     a.IssueCount,
			ReviewCount:    a.ReviewCount,
			TotalAdditions: a.TotalAdditions,
			TotalDeletions: a.TotalDeletions,
			FirstActivity:  a.FirstActivity.Format(time.RFC3339),
			LastActivity:   a.LastActivity.Format(time.RFC3339),
		})
	}
	return out
}

// toRoleTransitions maps domain role-transition points to their GraphQL models.
func toRoleTransitions(points []domain.RoleTransitionPoint) []*model.RoleTransitionPoint {
	out := make([]*model.RoleTransitionPoint, 0, len(points))
	for _, p := range points {
		out = append(out, &model.RoleTransitionPoint{
			Year:        p.Year,
			PrCreated:   p.PRCreated,
			ReviewCount: p.ReviewCount,
			Ratio:       p.Ratio,
			Description: p.Description,
		})
	}
	return out
}

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemberRepoDayStat holds the per-member, per-repository, per-day metrics for a
// single snapshot. It is the shared foundation for two time-series comparisons:
// summing across members yields a repository's daily series (repository-axis
// overlay), and filtering by repository yields each member's daily series within
// that repository (member-axis overlay). The day is stored as an ISO
// "2006-01-02" string normalized to UTC, matching MemberDayStat, so range
// filtering and bucketing (done on the frontend) are timezone-independent.
//
// Only (login, name_with_owner, day) combinations with activity get a row, so
// the table stays sparse despite its high theoretical cardinality.
type MemberRepoDayStat struct {
	ent.Schema
}

// Fields of the MemberRepoDayStat.
func (MemberRepoDayStat) Fields() []ent.Field {
	return []ent.Field{
		field.String("login").
			NotEmpty(),
		field.String("name_with_owner").
			NotEmpty(),
		field.String("day").
			NotEmpty(),
		field.Int("commit_count").
			Default(0),
		field.Int("pr_created").
			Default(0),
		field.Int("pr_merged").
			Default(0),
		field.Int("issue_count").
			Default(0),
		field.Int("review_count").
			Default(0),
		field.Int("additions").
			Default(0),
		field.Int("deletions").
			Default(0),
	}
}

// Edges of the MemberRepoDayStat.
func (MemberRepoDayStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("snapshot", Snapshot.Type).
			Ref("member_repo_day_stats").
			Unique().
			Required(),
	}
}

// Indexes of the MemberRepoDayStat.
func (MemberRepoDayStat) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("snapshot").
			Fields("login", "name_with_owner", "day").
			Unique(),
		// Repository-axis aggregation scans by repository within a snapshot.
		index.Edges("snapshot").
			Fields("name_with_owner"),
	}
}

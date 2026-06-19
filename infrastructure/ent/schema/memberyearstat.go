package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemberYearStat holds the per-member, per-year metrics for a single snapshot.
// It backs the yearly trend shown on a member's drill-down view.
type MemberYearStat struct {
	ent.Schema
}

// Fields of the MemberYearStat.
func (MemberYearStat) Fields() []ent.Field {
	return []ent.Field{
		field.String("login").
			NotEmpty(),
		field.Int("year"),
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

// Edges of the MemberYearStat.
func (MemberYearStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("snapshot", Snapshot.Type).
			Ref("member_year_stats").
			Unique().
			Required(),
	}
}

// Indexes of the MemberYearStat.
func (MemberYearStat) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("snapshot").
			Fields("login", "year").
			Unique(),
	}
}

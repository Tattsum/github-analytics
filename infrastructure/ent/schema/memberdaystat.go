package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemberDayStat holds the per-member, per-day metrics for a single snapshot.
// It backs the day-level time-series trend and the arbitrary date-range filter
// shown on the team overview and member drill-down views. The day is stored as
// an ISO "2006-01-02" string normalized to UTC, so range filtering and bucketing
// (done on the frontend) are timezone-independent.
type MemberDayStat struct {
	ent.Schema
}

// Fields of the MemberDayStat.
func (MemberDayStat) Fields() []ent.Field {
	return []ent.Field{
		field.String("login").
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

// Edges of the MemberDayStat.
func (MemberDayStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("snapshot", Snapshot.Type).
			Ref("member_day_stats").
			Unique().
			Required(),
	}
}

// Indexes of the MemberDayStat.
func (MemberDayStat) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("snapshot").
			Fields("login", "day").
			Unique(),
	}
}

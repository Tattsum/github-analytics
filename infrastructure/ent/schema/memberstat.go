package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemberStat holds the member-level scalar metrics for a single snapshot. These
// are the cross-member comparable values used for ranking and comparison on the
// frontend.
type MemberStat struct {
	ent.Schema
}

// Fields of the MemberStat.
func (MemberStat) Fields() []ent.Field {
	return []ent.Field{
		field.String("login").
			NotEmpty(),
		field.Int("total_commits").
			Default(0),
		field.Int("total_pr_created").
			Default(0),
		field.Int("total_pr_merged").
			Default(0),
		field.Int("total_issues").
			Default(0),
		field.Int("total_reviews").
			Default(0),
		field.Int("total_additions").
			Default(0),
		field.Int("total_deletions").
			Default(0),
		// firstActivityYear is 0 when the member has no recorded activity.
		field.Int("first_activity_year").
			Default(0),
		field.Int("peak_activity_year").
			Default(0),
		field.Int("peak_activity_commits").
			Default(0),
		field.Float("pr_to_review_ratio").
			Default(0),
	}
}

// Edges of the MemberStat.
func (MemberStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("snapshot", Snapshot.Type).
			Ref("member_stats").
			Unique().
			Required(),
	}
}

// Indexes of the MemberStat.
func (MemberStat) Indexes() []ent.Index {
	return []ent.Index{
		// Stats are always fetched per snapshot; login is unique within one.
		index.Edges("snapshot").
			Fields("login").
			Unique(),
	}
}

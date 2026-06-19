package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// MemberRepoStat holds the per-member, per-repository metrics for a single
// snapshot. Stored for all repositories a member contributed to, it backs both
// the member drill-down (top repositories) and the repository-axis aggregation.
type MemberRepoStat struct {
	ent.Schema
}

// Fields of the MemberRepoStat.
func (MemberRepoStat) Fields() []ent.Field {
	return []ent.Field{
		field.String("login").
			NotEmpty(),
		field.String("name_with_owner").
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

// Edges of the MemberRepoStat.
func (MemberRepoStat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("snapshot", Snapshot.Type).
			Ref("member_repo_stats").
			Unique().
			Required(),
	}
}

// Indexes of the MemberRepoStat.
func (MemberRepoStat) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("snapshot").
			Fields("login", "name_with_owner").
			Unique(),
		// Repository-axis aggregation scans by repository within a snapshot.
		index.Edges("snapshot").
			Fields("name_with_owner"),
	}
}

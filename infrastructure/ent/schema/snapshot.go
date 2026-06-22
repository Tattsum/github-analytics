// Package schema declares the ent entity schemas used to persist aggregated
// GitHub analytics snapshots in PostgreSQL.
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Snapshot represents a single batch run: one point-in-time capture of the
// aggregated team analytics. All per-member and per-repository stats reference
// the snapshot they belong to, so the web layer can read the latest snapshot.
type Snapshot struct {
	ent.Schema
}

// Fields of the Snapshot.
func (Snapshot) Fields() []ent.Field {
	return []ent.Field{
		field.Time("captured_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Snapshot.
func (Snapshot) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("member_stats", MemberStat.Type),
		edge.To("member_year_stats", MemberYearStat.Type),
		edge.To("member_day_stats", MemberDayStat.Type),
		edge.To("member_repo_stats", MemberRepoStat.Type),
	}
}

// Indexes of the Snapshot.
func (Snapshot) Indexes() []ent.Index {
	return []ent.Index{
		// The web layer selects the latest snapshot by captured_at.
		index.Fields("captured_at"),
	}
}

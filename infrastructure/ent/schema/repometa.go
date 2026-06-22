package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RepoMeta holds the owner metadata for one repository in a single snapshot.
// It is the authoritative source for "is this repository owned by an
// organization?" so the frontend can filter cross-repository comparisons to
// org-internal repositories. Owner type is a property of the repository, not of
// any (member, repository, day) row, so it lives here as one row per repository
// rather than being denormalized onto the high-cardinality stat buckets.
type RepoMeta struct {
	ent.Schema
}

// Fields of the RepoMeta.
func (RepoMeta) Fields() []ent.Field {
	return []ent.Field{
		field.String("name_with_owner").
			NotEmpty(),
		field.String("owner").
			NotEmpty(),
		// owner_type is GitHub's owner __typename: "Organization" or "User".
		// Empty when the collector could not determine it.
		field.String("owner_type").
			Default(""),
	}
}

// Edges of the RepoMeta.
func (RepoMeta) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("snapshot", Snapshot.Type).
			Ref("repo_metas").
			Unique().
			Required(),
	}
}

// Indexes of the RepoMeta.
func (RepoMeta) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("snapshot").
			Fields("name_with_owner").
			Unique(),
	}
}

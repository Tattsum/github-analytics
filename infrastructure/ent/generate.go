// Package ent contains the generated ent ORM client and the schema definitions
// that drive its code generation. ent and database access are confined to the
// infrastructure layer; domain and application layers depend only on repository
// interfaces and never import this package.
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema

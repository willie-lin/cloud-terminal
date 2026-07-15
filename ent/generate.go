package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature privacy,entql,namedges,bidiedges,schema/snapshot,sql/schemaconfig,sql/execquery,sql/upsert,sql/versioned-migration,sql/modifier,sql/lock,sql/aggregate,sql/json,sql/transaction  ./schema

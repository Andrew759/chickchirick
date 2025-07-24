package abstraction

import (
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
)

type Migrator interface {
	// Tables
	CreateTables(entities map[string][]migratorDto.MigratorInfo) error
	CreateTable(entity migratorDto.MigratorInfo) error
	DropTable(entity migratorDto.MigratorInfo) error

	// Constraints
	CreateConstraint(entity migratorDto.MigratorInfo, constraintName string) error
	DropConstraint(entity migratorDto.MigratorInfo, constraintName string) error

	// Indexes
	CreateIndex(entity migratorDto.MigratorInfo, indexName string) error
	DropIndex(entity migratorDto.MigratorInfo, indexName string) error

	// Schema
	DropSchema(schemaName string) error
}

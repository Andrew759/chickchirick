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
	CreateConstraint(dst interface{}, name string) error
	DropConstraint(dst interface{}, name string) error

	// Indexes
	CreateIndex(dst interface{}, name string) error
	DropIndex(dst interface{}, name string) error
}

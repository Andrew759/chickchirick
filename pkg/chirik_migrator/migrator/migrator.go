package migrator

import (
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"chickChirick/pkg/chirik_migrator/migrator/service"
)

type Migrator struct {
	TableCreator   service.TableCreator
	FixtureCreator service.FixtureCreator
	TableDropper   service.TableDropper
}

func (m Migrator) CreateTables(migratorEntities map[string][]migratorDto.MigratorInfo) error {
	processedTablesSqlMeta, err := m.TableCreator.CreateTables(migratorEntities)
	if err != nil {
		return err
	}
	if m.FixtureCreator.FixtureCount > 0 {
		for _, processedSqlMetas := range processedTablesSqlMeta {
			for _, processedSqlMeta := range processedSqlMetas {
				err = m.FixtureCreator.InsertFixtures(processedSqlMeta)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m Migrator) CreateTable(migratorInfo migratorDto.MigratorInfo) error {
	processedSqlMetas, err := m.TableCreator.CreateTable(migratorInfo)
	if err != nil {
		return err
	}
	if m.FixtureCreator.FixtureCount > 0 {
		for _, processedSqlMeta := range processedSqlMetas {
			err = m.FixtureCreator.InsertFixtures(processedSqlMeta)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DropTable TODO: implement this
func (m Migrator) DropTable(entity migratorDto.MigratorInfo) error {
	return nil
}

// CreateConstraint TODO: implement this
func (m Migrator) CreateConstraint(entity migratorDto.MigratorInfo, constraintName string) error {
	return nil
}

// DropConstraint TODO: implement this
func (m Migrator) DropConstraint(entity migratorDto.MigratorInfo, constraintName string) error {
	return nil
}

// CreateIndex TODO: implement this
func (m Migrator) CreateIndex(entity migratorDto.MigratorInfo, indexName string) error {
	return nil
}

// DropIndex TODO: implement this
func (m Migrator) DropIndex(entity migratorDto.MigratorInfo, indexName string) error {
	return nil
}

func (m Migrator) DropSchema(schemaName string) error {
	return m.TableDropper.DropSchema(schemaName)
}

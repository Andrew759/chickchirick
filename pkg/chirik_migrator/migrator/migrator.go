package migrator

import (
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"chickChirick/pkg/chirik_migrator/migrator/service"
	"context"

	"golang.org/x/sync/errgroup"
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

	//TODO: сейчас фикстуры не поддерживают целостность данных, поэтому использование горутин легко реализуется
	if m.FixtureCreator.FixtureCount > 0 {
		g, gCtx := errgroup.WithContext(context.Background())

		for _, processedSqlMetas := range processedTablesSqlMeta {
			for _, meta := range processedSqlMetas {
				meta := meta // Важно для старых версий Go (до 1.22)

				g.Go(func() error {
					select {
					case <-gCtx.Done():
						return gCtx.Err()
					default:
						// Идеоматично здесь передать gCtx в InsertFixtures, но метод этого не поддерживает
						return m.FixtureCreator.InsertFixtures(meta)
					}
				})
			}
		}

		// Wait вернет первую возникшую ошибку
		if err := g.Wait(); err != nil {
			return err
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

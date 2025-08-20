package service

import (
	"chickChirick/pkg/chirik_faker"
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	"fmt"
)

type FixtureCreator struct {
	SqlProcessor
}

func (fc FixtureCreator) InsertFixtures(sqlMeta dto.Meta) error {
	var errList []error
	for i := 0; i < fc.FixtureCount; i++ {
		sqlMetaIterationCopy := sqlMeta
		err := fc.writeFixtureToSqlMeta(&sqlMetaIterationCopy)
		if err != nil {
			errList = append(errList, err)
		}
		err = fc.ProcessSQLMeta(sqlMetaIterationCopy)
		if err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		//TODO: доработать
		return fmt.Errorf("%s", errList)
	}

	return nil
}

func (fc FixtureCreator) writeFixtureToSqlMeta(sqlMeta *dto.Meta) error {
	var sqlFieldList []string
	sqlFieldList = append(sqlFieldList, "INSERT INTO ? (")

	var sqlValues []dto.ValueMeta
	sqlValues = append(sqlValues, dto.ValueMeta{
		Value:  sqlMeta.TableName,
		IsSafe: true,
	})

	processedCount := 0
	for _, valueMeta := range sqlMeta.SqlValues {
		if !valueMeta.IsValueStore {
			continue
		}
		valName := fmt.Sprintf("%v", valueMeta.Value)
		//При установке фикстуры устанавливается флаг - не безопасно
		valueMeta.IsSafe = false

		//sqlValues = append(sqlValues, valueMeta)
		sqlFieldList = append(sqlFieldList, fmt.Sprintf("%v", valueMeta.Value))
		processedCount++

		if processedCount < sqlMeta.FieldCount {
			sqlFieldList = append(sqlFieldList, ", ")
		} else if processedCount == sqlMeta.FieldCount {
			sqlFieldList = append(sqlFieldList, ")")
		}

		fixtureValue, err := chirik_faker.FakeSQLValue(valName, valueMeta.Type)
		if err != nil {
			return err
		}

		//Если поле сконфигурировано так, чтобы игнорировать фикстуры - в качестве значения устанавливается nil
		if _, hasValue := fc.FixtureNilColumns[valName]; hasValue {
			fixtureValue = nil
		}

		sqlValues = append(sqlValues, dto.ValueMeta{
			Value:  fixtureValue,
			IsSafe: false,
		})

	}

	sqlFieldList = append(sqlFieldList, " VALUES (")
	for i := 1; i <= sqlMeta.FieldCount; i++ {
		if i != sqlMeta.FieldCount {
			sqlFieldList = append(sqlFieldList, "?, ")
		} else {
			sqlFieldList = append(sqlFieldList, "?")
		}
	}
	sqlFieldList = append(sqlFieldList, ");")

	sqlMeta.SqlFieldList = sqlFieldList
	sqlMeta.SqlValues = sqlValues
	sqlMeta.FieldCommentList = []string{}
	sqlMeta.FieldCommentValues = []string{}

	if fc.FixturePrefix != "" {
		sqlMeta.MigrationPrefix = fc.FixturePrefix + "_" + sqlMeta.MigrationPrefix
	}

	return nil
}

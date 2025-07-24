package service

import (
	"chickChirick/pkg/chirik_migrator/file"
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	"chickChirick/pkg/chirik_migrator/migrator/helper"
)

type SqlProcessor struct {
	dto.MigratorDiContainer
}

func (sp SqlProcessor) ProcessSQLMeta(sqlMeta dto.Meta) error {
	var resultSQL string
	for _, sqlField := range sqlMeta.SqlFieldList {
		resultSQL += sqlField
	}

	var err error
	resultSQL, err = helper.BuildRawSql(resultSQL, sqlMeta)
	if err != nil {
		return err
	}

	_, err = sp.NativeDB().Exec(resultSQL)
	if err != nil {
		return err
	}

	return file.WriteSQLToFile(resultSQL, sqlMeta.MigrationPrefix, sp.MConfig.MigrationFilesPath)
}

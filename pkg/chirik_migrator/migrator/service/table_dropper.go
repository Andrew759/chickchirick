package service

type TableDropper struct {
	SqlProcessor
}

func (td TableDropper) DropSchema(schemaName string) error {
	_, err := td.DBDecorator.NativeDB().Exec("DROP SCHEMA " + schemaName + "CASCADE; CREATE SCHEMA " + schemaName + ";")

	return err
}

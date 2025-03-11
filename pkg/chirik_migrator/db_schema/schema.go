package db_schema

type Schema struct {
	Name       string
	Table      string
	PrimaryKey *Field
	ForeignKey []*Field
	Fields     []*Field
	//TODO: FieldsByName  map[string]*Field
	Relationships Relationships
}

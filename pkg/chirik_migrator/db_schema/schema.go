package db_schema

type Schema struct {
	Name       string
	Table      string
	PrimaryKey *Field
	ForeignKey []*Field
	Fields     []*Field
	//TODO: идея нуждается в доработке. Индексы пока что не реализованы
	Indexes       []*Field
	Relationships Relationships
}

func (s Schema) HasPrimaryKey() bool {
	return s.PrimaryKey != nil
}

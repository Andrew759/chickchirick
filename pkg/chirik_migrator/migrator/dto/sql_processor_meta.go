package dto

type Meta struct {
	TableName          string
	SqlFieldList       []string
	SqlValues          []ValueMeta
	FieldCommentList   []string
	FieldCommentValues []string
	FieldCount         int
	MigrationPrefix    string
}

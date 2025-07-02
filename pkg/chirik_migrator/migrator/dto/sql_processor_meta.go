package dto

type Meta struct {
	TableName          string
	SqlFieldList       []string
	SqlValues          []any
	FieldCommentList   []string
	FieldCommentValues []string
	FieldMetas         []FieldMeta
	FieldCount         int
}

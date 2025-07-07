package postgres

const (
	Bool                     Type = "BOOLEAN"
	Smallint                 Type = "SMALLINT"
	Int                      Type = "INTEGER"
	Bigint                   Type = "BIGINT"
	BigSerial                Type = "BIGSERIAL"
	Real                     Type = "REAL"
	DoublePrecision          Type = "DOUBLE PRECISION"
	Varchar                  Type = "VARCHAR"
	Text                     Type = "TEXT"
	TimestampWithTimezone    Type = "TIMESTAMP WITH TIME ZONE"
	TimestampWithoutTimezone Type = "TIMESTAMP WITHOUT TIME ZONE"
	Bytes                    Type = "SMALLINT"
	Uuid                     Type = "UUID"
	Json                     Type = "JSON"
	Jsonb                    Type = "JSONB"
	Null                     Type = "NULL"
)

type Type string

func (d Type) String() string {
	return string(d)
}

func (d Type) IsJson() bool {
	return d == Json || d == Jsonb
}

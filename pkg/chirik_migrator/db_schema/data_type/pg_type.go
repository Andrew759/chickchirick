package data_type

const (
	boolean                  PgType = "BOOLEAN"
	smallint                 PgType = "SMALLINT"
	integer                  PgType = "INTEGER"
	bigInt                   PgType = "BIGINT"
	bigSerial                PgType = "BIGSERIAL"
	realT                    PgType = "REAL"
	doublePrecision          PgType = "DOUBLE PRECISION"
	varchar                  PgType = "VARCHAR"
	text                     PgType = "TEXT"
	timestampWithTimezone    PgType = "TIMESTAMP WITH TIME ZONE"
	timestampWithoutTimezone PgType = "TIMESTAMP WITHOUT TIME ZONE"
	bytes                    PgType = "SMALLINT"
	uuid                     PgType = "UUID"
	json                     PgType = "JSON"
	jsonb                    PgType = "JSONB"
	null                     PgType = "NULL"
)

type PgType string

func (d PgType) String() string {
	return string(d)
}

func (d PgType) IsHasSetValue() bool {
	return string(d) != ""
}

func (d PgType) IsJson() bool {
	return d == json || d == jsonb
}

func (d PgType) Bool() Type {
	return boolean
}

func (d PgType) SmallInt() Type {
	return smallint
}

func (d PgType) Int() Type {
	return integer
}

func (d PgType) BigInt() Type {
	return bigInt
}

func (d PgType) BigSerial() Type {
	return bigSerial
}

func (d PgType) Real() Type {
	return realT
}

func (d PgType) DoublePrecision() Type {
	return doublePrecision
}

func (d PgType) Varchar() Type {
	return varchar
}

func (d PgType) Text() Type {
	return text
}

func (d PgType) TimestampWithTimezone() Type {
	return timestampWithTimezone
}

func (d PgType) TimestampWithoutTimezone() Type {
	return timestampWithoutTimezone
}

func (d PgType) Bytes() Type {
	return bytes
}

func (d PgType) Uuid() Type {
	return uuid
}

func (d PgType) Json() Type {
	return json
}

func (d PgType) Jsonb() Type {
	return jsonb
}

func (d PgType) Null() Type {
	return null
}

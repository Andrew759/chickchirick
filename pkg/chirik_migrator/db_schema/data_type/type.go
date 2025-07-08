package data_type

type Type interface {
	String() string
	IsHasSetValue() bool
	IsJson() bool
	Bool() Type
	SmallInt() Type
	Int() Type
	BigInt() Type
	BigSerial() Type
	Real() Type
	DoublePrecision() Type
	Varchar() Type
	Text() Type
	TimestampWithTimezone() Type
	TimestampWithoutTimezone() Type
	Bytes() Type
	Uuid() Type
	Json() Type
	Jsonb() Type
	Null() Type
}

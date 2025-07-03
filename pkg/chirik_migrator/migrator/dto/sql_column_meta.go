package dto

type ValueMeta struct {
	Value  string
	Type   string
	IsSafe bool //Требуется ли экранирование для значения?
}

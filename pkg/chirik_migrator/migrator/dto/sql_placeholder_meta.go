package dto

import (
	"chickChirick/pkg/chirik_migrator/db_schema/data_type"
)

type ValueMeta struct {
	Value        any
	Type         data_type.Type
	IsSafe       bool //Требуется ли экранирование для значения?
	IsValueStore bool //Хранит ли данный элемент какое-либо значение? Пример - столбец таблицы
}

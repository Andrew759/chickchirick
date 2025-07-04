package dto

type ValueMeta struct {
	Value        any
	IsSafe       bool //Требуется ли экранирование для значения?
	IsStoreValue bool //Хранит ли данный элемент какое-либо значение? Пример - столбец таблицы
}

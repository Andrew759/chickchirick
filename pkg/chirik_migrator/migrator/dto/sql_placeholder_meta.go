package dto

type ValueMeta struct {
	Value        any
	Type         string
	IsSafe       bool //Требуется ли экранирование для значения?
	IsValueStore bool //Хранит ли данный элемент какое-либо значение? Пример - столбец таблицы
}

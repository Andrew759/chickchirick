package message

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Id        int     `json:"id" gorm:"type:int; unique; primaryKey; autoIncrement"`
	MessageId int     `json:"message_id" gorm:"type:int"`
	Message   Message `json:"message" gorm:"references:MessageId"`
	//TODO: удалить автогенерацию и мок, когда будет реализован сервис файлов
	FileUuid string `json:"file_uuid" gorm:"type:uuid; default:gen_random_uuid()"`
}

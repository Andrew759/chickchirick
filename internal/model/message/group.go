package message

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	MessageId int     `json:"message_id" gorm:"type:int"`
	Message   Message `json:"message" gorm:"references:MessageId"`
	//TODO: удалить автогенерацию и мок, когда будет реализован сервис групп
	GroupId string `json:"GroupId" gorm:"type:uuid; default:gen_random_uuid()"`
}

package message

import (
	"gorm.io/gorm"
)

type Personal struct {
	gorm.Model        `c_migrator:"enabled"`
	Id                int          `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	MessageId         int          `json:"message_id" gorm:"type:int"`
	Message           Message      `json:"message" gorm:"references:MessageId"`
	SenderId          int          `json:"sender_id" gorm:"type:int"`
	SenderRelation    UserRelation `json:"sender_relation" gorm:"references:SenderId"`
	RecipientId       int          `json:"recipient_id" gorm:"type:int"`
	RecipientRelation UserRelation `json:"recipient_relation" gorm:"references:RecipientId"`
}

package message

import (
	"gorm.io/gorm"
)

type Deleted struct {
	gorm.Model
	MessageId    int          `json:"message_id" gorm:"type:int"`
	Message      Message      `json:"message" gorm:"references:MessageId"`
	UserId       int          `json:"user_id" gorm:"type:int"`
	UserRelation UserRelation `json:"user_relation" gorm:"references:UserId"`
}

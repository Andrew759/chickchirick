package message

import (
	"gorm.io/gorm"
	"time"
)

type Status struct {
	gorm.Model       `c_migrator:"enabled" c_migrator_orm:"gorm"`
	Id               int          `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	MessageId        int          `json:"message_id" gorm:"type:int"`
	Message          Message      `json:"message" gorm:"references:MessageId"`
	ReadUserId       int          `json:"read_user_id" gorm:"type:int"`
	ReadUserRelation UserRelation `json:"read_user_relation" gorm:"references:ReadUserId"`
	Date             time.Time    `json:"date" gorm:"type:timestamp without time zone"`
}

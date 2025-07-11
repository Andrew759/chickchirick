package message

import (
	"gorm.io/gorm"
	"time"
)

type Status struct {
	gorm.Model       `c_migrator:"enabled"`
	Id               int          `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	MessageId        int          `json:"message_id" gorm:"type:int"`
	Message          Message      `json:"message" gorm:"references:MessageId"`
	ReadUserId       int          `json:"read_user_id" gorm:"type:int"`
	ReadUserRelation UserRelation `json:"read_user_relation" gorm:"references:ReadUserId"`
	Date             time.Time    `json:"date" gorm:"type:timestamp without time zone"`
}

func CreateStatus(db *gorm.DB, s *Status) error {
	return db.Create(s).Error
}

func GetStatus(db *gorm.DB) ([]Status, error) {
	var status []Status
	result := db.Find(&status)

	return status, result.Error
}

func GetStatusById(db *gorm.DB, id int) (Status, error) {
	var status Status
	result := db.First(&status, id)

	return status, result.Error
}

func DeleteStatusById(db *gorm.DB, id int) error {
	return db.Delete(&Status{}, id).Error
}

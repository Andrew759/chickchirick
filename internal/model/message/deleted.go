package message

import (
	"gorm.io/gorm"
)

type Deleted struct {
	gorm.Model   `c_migrator:"enabled"`
	MessageId    int          `json:"message_id" gorm:"type:int"`
	Message      Message      `json:"message" gorm:"references:MessageId"`
	UserId       int          `json:"user_id" gorm:"type:int"`
	UserRelation UserRelation `json:"user_relation" gorm:"references:UserId"`
}

func CreateDeleted(db *gorm.DB, d *Deleted) error {
	return db.Create(d).Error
}

func GetDeleted(db *gorm.DB) ([]Deleted, error) {
	var deleted []Deleted
	result := db.Find(&deleted)

	return deleted, result.Error
}

func GetDeletedById(db *gorm.DB, id int) (Deleted, error) {
	var deleted Deleted
	result := db.First(&deleted, id)

	return deleted, result.Error
}

func DeleteDeletedById(db *gorm.DB, id int) error {
	return db.Delete(&Deleted{}, id).Error
}

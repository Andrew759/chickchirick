package auth

import (
	"gorm.io/gorm"
	"time"
)

type Code struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Code       int8      `json:"code" gorm:"type:smallint"`
	SessionId  int       `json:"session_id" gorm:"type:int"`
	Session    Session   `json:"session" gorm:"references:SessionId"`
	ExpiresAt  time.Time `json:"expires_at" gorm:"type:timestamp without time zone"`
}

func CreateCode(db *gorm.DB, c *Code) error {
	return db.Create(c).Error
}

func UpdateCode(db *gorm.DB, c *Code) error {
	return db.Save(c).Error
}

func GetCodes(db *gorm.DB) ([]Code, error) {
	var codes []Code
	result := db.Find(&codes)

	return codes, result.Error
}

func GetCodeById(db *gorm.DB, id int) (Code, error) {
	var codes Code
	result := db.First(&codes, id)

	return codes, result.Error
}

func DeleteCodeById(db *gorm.DB, id int) error {
	return db.Delete(&Code{}, id).Error
}

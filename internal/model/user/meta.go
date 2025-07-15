package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Meta struct {
	gorm.Model `c_migrator:"enabled"  c_migrator_t_name:"user_meta"`
	UserUuid   uuid.UUID `json:"user_uuid" gorm:"type:uuid;default:gen_random_uuid()"`
	UserId     int       `json:"user_id" gorm:"type:int"`
	User       User      `json:"user" gorm:"references:UserId"`
}

func (Meta) TableName() string {
	return "user_meta"
}

func CreateMeta(db *gorm.DB, b *Meta) error {
	return db.Create(b).Error
}

func UpdateMeta(db *gorm.DB, m *Meta) error {
	return db.Save(m).Error
}

func GetMetas(db *gorm.DB) ([]Meta, error) {
	var metas []Meta
	result := db.Find(&metas)

	return metas, result.Error
}

func GetMetaById(db *gorm.DB, id int) (Meta, error) {
	var meta Meta
	result := db.First(&meta, id)

	return meta, result.Error
}

func DeleteMetaById(db *gorm.DB, id int) error {
	return db.Delete(&Meta{}, id).Error
}

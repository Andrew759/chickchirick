package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Meta struct {
	gorm.Model `c_migrator:"enabled"  c_migrator_t_name:"user_meta"`
	UserUuid   uuid.UUID `json:"user_uuid" gorm:"type:uuid;default:gen_random_uuid()"`
	UserId     int       `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User       User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var MetaNotFoundErr = errors.New("user not found")

func (Meta) TableName() string {
	return "user_meta"
}

func CreateMeta(db *gorm.DB, b *Meta) error {
	return db.Create(b).Error
}

func UpdateMetaById(db *gorm.DB, m *Meta, id int) error {
	var meta Meta
	result := db.First(&meta, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return MetaNotFoundErr
	}

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
	var meta Meta
	result := db.First(&meta, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return MetaNotFoundErr
	}

	return db.Delete(&Meta{}, id).Error
}

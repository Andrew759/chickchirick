package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Meta struct {
	gorm.Model `c_migrator:"enabled"  c_migrator_t_name:"user_meta"`
	//TODO: мигратор не обрабатывает поле UUID. Исправить!
	Id        int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid  uuid.UUID `json:"user_uuid" gorm:"type:uuid;default:gen_random_uuid()"`
	UserId    int       `json:"user_id" gorm:"type:int;unique;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.TimestampWithTimeZoneMicro
	UpdatedAt time.TimestampWithTimeZoneMicro
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var MetaNotFoundErr = errors.New("user not found")

func (Meta) TableName() string {
	return "user_meta"
}

func CreateMeta(ctx context.Context, db *gorm.DB, b *Meta) error {
	return db.WithContext(ctx).Create(b).Error
}

func UpdateMetaByUserId(ctx context.Context, db *gorm.DB, m *Meta, userId int) error {
	var meta Meta
	tx := db.WithContext(ctx)

	result := tx.Model(&Meta{}).Where("user_id = ?", userId).Take(&meta)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return MetaNotFoundErr
	}

	return tx.Save(m).Error
}

func GetMetas(ctx context.Context, db *gorm.DB) ([]Meta, error) {
	var metas []Meta
	result := db.WithContext(ctx).Find(&metas)

	return metas, result.Error
}

func GetMetaByUserId(ctx context.Context, db *gorm.DB, userId int) (Meta, error) {
	var meta Meta
	result := db.WithContext(ctx).Model(&Meta{}).Where("user_id = ?", userId).Take(&meta)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return meta, MetaNotFoundErr
	}

	return meta, result.Error
}

func DeleteMetaByUserId(ctx context.Context, db *gorm.DB, userId int) error {
	tx := db.WithContext(ctx)

	var meta Meta
	result := tx.Model(&Meta{}).Where("user_id = ?", userId).Take(&meta)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return MetaNotFoundErr
	}

	return tx.Delete(&Meta{}, meta.Id).Error
}

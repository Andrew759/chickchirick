package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	Id        int                             `json:"id" gorm:"column:id;type:integer;primaryKey;autoIncrement"`
	FileUuid  uuid.UUID                       `json:"file_uuid" gorm:"column:file_uuid;type:uuid;not null"`
	UserId    int                             `json:"user_id" gorm:"column:user_id;type:integer;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      User                            `json:"user,omitempty" gorm:"foreignKey:UserId;references:Id"`
	MediaKind string                          `json:"media_kind" gorm:"column:media_kind;type:varchar(20);default:'file'"`
	CreatedAt time.TimestampWithTimeZoneMicro `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.TimestampWithTimeZoneMicro `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt                  `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

func (Photo) TableName() string {
	return "photos"
}

var PhotoNotFoundErr = errors.New("photo not found")

func CreatePhoto(ctx context.Context, db *gorm.DB, b *Photo) error {
	return db.WithContext(ctx).Create(b).Error
}

func UpdatePhotoById(ctx context.Context, db *gorm.DB, p *Photo, id int) error {
	var photo Photo
	tx := db.WithContext(ctx)

	result := tx.First(&photo, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PhotoNotFoundErr
	}

	return tx.Save(p).Error
}

func GetPhotos(ctx context.Context, db *gorm.DB) ([]Photo, error) {
	var photos []Photo
	result := db.WithContext(ctx).Find(&photos)

	return photos, result.Error
}

func GetPhotoById(ctx context.Context, db *gorm.DB, id int) (Photo, error) {
	var photo Photo
	result := db.WithContext(ctx).First(&photo, id)

	return photo, result.Error
}

func GetPhotosByUserId(ctx context.Context, db *gorm.DB, userId int) ([]Photo, error) {
	var photos []Photo
	result := db.WithContext(ctx).Table("photos").Where("user_id = ? AND deleted_at IS NULL", userId).Order("id DESC").Find(&photos)

	return photos, result.Error
}

func DeletePhotoById(ctx context.Context, db *gorm.DB, id int) error {
	var photo Photo
	tx := db.WithContext(ctx)

	result := tx.First(&photo, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PhotoNotFoundErr
	}

	return tx.Delete(&Photo{}, id).Error
}

package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	FileUuid   uuid.UUID `gorm:"type:uuid;not null"`
	UserId     int       `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User       User      `json:"user" gorm:"foreignKey:UserId;references:Id"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
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
	result := db.WithContext(ctx).Where("user_id = ?", userId).Find(&photos)

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

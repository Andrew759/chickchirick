package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model `c_migrator:"enabled"`
	FieUuid    uuid.UUID `gorm:"type:uuid"`
	UserId     int       `json:"user_id" gorm:"type:int"`
	User       User      `json:"user" gorm:"references:UserId"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func CreatePhoto(db *gorm.DB, b *Photo) error {
	return db.Create(b).Error
}

func UpdatePhoto(db *gorm.DB, p *Photo) error {
	return db.Save(p).Error
}

func GetPhotos(db *gorm.DB) ([]Photo, error) {
	var photos []Photo
	result := db.Find(&photos)

	return photos, result.Error
}

func GetPhotoById(db *gorm.DB, id int) (Photo, error) {
	var photo Photo
	result := db.First(&photo, id)

	return photo, result.Error
}

func DeletePhotoById(db *gorm.DB, id int) error {
	return db.Delete(&Photo{}, id).Error
}

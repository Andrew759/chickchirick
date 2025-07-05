package user

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model `c_migrator:"enabled"`
	//TODO: удалить мок UUID, когда появится сервис файлов
	FieUuid pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserId  int         `json:"user_id" gorm:"type:int"`
	User    User        `json:"user" gorm:"references:UserId"`
}

func CreatePhoto(db *gorm.DB, b *Photo) error {
	return db.Create(b).Error
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

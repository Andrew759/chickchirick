package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model `c_migrator:"enabled"`
	UserId     int     `json:"user_id" gorm:"type:int"`
	User       User    `json:"user" gorm:"references:UserId"`
	Timezone   int8    `json:"timezone" gorm:"type:smallint;default:3"`
	Email      string  `json:"email" gorm:"type:varchar(256)"`
	Password   *string `json:"password" gorm:"type:varchar(1024)"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var PropertyNotFoundErr = errors.New("property not found")

func (Property) TableName() string {
	return "properties"
}

func CreateProperty(db *gorm.DB, b *Property) error {
	return db.Create(b).Error
}

func UpdatePropertyById(db *gorm.DB, p *Property, id int) error {
	var property Property
	result := db.First(&property, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PropertyNotFoundErr
	}

	return db.Save(p).Error
}

func GetProperties(db *gorm.DB) ([]Property, error) {
	var properties []Property
	result := db.Find(&properties)

	return properties, result.Error
}

func GetPropertyById(db *gorm.DB, id int) (Property, error) {
	var property Property
	result := db.First(&property, id)

	return property, result.Error
}

func DeletePropertyById(db *gorm.DB, id int) error {
	var property Property
	result := db.First(&property, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PropertyNotFoundErr
	}

	return db.Delete(&Property{}, id).Error
}

package user

import "gorm.io/gorm"

type Property struct {
	gorm.Model `c_migrator:"enabled"`
	UserId     int    `json:"user_id" gorm:"type:int"`
	User       User   `json:"user" gorm:"references:UserId"`
	Timezone   int8   `json:"timezone" gorm:"type:smallint"`
	Email      string `json:"email" gorm:"type:varchar(256)"`
}

func CreateProperty(db *gorm.DB, b *Property) error {
	return db.Create(b).Error
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
	return db.Delete(&Property{}, id).Error
}

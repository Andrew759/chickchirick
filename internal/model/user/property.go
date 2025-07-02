package user

import "gorm.io/gorm"

type Property struct {
	gorm.Model `c_migrator:"enabled" c_migrator_orm:"gorm"`
	UserId     int  `json:"user_id" gorm:"type:int"`
	User       User `json:"user" gorm:"references:UserId"`
	//TODO: Доработать таймзоны. Сейчас архитектура описана таким образом,
	// что всё время хранится без таймзон. Возможно потребуется отдельный сервис или
	// пакет для ресолва временной зоны пользователя
	Timezone int8   `json:"timezone" gorm:"type:smallint"`
	Email    string `json:"email" gorm:"type:varchar(256)"`
}

func CreateProperty(db *gorm.DB, b *Property) error {
	return db.Create(b).Error
}

func GetProperties(db *gorm.DB) ([]Property, error) {
	var properties []Property
	result := db.Find(&properties)

	return properties, result.Error
}

func GetProperty(db *gorm.DB, id int) (Property, error) {
	var property Property
	result := db.First(&property, id)

	return property, result.Error
}

func DeleteProperty(db *gorm.DB, id int) error {
	return db.Delete(&Property{}, id).Error
}

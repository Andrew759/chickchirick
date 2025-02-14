package user

import "gorm.io/gorm"

type Photos struct {
	gorm.Model
	//TODO: удалить мок UUID, когда появится сервис файлов
	FieUuid string `gorm:"type:uuid;default:gen_random_uuid()"`
	UserId  int    `json:"user_id" gorm:"type:int"`
	User    User   `json:"user" gorm:"references:UserId"`
}

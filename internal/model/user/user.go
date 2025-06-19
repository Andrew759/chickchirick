package user

import "gorm.io/gorm"

type User struct {
	gorm.Model `c_migrator:"enabled" c_migrator_orm:"gorm"`
	Id         int     `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Phone      int64   `json:"phone" gorm:"type:bigint"`
	Name       string  `json:"name" gorm:"type:varchar(256)"`
	Surname    string  `json:"surname" gorm:"type:varchar(256)"`
	Password   *string `json:"password" gorm:"type:varchar(1024)"`
}

// GetAllUsers Получение всех пользователей
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)

	return users, result.Error
}

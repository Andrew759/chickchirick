package user

import "gorm.io/gorm"

type User struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int     `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Phone      string  `json:"phone" gorm:"type:varchar(30)"`
	Name       string  `json:"name" gorm:"type:varchar(256)"`
	Login      string  `json:"login" gorm:"type:varchar(256)"`
	Surname    string  `json:"surname" gorm:"type:varchar(256)"`
	Password   *string `json:"password" gorm:"type:varchar(1024)"`
}

func CreateUser(db *gorm.DB, u *User) error {
	return db.Create(u).Error
}

func UpdateUser(db *gorm.DB, u *User) error {
	return db.Save(u).Error
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)

	return users, result.Error
}

func GetUserById(db *gorm.DB, id int) (User, error) {
	var user User
	result := db.First(&user, id)

	return user, result.Error
}

func DeleteUserById(db *gorm.DB, id int) error {
	return db.Delete(&User{}, id).Error
}

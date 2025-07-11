package user

import "gorm.io/gorm"

type Ban struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserId     int `json:"user_id" gorm:"type:int"`
	//User         User `json:"user" gorm:"references:UserId"`
	BannedUserId int `json:"banned_user_id" gorm:"type:int"`
	//BannedUser   User `json:"banned_user" gorm:"references:UserId"`
}

func CreateBan(db *gorm.DB, b *Ban) error {
	return db.Create(b).Error
}

func GetBans(db *gorm.DB) ([]Ban, error) {
	var bans []Ban
	result := db.Find(&bans)

	return bans, result.Error
}

func GetBanById(db *gorm.DB, id int) (Ban, error) {
	var ban Ban
	result := db.First(&ban, id)

	return ban, result.Error
}

func DeleteBanById(db *gorm.DB, id int) error {
	return db.Delete(&Ban{}, id).Error
}

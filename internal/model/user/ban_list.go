package user

import "gorm.io/gorm"

// BanList TODO: доработать связи
type BanList struct {
	gorm.Model `c_migrator:"enabled" c_migrator_orm:"gorm"`
	Id         int `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserId     int `json:"user_id" gorm:"type:int"`
	//User         User `json:"user" gorm:"references:UserId"`
	BannedUserId int `json:"banned_user_id" gorm:"type:int"`
	//BannedUser   User `json:"banned_user" gorm:"references:UserId"`
}

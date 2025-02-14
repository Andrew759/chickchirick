package user

import "gorm.io/gorm"

type Meta struct {
	gorm.Model
	UserUuid string `json:"user_uuid" gorm:"type:uuid;default:gen_random_uuid()"`
	UserId   int    `json:"user_id" gorm:"type:int"`
	User     User   `json:"user" gorm:"references:UserId"`
}

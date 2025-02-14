package message

import "gorm.io/gorm"

type UserRelation struct {
	gorm.Model
	UserId   int    `json:"id" gorm:"type:int; unique; primaryKey; autoIncrement"`
	UserUuid string `json:"user_uuid" gorm:"type:uuid"`
}

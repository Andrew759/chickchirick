package auth

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Id        int       `json:"id" gorm:"type:int; unique; primaryKey; autoIncrement"`
	UserUuid  string    `json:"user_uuid" gorm:"type:uuid"`
	Status    bool      `json:"status" gorm:"type:varchar(256)"`
	StartDate time.Time `json:"start_date" gorm:"type:timestamp without time zone"`
}

package auth

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	gorm.Model
	Id        int       `json:"id" gorm:"type:int; unique; primaryKey; autoIncrement"`
	SessionId int       `json:"session_id" gorm:"type:int"`
	Session   Session   `json:"session" gorm:"references:SessionId"`
	Token     string    `json:"token" gorm:"type:varchar(256)"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamp without time zone"`
}

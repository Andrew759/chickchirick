package auth

import (
	"gorm.io/gorm"
	"time"
)

type Code struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Code       int8      `json:"code" gorm:"type:smallint"`
	SessionId  int       `json:"session_id" gorm:"type:int"`
	Session    Session   `json:"session" gorm:"references:SessionId"`
	ExpiresAt  time.Time `json:"expires_at" gorm:"type:timestamp without time zone"`
}

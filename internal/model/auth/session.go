package auth

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int         `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid   pgtype.UUID `json:"user_uuid" gorm:"type:uuid"`
	Status     bool        `json:"status" gorm:"type:boolean"`
	StartDate  time.Time   `json:"start_date" gorm:"type:timestamp without time zone"`
}

package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type UserRelation struct {
	gorm.Model `c_migrator:"enabled"`
	UserId     int         `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid   pgtype.UUID `json:"user_uuid" gorm:"type:uuid"`
}

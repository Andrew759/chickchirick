package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type UserRelation struct {
	gorm.Model
	UserId   int         `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid pgtype.UUID `json:"user_uuid" gorm:"type:uuid"`
}

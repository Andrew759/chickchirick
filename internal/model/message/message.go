package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int         `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Text       pgtype.UUID `json:"text" gorm:"type:uuid"`
	Date       time.Time   `json:"start_date" gorm:"type:timestamp without time zone"`
}

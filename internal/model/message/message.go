package message

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	Id   int       `json:"id" gorm:"type:int; unique; primaryKey; autoIncrement"`
	Text string    `json:"text" gorm:"type:uuid"`
	Date time.Time `json:"start_date" gorm:"type:timestamp without time zone"`
}

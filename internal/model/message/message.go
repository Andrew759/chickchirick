package message

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Text       uuid.UUID `json:"text" gorm:"type:text;not null"`
	Date       time.Time `json:"start_date" gorm:"type:timestamp without time zone"`
}

func CreateMessage(db *gorm.DB, m *Message) error {
	return db.Create(m).Error
}

func UpdateMessage(db *gorm.DB, m *Message) error {
	return db.Save(m).Error
}

func GetMessages(db *gorm.DB) ([]Message, error) {
	var messages []Message
	result := db.Find(&messages)

	return messages, result.Error
}

func GetMessageById(db *gorm.DB, id int) (Message, error) {
	var message Message
	result := db.First(&message, id)

	return message, result.Error
}

func DeleteMessageById(db *gorm.DB, id int) error {
	return db.Delete(&Message{}, id).Error
}

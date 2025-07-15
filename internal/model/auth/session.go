package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid   uuid.UUID `json:"user_uuid" gorm:"type:uuid"`
	Status     bool      `json:"status" gorm:"type:boolean"`
	StartDate  time.Time `json:"start_date" gorm:"type:timestamp without time zone"`
}

func CreateSession(db *gorm.DB, s *Session) error {
	s.CreatedAt = time.Now()

	return db.Create(s).Error
}

func UpdateSession(db *gorm.DB, s *Session) error {
	s.UpdatedAt = time.Now()

	return db.Save(s).Error
}

func GetSessions(db *gorm.DB) ([]Session, error) {
	var sessions []Session
	result := db.Find(&sessions)

	return sessions, result.Error
}

func GetSessionById(db *gorm.DB, id int) (Session, error) {
	var session Session
	result := db.First(&session, id)

	return session, result.Error
}

func DeleteSessionById(db *gorm.DB, id int) error {
	return db.Delete(&Session{}, id).Error
}

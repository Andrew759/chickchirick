package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Meta struct {
	gorm.Model       `c_migrator:"enabled" c_migrator_t_name:"message_meta"`
	MessageUuid      pgtype.UUID `json:"message_uuid" gorm:"type:uuid;default:gen_random_uuid()"`
	MessageId        int         `json:"message_id" gorm:"type:int"`
	Message          Message     `json:"message" gorm:"references:MessageId"`
	MessageStatusId  int         `json:"message_status_id" gorm:"type:int"`
	Status           Status      `json:"status" gorm:"references:MessageStatusId"`
	RespondMessageId *int        `json:"respond_message_id" gorm:"type:int"`
	RespondMessage   *Message    `json:"RespondMessage" gorm:"references:RespondMessageId"`
}

func (Meta) TableName() string {
	return "message_meta"
}

func CreateMeta(db *gorm.DB, m *Meta) error {
	return db.Create(m).Error
}

func UpdateMeta(db *gorm.DB, m *Meta) error {
	return db.Save(m).Error
}

func GetMetas(db *gorm.DB) ([]Meta, error) {
	var meta []Meta
	result := db.Find(&meta)

	return meta, result.Error
}

func GetMetaById(db *gorm.DB, id int) (Meta, error) {
	var meta Meta
	result := db.First(&meta, id)

	return meta, result.Error
}

func DeleteMetaById(db *gorm.DB, id int) error {
	return db.Delete(&Meta{}, id).Error
}

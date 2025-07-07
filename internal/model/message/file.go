package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int     `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	MessageId  int     `json:"message_id" gorm:"type:int"`
	Message    Message `json:"message" gorm:"references:MessageId"`
	//TODO: удалить автогенерацию и мок, когда будет реализован сервис файлов
	FileUuid pgtype.UUID `json:"file_uuid" gorm:"type:uuid;default:gen_random_uuid()"`
}

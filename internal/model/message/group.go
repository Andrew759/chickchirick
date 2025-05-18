package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model `c_migrator:"enabled" c_migrator_orm:"gorm"`
	MessageId  int     `json:"message_id" gorm:"type:int"`
	Message    Message `json:"message" gorm:"references:MessageId"`
	//TODO: удалить автогенерацию и мок, когда будет реализован сервис групп
	GroupId pgtype.UUID `json:"GroupId" gorm:"type:uuid;default:gen_random_uuid()"`
}

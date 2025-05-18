package user

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Photos struct {
	gorm.Model `c_migrator:"enabled" c_migrator_orm:"gorm"`
	//TODO: удалить мок UUID, когда появится сервис файлов
	FieUuid pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserId  int         `json:"user_id" gorm:"type:int"`
	User    User        `json:"user" gorm:"references:UserId"`
}

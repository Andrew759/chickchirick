package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model `c_migrator:"enabled"`
	MessageId  int     `json:"message_id" gorm:"type:int"`
	Message    Message `json:"message" gorm:"references:MessageId"`
	//TODO: удалить автогенерацию и мок, когда будет реализован сервис групп
	GroupId pgtype.UUID `json:"GroupId" gorm:"type:uuid;default:gen_random_uuid()"`
}

func CreateGroup(db *gorm.DB, g *Group) error {
	return db.Create(g).Error
}

func GetGroups(db *gorm.DB) ([]Group, error) {
	var groups []Group
	result := db.Find(&groups)

	return groups, result.Error
}

func GetGroupById(db *gorm.DB, id int) (Group, error) {
	var group Group
	result := db.First(&group, id)

	return group, result.Error
}

func DeleteGroupById(db *gorm.DB, id int) error {
	return db.Delete(&Group{}, id).Error
}

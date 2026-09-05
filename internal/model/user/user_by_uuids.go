package user

import (
	"context"

	"gorm.io/gorm"
)

type UserWithUuid struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
	UserUuid string `json:"user_uuid"`
}

// GetUsersByUuids возвращает пользователей по списку UUID из user_meta
func GetUsersByUuids(ctx context.Context, db *gorm.DB, uuids []string) ([]UserWithUuid, error) {
	if len(uuids) == 0 {
		return []UserWithUuid{}, nil
	}

	var rows []UserWithUuid
	err := db.WithContext(ctx).
		Table("users AS u").
		Select("u.id, u.name, u.surname, u.login, user_meta.user_uuid::text AS user_uuid").
		Joins("JOIN user_meta ON user_meta.user_id = u.id").
		Where("user_meta.user_uuid IN ?", uuids).
		Where("u.deleted_at IS NULL").
		Scan(&rows).Error

	return rows, err
}

package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type UserRelation struct {
	gorm.Model `c_migrator:"enabled"`
	UserId     int         `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserUuid   pgtype.UUID `json:"user_uuid" gorm:"type:uuid"`
}

func CreateUserRelation(db *gorm.DB, ur *UserRelation) error {
	return db.Create(ur).Error
}

func UpdateUserRelation(db *gorm.DB, ur *UserRelation) error {
	return db.Save(ur).Error
}

func GetUserRelation(db *gorm.DB) ([]UserRelation, error) {
	var userRelation []UserRelation
	result := db.Find(&userRelation)

	return userRelation, result.Error
}

func GetUserRelationById(db *gorm.DB, id int) (UserRelation, error) {
	var userRelation UserRelation
	result := db.First(&userRelation, id)

	return userRelation, result.Error
}

func DeleteUserRelationById(db *gorm.DB, id int) error {
	return db.Delete(&UserRelation{}, id).Error
}

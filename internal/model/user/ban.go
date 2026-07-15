package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"context"
	"errors"

	"gorm.io/gorm"
)

type Ban struct {
	gorm.Model   `c_migrator:"enabled"`
	Id           int  `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	UserId       int  `json:"user_id" gorm:"type:int;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User         User `json:"user" gorm:"foreignKey:UserId;references:Id"`
	BannedUserId int  `json:"banned_user_id" gorm:"type:int"`
	BannedUser   User `json:"banned_user" gorm:"foreignKey:BannedUserId;references:Id"`
	Type         int  `json:"type" gorm:"type:int;not null"`
	CreatedAt    time.TimestampWithTimeZoneMicro
	UpdatedAt    time.TimestampWithTimeZoneMicro
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

var BanNotFoundErr = errors.New("ban not found")

func CreateBan(ctx context.Context, db *gorm.DB, b *Ban) error {
	return db.WithContext(ctx).Create(b).Error
}

func UpdateBanById(ctx context.Context, db *gorm.DB, b *Ban, id int) error {
	var ban Ban
	tx := db.WithContext(ctx)

	result := tx.First(&ban, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return BanNotFoundErr
	}

	return tx.Save(b).Error
}

func GetBans(ctx context.Context, db *gorm.DB) ([]Ban, error) {
	var bans []Ban
	result := db.WithContext(ctx).Find(&bans)

	return bans, result.Error
}

func GetBanById(ctx context.Context, db *gorm.DB, id int) (Ban, error) {
	var ban Ban
	result := db.WithContext(ctx).First(&ban, id)

	return ban, result.Error
}

func GetBansByUserId(ctx context.Context, db *gorm.DB, userId int) ([]Ban, error) {
	var bans []Ban
	result := db.WithContext(ctx).Model(&Ban{}).Where("user_id = ?", userId).Find(&bans)

	return bans, result.Error
}

func DeleteBanById(ctx context.Context, db *gorm.DB, id int) error {
	var ban Ban
	tx := db.WithContext(ctx)

	result := tx.First(&ban, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return BanNotFoundErr
	}

	return tx.Delete(&Ban{}, id).Error
}

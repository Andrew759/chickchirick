package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model `c_migrator:"enabled"`
	Id         int       `json:"id" gorm:"type:int;unique;primaryKey;autoIncrement"`
	Phone      string    `json:"phone" gorm:"type:varchar(30)"`
	Name       string    `json:"name" gorm:"type:varchar(256);not null"`
	Surname    string    `json:"surname" gorm:"type:varchar(256);not null"`
	Login      string    `json:"login" gorm:"type:varchar(256);unique; not null"`
	Property   *Property `json:"property" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Meta       *Meta     `json:"meta" gorm:"foreignKey:UserId"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type UserAlreadyExistErr struct {
	Field string
	Value string
}

func (e *UserAlreadyExistErr) Error() string {
	return fmt.Sprintf("user with %s '%s' already exists", e.Field, e.Value)
}

var UserNotFoundErr = errors.New("user not found")

func CreateUser(ctx context.Context, db *gorm.DB, u *User) error {
	var existing User
	tx := db.WithContext(ctx)

	err := tx.Where("login = ? OR phone = ?", u.Login, u.Phone).First(&existing).Error

	if err == nil {
		if existing.Login == u.Login {
			return &UserAlreadyExistErr{Field: "login", Value: u.Login}
		}
		if existing.Phone == u.Phone {
			return &UserAlreadyExistErr{Field: "phone", Value: u.Phone}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return tx.Create(u).Error
}

func UpdateUserById(ctx context.Context, db *gorm.DB, u *User, id int) error {
	var user User
	tx := db.WithContext(ctx)

	result := tx.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserNotFoundErr
	}

	return tx.Save(u).Error
}

func GetAllUsers(ctx context.Context, db *gorm.DB) ([]User, error) {
	var users []User
	result := db.WithContext(ctx).Find(&users)

	return users, result.Error
}

func GetUserById(ctx context.Context, db *gorm.DB, id int) (User, error) {
	var user User
	result := db.WithContext(ctx).First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, UserNotFoundErr
	}

	return user, result.Error
}

func DeleteUserById(ctx context.Context, db *gorm.DB, id int) error {
	var user User
	tx := db.WithContext(ctx)

	result := tx.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserNotFoundErr
	}

	return tx.Delete(&User{}, id).Error
}

func GetUserByUuid(ctx context.Context, db *gorm.DB, uuid string) (*User, error) {
	var user User

	err := db.WithContext(ctx).
		Joins("JOIN user_meta ON user_meta.user_id = users.id").
		Where("user_meta.user_uuid = ?", uuid).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, UserNotFoundErr
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByLoginOrPhone ищет пользователя по логину или номеру телефона.
func GetUserByLoginOrPhone(ctx context.Context, db *gorm.DB, loginOrPhone string) (*User, error) {
	var u User

	err := db.WithContext(ctx).
		Where("login = ? OR phone = ?", loginOrPhone, loginOrPhone).
		First(&u).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, UserNotFoundErr
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// SearchUsersByLoginOrPhone ищет пользователей по подстроке логина или телефона.
// excludeUserId — id текущего пользователя (0 = не исключать).
// limit — макс. число результатов (если <= 0, берётся 20).
func SearchUsersByLoginOrPhone(ctx context.Context, db *gorm.DB, query string, excludeUserId int, limit int) ([]User, error) {
	if limit <= 0 {
		limit = 20
	}
	q := strings.TrimSpace(query)
	if q == "" {
		return []User{}, nil
	}

	pattern := "%" + q + "%"
	tx := db.WithContext(ctx).Model(&User{}).
		Where("login ILIKE ? OR phone ILIKE ?", pattern, pattern)

	if excludeUserId > 0 {
		tx = tx.Where("id <> ?", excludeUserId)
	}

	var users []User
	err := tx.Order("login ASC").Limit(limit).Find(&users).Error
	return users, err
}

// SearchUsersWithUuid ищет по логину/телефону и возвращает id + uuid из user_meta.
func SearchUsersWithUuid(ctx context.Context, db *gorm.DB, query string, excludeUuid string, limit int) ([]UserWithUuid, error) {
	if limit <= 0 {
		limit = 20
	}
	q := strings.TrimSpace(query)
	if q == "" {
		return []UserWithUuid{}, nil
	}

	pattern := "%" + q + "%"
	tx := db.WithContext(ctx).
		Table("users AS u").
		Select("u.id, u.name, u.surname, u.login, user_meta.user_uuid::text AS user_uuid").
		Joins("JOIN user_meta ON user_meta.user_id = u.id").
		Where("u.login ILIKE ? OR u.phone ILIKE ?", pattern, pattern).
		Where("u.deleted_at IS NULL")

	if excludeUuid != "" {
		tx = tx.Where("user_meta.user_uuid::text <> ?", excludeUuid)
	}

	var rows []UserWithUuid
	err := tx.Order("u.login ASC").Limit(limit).Scan(&rows).Error
	return rows, err
}

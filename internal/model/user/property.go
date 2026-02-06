package user

import (
	"chickChirick/pkg/chirik_gorm_tweaks/time"
	"errors"

	"gorm.io/gorm"
)

type Property struct {
	//TODO: нужно указывать c_migrator_t_name сразу в двух местах при использовании GORM. Переделать
	gorm.Model `c_migrator:"enabled" c_migrator_t_name:"properties"`
	UserId     int     `json:"user_id" gorm:"unique;not null"`
	User       User    `json:"user" gorm:"foreignKey:UserId;references:Id"`
	Timezone   int8    `json:"timezone" gorm:"type:smallint;default:3"`
	Email      string  `json:"email" gorm:"type:varchar(256)"`
	Password   *string `json:"password" gorm:"type:varchar(1024)"`
	CreatedAt  time.TimestampWithTimeZoneMicro
	UpdatedAt  time.TimestampWithTimeZoneMicro
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

var PropertyNotFoundErr = errors.New("property not found")

var PropertyForUserAlreadyExistsErr = errors.New("property for user already exists")

func (p *Property) TableName() string {
	return "properties"
}

func (p *Property) SetPassword(password string) {
	p.Password = &password
}

func CreateProperty(db *gorm.DB, b *Property) error {
	hasProperty, err := HasPropertyByUserId(db, b.UserId)
	if err != nil {
		return err
	}
	if hasProperty {
		return PropertyForUserAlreadyExistsErr
	}

	return db.Create(b).Error
}

func UpdatePropertyById(db *gorm.DB, p *Property, id int) error {
	var property Property
	result := db.First(&property, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PropertyNotFoundErr
	}

	return db.Save(p).Error
}

func GetProperties(db *gorm.DB) ([]Property, error) {
	var properties []Property
	result := db.Find(&properties)

	return properties, result.Error
}

func GetPropertyById(db *gorm.DB, id int) (Property, error) {
	//TODO: доработать на подобии модели пользователя - добавить выбрасывание конкретного исключения,
	// что свойство не найдено. Затем поправить все остальные модели и контроллеры
	var property Property
	result := db.First(&property, id)

	return property, result.Error
}

func HasPropertyByUserId(db *gorm.DB, userId int) (bool, error) {
	var count int64
	err := db.Model(&Property{}).Where("user_id = ?", userId).Count(&count).Error

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func DeletePropertyById(db *gorm.DB, id int) error {
	var property Property
	result := db.First(&property, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return PropertyNotFoundErr
	}

	return db.Delete(&Property{}, id).Error
}

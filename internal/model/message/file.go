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

func CreateFile(db *gorm.DB, f *File) error {
	return db.Create(f).Error
}

func UpdateFile(db *gorm.DB, f *File) error {
	return db.Save(f).Error
}

func GetFiles(db *gorm.DB) ([]File, error) {
	var files []File
	result := db.Find(&files)

	return files, result.Error
}

func GetFileById(db *gorm.DB, id int) (File, error) {
	var file File
	result := db.First(&file, id)

	return file, result.Error
}

func DeleteFileById(db *gorm.DB, id int) error {
	return db.Delete(&File{}, id).Error
}

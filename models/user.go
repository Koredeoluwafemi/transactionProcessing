package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uint `gorm:"primary_key"`
	Firstname       string
	Lastname        string
	Email           string `gorm:"type:varchar(255);uniqueIndex"`
	Phone           string
	Password        string
	EmailVerifiedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (t *User) Create(db *gorm.DB) error {
	row := db.Create(&t)
	return row.Error
}
func (t *User) Update(db *gorm.DB, whereColumn string, value any) error {
	row := db.Model(&User{}).Where(whereColumn+" = ?", value).Updates(&t)
	return row.Error
}
func GetUser(db *gorm.DB, whereColumn map[string]any, preloads ...string) (User, error) {
	var record User
	rows := db
	if len(whereColumn) > 0 {
		for key, value := range whereColumn {
			rows = rows.Where(key+" = ?", value)
		}
	}
	if len(preloads) > 0 {
		for _, item := range preloads {
			rows = rows.Preload(item)
		}
	}
	rows.Last(&record)
	return record, rows.Error
}

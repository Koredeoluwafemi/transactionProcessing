package models

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID                 uint `gorm:"primary_key"`
	UserID             uint
	User               User
	OriginAccountID    uint
	OriginAccount      LedgerAccount
	RecipientAccountID uint
	RecipientAccount   LedgerAccount
	Amount             uint
	CompletedAt        time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (t *Transaction) Create(db *gorm.DB) error {
	row := db.Create(&t)
	return row.Error
}
func UpdateTransaction(db *gorm.DB, whereColumn string, value any, updateModel Transaction) error {
	row := db.Model(&Transaction{}).Where(whereColumn+" = ?", value).Updates(&updateModel)
	return row.Error
}

package models

import (
	"gorm.io/gorm"
	"time"
)

type LedgerAccount struct {
	ID            uint `gorm:"primary_key"`
	Balance       uint
	UserID        uint
	User          User
	AccountName   string
	AccountNumber string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (t *LedgerAccount) Create(db *gorm.DB) error {
	row := db.Create(&t)
	return row.Error
}
func UpdateAccountBalance(db *gorm.DB, id any, balance uint) error {
	row := db.Debug().Model(&LedgerAccount{}).Where("id = ?", id).UpdateColumn("balance", balance)
	return row.Error
}
func GetLedgerAccount(db *gorm.DB, whereValuesColumn map[string]any, preloads ...string) (LedgerAccount, error) {
	var record LedgerAccount
	rows := db.Debug()
	if len(whereValuesColumn) > 0 {
		for key, value := range whereValuesColumn {
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

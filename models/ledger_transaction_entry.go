package models

import (
	"gorm.io/gorm"
	"time"
)

type LedgerTransactionEntry struct {
	ID                  uint `gorm:"primary_key"`
	LedgerTransactionID uint
	LedgerTransaction   LedgerTransaction
	AccountID           uint
	Account             LedgerAccount
	Debit               uint
	Credit              uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (t *LedgerTransactionEntry) Create(db *gorm.DB) error {
	row := db.Create(&t)
	return row.Error
}

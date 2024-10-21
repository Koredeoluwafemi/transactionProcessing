package models

import (
	"gorm.io/gorm"
	"time"
)

type LedgerTransaction struct {
	ID            uint `gorm:"primary_key"`
	TransactionID uint
	Transaction   Transaction
	Name          string
	Debits        []LedgerTransactionEntry
	Credits       []LedgerTransactionEntry
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (t *LedgerTransaction) Create(db *gorm.DB) error {
	row := db.Create(&t)
	return row.Error
}

func GetLedgerTransaction(db *gorm.DB, whereColumn map[string]any, preloads ...string) (LedgerTransaction, error) {
	var record LedgerTransaction
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

package database

import (
	"log"
	"log/slog"
	"transpro/models"
)

func Migrate() {
	// Migrate the schema
	db := DB

	err := db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.LedgerAccount{},
		&models.LedgerTransaction{},
		&models.LedgerTransactionEntry{},
	)

	if err != nil {
		log.Println("migration", slog.Any("error", err))
	}

	Seed(db)
}

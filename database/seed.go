package database

import "gorm.io/gorm"

func Seed(db *gorm.DB) {

	userSeeder(db)
	ledgerAccountSeeder(db)

}

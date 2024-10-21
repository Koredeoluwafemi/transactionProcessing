package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
	"transpro/models"
)

func userSeeder(db *gorm.DB) {

	row := db.Where("email = ?", "user@transpro.com").First(&models.User{})
	if row.RowsAffected == 0 {
		//create password hash
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

		// register user
		db.Create(&models.User{
			Firstname:       "John",
			Lastname:        "Doe",
			Email:           "user@transpro.com",
			EmailVerifiedAt: time.Now(),
			Password:        string(passwordHash),
		})
	}

	row = db.Where("email = ?", "user2@transpro.com").First(&models.User{})
	if row.RowsAffected == 0 {
		//create password hash
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

		// register user
		db.Create(&models.User{
			Firstname:       "Seyi",
			Lastname:        "Man",
			Email:           "user2@transpro.com",
			EmailVerifiedAt: time.Now(),
			Password:        string(passwordHash),
		})
	}
}

func ledgerAccountSeeder(db *gorm.DB) {
	saveControl := models.LedgerAccount{AccountNumber: "123456", AccountName: "trader", UserID: 1, Balance: 100000}
	db.Where(saveControl).FirstOrCreate(&saveControl)
	saveControl = models.LedgerAccount{AccountNumber: "1234567", AccountName: "barber", UserID: 2, Balance: 500000}
	db.Where(saveControl).FirstOrCreate(&saveControl)
}

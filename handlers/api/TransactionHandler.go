package api

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
	"time"
	"transpro/database"
	"transpro/helper"
	"transpro/lock"
	"transpro/models"
)

type Account struct {
	AccountID          uint
	RecipientAccountID uint
	Amount             uint
	AvailableBalance   uint
}

type transactionEntry struct {
	Name   string
	Debit  insertEntry
	Credit insertEntry
}
type insertEntry struct {
	Amount    uint
	AccountID uint
}

// insertLedgerEntry ...
func insertLedgerEntry(db *gorm.DB, input transactionEntry, ledgerTransactionID uint) error {
	//insert debit leg
	debitTransactionEntry := models.LedgerTransactionEntry{
		LedgerTransactionID: ledgerTransactionID,
		AccountID:           input.Debit.AccountID,
		Debit:               input.Debit.Amount,
	}
	err := debitTransactionEntry.Create(db)
	if err != nil {
		return fmt.Errorf("unable to insert debit entry: %s", err)
	}

	//insert credit leg
	creditTransactionEntry := models.LedgerTransactionEntry{
		LedgerTransactionID: ledgerTransactionID,
		AccountID:           input.Credit.AccountID,
		Credit:              input.Credit.Amount,
	}
	err = creditTransactionEntry.Create(db)
	if err != nil {
		return fmt.Errorf("unable to insert credit entry: %s", err)
	}

	return nil
}

type TransactionInput struct {
	SenderAccountNumber    string `json:"sender_account_number"`
	RecipientAccountNumber string `json:"recipient_account_number"`
	Amount                 uint   `json:"amount"`
}

func (s TransactionInput) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.SenderAccountNumber, validation.Required),
		validation.Field(&s.RecipientAccountNumber, validation.Required),
		validation.Field(&s.Amount, validation.Required),
	)
}

func Transaction(c *fiber.Ctx) error {

	var input TransactionInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}
	if err := input.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}

	db := database.DB
	userID := helper.GetUserID(c)

	//check account balance
	senderAccount, err := models.GetLedgerAccount(db, map[string]any{"account_number": input.SenderAccountNumber, "user_id": userID})
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "account not found, permission denied!", "data": err})
	}

	recipientAccount, err := models.GetLedgerAccount(db, map[string]any{"account_number": input.RecipientAccountNumber})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "account not found", "data": err})
	}

	transaction := models.Transaction{
		UserID: userID, OriginAccountID: senderAccount.ID,
		RecipientAccountID: recipientAccount.ID, Amount: input.Amount,
	}
	err = transaction.Create(db)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}

	account := &Account{AccountID: senderAccount.ID, Amount: input.Amount, RecipientAccountID: recipientAccount.ID}
	err = account.ProcessAccount(db, transaction.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}

	err = models.UpdateTransaction(db, "id", transaction.ID, models.Transaction{CompletedAt: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": true, "message": "transaction successful", "data": ""})
}

func (a *Account) ProcessAccount(db *gorm.DB, transactionID uint) error {

	senderLock := lock.GetAccountLock(a.AccountID)
	recipientLock := lock.GetAccountLock(a.RecipientAccountID)

	senderLock.Lock()
	defer senderLock.Unlock()

	recipientLock.Lock()
	defer recipientLock.Unlock()

	tx := db.Begin()

	if a.AccountID == a.RecipientAccountID {
		return errors.New("credit account cannot be debit account")
	}

	originAccount, err := models.GetLedgerAccount(tx, map[string]any{"id": a.AccountID})
	if err != nil {
		return err
	}

	destinationAccount, err := models.GetLedgerAccount(tx, map[string]any{"id": a.RecipientAccountID})
	if err != nil {
		return err
	}

	if originAccount.Balance < a.Amount {
		return errors.New("insufficient balance")
	}

	debitEntry := insertEntry{Amount: a.Amount, AccountID: a.AccountID}
	creditEntry := insertEntry{Amount: a.Amount, AccountID: a.RecipientAccountID}
	ledgerEntry := transactionEntry{Debit: debitEntry, Credit: creditEntry}

	ledgerTransaction := models.LedgerTransaction{Name: "Money Transaction", TransactionID: transactionID}
	err = ledgerTransaction.Create(tx)
	if err != nil {
		tx.Rollback()
		return errors.New("unable to create ledger transaction")
	}

	err = insertLedgerEntry(tx, ledgerEntry, ledgerTransaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	//update source balance
	originAccount.Balance -= a.Amount
	err = models.UpdateAccountBalance(tx, a.AccountID, originAccount.Balance)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	//update recipient balance
	destinationAccount.Balance += a.Amount
	err = models.UpdateAccountBalance(tx, destinationAccount.ID, destinationAccount.Balance)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil

}

// BalanceResponse represents the balance response to be sent back
type BalanceResponse struct {
	AccountNumber string `json:"account_number"`
	Balance       uint   `json:"balance"`
}

// GetBalance handles the endpoint to get the balance of an account
func GetBalance(c *fiber.Ctx) error {

	accountNumber := c.Params("num")

	// Get the user ID (assuming user is authenticated)
	db := database.DB
	userID := helper.GetUserID(c)

	// Retrieve the account from the database
	account, err := models.GetLedgerAccount(db, map[string]any{"account_number": accountNumber, "user_id": userID})
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Account not found", "data": err.Error()})
	}

	// Return the account balance
	response := BalanceResponse{
		AccountNumber: account.AccountNumber,
		Balance:       account.Balance,
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": true, "message": "Balance retrieved", "data": response})
}

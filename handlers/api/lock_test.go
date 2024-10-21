package api

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"transpro/config"
	"transpro/database"
	_ "transpro/lock"
	"transpro/models"
)

// Test for Locking Mechanism in the Transaction Endpoint
func TestTransaction_LockingMechanism(t *testing.T) {

	// Initialize the Fiber app
	app := fiber.New()
	database.Start()
	db := database.DB
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk1MzA0OTQsIm5hbWUiOiJKb2huIERvZSIsInVpZCI6MX0.-RQd_edN-Nv11hNtN6svP1P8spRYEVMFZlGizB-5Zgk"

	// Define the routes to test
	jwtToken := jwtware.New(jwtware.Config{
		SigningKey: []byte(config.App.JWTKey),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Missing or malformed JWT", "status": false})
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid or expired JWT", "status": false})
			}
		},
	})
	app.Post("/transaction", jwtToken, Transaction)

	// Prepare concurrent transactions
	var wg sync.WaitGroup
	transactionRequest := `{
		"sender_account_number": "123456",
		"recipient_account_number": "1234567",
		"amount": 100
	}`

	senderAccount, _ := models.GetLedgerAccount(db, map[string]any{"account_number": "123456", "user_id": 1})
	recipientAccount, _ := models.GetLedgerAccount(db, map[string]any{"account_number": "1234567"})

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodPost, "/transaction", strings.NewReader(transactionRequest))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			// Simulate the HTTP request
			resp, err := app.Test(req, -1)

			// Assert no error occurred during the request
			assert.Nil(t, err)

			// Assert the response is OK (status code 200)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Retrieve the updated balances after the concurrent transactions
	updatedSenderAccount, _ := models.GetLedgerAccount(db, map[string]any{"account_number": "123456", "user_id": 1})
	updatedRecipientAccount, _ := models.GetLedgerAccount(db, map[string]any{"account_number": "1234567"})

	// Assert that the balances are updated correctly
	expectedSenderBalance := senderAccount.Balance - 10*100 // 10 transactions of 100 units each
	expectedRecipientBalance := recipientAccount.Balance + 10*100

	assert.Equal(t, expectedSenderBalance, updatedSenderAccount.Balance)
	assert.Equal(t, expectedRecipientBalance, updatedRecipientAccount.Balance)
}

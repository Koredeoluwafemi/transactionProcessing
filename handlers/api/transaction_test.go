package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"transpro/config"
	"transpro/database"
)

func TestTransaction(t *testing.T) {
	// Create a Fiber app
	app := fiber.New()
	database.Start()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk1MTkzNjIsIm5hbWUiOiJKb2huIERvZSIsInVpZCI6MX0.E9EZNJqtdrHnAYPibESYDzJ_HQWpqO1y9M8GtHa0n7E"

	// Test cases
	tests := []struct {
		name        string
		input       TransactionInput
		statusCode  int
		BearerToken string
	}{
		{
			name: "Valid transaction",
			input: TransactionInput{
				SenderAccountNumber:    "123456",
				RecipientAccountNumber: "1234567",
				Amount:                 500,
			},
			statusCode:  http.StatusOK,
			BearerToken: token,
		},
		{
			name: "Invalid sender account",
			input: TransactionInput{
				SenderAccountNumber:    "invalid",
				RecipientAccountNumber: "1234567",
				Amount:                 500,
			},
			statusCode:  http.StatusUnauthorized,
			BearerToken: token,
		},
		{
			name: "Invalid recipient account",
			input: TransactionInput{
				SenderAccountNumber:    "123456",
				RecipientAccountNumber: "invalid",
				Amount:                 500,
			},
			statusCode:  http.StatusBadRequest,
			BearerToken: token,
		},
		{
			name: "Insufficient balance",
			input: TransactionInput{
				SenderAccountNumber:    "123456",
				RecipientAccountNumber: "1234567",
				Amount:                 100000000,
			},
			statusCode:  http.StatusBadRequest,
			BearerToken: token,
		},
	}

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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Marshal the input to JSON
			jsonBody, err := json.Marshal(test.input)
			assert.NoError(t, err)

			// Create a test request
			req := httptest.NewRequest("POST", "/transaction", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+test.BearerToken)

			// Record the response
			res, err := app.Test(req)
			assert.NoError(t, err)

			// Check the status code
			assert.Equal(t, test.statusCode, res.StatusCode)
		})
	}
}

func TestGetBalance(t *testing.T) {
	// Create a Fiber app
	app := fiber.New()
	database.Start()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk1MTkzNjIsIm5hbWUiOiJKb2huIERvZSIsInVpZCI6MX0.E9EZNJqtdrHnAYPibESYDzJ_HQWpqO1y9M8GtHa0n7E"

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

	// Define the route
	app.Get("/balance/:num", jwtToken, GetBalance)

	// Test cases
	tests := []struct {
		name        string
		accountNum  string
		statusCode  int
		BearerToken string
	}{
		{
			name:        "Valid account",
			accountNum:  "123456",
			statusCode:  http.StatusOK,
			BearerToken: token,
		},
		{
			name:        "Invalid account",
			accountNum:  "invalid",
			statusCode:  http.StatusNotFound,
			BearerToken: token,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test request
			req := httptest.NewRequest("GET", "/balance/"+test.accountNum, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+test.BearerToken)

			// Record the response
			res, err := app.Test(req)
			assert.NoError(t, err)

			// Check the status code
			assert.Equal(t, test.statusCode, res.StatusCode)
		})
	}
}

//func TestAccountLocking(t *testing.T) {
//	// Create an account
//	account := &Account{
//		AccountID:        1,
//		AvailableBalance: 1000,
//	}
//
//	// Test concurrent updates
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			account.Mu.Lock()
//			account.AvailableBalance += 100
//			account.Mu.Unlock()
//		}()
//	}
//
//	wg.Wait()
//
//	// Verify the final balance
//	if account.AvailableBalance != 2000 {
//		t.Errorf("Expected balance 2000, got %d", account.AvailableBalance)
//	}
//}
//
//func TestAccountLocking_ConcurrentReads(t *testing.T) {
//	// Create an account
//	account := &Account{
//		AccountID:        1,
//		AvailableBalance: 1000,
//	}
//
//	// Test concurrent reads
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			account.Mu.RLock()
//			_ = account.AvailableBalance
//			account.Mu.RUnlock()
//		}()
//	}
//
//	wg.Wait()
//}
//
//func TestAccountLocking_ConcurrentReadWrite(t *testing.T) {
//	// Create an account
//	account := &Account{
//		AccountID:        1,
//		AvailableBalance: 1000,
//	}
//
//	// Test concurrent read-write
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			if i%2 == 0 {
//				// Write
//				account.Mu.Lock()
//				account.AvailableBalance += 100
//				account.Mu.Unlock()
//			} else {
//				// Read
//				account.Mu.RLock()
//				_ = account.AvailableBalance
//				account.Mu.RUnlock()
//			}
//		}()
//	}
//
//	wg.Wait()
//}
//
//func TestAccountLocking_Deadlock(t *testing.T) {
//	// Create two accounts
//	account1 := &Account{
//		AccountID:        1,
//		AvailableBalance: 1000,
//	}
//	account2 := &Account{
//		AccountID:        2,
//		AvailableBalance: 2000,
//	}
//
//	// Test deadlock scenario
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			account1.Mu.Lock()
//			time.Sleep(10 * time.Millisecond)
//			account2.Mu.Lock()
//			account2.Mu.Unlock()
//			account1.Mu.Unlock()
//		}()
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			account2.Mu.Lock()
//			time.Sleep(10 * time.Millisecond)
//			account1.Mu.Lock()
//			account1.Mu.Unlock()
//			account2.Mu.Unlock()
//		}()
//	}
//
//	wg.Wait()
//}

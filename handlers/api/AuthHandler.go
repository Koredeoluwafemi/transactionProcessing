package api

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"transpro/config"
	"transpro/database"
	"transpro/models"
)

type loginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

func (s loginRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required),
		validation.Field(&s.Password, validation.Required),
	)
}

func Login(c *fiber.Ctx) error {

	var request loginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}
	if err := request.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}

	db := database.DB
	password := strings.TrimSpace(request.Password)
	email := strings.TrimSpace(request.Email)

	user, err := models.GetUser(db, map[string]any{"email": email})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": err.Error(), "data": err})
	}

	//user has been verified
	hash := []byte(user.Password)
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Unable to login; invalid email or password", "data": err})
	}

	// Create token
	tokenResponse, err := setAccessToken(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Unable to generate token", "data": err})
	}

	output := fiber.Map{
		"uid":                    user.ID,
		"firstname":              user.Firstname,
		"lastname":               user.Lastname,
		"email":                  user.Email,
		"phone":                  user.Phone,
		"token":                  tokenResponse.AccessToken,
		"token_expiry_timestamp": tokenResponse.ExpiryTime.Format("2006-01-02 03:04:05"),
		"token_expiry_unix":      tokenResponse.ExpiryTime.Unix(),
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": true, "message": "success", "data": output})

}

type tokenBox struct {
	AccessToken string
	ExpiryTime  time.Time
}

func setAccessToken(user models.User) (tokenBox, error) {

	var response tokenBox
	accessToken := jwt.New(jwt.SigningMethodHS256)

	expiryTime := time.Now().Add(time.Minute * time.Duration(180))
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["name"] = user.Firstname + " " + user.Lastname
	claims["exp"] = expiryTime.Unix()

	// Generate token and send as response.
	accessTokenString, err := accessToken.SignedString([]byte(config.App.JWTKey))
	if err != nil {
		return response, err
	}

	response.AccessToken = accessTokenString
	response.ExpiryTime = expiryTime

	return response, nil
}

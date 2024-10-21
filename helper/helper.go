package helper

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"strconv"
	"strings"
	"transpro/config"
)

func GetUserID(c *fiber.Ctx) uint {
	if c.Get("Authorization") == "" {
		return 0
	}
	tokenString := strings.ReplaceAll(c.Get("Authorization"), "Bearer ", "")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App.JWTKey), nil
	})
	if err != nil {
		return 0
	}

	claims := token.Claims.(jwt.MapClaims)

	if value, exist := claims["uid"]; exist {
		stringData := TransToString(value)
		uid, err := strconv.Atoi(stringData)
		if err != nil {
			return 0
		}
		return uint(uid)
	}

	return 0
}

func TransToString(data interface{}) (res string) {
	switch v := data.(type) {
	case float64:
		res = strconv.FormatFloat(data.(float64), 'f', 0, 64)
	case float32:
		res = strconv.FormatFloat(float64(data.(float32)), 'f', 6, 32)
	case int:
		res = strconv.FormatInt(int64(data.(int)), 10)
	case int64:
		res = strconv.FormatInt(data.(int64), 10)
	case uint:
		res = strconv.FormatUint(uint64(data.(uint)), 10)
	case uint64:
		res = strconv.FormatUint(data.(uint64), 10)
	case uint32:
		res = strconv.FormatUint(uint64(data.(uint32)), 10)
	case json.Number:
		res = data.(json.Number).String()
	case string:
		res = data.(string)
	case []byte:
		res = string(v)
	default:
		res = ""
	}
	return
}

func GetRoot() string {
	//get current directory
	path, err := os.Getwd()

	//get parent path
	if err != nil {
		log.Println(err)
	}

	return path + "/" + "resources"
}

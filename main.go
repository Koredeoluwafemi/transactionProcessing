package main

import (
	"github.com/gofiber/fiber/v2"
	"transpro/config"
	"transpro/database"
	"transpro/helper"
	_ "transpro/lock"
	"transpro/routes"
)

func main() {

	//get root directory
	resourcesPath := helper.GetRoot()

	database.Start()
	database.Migrate()

	fiberConfig := fiber.Config{}
	app := fiber.New(fiberConfig)
	routes.Routes(app)

	app.Static("/", resourcesPath)

	port := config.App.Port

	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}

package main

import (
	"./database"
	"./routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	db := database.Initialize("./database/db.sqlite")
	database.Migrate(db)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/settings", routes.GetDeviceSetting(db))
	e.POST("/settings", routes.SaveDeviceSettings(db))
	e.GET("/simulate", routes.SimulatePriceChanges(db))

	e.Logger.Fatal(e.Start(":9000"))
}

package main

import (
	"application/config"
	"application/handlers"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}
	e := echo.New()
	handlers.InitRoutes(e)
	e.Logger.Debug(e.Start(":8000"))
}

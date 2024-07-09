package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"vortex/controllers"
	"vortex/migrations"
	"vortex/repository"
	"vortex/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var config Config

func main() {
	_, exists := os.LookupEnv("PROD")

	if !exists {
		err := godotenv.Load("./dev.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	config.Load()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	migrations.Migrate(db)
	repository.Init(db)

	var clientService services.ClientService
	clientService.InitService()
	var checkScheduleService services.CheckScheduleService
	checkScheduleService.InitService(&clientService)
	go checkScheduleService.StartWatcher()

	var clientsController controllers.ClientController
	clientsController.Init(&clientService)
	e := echo.New()
	e.POST("/", clientsController.Add)
	e.PUT("/", clientsController.Update)
	e.DELETE("/", clientsController.Delete)
	e.GET("/", clientsController.UpdateStatuses)
	e.Logger.Fatal(e.Start(":" + config.HTTPPort))
}

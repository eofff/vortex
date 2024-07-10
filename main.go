package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"vortex/controllers"
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

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var algorithmStatusRepository repository.IAlgorithmStatusRepository = &repository.AlgorithmStatusRepository{}
	algorithmStatusRepository.Init(db)

	var algorithmStatusService services.IAlgorithmStatusService = &services.AlgorithmStatusService{}
	algorithmStatusService.Init(algorithmStatusRepository)

	var clientService services.IClientService = &services.ClientService{}
	clientService.Init(algorithmStatusService)
	var checkScheduleService services.CheckScheduleService
	checkScheduleService.Init(clientService)
	checkScheduleService.StartWatcher()

	e := echo.New()
	var clientsController controllers.ClientController
	clientsController.Init(clientService, e)
	e.Logger.Fatal(e.Start(":" + config.HTTPPort))
}

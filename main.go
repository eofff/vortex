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
	// load env values from file if env variable PROD not exists
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

	// check db connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var algorithmStatusRepository repository.IAlgorithmStatusRepository = &repository.AlgorithmStatusRepository{}
	algorithmStatusRepository.Init(db)

	// initialize services
	var algorithmStatusService services.IAlgorithmStatusService = &services.AlgorithmStatusService{}
	algorithmStatusService.Init(algorithmStatusRepository)

	var deployerService services.Deployer = &services.DeployerService{}

	var clientService services.IClientService = &services.ClientService{}
	clientService.Init(algorithmStatusService, deployerService)
	var checkScheduleService services.CheckScheduleService
	checkScheduleService.Init(clientService)

	// start 5 minutes autochecker
	checkScheduleService.StartWatcher()

	// initialize instance of echo http server https://echo.labstack.com/
	e := echo.New()
	var clientsController controllers.ClientController
	clientsController.Init(clientService, e)
	e.Logger.Fatal(e.Start(":" + config.HTTPPort))
}

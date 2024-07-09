package main

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort   string
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
}

func (c *Config) Load() {
	var exists bool
	var err error
	c.DbHost, exists = os.LookupEnv("DBHOST")
	if !exists {
		log.Fatal("there is no DBHOST var")
	}

	port, exists := os.LookupEnv("DBPORT")
	if !exists {
		log.Fatal("there is no DBPORT var")
	}

	c.DbPort, err = strconv.Atoi(port)
	if err != nil {
		log.Fatal("DBPORT var is incorrect")
	}

	c.DbUser, exists = os.LookupEnv("DBUSER")
	if !exists {
		log.Fatal("there is no DBUSER var")
	}

	c.DbPassword, exists = os.LookupEnv("DBPASSWORD")
	if !exists {
		log.Fatal("there is no DBPASSWORD var")
	}

	c.DbName, exists = os.LookupEnv("DBNAME")
	if !exists {
		log.Fatal("there is no DBNAME var")
	}

	c.HTTPPort, exists = os.LookupEnv("HTTP_PORT")
	if !exists {
		log.Fatal("there is no HTTP_PORT var")
	}
}

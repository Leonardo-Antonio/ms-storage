package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env Enviroment

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	Env = Enviroment{
		Micro: Micro{
			Name: os.Getenv("APP_NAME"),
			Port: os.Getenv("APP_PORT"),
			Host: os.Getenv("APP_HOST"),
		},
	}
}

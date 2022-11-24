package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"test/pkg/handler"
	"test/pkg/repository"
	"test/pkg/service"
)

// @title          Server API
// @version        0.0.1
// @description    Work with server.
// @termsOfService http://swagger.io/terms/

// @contact.name Tairyoku
// @contact.url  http://www.swagger.io/tairyoku

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath /

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := repository.NewRepositoryDB(repository.Config{
		Username: os.Getenv("DBUsername"),
		Password: os.Getenv("DBPassword"),
		Host:     os.Getenv("DBHost"),
		Url:      os.Getenv("DBUrl"),
		DBName:   os.Getenv("DBName"),
	})
	if err != nil {
		log.Fatalf("error %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(service.Server)
	if err := server.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error %s", err.Error())
	}

}

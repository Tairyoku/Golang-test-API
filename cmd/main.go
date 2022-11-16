package main

import (
	"log"
	"test"
	"test/pkg/handler"
	"test/pkg/repository"
	"test/pkg/service"
)

func main() {

	db, err := repository.NewMysqlDB(repository.Config{
		Username: "root",
		Password: "",
		Host:     "tcp",
		Url:      "127.0.0.1:3306",
		DBName:   "myFirstDB",
	})
	if err != nil {
		log.Fatalf("error %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(test.Server)
	if err := server.Run(":8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error %s", err.Error())
	}

}

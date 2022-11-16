package main

import (
	"log"
	"test"
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

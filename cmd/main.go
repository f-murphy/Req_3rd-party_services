package main

import (
	"req3rdPartyServices/handler"
	"req3rdPartyServices/repository"
	"req3rdPartyServices/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func main() {
 	db, err := sqlx.Connect("postgres", "postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable")
	if err != nil {
		logrus.Fatal("failed to initialize db: ", err.Error())
	}

	repos := repository.NewTaskRepository(db)
	services := service.NewTaskService(repos)
	handlers := handler.NewTaskHandler(services)

	r := gin.Default()
	r.POST("/task", handlers.CreateTask)
	r.GET("/task/:id", handlers.GetTask)
	r.Run(":8080")
}

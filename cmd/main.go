package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"req3rdPartyServices/task"
)

func main() {
	r := gin.Default()
	fmt.Println("Listening on port 8080")

	r.POST("/task", task.PostTask)
	r.GET("/task/:id", task.GetTaskStatus)
	
	err := r.Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}

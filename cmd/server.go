package cmd

import (
	"fmt"
	"log"
	"net/http"
	"req3rdPartyServices/internal/service"
)

func StartServer() {
	fmt.Println("Listening on port 8080")
	http.HandleFunc("/tasks", service.TaskHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

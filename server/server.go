package server

import (
	"fmt"
	"log"
	"net/http"
	"req3rdPartyServices/task"
)

func StartServer() {
	fmt.Println("Listening on port 8080")
	http.HandleFunc("/task", task.RouteRedirection)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

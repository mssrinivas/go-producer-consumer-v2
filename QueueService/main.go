package main

import (
	"log"
	"net/http"
)

var logger *log.Logger

func HttpKeepAlive(port string) {
	errChan := make(chan error)
	go func() {
		log.Println("HTTP KeepAlive :transport", "HTTP", "started on port", port)
		errChan <- http.ListenAndServe(port, nil)
	}()
	log.Fatal("exit", <-errChan)
}

func main() {
	port := GetValueFromEnvVariable("ENV_PORT", ":8085")
	consumer := GetValueFromEnvVariable("QUEUE_URL", "http://localhost:8085")
	queueservice := svc.QueueService{
		StopChan: make(chan bool)
	}
	queueservice.Initialize()
	HttpKeepAlive(port)
}
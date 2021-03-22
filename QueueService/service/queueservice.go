package service

import (
	v1 "QueueService/contracts"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var ConsumerBuffer chan *v1.Task

type QueueService struct {
	Router   *mux.Router
	StopChan chan bool
	Handler  RequestHandler
}

func (q *QueueService) Initialize() {
	requestCollector := NewRequestHandler(5)
	ConsumerBuffer = make(chan *v1.Task, 5)
	http.HandleFunc("/queue/start", q.InitializeQueueService)
	http.HandleFunc("/queue/stop", q.StopQueue)

	http.HandleFunc("/queue/produce", requestCollector.AddTask)
	http.HandleFunc("/queue/status", requestCollector.CheckQueue)
	http.HandleFunc("/queue/ack", requestCollector.ReceiveAckStatus)
}

func (q *QueueService) InitializeQueueService(_ http.ResponseWriter, _ *http.Request) {
	log.Print("Queue Started")
	go q.StartQueue()
}

func (q *QueueService) StartQueue() {
	for {
		select {
		case <-q.StopChan:
			log.Print("Queue Stopped")
			break
		}
	}
}

func (q *QueueService) StopQueue(_ http.ResponseWriter, _ *http.Request) {
	go func() {
		q.StopChan <- true
	}()
}

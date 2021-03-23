package service

import (
	v1 "Consumer/contracts"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var ctx = context.Background()

type ConsumerService struct {
	Collector Collector
	Router    *mux.Router
	StopChan  chan bool
	Queue     string
	Redis     *redis.Client
}

func (csc *ConsumerService) Initialize() {
	csc.Collector = NewCollector()
	http.HandleFunc("/consumer/start", csc.InitializeConsumer)
	http.HandleFunc("/consumer/stop", csc.StopConsumer)
}

// Function to capture the task-requests in buffered channels and produce to Consumer
func (csc *ConsumerService) InitializeConsumer(_ http.ResponseWriter, _ *http.Request) {
	log.Print("Consumer Started")
	go csc.StartConsumer()
	go csc.Collector.ConsumeTasks()
}

func (csc *ConsumerService) StartConsumer() {
	for {
		select {
		case task := <-TaskChan:
			// if successful then do nothing print success
			log.Print("task", task.TaskName, "consumed successfully")
			resp := csc.WriteData(task)
			if resp == false {
				csc.QueueClient(v1.ConsumerResponse{
					TaskName: task.TaskName,
					Status:   false,
				})
				log.Print("task", task.TaskName, "failed to be consumed", task.TaskName)
			} else {
				csc.QueueClient(v1.ConsumerResponse{
					TaskName: task.TaskName,
					Status:   true,
				})
				log.Print("task", task.TaskName, "consumed successfully")
			}
		case <-csc.StopChan:
			log.Print("Consumer Stopped")
			break
		}
	}
}

func (csc *ConsumerService) StopConsumer(_ http.ResponseWriter, _ *http.Request) {
	go func() {
		csc.StopChan <- true
	}()
}

func (csc *ConsumerService) QueueClient(response v1.ConsumerResponse) error {
	client := http.Client{}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(response)
	request, err := http.NewRequest("POST", "http://localhost:8088/queue/ack", bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		log.Fatal("Unable to POST task to QueueService")
		return err
	}

	_, err = client.Do(request)
	if err != nil {
		log.Print("Unexpected response from QueueService")
		return err
	}
	return nil
}

func (csc *ConsumerService) WriteData(task v1.Task) bool {
	key := task.TaskName + "_" + task.LastUpdateTime
	value := v1.Task{
		TaskName:       task.TaskName,
		TaskType:       task.TaskType,
		LastUpdateTime: task.LastUpdateTime,
		ScheduledTime:  task.ScheduledTime,
		Periodicity:    task.Periodicity,
		TaskStatus:     task.TaskStatus,
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(value)
	err := csc.Redis.Set(ctx, key, reqBodyBytes.Bytes(), 0).Err()
	if err != nil {
		return false
	}
	log.Print(key, "successfully added to Redis")
	return true
}

// TaskDispatcher fetches data from redis and updates triggered tasks
func (csc *ConsumerService) GetTaskData() {
	// return the respective record for ex:
	return
}

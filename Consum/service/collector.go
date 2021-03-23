
package service

import (
	v1 "Consumer/contracts"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	ErrIncorrectTaskFields = "task fields cannot be empty"
	ErrInvalidTimeFormat   = "invalid time format in request"
	ErrInvalidPeriodicity  = "invalid periodicity value"
)

type Collector struct {
	logger *log.Logger
}

func NewCollector() Collector {
	return Collector{}
}

// A buffered channel that captures tasks to be produced to Consumer
var TaskChan = make(chan v1.Task, 10)

// Function to capture the task-requests in buffered channels and produce to Consumer
func (c *Collector) ConsumeTasks() {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
			case _ = <-ticker.C:
				RequestColl()
			}
		}
	}()
	<-done
}

func RequestColl() {
	taskRequest := v1.Task{}
	client := http.Client{}
	request, err := http.NewRequest("GET", "http://localhost:8088/queue/status", nil)
	if err != nil {
		log.Print("Unable to POST task to Queue")
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Print("Unexpected response from QueueService", err)
	}

	respBody := extractResponseBody(resp)
	if len(respBody) > 0 {
		err = json.Unmarshal(respBody, &taskRequest)
		TaskChan <- taskRequest
	} else {
		log.Print("Empty Queue")
	}
}

func extractResponseBody(resp *http.Response) []byte {
	body := ""
	if resp.Body != nil {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			body = string(bytes)
		}
	}
	return []byte(body)
}

// Function to validate the request from the client
func (ctr *Collector) ValidateRequest(task v1.Task) error {
	if task.TaskName == "" || task.TaskType == "" || task.LastUpdateTime == "" ||
		task.ScheduledTime == "" {
		ctr.logger.Fatal(
			"service", "Consumer",
			"method", "ValidateRequest",
			"error", "task fields cannot be empty")
		return errors.New(ErrIncorrectTaskFields)
	}

	//err := ctr.ValidateRequestTimeFields(task.LastUpdateTime, "lastUpdateTime")
	//if err != nil {
	//	return err
	//}
	//
	//err = ctr.ValidateRequestTimeFields(task.ScheduledTime, "scheduledTime")
	//if err != nil {
	//	return err
	//}

	if task.Periodicity < 0 {
		ctr.logger.Fatal(
			"service", "Consumer",
			"method", "ValidateRequest",
			"error", "periodicity should be positive")
		return errors.New(ErrInvalidPeriodicity)
	}

	return nil
}

func (ctr *Collector) ValidateRequestTimeFields(req string, timerType string) error {
	_, err := time.Parse(time.RFC3339, req)
	if err != nil {
		ctr.logger.Fatal(
			"service", "Consumer",
			"method", "ValidateRequest",
			"error", "Invalid format", timerType)
		return errors.New(ErrInvalidTimeFormat)
	}
	return nil
}

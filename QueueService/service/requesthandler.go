package service

import (
	v1 "QueueService/contracts"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestHandler struct {
	Queue  TaskQueue
	logger *log.Logger
}

func NewRequestHandler(maxBuffer int) RequestHandler {
	taskQueue := NewTaskQueue(maxBuffer)
	return RequestHandler{
		Queue: taskQueue,
	}
}

// Function to capture the task-requests into circular queue
func (r *RequestHandler) AddTask(_ http.ResponseWriter, req *http.Request) {
	task := v1.Task{}
	requestBody := extractRequestBody(req)
	err := json.Unmarshal(requestBody, &task)
	if err != nil {
		return
	}

	err = ValidateRequest(task)
	if err != nil {
		return
	}

	r.Queue.Enqueue(&task)

	return
}

// Function to check for available tasks in circular queue
func (r *RequestHandler) CheckQueue(resp http.ResponseWriter, _ *http.Request) {
	// Add code
	// Check if a task exists in the queue
	if r.Queue.TaskList[r.Queue.RearIndex] != nil {
		str := fmt.Sprintf("%v", *r.Queue.TaskList[r.Queue.RearIndex])
		r.Queue.RearIndex++
		resp.Write([]byte(str))
	} else {
		// if successful then do nothing print success
		resp.Write(nil)
	}
}

// Function to capture the task-consumption acknowledgments
func (r *RequestHandler) ReceiveAckStatus(_ http.ResponseWriter, req *http.Request) {
	reqBody := extractRequestBody(req)
	consumerResponse := v1.ConsumerResponse{}
	err := json.Unmarshal(reqBody, &consumerResponse)
	if err != nil {
		return
	}
	// Release the lock only after a successful status
	if consumerResponse.Status {
		r.Queue.DeQueue(consumerResponse.TaskName)
	}
	return
}


// Helper functions
func extractRequestBody(req *http.Request) []byte {
	body := ""
	if req.Body != nil {
		bytes, err := ioutil.ReadAll(req.Body)
		if err == nil {
			body = string(bytes)
		}
	}
	return []byte(body)
}

// Function to validate the request from the producer
func ValidateRequest(task v1.Task) error {
	return nil
}

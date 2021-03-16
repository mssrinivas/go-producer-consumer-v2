package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RequestHandler struct {
	Queue  TaskQueue
	logger *log.Logger
}

func NewRequestHandler() Collector {
	taskQueue := NewTaskQueue()
	return RequestHandler{
		Queue: taskQueue,
	}
}

// Function to capture the task-requests into circular queue
func (r *RequestHandler) AddTask(_ http.ResponseWriter, _ *http.Request) {
	task := v1.Task{}
	requestBody := extractRequestBody(req)
	err := json.Unmarshal(requestBody, &task)
	if err != nil {
		//log error
		return
	}

	err = ValidateRequest(task)
	if err != nil {
		return
	}

	r.Queue.Enqueue(task)
	return
}

// Function called when a task is consumed by the consumer
func (r *RequestHandler) ConsumeTask(_ http.ResponseWriter, _ *http.Request) http.Response {
	// Add code
	task := r.Queue.Dequeue()
	// return a HTTP response with the task marshalled in the body
}

// Function to check for available tasks in circular queue
func (r *RequestHandler) CheckQueueStatus(_ http.ResponseWriter, _ *http.Request) {
	// Add code
	// Check if a task exists in the queue
	return

}

// Function to capture the task-consumption acknowledgments
func (r *RequestHandler) CheckAckStatus(_ http.ResponseWriter, req *http.Request) {
	// If requestBody contains Ack true then remove from Queue
	reqBody := extractRequestBody(req)
	err := json.Unmarshal(requestBody, &ConsumerResponse)
	if err != nil {
		//log error
		return
	}
	// Release the lock only after a successful status
	if ConsumerResponse.Status {
		r.Queue.MutexMap[task.taskName].UnLock()
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
}

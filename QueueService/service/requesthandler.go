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
  Queue     TaskQueue 
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

	q.Queue.Enqueue(task)
  MutexMap[task.taskName].lock()
	return

}

//
func (r *RequestHandler) ConsumeTask(_ http.ResponseWriter, _ *http.Request) v1.Task {
  task := r.Queue.Dequeue()
	return v1.Task
}


// Function to capture the task-requests into circular queue
func (r *RequestHandler) CheckQueueStatus(_ http.ResponseWriter, _ *http.Request) {
	// Check if rear != front to conclude that a task exists in queue
	return

}

// Function to capture the task-requests into circular queue
func (r *RequestHandler) CheckAckStatus(_ http.ResponseWriter, req *http.Request) {
	// If requestBody contains Ack true then remove from Queue
  reqBody := extractRequestBody(req)
  err := json.Unmarshal(requestBody, &ConsumerResponse)
	if err != nil {
		//log error
		return
	}

  if ConsumerResponse.staus {
    MutexMap[task.taskName].unlock(task.taskName)
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

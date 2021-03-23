package service

import (
	v1 "QueueService/contracts"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestHandler struct {
	Queue  *TaskQueue
	logger *log.Logger
}

func NewRequestHandler(maxBuffer int) RequestHandler {
	taskQueue := NewTaskQueue(maxBuffer)
	return RequestHandler{
		Queue: &taskQueue,
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
	// Check if a task exists in the queue
	peekElement, err := r.Queue.PeekQueue()
	if err != nil {
		r.logger.Fatal("Error fetching peek of the queue")
	}
	if peekElement != nil {
		jData, err := json.Marshal(*peekElement)
		if err != nil {
			// handle error
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(jData)
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

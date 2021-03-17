package service

import (
	v1 "QueueService/contracts"
)

type TaskQueue struct {
	FrontIndex int
	RearIndex  int
	TaskList   []v1.Task
	MutexMap   map[string]sync.Mutex
	IndexMap.  map[string]int
}

func NewTaskQueue(maxBuffer int) TaskQueue {
	queue := TaskQueue{
		FrontIndex: 0,
		RearIndex:  0,
		TaskList:   make([]v1.Task, maxBuffer),
		MutexMap:   make([string]sync.Mutex, maxBuffer),
	}
	return queue
}

// Method to add a task into the queue
func (queue TaskQueue) Enqueue(task v1.Task) error {
	// Add code
	// Perform the lock on the particular task
	MutexMap[task.TaskName].Lock()
	return
}

// Method to check for the next element to be processed in the queue
func (queue TaskQueue) PeekQueue() (task, error) {
	// Add code)
	return
}

// Method to remove a processed task from the queue
func (queue TaskQueue) DeQueue() v1.Task {
	// Add code
	return
}

package service

import (
	v1 "QueueService/contracts"
	"errors"
	"sync"
)

var Mutex *sync.Mutex

type TaskQueue struct {
	FrontIndex int
	RearIndex  int
	TaskList   []*v1.Task
	IndexMap   map[string]int
}

func NewTaskQueue(maxBuffer int) TaskQueue {
	Mutex = &sync.Mutex{}
	queue := TaskQueue{
		FrontIndex: 0,
		RearIndex:  0,
		TaskList:   make([]*v1.Task, maxBuffer),
		IndexMap:   make(map[string]int, maxBuffer),
	}
	return queue
}

// Method to add a task into the queue
func (q *TaskQueue) Enqueue(task *v1.Task) (bool, error) {
	// Add code
	// Perform the lock on the particular task
	if q.FrontIndex == len(q.TaskList)-1 {
		if q.TaskList[q.FrontIndex] != nil {
			err := errors.New("queue buffer full")
			return false, err
		} else {
			q.FrontIndex = 0
		}
	}

	Mutex.Lock()
	q.TaskList[q.FrontIndex] = task
	q.IndexMap[task.TaskName] = q.FrontIndex
	q.FrontIndex++
	Mutex.Unlock()
	return true, nil
}

// Method to check for the next element to be processed in the queue
func (q *TaskQueue) PeekQueue() (*v1.Task, error) {
	var task *v1.Task
	Mutex.Lock()
	if q.RearIndex == len(q.TaskList)-1 {
		q.RearIndex = 0
	} else {
		q.RearIndex++
	}
	if q.TaskList[q.RearIndex] != nil {
		task = q.TaskList[q.RearIndex]
	}
	Mutex.Unlock()
	return task, nil
}

// Method to remove a processed task from the queue
func (q *TaskQueue) DeQueue(taskName string) {
	Mutex.Lock()
	if val, ok := q.IndexMap[taskName]; ok {
		q.TaskList[val] = nil
		delete(q.IndexMap, taskName)
	}
	Mutex.Unlock()
	return
}

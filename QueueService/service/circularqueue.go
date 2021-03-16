import(
  v1 "QueueService/contracts"
)

type TaskQueue struct {
	 FrontIndex int 
	 RearIndex  int
   TaskList   []v1.Task
   MutexMap   map[string]&sync.Mutex{}
}

func NewTaskQueue(maxBuffer int) TaskQueue{
  queue := TaskQueue{
    FrontIndex: 0,
    RearIndex: 0,
    TaskList: make([]v1.Task, maxBuffer),
    MutexMap:  make([string]&sync.Mutex{}, maxBuffer),
  }
  
  return queue
}
 
// Method to add a task into the queue 
func (queue TaskQueue) Enqueue(task v1.Task) error {
}
 
// Method to check for the next element to be processed in the queue
func (queue TaskQueue) PeekQueue() task, error {
}

// Method to remove a processed task from the queue
func (queue TaskQueue) DeQueue() v1.Task {

}
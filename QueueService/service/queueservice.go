package service

type QueueService struct {
	Router   *mux.Router
	StopChan chan bool
	Handler  RequestHandler
}

func (q *QueueService) Initialize() {
	q.Queue = NewTaskQueue()
	http.HandleFunc("/queue/start", q.InitializeQueueService)
	http.HandleFunc("/queue/stop", q.StopQueue)

	http.HandleFunc("/queue/produce", q.Handler.AddTask)
	http.HandleFunc("/queue/consume", q.Handler.ConsumeTask)

	http.HandleFunc("/queue/status", q.Handler.CheckQueueStatus)
	http.HandleFunc("/queue/ack", q.Handler.CheckAckStatus)
}

func (q *QueueService) InitalizeQueueService(_ http.ResponseWriter, _ *http.Request) {
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

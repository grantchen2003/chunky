package workerpool

type WorkerPool struct {
	jobQueue      chan Job
	workerCount   int
	workers       []*worker
	doneChannel   chan bool
	resultChannel chan jobResult
}

func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		jobQueue:      make(chan Job),
		workerCount:   workerCount,
		workers:       make([]*worker, workerCount),
		doneChannel:   make(chan bool),
		resultChannel: make(chan jobResult),
	}

	for i := range workerCount {
		worker := &worker{
			jobQueue:    wp.jobQueue,
			doneChannel: wp.doneChannel,
			resultChan:  wp.resultChannel,
		}
		wp.workers[i] = worker
		go worker.Start()
	}

	return wp
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobQueue <- job
}

func (wp *WorkerPool) Wait() {
	close(wp.jobQueue)
	for range wp.workerCount {
		<-wp.doneChannel
	}
	close(wp.resultChannel)
	close(wp.doneChannel)
}

func (wp *WorkerPool) ResultsChannel() <-chan jobResult {
	return wp.resultChannel
}

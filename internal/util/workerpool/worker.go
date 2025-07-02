package workerpool

type worker struct {
	jobQueue    chan Job
	doneChannel chan bool
	resultChan  chan jobResult
}

func (w *worker) Start() {
	for job := range w.jobQueue {
		value, err := job.Execute()
		w.resultChan <- jobResult{Value: value, Err: err}
	}

	w.doneChannel <- true
}

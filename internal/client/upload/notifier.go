package upload

type Notifier struct {
	ProgressChan chan Progress
	ResultChan   chan Result
	StatusChan   chan Status
}

func NewNotifier() *Notifier {
	return &Notifier{
		ProgressChan: make(chan Progress),
		ResultChan:   make(chan Result),
		StatusChan:   make(chan Status),
	}
}

func (n *Notifier) Close() {
	close(n.ProgressChan)
	close(n.ResultChan)
	close(n.StatusChan)
}

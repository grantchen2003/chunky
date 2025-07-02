package workerpool

type Job struct {
	task func() (any, error)
}

func NewJob(task func() (any, error)) *Job {
	return &Job{task: task}
}

func (j *Job) Execute() (any, error) {
	return j.task()
}

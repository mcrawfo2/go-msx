package types

import (
	"context"
	"github.com/pkg/errors"
)

var ErrNoJobResult = errors.New("No result from job")

type JobChan chan *Job

type Job struct {
	Result chan error
	Action ActionFunc
}

func (j *Job) Execute(ctx context.Context) {
	defer close(j.Result)
	j.Result <- j.Action(ctx)
}

func NewJob(action ActionFunc) *Job {
	return &Job{
		Result: make(chan error, 1),
		Action: action,
	}
}

type Worker struct {
	jobs JobChan
	ctx  context.Context
}

func (w *Worker) job(job *Job) {
	job.Execute(w.ctx)
}

func (w *Worker) pump() {
	for job := range w.jobs {
		w.job(job)
	}
}

func (w *Worker) Run(action ActionFunc) error {
	job := NewJob(action)
	w.jobs <- job
	for err := range job.Result {
		return err
	}
	return ErrNoJobResult
}

func (w *Worker) Stop() {
	close(w.jobs)
}

func NewWorker(ctx context.Context) *Worker {
	result := &Worker{
		ctx:  ctx,
		jobs: make(chan *Job, 1),
	}
	go result.pump()
	return result
}

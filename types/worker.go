// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"context"
	"github.com/pkg/errors"
)

var ErrNoJobResult = errors.New("No result from action")

type JobChan chan *Job

// Job implements a unit of work
type Job struct {
	Result    chan error
	Context   context.Context
	Operation Operation
}

func (j *Job) Execute(ctx context.Context) {
	defer close(j.Result)

	if j.Context != nil {
		ctx = j.Context
	}

	j.Result <- j.Operation.Run(ctx)
}

func NewJob(action ActionFunc, options ...JobOption) *Job {
	result := &Job{
		Result:    make(chan error, 1),
		Operation: NewOperation(action),
	}
	for _, option := range options {
		option(result)
	}
	return result
}

type JobOption func(j *Job)

func JobDecorator(deco ActionFuncDecorator) JobOption {
	return func(j *Job) {
		j.Operation = j.Operation.WithDecorator(deco)
	}
}

func JobFilter(filter ActionFilter) JobOption {
	return func(j *Job) {
		j.Operation = j.Operation.WithFilter(filter)
	}
}

func JobContext(ctx context.Context) JobOption {
	return func(j *Job) {
		j.Context = ctx
	}
}

type WorkQueue interface {
	Schedule(action ActionFunc, options ...JobOption) chan error
	Run(action ActionFunc, options ...JobOption) error
}

// Worker implements a serial work queue
type Worker struct {
	jobs JobChan
	ctx  context.Context
}

func (w *Worker) pump() {
	for {
		select {
		case job := <-w.jobs:
			if job == nil {
				// channel was closed
				return
			}

			job.Execute(w.ctx)

		case <-w.ctx.Done():
			// context was cancelled
			return
		}
	}
}

// Schedule asynchronously executes a single action
func (w *Worker) Schedule(action ActionFunc, options ...JobOption) chan error {
	job := NewJob(action, options...)
	w.jobs <- job
	return job.Result
}

// Run synchronously executes a single action
func (w *Worker) Run(action ActionFunc, options ...JobOption) error {
	jobResult := w.Schedule(action, options...)
	for err := range jobResult {
		return err
	}
	return ErrNoJobResult
}

// Stop terminates the pump
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

// WorkerPool implements a parallel work queue
type WorkerPool struct {
	jobs    chan *Job
	workers []*Worker
}

// Schedule asynchronously executes a single action
func (p *WorkerPool) Schedule(action ActionFunc, options ...JobOption) chan error {
	job := NewJob(action, options...)
	p.jobs <- job
	return job.Result
}

// Run synchronously executes a single action
func (p *WorkerPool) Run(action ActionFunc, options ...JobOption) error {
	jobResult := p.Schedule(action, options...)
	for err := range jobResult {
		return err
	}
	return ErrNoJobResult
}

func (p *WorkerPool) Stop() {
	close(p.jobs)
}

func NewWorkerPool(ctx context.Context, workers int) (*WorkerPool, error) {
	if workers <= 0 {
		return nil, errors.New("Minimum pool size is 1")
	}

	pool := &WorkerPool{
		jobs:    make(chan *Job),
		workers: make([]*Worker, workers),
	}

	for i := 0; i < workers; i++ {
		worker := &Worker{
			jobs: pool.jobs,
			ctx:  ctx,
		}

		pool.workers[i] = worker
		go worker.pump()
	}

	return pool, nil
}

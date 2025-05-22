package publish

import (
	"context"
	"sync"
)

const queueSize = 100

type JobFn func(ctx context.Context, workerID int)

// WorkerPool manages a pool of worker goroutines
type WorkerPool struct {
	numWorkers int
	jobQueue   chan JobFn
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan JobFn, queueSize),
		wg:         sync.WaitGroup{},
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := range wp.numWorkers {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

// Stop stops the worker pool and waits for all workers to finish
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}

// SubmitJob submits a job to the worker pool
func (wp *WorkerPool) SubmitJob(job JobFn) {
	wp.jobQueue <- job
}

// worker is a worker goroutine that processes jobs
func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
	defer wp.wg.Done()
	for job := range wp.jobQueue {
		job(ctx, workerID)
	}
}

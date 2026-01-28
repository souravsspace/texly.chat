package queue

import (
	"context"
	"fmt"
	"sync"
)

/*
* Job represents a scraping job to be processed
 */
type Job struct {
	SourceID string
	BotID    string
	URL      string
}

/*
* JobQueue interface defines methods for managing jobs
 */
type JobQueue interface {
	Enqueue(job Job) error
	Start(ctx context.Context, handler JobHandler)
	Stop()
}

/*
* JobHandler is a function that processes a job
 */
type JobHandler func(job Job) error

/*
* InMemoryQueue implements JobQueue using Go channels
 */
type InMemoryQueue struct {
	jobs       chan Job
	workerPool int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

/*
* NewInMemoryQueue creates a new in-memory job queue
* bufferSize: how many jobs can be queued
* workerPool: number of concurrent workers
 */
func NewInMemoryQueue(bufferSize, workerPool int) *InMemoryQueue {
	ctx, cancel := context.WithCancel(context.Background())
	return &InMemoryQueue{
		jobs:       make(chan Job, bufferSize),
		workerPool: workerPool,
		ctx:        ctx,
		cancel:     cancel,
	}
}

/*
* Enqueue adds a job to the queue
 */
func (q *InMemoryQueue) Enqueue(job Job) error {
	select {
	case q.jobs <- job:
		return nil
	case <-q.ctx.Done():
		return fmt.Errorf("queue is stopped")
	default:
		return fmt.Errorf("queue is full")
	}
}

/*
* Start begins processing jobs with the given handler
 */
func (q *InMemoryQueue) Start(ctx context.Context, handler JobHandler) {
	for i := 0; i < q.workerPool; i++ {
		q.wg.Add(1)
		go q.worker(ctx, handler, i)
	}
}

/*
* worker is the background goroutine that processes jobs
 */
func (q *InMemoryQueue) worker(ctx context.Context, handler JobHandler, workerID int) {
	defer q.wg.Done()
	
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopping\n", workerID)
			return
		case <-q.ctx.Done():
			fmt.Printf("Worker %d stopping (queue shutdown)\n", workerID)
			return
		case job, ok := <-q.jobs:
			if !ok {
				fmt.Printf("Worker %d: job channel closed\n", workerID)
				return
			}
			
			// Process the job
			if err := handler(job); err != nil {
				fmt.Printf("Worker %d: error processing job %s: %v\n", workerID, job.SourceID, err)
			} else {
				fmt.Printf("Worker %d: successfully processed job %s\n", workerID, job.SourceID)
			}
		}
	}
}

/*
* Stop gracefully stops the queue and waits for workers to finish
 */
func (q *InMemoryQueue) Stop() {
	q.cancel()
	close(q.jobs)
	q.wg.Wait()
	fmt.Println("All workers stopped")
}

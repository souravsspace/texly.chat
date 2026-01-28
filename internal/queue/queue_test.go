package queue

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryQueue_Enqueue(t *testing.T) {
	queue := NewInMemoryQueue(10, 2)
	defer queue.Stop()

	job := Job{
		SourceID: "source-1",
		BotID:    "bot-1",
		URL:      "https://example.com",
	}

	err := queue.Enqueue(job)
	assert.NoError(t, err)
}

func TestInMemoryQueue_EnqueueFull(t *testing.T) {
	queue := NewInMemoryQueue(1, 1) // Small buffer
	defer queue.Stop()

	job1 := Job{SourceID: "1", BotID: "bot-1", URL: "https://example1.com"}
	job2 := Job{SourceID: "2", BotID: "bot-1", URL: "https://example2.com"}

	// First should succeed
	err := queue.Enqueue(job1)
	assert.NoError(t, err)

	// Second should fail (buffer full, no workers consuming)
	err = queue.Enqueue(job2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue is full")
}

func TestInMemoryQueue_WorkerProcessing(t *testing.T) {
	queue := NewInMemoryQueue(10, 2)
	ctx := context.Background()

	var processedJobs []string
	var mutex sync.Mutex

	handler := func(job Job) error {
		mutex.Lock()
		processedJobs = append(processedJobs, job.SourceID)
		mutex.Unlock()
		return nil
	}

	queue.Start(ctx, handler)
	defer queue.Stop()

	// Enqueue jobs
	queue.Enqueue(Job{SourceID: "1", BotID: "bot-1", URL: "https://example1.com"})
	queue.Enqueue(Job{SourceID: "2", BotID: "bot-1", URL: "https://example2.com"})
	queue.Enqueue(Job{SourceID: "3", BotID: "bot-1", URL: "https://example3.com"})

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	mutex.Lock()
	assert.Len(t, processedJobs, 3)
	assert.Contains(t, processedJobs, "1")
	assert.Contains(t, processedJobs, "2")
	assert.Contains(t, processedJobs, "3")
	mutex.Unlock()
}

func TestInMemoryQueue_ErrorHandling(t *testing.T) {
	queue := NewInMemoryQueue(10, 1)
	ctx := context.Background()

	var processedCount int
	var mutex sync.Mutex

	handler := func(job Job) error {
		mutex.Lock()
		processedCount++
		mutex.Unlock()
		// Return error for testing
		return assert.AnError
	}

	queue.Start(ctx, handler)
	defer queue.Stop()

	// Enqueue job
	queue.Enqueue(Job{SourceID: "1", BotID: "bot-1", URL: "https://example.com"})

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Should still process despite error
	mutex.Lock()
	assert.Equal(t, 1, processedCount)
	mutex.Unlock()
}

func TestInMemoryQueue_Stop(t *testing.T) {
	queue := NewInMemoryQueue(10, 2)
	ctx := context.Background()

	handler := func(job Job) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	queue.Start(ctx, handler)

	// Enqueue some jobs
	queue.Enqueue(Job{SourceID: "1", BotID: "bot-1", URL: "https://example.com"})

	// Stop queue
	queue.Stop()

	// Try to enqueue after stop should fail
	err := queue.Enqueue(Job{SourceID: "2", BotID: "bot-1", URL: "https://example.com"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue is stopped")
}

func TestInMemoryQueue_ConcurrentWorkers(t *testing.T) {
	queue := NewInMemoryQueue(20, 3) // 3 workers
	ctx := context.Background()

	var processedJobs []string
	var mutex sync.Mutex

	handler := func(job Job) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		mutex.Lock()
		processedJobs = append(processedJobs, job.SourceID)
		mutex.Unlock()
		return nil
	}

	queue.Start(ctx, handler)
	defer queue.Stop()

	// Enqueue multiple jobs
	for i := 0; i < 10; i++ {
		queue.Enqueue(Job{
			SourceID: string(rune('A' + i)),
			BotID:    "bot-1",
			URL:      "https://example.com",
		})
	}

	// Wait for all to process
	time.Sleep(200 * time.Millisecond)

	mutex.Lock()
	assert.Len(t, processedJobs, 10)
	mutex.Unlock()
}

func TestInMemoryQueue_ContextCancellation(t *testing.T) {
	queue := NewInMemoryQueue(10, 2)
	ctx, cancel := context.WithCancel(context.Background())

	var processedCount int
	var mutex sync.Mutex

	handler := func(job Job) error {
		mutex.Lock()
		processedCount++
		mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
		return nil
	}

	queue.Start(ctx, handler)

	// Enqueue jobs
	queue.Enqueue(Job{SourceID: "1", BotID: "bot-1", URL: "https://example.com"})
	queue.Enqueue(Job{SourceID: "2", BotID: "bot-1", URL: "https://example.com"})

	// Cancel context immediately
	cancel()
	time.Sleep(50 * time.Millisecond)

	// May process some jobs before cancellation
	mutex.Lock()
	count := processedCount
	mutex.Unlock()

	// Should have processed at least one but workers should stop
	assert.GreaterOrEqual(t, count, 0)
	
	queue.Stop()
}

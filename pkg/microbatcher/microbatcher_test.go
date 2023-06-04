package main_test

import (
	"sync"
	"testing"
	"time"

	main "microbatcher"
)

type MockBatchProcessor struct{}

func (m MockBatchProcessor) ProcessBatch(jobs []main.Job) []main.JobResult {
	var results []main.JobResult
	for _, job := range jobs {
		results = append(results, main.JobResult{
			Job:  job,
			Data: job.Data.(int) * 2,
		})
	}
	return results
}

func TestSubmitJob(t *testing.T) {
	batcher := main.NewBatcher(MockBatchProcessor{}, 10, 1*time.Second)
	defer batcher.Shutdown()

	job := main.Job{
		Data: 2,
	}
	result := batcher.SubmitJob(job)

	if result.Data != 4 {
		t.Errorf("Got %v, expected %v", result.Data, 4)
	}
}

func TestShutdown(t *testing.T) {
	batcher := main.NewBatcher(MockBatchProcessor{}, 10, 1*time.Second)

	var wg sync.WaitGroup
	jobCount := 10

	results := make([]main.JobResult, jobCount)

	// Submit jobs
	for i := 0; i < jobCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			job := main.Job{
				Data: i,
			}
			results[i] = batcher.SubmitJob(job)
		}(i)
	}

	// Wait until all jobs are submitted before calling Shutdown
	wg.Wait()
	batcher.Shutdown()

	for i, result := range results {
		if result.Data != i*2 {
			t.Errorf("Got %v, expected %v", result.Data, i*2)
		}
	}
}

package microbatcher

import (
	"sync"
	"time"
)

// Job: represents a job that needs to be processed.
type Job struct {
	Data interface{}
}

// JobResult: result of a processed job.
type JobResult struct {
	Job  Job
	Data interface{}
}

// BatchProcessor: batch processing system should implement.
type BatchProcessor interface {
	ProcessBatch(jobs []Job) []JobResult
}

// Batcher: struct that represents a system for processing jobs in micro-batches.
type Batcher struct {
	jobs         chan Job
	jobResults   chan JobResult
	shutdown     chan bool
	batchSize    int
	batchTimeout time.Duration
	wg           sync.WaitGroup
	processor    BatchProcessor
}

// NewBatcher is a constructor for Batcher. It starts the batchJobs goroutine.
func NewBatcher(processor BatchProcessor, batchSize int, batchTimeout time.Duration) *Batcher {
	batcher := &Batcher{
		jobs:         make(chan Job),
		jobResults:   make(chan JobResult),
		shutdown:     make(chan bool),
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
		processor:    processor,
	}
	batcher.wg.Add(1)
	go batcher.batchJobs()
	return batcher
}

// SubmitJob allows a caller to submit a job to be processed. It returns the result of the job.
func (b *Batcher) SubmitJob(job Job) JobResult {
	b.jobs <- job
	return <-b.jobResults
}

// Shutdown stops the Batcher from accepting more jobs and processes all jobs that have been accepted but not yet processed.
func (b *Batcher) Shutdown() {
	close(b.jobs)
	b.wg.Wait()
	close(b.jobResults)
}

// batchJobs is a goroutine that runs in the background, processing jobs in micro-batches.
func (b *Batcher) batchJobs() {
	defer b.wg.Done()
	ticker := time.NewTicker(b.batchTimeout)
	defer ticker.Stop()
	var batch []Job
	for {
		select {
		case job, ok := <-b.jobs:
			if !ok {
				if len(batch) > 0 {
					b.processBatch(batch)
				}
				return
			}
			batch = append(batch, job)
			if len(batch) >= b.batchSize {
				b.processBatch(batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				b.processBatch(batch)
				batch = nil
			}
		}
	}
}

// processBatch takes a batch of jobs, processes them using the BatchProcessor, and sends the results to the jobResults channel.
func (b *Batcher) processBatch(jobs []Job) {
	results := b.processor.ProcessBatch(jobs)
	for _, result := range results {
		b.jobResults <- result
	}
}

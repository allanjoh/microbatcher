# Microbatcher

Microbatcher is a Go library for processing tasks in small batches, which can help to improve throughput by reducing the number of requests made to a downstream system.

## Features

- Allows to submit a single job and returns a job result.
- Processes accepted jobs in batches using a custom batch processor.
- Provides a way to configure the batch size and processing frequency.
- Exposes a shutdown method which returns after all previously accepted jobs are processed.

## Usage

To use the library, you first need to implement the `BatchProcessor` interface in your code. This interface has a single method, `ProcessBatch`, which should include the logic for processing a batch of jobs.

Here's an example of how to use the library with a `BatchProcessor` implementation:

```go
package main

import (
	"time"

	"github.com/yourusername/yourproject/pkg/microbatcher"
)

type MyBatchProcessor struct{}

func (r MyBatchProcessor) ProcessBatch(jobs []microbatcher.Job) []microbatcher.JobResult {
	var results []microbatcher.JobResult
	for _, job := range jobs {
		results = append(results, microbatcher.JobResult{
			Job:  job,
			Data: job.Data,
		})
	}
	return results
}

func main() {
	batcher := microbatcher.NewBatcher(MyBatchProcessor{}, 10, 1*time.Second)

	job := microbatcher.Job{
		Data: "my job data",
	}
	result := batcher.SubmitJob(job)

	// use result...

	batcher.Shutdown()
}


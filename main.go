package main

import (
	"fmt"
	"time"

	"github.com/allanjoh/microbatcher/pkg/microbatcher"
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
	jobs := 10
	batcher := microbatcher.NewBatcher(MyBatchProcessor{}, jobs, 1*time.Second)
	for i := 0; i <= jobs; i++ {
		job := microbatcher.Job{
			Data: "my job data",
		}
		result := batcher.SubmitJob(job)

		fmt.Println(result)
	}

	batcher.Shutdown()
}

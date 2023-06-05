package microbatcher_test

import (
	"time"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"

	"github.com/allanjoh/microbatcher/pkg/microbatcher"
	"github.com/allanjoh/microbatcher/pkg/microbatcher/mocks"
)

type BatcherSuite struct {
	suite.Suite
	batcher         *microbatcher.Batcher
	mockProcessor   *mocks.BatchProcessor
}

func (suite *BatcherSuite) SetupTest() {
	suite.mockProcessor = new(mocks.BatchProcessor)
	suite.batcher = microbatcher.NewBatcher(suite.mockProcessor, 10, 1*time.Second)
}

func (suite *BatcherSuite) TestSubmitJob() {
	job := microbatcher.Job{Data: "job1"}
	result := microbatcher.JobResult{Job: job, Data: "result1"}

	suite.mockProcessor.On("ProcessBatch", mock.AnythingOfType("[]microbatcher.Job")).Return([]microbatcher.JobResult{result}).Once()

	go func() {
		time.Sleep(1 * time.Second)
		suite.batcher.SubmitJob(job)
	}()

	res := suite.batcher.SubmitJob(job)

	suite.Equal(result, res)
	suite.mockProcessor.AssertExpectations(suite.T())
}

// TestShutdown verifies that the Shutdown method stops accepting more jobs and processes all pending jobs.
func (suite *BatcherSuite) TestShutdown() {
	job := microbatcher.Job{Data: "job"}
	result := microbatcher.JobResult{Job: job, Data: "result"}
	suite.mockProcessor.On("ProcessBatch", mock.AnythingOfType("[]microbatcher.Job")).Return([]microbatcher.JobResult{result}).Once()
	suite.batcher.SubmitJob(job)

	suite.batcher.Shutdown()
	suite.mockProcessor.AssertExpectations(suite.T())

	assert.Panics(suite.T(), func() { suite.batcher.SubmitJob(job) })
}

func TestBatcherSuite(t *testing.T) {
	suite.Run(t, new(BatcherSuite))
}

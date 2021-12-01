package ops

import (
	"context"
	"expinc/sunagent/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunJob_Successful(t *testing.T) {
	// start job
	info, err := StartJob(context.Background(), "dummy", nil)
	assert.Equal(t, nil, err)

	// check in progress
	time.Sleep(3 * time.Second)
	info, err = GetJobInfo(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, JobStatusExecuting, info.Status)
	assert.Less(t, info.Progress, 100)
	assert.Greater(t, info.Progress, 0)

	// check finshed
	time.Sleep(10 * time.Second)
	info, err = GetJobInfo(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, JobStatusSuccessful, info.Status)
	assert.Equal(t, 100, info.Progress)
	assert.Equal(t, "Successful", info.Result)
}

func TestRunJob_Failed(t *testing.T) {
	// start job
	params := map[string]interface{}{
		"expectedResult": JobStatusFailed,
	}
	info, err := StartJob(context.Background(), "dummy", params)
	assert.Equal(t, nil, err)

	// check result
	time.Sleep(3 * time.Second)
	info, err = GetJobInfo(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, JobStatusFailed, info.Status)
	resultErr, ok := info.Result.(common.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, common.ErrorUnknown, resultErr.Code())
}

func TestCancelJob_Ordinary(t *testing.T) {
	// start job
	info, err := StartJob(context.Background(), "dummy", nil)
	assert.Equal(t, nil, err)

	// check in progress
	time.Sleep(3 * time.Second)
	info, err = GetJobInfo(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, JobStatusExecuting, info.Status)
	assert.Less(t, info.Progress, 100)
	assert.Greater(t, info.Progress, 0)

	// cancel job and check
	info, err = CancelJob(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	time.Sleep(2 * time.Second)
	info, err = GetJobInfo(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, JobStatusCanceled, info.Status)
	assert.Less(t, info.Progress, 100)
}

func TestCancelJob_Finished(t *testing.T) {
	// start job
	params := map[string]interface{}{
		"expectedResult": JobStatusFailed,
		"execSeconds":    1,
	}
	info, err := StartJob(context.Background(), "dummy", params)
	assert.Equal(t, nil, err)

	// cancel job after finish
	time.Sleep(2 * time.Second)
	info, err = CancelJob(context.Background(), info.Id)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, common.ErrorNotAllowed, err.(common.Error).Code())
}

func TestStartJob_InvalidType(t *testing.T) {
	_, err := StartJob(context.Background(), "", nil)
	internalErr, ok := err.(common.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, common.ErrorInvalidParameter, internalErr.Code())
}

func TestRunJob_Panic(t *testing.T) {
	// start job
	params := map[string]interface{}{
		"expectedResult": "panic",
	}
	info, err := StartJob(context.Background(), "dummy", params)
	assert.Equal(t, nil, err)

	// check result
	time.Sleep(3 * time.Second)
	info, err = GetJobInfo(context.Background(), info.Id)
	assert.Equal(t, nil, err)
	assert.Equal(t, JobStatusFailed, info.Status)
	resultErr, ok := info.Result.(common.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, common.ErrorUnknown, resultErr.Code())
}

func TestAccessJob_NonExist(t *testing.T) {
	_, err := GetJobInfo(context.Background(), "")
	internalErr, ok := err.(common.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, common.ErrorNotFound, internalErr.Code())

	_, err = CancelJob(context.Background(), "")
	internalErr, ok = err.(common.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, common.ErrorNotFound, internalErr.Code())
}

func TestListJobs(t *testing.T) {
	for i := 0; i < 3; i++ {
		StartJob(context.Background(), "dummy", nil)
	}

	jobs := ListJobInfo(context.Background())
	assert.Equal(t, len(id2Jobs), len(jobs))
	for _, job := range jobs {
		assert.NotEmpty(t, job.Id)
	}
}

func TestCleanFinishedJobs(t *testing.T) {
	// Clear jobs and set an invalid clean threshold, which makes threshold as 5
	id2Jobs = make(map[string]Job)
	SetJobCleanThreshold(context.Background(), 1)

	// Add 5 jobs to trigger a job clean.
	// There should be 2 jobs left.
	params := map[string]interface{}{
		"expectedResult": JobStatusFailed,
	}
	for i := 0; i < 5; i++ {
		StartJob(context.Background(), "dummy", params)
		time.Sleep(time.Second)
	}
	assert.Equal(t, 2, len(id2Jobs))

	// Add 6 jobs that will execute for a longer period.
	// The former finished jobs will be cleaned.
	// The 6 executing jobs will be left.
	params = map[string]interface{}{
		"execSeconds": 10,
	}
	for i := 0; i < 6; i++ {
		StartJob(context.Background(), "dummy", params)
	}
	time.Sleep(time.Second)
	assert.Equal(t, 6, len(id2Jobs))
}

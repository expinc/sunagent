package ops

import (
	"encoding/base64"
	"expinc/sunagent/common"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	JobStatusSpawned    = "SPAWNED"
	JobStatusExecuting  = "EXECUTING"
	JobStatusSuccessful = "SUCCESSFUL"
	JobStatusFailed     = "FAILED"
	JobStatusCanceled   = "CANCELED"
)

type JobInfo struct {
	Name      string      `json:"name"`
	Id        string      `json:"id"`
	Status    string      `json:"status"`
	BeginTime time.Time   `json:"beginTime"`
	EndTime   time.Time   `json:"endTime"`
	Progress  int         `json:"progress"`
	Result    interface{} `json:"result"`
}

type Job interface {
	// functions must be implemented by concrete jobs:
	execute() error
	cancel()
	dispose()

	// functions implemented by jobBase
	getInfo() *JobInfo
}

// all concrete job should combine jobBase
type jobBase struct {
	info   *JobInfo
	params map[string]interface{}
}

func (base *jobBase) getInfo() *JobInfo {
	return base.info
}

// all jobs:
var id2Jobs map[string]Job
var jobMutex = sync.RWMutex{}

// a concrete job for test
type dummyJob struct {
	jobBase

	canceled chan bool
}

func (job *dummyJob) execute() error {
	job.canceled = make(chan bool, 1)

	for i := 0; i < 10; i++ {
		select {
		case <-job.canceled:
			break
		default:
			time.Sleep(1 * time.Second)
			job.info.Progress += 10
		}
	}
	return nil
}

func (job *dummyJob) cancel() {
	job.canceled <- true
}

func (job *dummyJob) dispose() {}

func generateJobId() string {
	uuidValue := uuid.New()
	return base64.URLEncoding.EncodeToString(uuidValue[:])
}

func StartJob(typ string, params map[string]interface{}) (info JobInfo, err error) {
	// create job
	var job Job
	info.Id = generateJobId()
	info.Status = JobStatusSpawned
	switch typ {
	case "dummy":
		info.Name = "dummy"
		job = &dummyJob{
			jobBase: jobBase{
				info:   &info,
				params: params,
			},
			canceled: make(chan bool),
		}
	default:
		errMsg := fmt.Sprintf("Invalid job type: %s", typ)
		err = common.NewError(common.ErrorInvalidParameter, errMsg)
		return
	}

	// execute job
	go func(job Job) {
		job.getInfo().Status = JobStatusExecuting
		job.getInfo().BeginTime = time.Now()

		err := job.execute()
		if nil != err {
			job.getInfo().Status = JobStatusFailed
			job.getInfo().Result = err
		} else if JobStatusCanceled != job.getInfo().Status {
			// if the status is CANCELED, the job should be canceled, no need to change the status
			// otherwise, the job should be finished successfully, the status should be set as SUCCESSFUL
			job.getInfo().Status = JobStatusSuccessful
		}
		job.getInfo().EndTime = time.Now()
	}(job)

	// add job to list
	jobMutex.Lock()
	id2Jobs[info.Id] = job
	// TODO: clean finished jobs
	jobMutex.Unlock()

	return
}

func GetJobInfo(id string) (info JobInfo, err error) {
	jobMutex.RLock()
	defer func() {
		jobMutex.RUnlock()
	}()

	job, ok := id2Jobs[id]
	if ok {
		info = *job.getInfo()
	} else {
		errMsg := fmt.Sprintf("No job with ID=%s", id)
		err = common.NewError(common.ErrorNotFound, errMsg)
	}
	return
}

func ListJobInfo() []JobInfo {
	jobMutex.RLock()
	defer func() {
		jobMutex.RUnlock()
	}()

	var result []JobInfo
	for _, job := range id2Jobs {
		result = append(result, *job.getInfo())
	}
	return result
}

func CancelJob(id string) (info JobInfo, err error) {
	jobMutex.RLock()
	defer func() {
		jobMutex.RUnlock()
	}()

	job, ok := id2Jobs[id]
	if ok {
		job.cancel()
		info = *job.getInfo()
	} else {
		errMsg := fmt.Sprintf("No job with ID=%s", id)
		err = common.NewError(common.ErrorNotFound, errMsg)
	}
	return
}

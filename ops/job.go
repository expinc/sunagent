package ops

import (
	"context"
	"encoding/base64"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"fmt"
	"sort"
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

func IsFinshedJobStatus(status string) bool {
	switch status {
	case JobStatusSuccessful:
		return true
	case JobStatusFailed:
		return true
	case JobStatusCanceled:
		return true
	default:
		return false
	}
}

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
	ctx    context.Context
	info   *JobInfo
	params map[string]interface{}
}

func (base *jobBase) getInfo() *JobInfo {
	return base.info
}

// all jobs:
var id2Jobs map[string]Job = make(map[string]Job)
var jobMutex = sync.RWMutex{}
var jobCleanThreshold = 100

// a concrete job for test
type dummyJob struct {
	jobBase

	canceled chan bool
}

func (job *dummyJob) execute() error {
	job.canceled = make(chan bool, 1)
	execSeconds := 10

	if nil != job.params {
		expectedResult, ok := job.params["expectedResult"]
		if ok {
			if JobStatusFailed == expectedResult {
				return common.NewError(common.ErrorUnknown, "Fake error")
			} else if "panic" == expectedResult {
				panic(common.NewError(common.ErrorUnknown, "Fake panic"))
			}
		}

		execSecondsObj, ok := job.params["execSeconds"]
		if ok {
			execSeconds, _ = execSecondsObj.(int)
		}
	}

	for i := 0; i < execSeconds; i++ {
		select {
		case <-job.canceled:
			job.info.Status = JobStatusCanceled
			return nil
		default:
			time.Sleep(1 * time.Second)
			job.info.Progress += 10
		}
	}
	job.info.Result = "Successful"
	return nil
}

func (job *dummyJob) cancel() {
	job.canceled <- true
}

func (job *dummyJob) dispose() {
	close(job.canceled)
}

func generateJobId() string {
	uuidValue := uuid.New()
	return base64.URLEncoding.EncodeToString(uuidValue[:])
}

func cleanFinishedJobs() {
	var jobInfoList []JobInfo
	for _, job := range id2Jobs {
		jobInfoList = append(jobInfoList, *job.getInfo())
	}
	sort.Slice(jobInfoList, func(i, j int) bool {
		if jobInfoList[i].BeginTime.Before(jobInfoList[j].BeginTime) {
			return true
		}
		return false
	})

	cntRemain := len(jobInfoList)
	for _, info := range jobInfoList {
		if IsFinshedJobStatus(info.Status) {
			delete(id2Jobs, info.Id)
			cntRemain--
		}

		if cntRemain <= jobCleanThreshold/2 {
			break
		}
	}
}

func StartJob(ctx context.Context, typ string, params map[string]interface{}) (info JobInfo, err error) {
	// create job
	var job Job
	info.Id = generateJobId()
	log.InfoCtx(ctx, fmt.Sprintf("Creating job: type=%v, id=%v", typ, info.Id))
	info.Status = JobStatusSpawned
	switch typ {
	case "dummy":
		info.Name = "dummy"
		job = &dummyJob{
			jobBase: jobBase{
				ctx:    ctx,
				info:   &info,
				params: params,
			},
			canceled: make(chan bool),
		}
	default:
		errMsg := fmt.Sprintf("Invalid job type: %s", typ)
		log.ErrorCtx(ctx, errMsg)
		err = common.NewError(common.ErrorInvalidParameter, errMsg)
		return
	}

	// execute job
	go func(job Job) {
		job.getInfo().Status = JobStatusExecuting
		job.getInfo().BeginTime = time.Now()

		defer func(job Job) {
			job.getInfo().EndTime = time.Now()

			// if a panic occurs when executing, the job should be failed
			if p := recover(); nil != p {
				job.getInfo().Status = JobStatusFailed
				job.getInfo().Result = p
			}
		}(job)

		log.InfoCtx(ctx, fmt.Sprintf("Executing job: id=%v", job.getInfo().Id))
		err := job.execute()
		if nil != err {
			job.getInfo().Status = JobStatusFailed
			job.getInfo().Result = err
		} else if JobStatusCanceled != job.getInfo().Status {
			// if the status is CANCELED, the job should be canceled, no need to change the status
			// otherwise, the job should be finished successfully, the status should be set as SUCCESSFUL
			job.getInfo().Status = JobStatusSuccessful
		}
		log.InfoCtx(ctx, fmt.Sprintf("Finished job: id=%v, status=%v", job.getInfo().Id, job.getInfo().Status))
	}(job)

	// add job to list
	jobMutex.Lock()
	id2Jobs[info.Id] = job
	if jobCleanThreshold <= len(id2Jobs) {
		log.InfoCtx(ctx, fmt.Sprintf("Cleaning finished jobs: currentJobs=%v", len(id2Jobs)))
		cleanFinishedJobs()
		log.InfoCtx(ctx, fmt.Sprintf("Cleaned finished jobs: currentJobs=%v", len(id2Jobs)))
	}
	jobMutex.Unlock()

	return
}

func GetJobInfo(ctx context.Context, id string) (info JobInfo, err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Getting job: id=%v", id))
	jobMutex.RLock()
	defer func() {
		jobMutex.RUnlock()
	}()

	job, ok := id2Jobs[id]
	if ok {
		info = *job.getInfo()
	} else {
		errMsg := fmt.Sprintf("No job with ID=%s", id)
		log.ErrorCtx(ctx, errMsg)
		err = common.NewError(common.ErrorNotFound, errMsg)
	}
	return
}

func ListJobInfo(ctx context.Context) []JobInfo {
	log.InfoCtx(ctx, "Listing jobs")
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

func CancelJob(ctx context.Context, id string) (info JobInfo, err error) {
	log.InfoCtx(ctx, fmt.Sprintf("Canceling job: id=%v", id))
	jobMutex.RLock()
	defer func() {
		jobMutex.RUnlock()
	}()

	job, ok := id2Jobs[id]
	if ok {
		if JobStatusExecuting == job.getInfo().Status {
			job.cancel()
			info = *job.getInfo()
		} else {
			errMsg := fmt.Sprintf("Cannot cancel finished job: id=%v", id)
			log.ErrorCtx(ctx, errMsg)
			err = common.NewError(common.ErrorNotAllowed, errMsg)
		}
	} else {
		errMsg := fmt.Sprintf("No job with ID=%s", id)
		log.ErrorCtx(ctx, errMsg)
		err = common.NewError(common.ErrorNotFound, errMsg)
	}
	return
}

func SetJobCleanThreshold(ctx context.Context, num int) {
	log.InfoCtx(ctx, fmt.Sprintf("Setting job clean threshold as %v", num))
	jobMutex.Lock()
	defer func() {
		jobMutex.Unlock()
	}()

	jobCleanThreshold = num
	if jobCleanThreshold < 5 {
		jobCleanThreshold = 5
		log.WarnCtx(ctx, "Cannot set job clean threshold below 5. Fallback it as 5")
	}
}

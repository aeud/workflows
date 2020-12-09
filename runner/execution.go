package runner

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type JobState string

const (
	StateSucceeded JobState = JobState("SUCCEEDED")
	StateFailed    JobState = JobState("FAILED")
	StateLoading   JobState = JobState("LOADING")
	StateNew       JobState = JobState("NEW")
	// Define the check state and timeout duration when getting the Job State
	DefaultCheckStateDuration = 10 * time.Second
	DefaultTimeoutDuration    = 10 * time.Hour
)

type Execution struct {
	State             JobState
	ErrorMessage      string
	Message           string
	Task              *Task
	JobID             string
	asyncErrorHandler chan error
}

func (e *Execution) UpdateState(s JobState) {
	if e.State != s {
		e.State = s
		log.Printf("job %s (%s) changed to state %s", e.JobID, e.Task.Name, s)
	}
}

func NewExecution(t *Task) (*Execution, error) {
	e := &Execution{
		State:             StateNew,
		Task:              t,
		JobID:             "(not attributed yet)",
		asyncErrorHandler: make(chan error),
	}
	if err := e.sendToTaskRunnerEngine(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Execution) sendToTaskRunnerEngine() error {
	e.UpdateState(StateLoading)
	v, err := TaskRunnerNewJob(e.Task)
	if err != nil {
		return err
	}
	e.JobID = v.Job.JobID
	e.UpdateState(JobState(v.Job.State))
	log.Printf("new job inserted in the Task Runner (job_id: %s)", e.JobID)
	go asyncCheckState(e)
	return nil
}

func DurationgBetweenChecks() time.Duration {
	def := DefaultCheckStateDuration
	if s := os.Getenv("WORKFLOW_CHECK_STATE_DURATION_SEC"); s != "" {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Printf("Cannot parse the WORKFLOW_CHECK_STATE_DURATION_SEC %s value. Using default: %s", s, def)
		}
		return time.Duration(v) * time.Second
	}
	return def
}

func DurationgBeforeCheckTimeout() time.Duration {
	return DefaultTimeoutDuration
}

func asyncCheckState(e *Execution) {
	duration := DurationgBetweenChecks()
	t := time.NewTicker(duration)
	tErr := time.NewTicker(DurationgBeforeCheckTimeout())
	for {
		select {
		case <-t.C:
			v, err := TaskRunnerCheckExecution(e)
			if err != nil {
				e.asyncErrorHandler <- err
				return
			}
			e.UpdateState(JobState(v.Job.State))
			switch state := e.State; state {
			case StateSucceeded:
				e.asyncErrorHandler <- nil
				return
			case StateFailed:
				e.asyncErrorHandler <- fmt.Errorf("job %s (%s) failed", e.Task.Name, e.JobID)
				return
			default:
				log.Printf("checking state for job %s (%s) in %s", e.JobID, e.Task.Name, duration)
			}

		case <-tErr.C:
			e.asyncErrorHandler <- fmt.Errorf("timeout state check for job %s (%s)", e.JobID, e.Task.Name)
			return
		}
	}
}

func (e *Execution) Wait() error {
	return <-e.asyncErrorHandler
}

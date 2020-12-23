package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type JobState string

const (
	StateSucceeded   JobState = JobState("SUCCEEDED")
	StatePreparing   JobState = JobState("PREPARING")
	StateQueued      JobState = JobState("QUEUED")
	StateRunning     JobState = JobState("RUNNING")
	StateCancelling  JobState = JobState("CANCELLING")
	StateCancelled   JobState = JobState("CANCELLED")
	StateFailed      JobState = JobState("FAILED")
	StateUnspecified JobState = JobState("STATE_UNSPECIFIED")
	// Define the check state and timeout duration when getting the Job State
	DefaultCheckStateDuration = 10 * time.Second
	DefaultTimeoutDuration    = 10 * time.Hour
)

type Execution struct {
	State   JobState `json:"state"`
	JobID   string   `json:"jobId"`
	Message string   `json:"message"`
	// task              *Task    `json:"-"`
	asyncErrorHandler chan error
	async             bool
}

func (e *Execution) JSON() []byte {
	bs, _ := json.Marshal(e)
	return bs
}

func (e *Execution) UpdateState() error {
	if _, err := TaskRunnerCheckExecution(e); err != nil {
		return err
	}
	return nil
}

func NewExecution(t *Task) (*Execution, error) {
	resp, err := TaskRunnerNewJob(t)
	if err != nil {
		return nil, err
	}
	e := resp.Job
	log.Printf("new job inserted in the Task Runner (job_id: %s)", e.JobID)
	state := e.State
	if state == StateQueued {
		e.async = true
		e.asyncErrorHandler = make(chan error)
		go asyncCheckState(e)
	}
	return e, nil
}

// TODO to change (no env variables...)
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

// TODO to change (no env variables...) on the same model that before
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
			if err := e.UpdateState(); err != nil {
				e.asyncErrorHandler <- err
				return
			}
			switch state := e.State; state {
			case StateSucceeded:
				e.asyncErrorHandler <- nil
				return
			case StateFailed:
				e.asyncErrorHandler <- fmt.Errorf("job %s failed", e.JobID)
				return
			case StateCancelled:
				e.asyncErrorHandler <- fmt.Errorf("job %s failed", e.JobID)
				return
			default:
				log.Printf("checking state for job %s in %s", e.JobID, duration)
			}

		case <-tErr.C:
			e.asyncErrorHandler <- fmt.Errorf("timeout state check for job %s", e.JobID)
			return
		}
	}
}

func (e *Execution) Wait() error {
	if e.async {
		return <-e.asyncErrorHandler
	}
	return nil
}

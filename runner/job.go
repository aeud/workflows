package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func TaskRunnerAPIURL() string {
	return os.Getenv("TASK_RUNNER_URL")
}
func TaskRunnerAPIJob() string {
	return fmt.Sprintf("%s/%s", TaskRunnerAPIURL(), "job")
}
func TaskRunnerAPIJobStatus() string {
	return fmt.Sprintf("%s/%s", TaskRunnerAPIURL(), "jobstatus")
}

type TaskRunnerAPIResponse struct {
	Job *Execution `json:"job"`
}

// DecorateRequestWithAuthentification takes an http.Request, and add the
// calculated (or provided) JWT in the Authorization: Bearer header.
func DecorateRequestWithAuthentification(req *http.Request) error {
	token, err := generateJWT(TaskRunnerAPIURL())
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func SendRequest(method, url string, body io.Reader) (*http.Response, error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if err := DecorateRequestWithAuthentification(req); err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}

func TaskRunnerNewJob(t *Task) (*TaskRunnerAPIResponse, error) {
	jsonPayload, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := SendRequest("POST", TaskRunnerAPIJob(), bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error when sending the task to the TR engine: %s (%s)", t.Name, bs)
	}
	v := TaskRunnerAPIResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func TaskRunnerCheckExecution(e *Execution) (*TaskRunnerAPIResponse, error) {
	url := fmt.Sprintf("%s?job_id=%s", TaskRunnerAPIJobStatus(), e.JobID)
	resp, err := SendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	v := TaskRunnerAPIResponse{
		Job: e,
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

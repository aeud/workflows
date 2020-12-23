package runner

import (
	"fmt"
)

type Task struct {
	Name       string            `yaml:"name" json:"job_name"`
	ImageURI   string            `yaml:"imageUri" json:"imageUri"`
	Engine     string            `yaml:"engine" json:"engine"`
	URL        string            `yaml:"url" json:"url"`
	ScaleTier  string            `yaml:"scaleTier" json:"scaleTier"`
	MasterType string            `yaml:"masterType" json:"masterType"`
	Args       map[string]string `yaml:"args" json:"args"`
	Labels     map[string]string `yaml:"labels" json:"labels"`
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %s", t.Name)
}

func (t *Task) Run() (*Execution, error) {
	return NewExecution(t)
}

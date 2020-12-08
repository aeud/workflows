package runner

import "fmt"

type Task struct {
	Name     string   `yaml:"name" json:"job_name"`
	ImageURI string   `yaml:"imageUri" json:"imageUri"`
	Args     []string `yaml:"args" json:"args"`
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %s", t.Name)
}

func (t *Task) Run() (*Execution, error) {
	return NewExecution(t)
}

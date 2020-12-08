package flows

import (
	"encoding/json"
	"fmt"
	"strings"
	"workflows/runner"

	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"gopkg.in/yaml.v2"
)

type Flow struct {
	Name      string     `yaml:"name" json:"name"`
	Steps     []FlowStep `yaml:"steps" json:"steps"`
	directory map[string]FlowStep
}

func (f *Flow) JSON() []byte {
	bs, _ := json.Marshal(f)
	return bs
}

func (f *Flow) YAML() []byte {
	bs, _ := yaml.Marshal(f)
	return bs
}

func (f *Flow) Graph() dag.Graph {
	f.directory = make(map[string]FlowStep)
	graph := dag.Graph{}
	originRoot, targetRoot := FlowStep{
		ID:              "origin_root",
		ignoreExecution: true,
	}, FlowStep{
		ID:              "target_root",
		ignoreExecution: true,
	}
	graph.Add(originRoot)
	graph.Add(targetRoot)
	for _, step := range f.Steps {
		f.directory[step.ID] = step
		graph.Add(step)
		graph.Connect(dag.BasicEdge(step, originRoot))
		graph.Connect(dag.BasicEdge(targetRoot, step))
		if dependencies := step.DependsOn; dependencies != nil {
			for _, target := range *dependencies {
				graph.Connect(dag.BasicEdge(step, f.directory[target]))
			}
		}

	}
	return graph
}

func (f *Flow) DAG() dag.AcyclicGraph {
	return dag.AcyclicGraph{
		Graph: f.Graph(),
	}
}

func (f *Flow) Walk() error {
	g := f.DAG()
	diagnostics := g.Walk(walkFunction)
	errors := make([]string, 0)
	for _, d := range diagnostics {
		if severity := d.Severity(); severity == tfdiags.Error {
			errors = append(errors, d.Description().Detail)
		}
	}
	if len(errors) > 0 {
		errorMessage := strings.Join(errors, "\n")
		return fmt.Errorf(errorMessage)
	}
	return nil
}

type FlowStep struct {
	ID              string       `yaml:"id" json:"id"`
	Description     string       `yaml:"description" json:"description"`
	DependsOn       *[]string    `yaml:"dependsOn" json:"depends_on"`
	Task            *runner.Task `yaml:"task" json:"task"`
	ignoreExecution bool
	flowName        string
}

func (s *FlowStep) Run() error {
	if s.Task == nil || s.ignoreExecution {
		return nil
	}
	task := s.Task
	if task.Name == "" {
		task.Name = s.ID
	}
	execution, err := task.Run()
	if err != nil {
		return err
	}
	return execution.Wait()
}

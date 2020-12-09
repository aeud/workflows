package flows

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

const (
	FormatYAML = "yaml"
	FormatJSON = "json"
)

func NewFlowFromYAMLFile(filePath string) (Workflow, error) {
	return NewFlowFromFile(filePath, FormatYAML)
}

func NewFlowFromFile(filePath, format string) (Workflow, error) {
	flow := Workflow{}
	file, err := OpenFromPath(filePath)
	if err != nil {
		return flow, err
	}
	defer file.Close()
	var decoder interface {
		Decode(v interface{}) error
	}
	switch format {
	case FormatYAML:
		decoder = yaml.NewDecoder(file)
	case FormatJSON:
		decoder = json.NewDecoder(file)
	default:
		return flow, fmt.Errorf("cannot parse the format %s", format)
	}
	if err := decoder.Decode(&flow); err != nil {
		return flow, err
	}
	if flow.Name != "" {
		for _, step := range flow.Steps {
			step.flowName = flow.Name
		}
	}
	return flow, nil
}

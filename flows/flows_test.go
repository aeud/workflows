package flows

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	flow := Flow{
		Steps: []FlowStep{
			FlowStep{ID: "a", ImageURI: "image_a"},
			FlowStep{ID: "b", ImageURI: "image_b", DependsOn: &[]string{"a"}},
		},
	}
	graph := flow.DAG()
	fmt.Println(graph)
}

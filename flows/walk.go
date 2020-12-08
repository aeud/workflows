package flows

import (
	"fmt"

	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
)

func walkFunction(v dag.Vertex) tfdiags.Diagnostics {
	s := v.(FlowStep)
	if err := s.Run(); err != nil {
		return tfdiags.Diagnostics{
			tfdiags.Sourceless(tfdiags.Error, fmt.Sprintf("task %s failed", s.ID), fmt.Sprintf("error message: %s", err.Error())),
		}
	}
	return tfdiags.Diagnostics{
		tfdiags.SimpleWarning(fmt.Sprintf("task %s succeeded", s.ID)),
	}
}

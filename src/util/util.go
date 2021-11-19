package util

import gen "github.com/sonraisecurity/sonrai-asearch/src/proto"

func FindStepById(stepId string, search *gen.Search) *gen.SearchStep {
	nextSteps := search.Steps
	for {
		if len(nextSteps) == 0 {
			break
		}
		steps := nextSteps
		nextSteps = make([]*gen.SearchStep, 0)
		for _, step := range steps {
			if step.Id == stepId {
				return step
			}

			for i := range step.NextSteps {
				nextSteps = append(nextSteps, step.NextSteps[i])
			}
		}
	}
	return nil
}
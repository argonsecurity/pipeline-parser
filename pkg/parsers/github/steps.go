package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"github.com/mitchellh/mapstructure"
)

func parseSteps(steps *[]githubModels.Step) ([]models.Step, error) {
	if steps == nil {
		return nil, nil
	}

	var stepsSlice []models.Step
	if err := mapstructure.Decode(steps, &stepsSlice); err != nil {
		return nil, err
	}
	return stepsSlice, nil
}

package triggers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func ParseControls(controls *job.Controls, isDeny bool) *models.Condition {
	if controls == nil {
		return nil
	}

	return &models.Condition{
		Allow:     utils.GetPtr(!isDeny),
		Branches:  generateFilter(controls.Refs, isDeny),
		Paths:     generateFilter(controls.Changes, isDeny),
		Variables: parseVariables(controls.Variables),
	}
}

func parseVariables(expressions []string) map[string]string {
	if expressions == nil {
		return nil
	}
	variables := make(map[string]string)
	for _, expression := range expressions {
		comparisons := getComparisons(expression)
		for _, comparison := range comparisons {
			if comparison.IsPositive() {
				variables[comparison.Variable] = comparison.Value
			}
		}
	}
	return variables
}

func generateFilter(list []string, isDeny bool) *models.Filter {
	if list == nil {
		return nil
	}

	if isDeny {
		return &models.Filter{
			DenyList: list,
		}
	}

	return &models.Filter{
		AllowList: list,
	}
}

package triggers

import (
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	refKeywords = []string{
		"api",
		"branches",
		"chat",
		"external",
		"external_pull_requests",
		"merge_requests",
		"pipelines",
		"pushes",
		"schedules",
		"tags",
		"triggers",
		"web",
	}

	refToEventMapping = map[string]models.EventType{
		"pushes":                 models.PushEvent,
		"merge_requests":         models.PullRequestEvent,
		"external_pull_requests": models.PullRequestEvent,
		"pipelines":              models.PipelineRunEvent,
		"trigger":                models.PipelineTriggerEvent,
		"schedules":              models.ScheduledEvent,
	}
)

func ParseControls(controls *job.Controls, isDeny bool) *models.Condition {
	if controls == nil {
		return nil
	}

	branches, events := getBranchesAndEvents(controls.Refs)

	return &models.Condition{
		Allow:     utils.GetPtr(!isDeny),
		Branches:  generateFilter(branches, isDeny),
		Paths:     generateFilter(controls.Changes, isDeny),
		Events:    events,
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

func getBranchesAndEvents(refs []string) ([]string, []models.EventType) {
	branches := []string{}
	events := []models.EventType{}

	for _, ref := range refs {
		if event, ok := refToEventMapping[ref]; ok {
			events = append(events, event)
		} else if utils.SliceContains(refKeywords, ref) {
			events = append(events, models.EventType(ref))
		} else {
			branches = append(branches, ref)
		}
	}

	return branches, events
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

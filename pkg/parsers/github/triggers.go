package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const (
	pushEvent              = "push"
	forkEvent              = "fork"
	workflowDispatchEvent  = "workflow_dispatch"
	pullRequestEvent       = "pull_request"
	pullRequestTargetEvent = "pull_request_target"
)

var (
	githubEventToModelEvent = map[string]models.EventType{
		pushEvent:             models.PushEvent,
		forkEvent:             models.ForkEvent,
		workflowDispatchEvent: models.ManualEvent,
		pullRequestEvent:      models.PullRequestEvent,
	}
)

func parseWorkflowTriggers(workflow *githubModels.Workflow) (*[]models.Trigger, error) {
	if workflow.On == nil {
		return nil, nil
	}

	// Handle workflow.on if it is a list of event names
	if events, isEventListFormat := utils.ToSlice[string](workflow.On); isEventListFormat {
		return utils.GetPtr(generateTriggersFromEvents(events)), nil
	}

	// Handle workflow.on if each event has a specific configuration
	on := workflow.On
	triggers := []models.Trigger{}
	events := utils.Filter(utils.GetMapKeys(on.Events), func(event string) bool {
		_, ok := githubEventToModelEvent[event]
		return !ok
	})
	triggers = append(triggers, generateTriggersFromEvents(events)...)

	if on.Push != nil {
		triggers = append(triggers, parseRef(on.Push, models.PushEvent))
	}

	if on.PullRequest != nil {
		triggers = append(triggers, parseRef(on.PullRequest, models.PullRequestEvent))
	}

	if on.PullRequestTarget != nil {
		triggers = append(triggers, parseRef(on.PullRequestTarget, models.EventType(pullRequestTargetEvent)))
	}

	if on.WorkflowDispatch != nil {
		triggers = append(triggers, parseWorkflowDispatch(on.WorkflowDispatch))
	}

	if on.WorkflowCall != nil {
		triggers = append(triggers, parseWorkflowCall(on.WorkflowCall))
	}

	if on.WorkflowRun != nil {
		triggers = append(triggers, parseWorkflowRun(on.WorkflowRun))
	}

	if on.Schedule != nil {
		triggers = append(triggers, parseSchedule(on.Schedule)...)
	}
	return &triggers, nil
}

func parseWorkflowRun(workflowRun *githubModels.WorkflowRun) models.Trigger {
	trigger := models.Trigger{
		Event:     models.PipelineRunEvent,
		Pipelines: workflowRun.Workflows,
		Branches: &models.Filter{
			AllowList: workflowRun.Branches,
			DenyList:  workflowRun.BranchesIgnore,
		},
	}

	if workflowRun.Types != nil {
		trigger.Filters = map[string]any{"types": workflowRun.Types}
	}

	return trigger
}

func parseWorkflowCall(workflowCall *githubModels.WorkflowCall) models.Trigger {
	return models.Trigger{
		Event:     models.PipelineTriggerEvent,
		Paramters: parseInputs(workflowCall.Inputs),
	}
}

func parseInputs(inputs githubModels.Inputs) []models.Parameter {
	parameters := []models.Parameter{}
	if inputs != nil {
		for k, v := range inputs {
			parameters = append(parameters, models.Parameter{
				Name:        &k,
				Description: &v.Description,
				Default:     v.Default,
			})
		}
	}
	return parameters
}

func parseWorkflowDispatch(workflowDispatch *githubModels.WorkflowDispatch) models.Trigger {
	trigger := models.Trigger{
		Event: models.ManualEvent,
	}

	if workflowDispatch.Inputs != nil {
		trigger.Paramters = parseInputs(workflowDispatch.Inputs)
	}
	return trigger
}

func parseRef(ref *githubModels.Ref, event models.EventType) models.Trigger {
	trigger := models.Trigger{
		Event: event,
	}

	if len(ref.Paths)+len(ref.PathsIgnore) > 0 {
		trigger.Paths = &models.Filter{}
	}

	if len(ref.Branches)+len(ref.BranchesIgnore) > 0 {
		trigger.Branches = &models.Filter{}
	}

	for _, path := range ref.Paths {
		trigger.Paths.AllowList = append(trigger.Paths.AllowList, path)
	}
	for _, path := range ref.PathsIgnore {
		trigger.Paths.DenyList = append(trigger.Paths.DenyList, path)
	}
	for _, branch := range ref.Branches {
		trigger.Branches.AllowList = append(trigger.Branches.AllowList, branch)
	}
	for _, branch := range ref.BranchesIgnore {
		trigger.Branches.DenyList = append(trigger.Branches.DenyList, branch)
	}

	return trigger
}

func parseSchedule(schedule *[]githubModels.Cron) []models.Trigger {
	return utils.Map(*schedule, func(cron githubModels.Cron) models.Trigger {
		return models.Trigger{
			Event:     models.ScheduledEvent,
			Scheduled: utils.GetPtr(cron.Cron),
		}
	})

}

func generateTriggersFromEvents(events []string) []models.Trigger {
	return utils.Map(events, generateTriggerFromEvent)
}

func generateTriggerFromEvent(event string) models.Trigger {
	modelEvent, ok := githubEventToModelEvent[event]
	if !ok {
		modelEvent = models.EventType(event)
	}
	return models.Trigger{
		Event: modelEvent,
	}
}

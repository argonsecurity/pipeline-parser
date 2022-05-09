package github

import (
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
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

func parseWorkflowTriggers(workflow *githubModels.Workflow) *models.Triggers {
	if workflow.On == nil {
		return nil
	}

	// Handle workflow.on if it is a list of event names
	if events, isEventListFormat := utils.ToSlice[string](workflow.On); isEventListFormat {
		return &models.Triggers{
			Triggers:      utils.GetPtr(generateTriggersFromEvents(events)),
			FileReference: workflow.On.FileReference,
		}
	}

	// Handle workflow.on if each event has a specific configuration
	on := workflow.On
	triggerSlice := []models.Trigger{}

	if on.Push != nil {
		triggerSlice = append(triggerSlice, parseRef(on.Push, models.PushEvent))
	}

	if on.PullRequest != nil {
		triggerSlice = append(triggerSlice, parseRef(on.PullRequest, models.PullRequestEvent))
	}

	if on.PullRequestTarget != nil {
		triggerSlice = append(triggerSlice, parseRef(on.PullRequestTarget, models.EventType(pullRequestTargetEvent)))
	}

	if on.WorkflowDispatch != nil {
		triggerSlice = append(triggerSlice, parseWorkflowDispatch(on.WorkflowDispatch))
	}

	if on.WorkflowCall != nil {
		triggerSlice = append(triggerSlice, parseWorkflowCall(on.WorkflowCall))
	}

	if on.WorkflowRun != nil {
		triggerSlice = append(triggerSlice, parseWorkflowRun(on.WorkflowRun))
	}

	if on.Schedule != nil {
		triggerSlice = append(triggerSlice, parseSchedule(on.Schedule))
	}

	if len(on.Events) > 0 {
		triggerSlice = append(triggerSlice, parseEvents(on.Events)...)
	}

	return &models.Triggers{
		Triggers:      &triggerSlice,
		FileReference: workflow.On.FileReference,
	}
}

func parseEvents(events githubModels.Events) []models.Trigger {
	return utils.MapToSlice(events, func(eventName string, event githubModels.Event) models.Trigger {
		trigger := models.Trigger{
			Event: models.EventType(eventName),
			Filters: map[string]any{
				"types": event.Types,
			},
			FileReference: event.FileReference,
		}
		return trigger
	})
}

func parseWorkflowRun(workflowRun *githubModels.WorkflowRun) models.Trigger {
	trigger := models.Trigger{
		FileReference: workflowRun.FileReference,
		Event:         models.PipelineRunEvent,
		Pipelines:     workflowRun.Workflows,
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
		Event:         models.PipelineTriggerEvent,
		Parameters:    parseInputs(workflowCall.Inputs),
		FileReference: workflowCall.FileReference,
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
		Event:         models.ManualEvent,
		FileReference: workflowDispatch.FileReference,
	}

	if workflowDispatch.Inputs != nil {
		trigger.Parameters = parseInputs(workflowDispatch.Inputs)
	}
	return trigger
}

func parseRef(ref *githubModels.Ref, event models.EventType) models.Trigger {
	trigger := models.Trigger{
		Event:         event,
		FileReference: ref.FileReference,
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

func parseSchedule(schedule *githubModels.Schedule) models.Trigger {
	return models.Trigger{
		Event: models.ScheduledEvent,
		Schedules: utils.GetPtr(utils.Map(*schedule.Crons, func(cron githubModels.Cron) string {
			return cron.Cron
		})),
		FileReference: schedule.FileReference,
	}

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

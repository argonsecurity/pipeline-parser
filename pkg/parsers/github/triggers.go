package github

import (
	"errors"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/mitchellh/mapstructure"
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

type gitevent struct {
	Paths          []string `mapstructure:"paths"`
	PathsIgnore    []string `mapstructure:"paths-ignore"`
	Branches       []string `mapstructure:"branches"`
	BranchesIgnore []string `mapstructure:"branches-ignore"`
}

type inputs struct {
	Description string      `mapstructure:"description"`
	Default     interface{} `mapstructure:"default"`
	Required    bool        `mapstructure:"required"`
	Type        string      `mapstructure:"type"`
	Options     []string    `mapstructure:"options,omitempty"`
}

type workflowdispatch struct {
	Inputs map[string]inputs `mapstructure:"inputs"`
}

type On struct {
	Push              *gitevent `mapstructure:"push"`
	PullRequest       *gitevent `mapstructure:"pull_request"`
	PullRequestTarget *gitevent `mapstructure:"pull_request_target"`
	WorkflowCall      *struct {
		Inputs  *inputs `mapstructure:"inputs"`
		Outputs map[string]*struct {
			Description string `mapstructure:"description"`
			Value       string `mapstructure:"value"`
		}
		Secrets map[string]*struct {
			Description string `mapstructure:"description"`
			Required    bool   `mapstructure:"required"`
		}
	} `mapstructure:"workflow_call"`
	Schedule    map[string]string `mapstructure:"schedule"`
	WorkflowRun *struct {
		Types          []string `mapstructure:"types"`
		Workflows      []string `mapstructure:"workflows"`
		Branches       []string `mapstructure:"branches"`
		BranchesIgnore []string `mapstructure:"branches-ignore"`
	} `mapstructure:"workflow_run"`
	WorkflowDispatch *workflowdispatch `mapstructure:"workflow_dispatch"`
	Events
}

type Events map[string]*struct {
	Types []string `mapstructure:"types"`
}

func parseWorkflowTriggers(workflow *Workflow) ([]models.Trigger, error) {
	triggers := []models.Trigger{}
	if workflow.On == nil {
		return nil, nil
	}

	var on On
	if events, ok := utils.ToSlice[string](workflow.On); ok {
		triggers = generateTriggersFromEvents(events)
	} else if err := mapstructure.Decode(workflow.On, &on); err == nil {
		events := utils.Filter(utils.GetMapKeys(on.Events), func(event string) bool {
			_, ok := githubEventToModelEvent[event]
			return !ok
		})
		triggers = append(triggers, generateTriggersFromEvents(events)...)
		if on.Push != nil {
			triggers = append(triggers, parseGitEvent(on.Push, models.PushEvent))
		}

		if on.PullRequest != nil {
			triggers = append(triggers, parseGitEvent(on.PullRequest, models.PullRequestEvent))
		}

		if on.PullRequestTarget != nil {
			triggers = append(triggers, parseGitEvent(on.PullRequestTarget, models.EventType(pullRequestTargetEvent)))
		}

		if on.WorkflowDispatch != nil {
			triggers = append(triggers, parseWorkflowDispatch(on.WorkflowDispatch))
		}
	} else {
		return nil, errors.New("failed to parse workflow triggers")
	}
	return triggers, nil
}

func parseWorkflowDispatch(workflowDispatch *workflowdispatch) models.Trigger {
	if workflowDispatch == nil {
		return models.Trigger{}
	}
	trigger := models.Trigger{
		Event: models.ManualEvent,
	}

	if workflowDispatch.Inputs != nil {
		for k, v := range workflowDispatch.Inputs {
			trigger.Paramters = append(trigger.Paramters, models.Parameter{
				Name:        &k,
				Description: &v.Description,
				Default:     v.Default,
			})
		}
	}
	return trigger
}

func parseGitEvent(gitevent *gitevent, event models.EventType) models.Trigger {
	if gitevent == nil {
		return models.Trigger{}
	}
	trigger := models.Trigger{
		Event: event,
		Paths: &models.Filter{
			AllowList: []string{},
			DenyList:  []string{},
		},
		Branches: &models.Filter{
			AllowList: []string{},
			DenyList:  []string{},
		},
	}

	for _, path := range gitevent.Paths {
		trigger.Paths.AllowList = append(trigger.Paths.AllowList, path)
	}
	for _, path := range gitevent.PathsIgnore {
		trigger.Paths.DenyList = append(trigger.Paths.DenyList, path)
	}
	for _, branch := range gitevent.Branches {
		trigger.Branches.AllowList = append(trigger.Branches.AllowList, branch)
	}
	for _, branch := range gitevent.BranchesIgnore {
		trigger.Branches.DenyList = append(trigger.Branches.DenyList, branch)
	}

	return trigger
}

func generateTriggersFromEvents(events []string) []models.Trigger {
	return utils.Map(events, func(event string) models.Trigger {
		modelEvent, ok := githubEventToModelEvent[event]
		if !ok {
			modelEvent = models.EventType(event)
		}
		return models.Trigger{
			Event: modelEvent,
		}
	})
}

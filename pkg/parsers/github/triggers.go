package github

import (
	"errors"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const (
	pushEvent             = "push"
	forkEvent             = "fork"
	workflowDispatchEvent = "workflow_dispatch"
	pullRequestEvent      = "pull_request"
)

// type Trigger struct {
// 	Branches    *Filter
// 	Paths       *Filter
// 	PullRequest *bool
// 	Manual      *bool
// 	Disabled    *bool
// 	Push        *bool
// 	Fork        *bool
// 	Scheduled   *string
// 	Events      *[]string
//     }

type gitevent struct {
	Paths       []string `yaml:"paths"`
	PathsIgnore []string `yaml:"paths-ignore"`
	branchFilters
}

type branchFilters struct {
	Branches       []string `yaml:"branches"`
	BranchesIgnore []string `yaml:"branches-ignore"`
}

type inputs struct {
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
	Required    bool        `yaml:"required"`
	Type        string      `yaml:"type"`
	Options     []string    `yaml:"options,omitempty"`
}

type On struct {
	Push              *gitevent `yaml:"push"`
	PullRequest       *gitevent `yaml:"pull_request"`
	PullRequestTarget *gitevent `yaml:"pull_request_target"`
	WorkflowCall      *struct {
		Inputs  *inputs `yaml:"inputs"`
		Outputs map[string]*struct {
			Description string `yaml:"description"`
			Value       string `yaml:"value"`
		}
		Secrets map[string]*struct {
			Description string `yaml:"description"`
			Required    bool   `yaml:"required"`
		}
	} `yaml:"workflow_call"`
	Schedule    map[string]string `yaml:"schedule"`
	WorkflowRun *struct {
		Types     []string `yaml:"types"`
		Workflows []string `yaml:"workflows"`
		branchFilters
	} `yaml:"workflow_run"`
	WorkflowDispatch *struct {
		Inputs *inputs `yaml:"inputs"`
	} `yaml:"workflow_dispatch"`
	Events
}

type Events map[string]*struct {
	Types []string `yaml:"types"`
}

func parseWorkflowTriggers(workflow *Workflow) (*models.Trigger, error) {
	trigger := &models.Trigger{}
	if workflow.On == nil {
		return nil, nil
	}

	if events, ok := utils.ToSlice[string](workflow.On); ok {
		trigger.Events = &events
		for _, event := range events {
			switch event {
			case pushEvent:
				trigger.Push = utils.GetPtr(true)
			case forkEvent:
				trigger.Fork = utils.GetPtr(true)
			case workflowDispatchEvent:
				trigger.Manual = utils.GetPtr(true)
			case pullRequestEvent:
				trigger.PullRequest = utils.GetPtr(true)
			}
		}
	} else if on, ok := workflow.On.(On); ok {
		events := utils.GetMapKeys(on.Events)
		trigger.Events = &events
	} else {
		return nil, errors.New("failed to parse workflow triggers")
	}
	return trigger, nil
}

package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Inputs map[string]struct {
	Description string
	Default     interface{}
	Required    bool
	Type        string
	Options     []string
}

type Outputs map[string]*struct {
	Description string
	Value       string
}

type WorkflowDispatch struct {
	Inputs       Inputs
	FileLocation *models.FileLocation
}

type WorkflowCall struct {
	Inputs  Inputs
	Outputs Outputs
	Secrets map[string]*struct {
		Description string
		Required    bool
	}
	FileLocation *models.FileLocation
}

type WorkflowRun struct {
	Types     []string
	Workflows []string
	Ref
	FileLocation *models.FileLocation
}

type Events map[string]Event

type Event struct {
	Types []string `mapstructure:"types"`
}

type Cron struct {
	Cron         string
	FileLocation *models.FileLocation
}

type On struct {
	Push              *Ref
	PullRequest       *Ref
	PullRequestTarget *Ref
	WorkflowCall      *WorkflowCall
	Schedule          *[]Cron
	WorkflowRun       *WorkflowRun
	WorkflowDispatch  *WorkflowDispatch
	FileLocation      *models.FileLocation
	Events
}

func (on *On) UnmarshalYAML(node *yaml.Node) error {
	events := make([]string, 0)
	if err := node.Decode(&events); err == nil {
		for _, event := range events {
			on.unmarshalKey(event, yaml.Node{})
		}
		return nil
	}

	var triggersMap map[string]yaml.Node
	if err := node.Decode(&triggersMap); err != nil {
		return err
	}

	for k, v := range triggersMap {
		if err := on.unmarshalKey(k, v); err != nil {
			return err
		}
	}
	on.FileLocation = loadersUtils.GetFileLocation(node)
	return nil
}

func (on *On) unmarshalKey(key string, node yaml.Node) error {
	var err error
	fileLocation := loadersUtils.GetFileLocation(&node)
	switch key {
	case "schedule":
		mapstructure.Decode(node, &on.Schedule)
	case "push":
		on.Push = &Ref{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.Push)
		}
	case "pull_request":
		on.PullRequest = &Ref{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.PullRequest)
		}
	case "pull_request_target":
		on.PullRequestTarget = &Ref{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.PullRequestTarget)
		}
	case "workflow_call":
		on.WorkflowCall = &WorkflowCall{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.WorkflowCall)
		}
	case "workflow_run":
		on.WorkflowRun = &WorkflowRun{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.WorkflowRun)
		}
	case "workflow_dispatch":
		on.WorkflowDispatch = &WorkflowDispatch{FileLocation: fileLocation}
		if !node.IsZero() {
			err = node.Decode(&on.WorkflowDispatch)
		}
	default:
		if on.Events == nil {
			on.Events = make(Events)
		}
		var event Event
		if err := node.Decode(&event); err != nil {
			return err
		}
		on.Events[key] = event
	}
	return err
}

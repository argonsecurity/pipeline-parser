package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
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
	Types        []string
	FileLocation *models.FileLocation
}

type Cron struct {
	Cron         string
	FileLocation *models.FileLocation
}

type Schedule struct {
	Crons        *[]Cron
	FileLocation *models.FileLocation
}

type On struct {
	Push              *Ref
	PullRequest       *Ref
	PullRequestTarget *Ref
	WorkflowCall      *WorkflowCall
	Schedule          *Schedule
	WorkflowRun       *WorkflowRun
	WorkflowDispatch  *WorkflowDispatch
	FileLocation      *models.FileLocation
	Events
}

func (on *On) UnmarshalYAML(node *yaml.Node) error {
	on.FileLocation = loadersUtils.GetFileLocation(node)
	if node.Tag == consts.SequenceTag {
		fileLocation := loadersUtils.GetFileLocation(node)
		for _, event := range node.Content {
			on.unmarshalString(event.Value, &yaml.Node{}, fileLocation)
		}
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		if err := on.unmarshalNode(node.Content[i], node.Content[i+1]); err != nil {
			return err
		}
	}
	on.FileLocation.StartRef.Line-- // The line of the "on" node is currently not accessible, this is a patch
	return nil
}

func (on *On) unmarshalNode(key, val *yaml.Node) error {
	fileLocation := loadersUtils.GetMapKeyFileLocation(key, val)
	return on.unmarshalString(key.Value, val, fileLocation)
}

func (on *On) unmarshalString(key string, val *yaml.Node, fileLocation *models.FileLocation) error {
	var err error
	switch key {
	case "schedule":
		on.Schedule = &Schedule{
			FileLocation: fileLocation,
		}
		err = val.Decode(&on.Schedule.Crons)
	case "push":
		on.Push = &Ref{FileLocation: fileLocation}
		if !val.IsZero() {
			err = val.Decode(&on.Push)
		}
	case "pull_request":
		on.PullRequest = &Ref{FileLocation: fileLocation}
		if !val.IsZero() {
			err = val.Decode(&on.PullRequest)
		}
	case "pull_request_target":
		on.PullRequestTarget = &Ref{FileLocation: fileLocation}
		if !val.IsZero() {
			err = val.Decode(&on.PullRequestTarget)
		}
	case "workflow_call":
		on.WorkflowCall = &WorkflowCall{FileLocation: fileLocation}
		if !val.IsZero() {
			err = val.Decode(&on.WorkflowCall)
		}
	case "workflow_run":
		on.WorkflowRun = &WorkflowRun{FileLocation: fileLocation}
		if !val.IsZero() {
			if err = val.Decode(&on.WorkflowRun); err == nil {
				err = val.Decode(&(on.WorkflowRun.Ref))
			}
		}
	case "workflow_dispatch":
		on.WorkflowDispatch = &WorkflowDispatch{FileLocation: fileLocation}
		if !val.IsZero() {
			err = val.Decode(&on.WorkflowDispatch)
		}
	default:
		if on.Events == nil {
			on.Events = make(Events)
		}
		event := Event{FileLocation: fileLocation}
		if err := val.Decode(&event); err != nil {
			return err
		}
		on.Events[key] = event
	}
	return err
}

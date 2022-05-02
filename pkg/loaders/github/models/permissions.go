package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

const (
	readAll  = "read-all"
	writeAll = "write-all"
)

type PermissionsEvent struct {
	Actions            string `yaml:"actions,omitempty"`
	Checks             string `yaml:"checks,omitempty"`
	Contents           string `yaml:"contents,omitempty"`
	Deployments        string `yaml:"deployments,omitempty"`
	Discussions        string `yaml:"discussions,omitempty"`
	IdToken            string `yaml:"id-token,omitempty"`
	Issues             string `yaml:"issues,omitempty"`
	Packages           string `yaml:"packages,omitempty"`
	Pages              string `yaml:"pages,omitempty"`
	PullRequests       string `yaml:"pull-requests,omitempty"`
	RepositoryProjects string `yaml:"repository-projects,omitempty"`
	SecurityEvents     string `yaml:"security-events,omitempty"`
	Statuses           string `yaml:"statuses,omitempty"`

	FileLocation *models.FileLocation
}

func createFullPermissions(permission string) *PermissionsEvent {
	return &PermissionsEvent{
		Actions:            permission,
		Checks:             permission,
		Contents:           permission,
		Deployments:        permission,
		Discussions:        permission,
		IdToken:            permission,
		Issues:             permission,
		Packages:           permission,
		Pages:              permission,
		PullRequests:       permission,
		RepositoryProjects: permission,
		SecurityEvents:     permission,
		Statuses:           permission,
	}
}

func (p *PermissionsEvent) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == "!!str" {
		switch node.Value {
		case readAll:
			*p = *createFullPermissions("read")
		case writeAll:
			*p = *createFullPermissions("write")
		}
		return nil
	}

	var tmpInterface any
	if err := node.Decode(&tmpInterface); err != nil {
		return err
	}

	if err := mapstructure.Decode(tmpInterface, &p); err != nil {
		return err
	}

	p.FileLocation = loadersUtils.GetFileLocation(node)
	return nil
}

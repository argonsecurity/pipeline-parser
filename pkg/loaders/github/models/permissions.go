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
	Actions            string `mapstructure:"actions,omitempty" yaml:"actions,omitempty"`
	Checks             string `mapstructure:"checks,omitempty" yaml:"checks,omitempty"`
	Contents           string `mapstructure:"contents,omitempty" yaml:"contents,omitempty"`
	Deployments        string `mapstructure:"deployments,omitempty" yaml:"deployments,omitempty"`
	Discussions        string `mapstructure:"discussions,omitempty" yaml:"discussions,omitempty"`
	IdToken            string `mapstructure:"id-token,omitempty" yaml:"id-token,omitempty"`
	Issues             string `mapstructure:"issues,omitempty" yaml:"issues,omitempty"`
	Packages           string `mapstructure:"packages,omitempty" yaml:"packages,omitempty"`
	Pages              string `mapstructure:"pages,omitempty" yaml:"pages,omitempty"`
	PullRequests       string `mapstructure:"pull-requests,omitempty" yaml:"pull-requests,omitempty"`
	RepositoryProjects string `mapstructure:"repository-projects,omitempty" yaml:"repository-projects,omitempty"`
	SecurityEvents     string `mapstructure:"security-events,omitempty" yaml:"security-events,omitempty"`
	Statuses           string `mapstructure:"statuses,omitempty" yaml:"statuses,omitempty"`

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
	p.FileLocation.StartRef.Line--
	return nil
}

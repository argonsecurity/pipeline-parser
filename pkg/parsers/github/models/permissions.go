package models

import (
	"reflect"

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
	var summarizedPermissions string
	if err := node.Decode(&summarizedPermissions); err == nil {
		if summarizedPermissions == readAll {
			*p = *createFullPermissions("read")
		} else if summarizedPermissions == writeAll {
			*p = *createFullPermissions("write")
		}
		return nil
	}

	var tmpInterface any
	if err := node.Decode(&tmpInterface); err != nil {
		return err
	}
	return mapstructure.Decode(tmpInterface, &p)
}

func DecodeTokenPermissionsHookFunc() mapstructure.DecodeHookFuncType {
	return func(f, t reflect.Type, data any) (any, error) {
		if t != reflect.TypeOf(PermissionsEvent{}) {
			return data, nil
		}

		if f.Kind() == reflect.String {
			if data == readAll {
				return createFullPermissions("read"), nil
			} else if data == writeAll {
				return createFullPermissions("write"), nil
			}
		}

		if f.Kind() == reflect.Map {
			permissions := PermissionsEvent{}
			mapstructure.Decode(data, &permissions)
			return permissions, nil
		}

		return data, nil
	}
}

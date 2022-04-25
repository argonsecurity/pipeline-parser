package github

import (
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/mitchellh/mapstructure"
)

const (
	readPermission  = "read"
	writePermission = "write"
)

var (
	customPermissionsMap = map[string]string{
		"actions":       models.RunPipelinePermission,
		"pull-requests": models.PullRequestPermission,
	}
)

func parseTokenPermissions(permissions *githubModels.PermissionsEvent) (*map[string]models.Permission, error) {
	if permissions == nil {
		return nil, nil
	}

	tokenPermissions := make(map[string]models.Permission)
	var permissionsMap map[string]string
	if err := mapstructure.Decode(permissions, &permissionsMap); err != nil {
		return nil, err
	}

	for permissionName, value := range permissionsMap {
		if customPermissionsMap[permissionName] != "" {
			permissionName = customPermissionsMap[permissionName]
		}
		tokenPermissions[permissionName] = parsePermissionValue(value)
	}
	return &tokenPermissions, nil
}

func parsePermissionValue(permission string) models.Permission {
	if permission == readPermission {
		return models.Permission{
			Read: true,
		}
	}
	if permission == writePermission {
		return models.Permission{
			Write: true,
		}
	}
	return models.Permission{}
}

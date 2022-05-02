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

func parseTokenPermissions(permissions *githubModels.PermissionsEvent) (*models.TokenPermissions, error) {
	if permissions == nil {
		return nil, nil
	}

	var permissionsMap map[string]any
	if err := mapstructure.Decode(permissions, &permissionsMap); err != nil {
		return nil, err
	}

	tokenPermissions := make(map[string]models.Permission)
	for permissionName, value := range permissionsMap {
		if val, ok := value.(string); ok {
			if customPermissionsMap[permissionName] != "" {
				permissionName = customPermissionsMap[permissionName]
			}
			tokenPermissions[permissionName] = parsePermissionValue(val)
		}
	}

	return &models.TokenPermissions{
		Permissions:  tokenPermissions,
		FileLocation: permissions.FileLocation,
	}, nil
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

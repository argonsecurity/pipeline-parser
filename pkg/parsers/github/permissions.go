package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
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

func parseTokenPermissions(permissions *githubModels.PermissionsEvent) *map[string]models.Permission {
	if permissions == nil {
		return nil
	}
	tokenPermissions := make(map[string]models.Permission)
	var permissionsMap map[string]string
	mapstructure.Decode(permissions, &permissionsMap)
	for permissionName, value := range permissionsMap {
		if customPermissionsMap[permissionName] != "" {
			permissionName = customPermissionsMap[permissionName]
		}
		tokenPermissions[permissionName] = parsePermissionValue(value)
	}
	return &tokenPermissions
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

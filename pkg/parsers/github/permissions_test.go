package github

import (
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestParseTokenPermissions(t *testing.T) {
	testCases := []struct {
		name                string
		permissions         *githubModels.PermissionsEvent
		expectedPermissions *models.TokenPermissions
	}{
		{
			name:                "Permissions is nil",
			permissions:         nil,
			expectedPermissions: nil,
		},
		{
			name: "All permissions keys",
			permissions: &githubModels.PermissionsEvent{
				Actions:       "read",
				Checks:        "write",
				Contents:      "read",
				Deployments:   "read",
				Issues:        "write",
				Pages:         "read",
				Statuses:      "read",
				Packages:      "nothing",
				FileReference: testutils.CreateFileReference(6, 7, 8, 9),
			},
			expectedPermissions: &models.TokenPermissions{
				Permissions: map[string]models.Permission{
					"checks": {
						Write: true,
					},
					"contents": {
						Read: true,
					},
					"deployments": {
						Read: true,
					},
					"issues": {
						Write: true,
					},
					"pages": {
						Read: true,
					},
					"run-pipeline": {
						Read: true,
					},
					"statuses": {
						Read: true,
					},
					"packages": {},
				},
				FileReference: testutils.CreateFileReference(6, 7, 8, 9),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseTokenPermissions(testCase.permissions)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedPermissions, got, testCase.name)
		})
	}
}

func TestParsePermissionValue(t *testing.T) {
	testCases := []struct {
		name        string
		permission  string
		expectedVal models.Permission
	}{
		{
			name:       "Read permission",
			permission: readPermission,
			expectedVal: models.Permission{
				Read: true,
			},
		},
		{
			name:       "Write permission",
			permission: writePermission,
			expectedVal: models.Permission{
				Write: true,
			},
		},
		{
			name:        "Mo read no write permissions",
			permission:  "no_permission",
			expectedVal: models.Permission{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parsePermissionValue(testCase.permission)

			assert.Equal(t, testCase.expectedVal, got, testCase.name)
		})
	}
}

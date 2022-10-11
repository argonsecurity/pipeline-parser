package gitlab

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseDefaults(t *testing.T) {
	testCases := []struct {
		name                  string
		gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration
		expectedDefaults      *models.Defaults
	}{
		{
			name:                  "Gitlab CI config is empty",
			gitlabCIConfiguration: &gitlabModels.GitlabCIConfiguration{},
			expectedDefaults:      &models.Defaults{},
		},
		{
			name: "Gitlab CI config with defaults data",
			gitlabCIConfiguration: &gitlabModels.GitlabCIConfiguration{
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				Image: &common.Image{
					Name:          "registry/namespace/image:tag",
					Entrypoint:    []string{"entrypoint"},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				AfterScript: &common.Script{
					Commands:      []string{"command-1", "command-2"},
					FileReference: testutils.CreateFileReference(1, 2, 1, 10),
				},
				BeforeScript: &common.Script{
					Commands:      []string{"command-1", "command-2"},
					FileReference: testutils.CreateFileReference(2, 3, 2, 20),
				},
				Default: &gitlabModels.Default{
					Artifacts: &gitlabModels.Artifacts{
						Reports: &gitlabModels.Reports{
							SecretDetection:    "",
							Sast:               "",
							DependencyScanning: "",
							Terraform:          "",
							LicenseScanning:    "",
						},
					},
				},
			},
			expectedDefaults: &models.Defaults{
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				Scans: &models.Scans{
					Secrets:      utils.GetPtr(true),
					SAST:         utils.GetPtr(true),
					Dependencies: utils.GetPtr(true),
					Iac:          utils.GetPtr(true),
					License:      utils.GetPtr(true),
				},
				Runner: &models.Runner{
					DockerMetadata: &models.DockerMetadata{
						Image:       utils.GetPtr("namespace/image"),
						Label:       utils.GetPtr("tag"),
						RegistryURL: utils.GetPtr("registry"),
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				PostSteps: []*models.Step{
					{
						Type: "shell",
						Shell: &models.Shell{
							Script: utils.GetPtr("command-1"),
						},
						FileReference: testutils.CreateFileReference(2, 4, 2, 19),
					},
					{
						Type: "shell",
						Shell: &models.Shell{
							Script: utils.GetPtr("command-2"),
						},
						FileReference: testutils.CreateFileReference(3, 4, 3, 19),
					},
				},
				PreSteps: []*models.Step{
					{
						Type: "shell",
						Shell: &models.Shell{
							Script: utils.GetPtr("command-1"),
						},
						FileReference: testutils.CreateFileReference(3, 5, 3, 29),
					},
					{
						Type: "shell",
						Shell: &models.Shell{
							Script: utils.GetPtr("command-2"),
						},
						FileReference: testutils.CreateFileReference(4, 5, 4, 29),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseDefaults(testCase.gitlabCIConfiguration)

			testutils.DeepCompare(t, testCase.expectedDefaults, got)
		})
	}

}

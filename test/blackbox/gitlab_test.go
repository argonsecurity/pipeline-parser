package blackbox

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestGitLab(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "gradle.yaml",
			Expected: &models.Pipeline{
				Jobs: SortJobs([]*models.Job{
					{
						ID:               utils.GetPtr("test"),
						Name:             utils.GetPtr("test"),
						ConcurrencyGroup: utils.GetPtr("test"),
						Steps: []*models.Step{
							{
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script: utils.GetPtr("gradle check"),
								},
								FileReference: testutils.CreateFileReference(35, 3, 35, 15),
							},
						},
						Metadata: models.Metadata{
							Test: true,
						},
					},
					{
						ID:               utils.GetPtr("build"),
						Name:             utils.GetPtr("build"),
						ConcurrencyGroup: utils.GetPtr("build"),
						Steps: []*models.Step{
							{
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script: utils.GetPtr("gradle --build-cache assemble"),
								},
								FileReference: testutils.CreateFileReference(25, 3, 25, 32),
							},
						},
						Metadata: models.Metadata{
							Build: true,
						},
					},
				}),
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("gradle"),
							Label: utils.GetPtr("alpine"),
						},
						FileReference: testutils.CreateFileReference(10, 1, 10, 8),
					},
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"GRADLE_OPTS": "-Dorg.gradle.daemon=false",
						},
						FileReference: testutils.CreateFileReference(16, 0, 17, 16),
					},
					PreSteps: []*models.Step{
						{
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Script: utils.GetPtr(`GRADLE_USER_HOME="$(pwd)/.gradle"`),
							},
							FileReference: testutils.CreateFileReference(20, 3, 20, 38),
						},
						{
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Script: utils.GetPtr(`export GRADLE_USER_HOME`),
							},
							FileReference: testutils.CreateFileReference(21, 3, 21, 28),
						},
					},
				},
			},
		},
	}

	executeTestCases(t, testCases, "gitlab", consts.GitLabPlatform)
}

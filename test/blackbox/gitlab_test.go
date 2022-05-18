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
						ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("test")),
						Steps: []*models.Step{
							{
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script: utils.GetPtr("gradle check"),
								},
								FileReference: testutils.CreateFileReference(35, 3, 35, 23),
							},
						},
						Metadata: models.Metadata{
							Test: true,
						},
						FileReference: testutils.CreateFileReference(33, 1, 35, 23),
					},
					{
						ID:               utils.GetPtr("build"),
						Name:             utils.GetPtr("build"),
						ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("build")),
						Steps: []*models.Step{
							{
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script: utils.GetPtr("gradle --build-cache assemble"),
								},
								FileReference: testutils.CreateFileReference(25, 3, 25, 40),
							},
						},
						Metadata: models.Metadata{
							Build: true,
						},
						FileReference: testutils.CreateFileReference(23, 1, 31, 16),
					},
				}),
				Defaults: &models.Defaults{
					Runner: &models.Runner{
						DockerMetadata: &models.DockerMetadata{
							Image: utils.GetPtr("gradle"),
							Label: utils.GetPtr("alpine"),
						},
						FileReference: testutils.CreateFileReference(10, 1, 10, 21),
					},
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"GRADLE_OPTS": "-Dorg.gradle.daemon=false",
						},
						FileReference: testutils.CreateFileReference(16, 1, 17, 43),
					},
					PreSteps: []*models.Step{
						{
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Script: utils.GetPtr(`GRADLE_USER_HOME="$(pwd)/.gradle"`),
							},
							FileReference: testutils.CreateFileReference(20, 3, 20, 61),
						},
						{
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Script: utils.GetPtr(`export GRADLE_USER_HOME`),
							},
							FileReference: testutils.CreateFileReference(21, 3, 21, 51),
						},
					},
				},
			},
		},
		{
			Filename: "terraform.yaml",
			Expected: &models.Pipeline{
				Imports: []string{
					"Terraform/Base.latest.gitlab-ci.yml",
					"Jobs/SAST-IaC.latest.gitlab-ci.yml",
				},
				Jobs: SortJobs([]*models.Job{
					{
						ID:            utils.GetPtr("fmt"),
						Name:          utils.GetPtr("fmt"),
						FileReference: testutils.CreateFileReference(16, 1, 18, 10),
					},
					{
						ID:            utils.GetPtr("validate"),
						Name:          utils.GetPtr("validate"),
						FileReference: testutils.CreateFileReference(20, 1, 22, 10),
					},
					{
						ID:            utils.GetPtr("build"),
						Name:          utils.GetPtr("build"),
						FileReference: testutils.CreateFileReference(24, 1, 25, 28),
						Metadata: models.Metadata{
							Build: true,
						},
					},
					{
						ID:   utils.GetPtr("deploy"),
						Name: utils.GetPtr("deploy"),
						Dependencies: []*models.JobDependency{
							{
								JobID: utils.GetPtr("build"),
							},
						},
						FileReference: testutils.CreateFileReference(27, 1, 32, 25),
					},
				}),
				Defaults: &models.Defaults{},
			},
		},
	}

	executeTestCases(t, testCases, "gitlab", consts.GitLabPlatform)
}

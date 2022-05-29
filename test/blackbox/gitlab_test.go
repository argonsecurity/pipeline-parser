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
						FileReference: testutils.CreateFileReference(10, 8, 10, 28),
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
		{
			Filename: "build-job.yaml",
			Expected: &models.Pipeline{
				Triggers: &models.Triggers{
					FileReference: testutils.CreateFileReference(22, 3, 25, 18),
					Triggers: []*models.Trigger{
						{
							Event:         models.PullRequestEvent,
							FileReference: testutils.CreateFileReference(23, 7, 23, 55),
						},
					},
				},
				Jobs: SortJobs([]*models.Job{
					{
						ID:   utils.GetPtr("python-build"),
						Name: utils.GetPtr("python-build"),
						Conditions: []*models.Condition{
							{
								Statement: "$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME =~ /^feature/",
								Allow:     utils.GetPtr(true),
							},
							{
								Branches: &models.Filter{
									AllowList: []string{"/^feature-.*/", "main"},
								},
								Allow:  utils.GetPtr(true),
								Events: []models.EventType{models.PullRequestEvent, models.EventType("api")},
							},
						},
						Steps: []*models.Step{
							{
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script: utils.GetPtr("cd requests"),
								},
								FileReference: testutils.CreateFileReference(15, 5, 15, 40),
							},
							{
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script: utils.GetPtr("python3 setup.py sdist"),
								},
								FileReference: testutils.CreateFileReference(16, 5, 16, 51),
							},
						},
						ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("build")),
						Metadata:         models.Metadata{Build: true},
						FileReference:    testutils.CreateFileReference(4, 1, 16, 29),
					},
				}),
				Defaults: &models.Defaults{
					Scans: &models.Scans{
						Secrets:      utils.GetPtr(true),
						SAST:         utils.GetPtr(true),
						Iac:          utils.GetPtr(true),
						Dependencies: utils.GetPtr(true),
						License:      utils.GetPtr(true),
					},
					PreSteps: []*models.Step{
						{
							Type: models.ShellStepType,
							Shell: &models.Shell{
								Script: utils.GetPtr(`echo "before_script"`),
							},
							FileReference: testutils.CreateFileReference(18, 1, 19, 25),
						},
					},
					Conditions: []*models.Condition{
						{
							Statement: "$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH",
							Allow:     utils.GetPtr(false),
						},
					},
				},
			},
		},
	}

	executeTestCases(t, testCases, "gitlab", consts.GitLabPlatform)
}

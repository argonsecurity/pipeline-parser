package blackbox

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func getAllGitHubPermissions(permission models.Permission) *models.TokenPermissions {
	allPermissions := map[string]models.Permission{
		"run-pipeline":        permission,
		"checks":              permission,
		"contents":            permission,
		"deployments":         permission,
		"discussions":         permission,
		"id-token":            permission,
		"issues":              permission,
		"packages":            permission,
		"pages":               permission,
		"pull-request":        permission,
		"repository-projects": permission,
		"security-events":     permission,
		"statuses":            permission,
	}
	return &models.TokenPermissions{
		Permissions: allPermissions,
	}
}

func TestGitHub(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "steps.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("steps"),
				Jobs: SortJobs([]*models.Job{
					{
						ID:   utils.GetPtr("job1"),
						Name: utils.GetPtr("Job 1"),
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("task without params"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
								},
								FileReference: testutils.CreateFileReference(7, 9, 8, 34),
							},
							{
								Name: utils.GetPtr("task with params"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
									Inputs: []*models.Parameter{
										{
											Name:          utils.GetPtr("repo"),
											Value:         "repository",
											FileReference: testutils.CreateFileReference(13, 11, 13, 21), // End column is supposed to be 27
										},
									},
								},
								FileReference: testutils.CreateFileReference(10, 9, 13, 21), // End column is supposed to be 27
							},
							{
								Name: utils.GetPtr("task with multiline params"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1"),
									VersionType: "tag",
									Inputs: []*models.Parameter{
										{
											Name:          utils.GetPtr("repos"),
											Value:         "repository1\nrepository2\n",
											FileReference: testutils.CreateFileReference(18, 11, 20, 11), // End column is supposed to be 29
										},
										{
											Name:          utils.GetPtr("input"),
											Value:         "value",
											FileReference: testutils.CreateFileReference(21, 11, 21, 16), // End column is supposed to be 23
										},
									},
								},
								FileReference: testutils.CreateFileReference(15, 9, 21, 16), // End column is supposed to be 23
							},
							{
								Name: utils.GetPtr("task with commit ID version"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("c44948622e1b6bb0eb0cec5b813c1ac561158e1e"),
									VersionType: "commit",
								},
								FileReference: testutils.CreateFileReference(23, 9, 24, 72),
							},
							{
								Name: utils.GetPtr("task with branch version"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("master"),
									VersionType: "branch",
								},
								FileReference: testutils.CreateFileReference(26, 9, 27, 38),
							},
							{
								Name: utils.GetPtr("task with tag version"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("v1.1.1"),
									VersionType: "tag",
								},
								FileReference: testutils.CreateFileReference(29, 9, 30, 38),
							},
							{
								Name: utils.GetPtr("shell"),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script:        utils.GetPtr("command line"),
									FileReference: testutils.CreateFileReference(33, 14, 33, 26),
								},
								FileReference: testutils.CreateFileReference(32, 9, 33, 26),
							},
							{
								Name: utils.GetPtr("custom shell"),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script:        utils.GetPtr("command line"),
									FileReference: testutils.CreateFileReference(37, 14, 37, 26),
									Type:          utils.GetPtr("cmd"),
								},
								FileReference: testutils.CreateFileReference(35, 9, 37, 26),
							}, {
								Name: utils.GetPtr("shell with break rows"),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Script:        utils.GetPtr("echo 1\necho 2\necho 3\n"),
									FileReference: testutils.CreateFileReference(40, 14, 43, 14),
								},
								FileReference: testutils.CreateFileReference(39, 9, 43, 14),
							},
						},
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
						FileReference:   testutils.CreateFileReference(4, 3, 43, 14),
					},
				}),
			},
		},
		{
			Filename: "dependant-jobs.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("dependable jobs"),
				Jobs: SortJobs([]*models.Job{
					{
						ID:              utils.GetPtr("dependable-job"),
						Name:            utils.GetPtr("Dependable Job"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						FileReference:   testutils.CreateFileReference(4, 3, 5, 25),
					},
					{
						ID:              utils.GetPtr("dependant-job"),
						Name:            utils.GetPtr("Dependant Job"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Dependencies: []*models.JobDependency{
							{
								JobID: utils.GetPtr("dependable-job"),
							},
						},
						FileReference: testutils.CreateFileReference(7, 3, 9, 27),
					},
				}),
			},
		},
		{
			Filename: "all-triggers.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("all-triggers"),
				Triggers: &models.Triggers{
					FileReference: testutils.CreateFileReference(2, 3, 30, 20),
					Triggers: SortTriggers([]*models.Trigger{
						{
							Event:         models.ScheduledEvent,
							Schedules:     &[]string{"30 2 * * *"},
							FileReference: testutils.CreateFileReference(3, 3, 4, 23),
						},
						{
							Event: models.PushEvent,
							Branches: &models.Filter{
								AllowList: []string{"master"},
							},
							FileReference: testutils.CreateFileReference(5, 3, 7, 15),
						},
						{
							Event: models.PipelineRunEvent,
							Branches: &models.Filter{
								DenyList: []string{"master"},
							},
							FileReference: testutils.CreateFileReference(26, 3, 28, 15),
						},
						{
							Event: models.PullRequestEvent,
							Paths: &models.Filter{
								DenyList: []string{"*/test/*"},
							},
							FileReference: testutils.CreateFileReference(8, 3, 10, 17),
						},
						{
							Event: "pull_request_target",
							Paths: &models.Filter{
								AllowList: []string{"*/test/*"},
							},
							FileReference: testutils.CreateFileReference(11, 3, 13, 17),
						},
						{
							Event: models.ManualEvent,
							Parameters: []models.Parameter{
								{
									Name:        utils.GetPtr("workflow-input"),
									Description: utils.GetPtr("The workflow input"),
									Default:     "default-value",
								},
							},
							FileReference: testutils.CreateFileReference(14, 3, 19, 23),
						},
						{
							Event: models.PipelineTriggerEvent,
							Parameters: []models.Parameter{
								{
									Name:        utils.GetPtr("workflow-input"),
									Description: utils.GetPtr("The workflow input"),
									Default:     "default-value",
								},
							},
							FileReference: testutils.CreateFileReference(20, 3, 25, 23),
						},
						{
							Event: "label",
							Filters: map[string]any{
								"types": []string{"created"},
							},
							FileReference: testutils.CreateFileReference(29, 3, 30, 20),
						},
					}),
				},
			},
		},
		{
			Filename: "token-permissions.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("permissions"),
				Jobs: SortJobs([]*models.Job{
					{
						FileReference:    testutils.CreateFileReference(8, 3, 10, 26),
						ID:               utils.GetPtr("job1"),
						Name:             utils.GetPtr("Job 1"),
						ContinueOnError:  utils.GetPtr(false),
						TokenPermissions: getAllGitHubPermissions(models.Permission{Read: true}),
						TimeoutMS:        utils.GetPtr(21600000),
					},
				}),
				Defaults: &models.Defaults{
					TokenPermissions: &models.TokenPermissions{
						FileReference: testutils.CreateFileReference(2, 3, 5, 22),
						Permissions: map[string]models.Permission{
							"run-pipeline": {
								Read: true,
							},
							"statuses": {
								Write: true,
							},
							"pull-request": {
								Read: true,
							},
						},
					},
				},
			},
		},
		{
			Filename: "runners.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("runners"),
				Jobs: SortJobs([]*models.Job{
					{
						ID:              utils.GetPtr("job1"),
						Name:            utils.GetPtr("Job 1"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Runner: &models.Runner{
							OS:            utils.GetPtr("linux"),
							Labels:        &[]string{"ubuntu-latest"},
							SelfHosted:    utils.GetPtr(false),
							FileReference: testutils.CreateFileReference(6, 14, 6, 27),
						},
						FileReference: testutils.CreateFileReference(4, 3, 6, 27),
					},
					{
						ID:              utils.GetPtr("job2"),
						Name:            utils.GetPtr("Job 2"),
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
						Runner: &models.Runner{
							OS:            utils.GetPtr("windows"),
							Labels:        &[]string{"self-hosted", "windows-latest"},
							SelfHosted:    utils.GetPtr(true),
							FileReference: testutils.CreateFileReference(9, 14, 9, 43),
						},
						FileReference: testutils.CreateFileReference(7, 3, 9, 42),
					},
					{
						ID:              utils.GetPtr("job3"),
						Name:            utils.GetPtr("Job 3"),
						TimeoutMS:       utils.GetPtr(21600000),
						ContinueOnError: utils.GetPtr(false),
						Runner: &models.Runner{
							OS:            utils.GetPtr("linux"),
							Arch:          utils.GetPtr("x64"),
							Labels:        &[]string{"self-hosted", "linux", "x64"},
							SelfHosted:    utils.GetPtr(true),
							FileReference: testutils.CreateFileReference(12, 14, 12, 39),
						},
						FileReference: testutils.CreateFileReference(10, 3, 12, 38),
					},
				}),
			},
		},
		{
			Filename: "environment-variables.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("environment-variables"),
				Jobs: SortJobs([]*models.Job{
					{
						ID:   utils.GetPtr("job1"),
						Name: utils.GetPtr("Job 1"),
						EnvironmentVariables: &models.EnvironmentVariablesRef{
							EnvironmentVariables: models.EnvironmentVariables{
								"STRING": "string",
								"NUMBER": 1,
							},
							FileReference: testutils.CreateFileReference(10, 7, 12, 16),
						},
						FileReference:   testutils.CreateFileReference(8, 3, 18, 20),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Step 1"),
								Type: "shell",
								Shell: &models.Shell{
									Script:        utils.GetPtr("command line"),
									FileReference: testutils.CreateFileReference(15, 14, 15, 26),
								},
								FileReference: testutils.CreateFileReference(14, 9, 18, 20),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"STRING": "string",
										"NUMBER": 1,
									},
									FileReference: testutils.CreateFileReference(16, 11, 18, 20),
								},
							},
						},
					},
				}),
				Defaults: &models.Defaults{
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"STRING": "string",
							"NUMBER": 1,
						},
						FileReference: testutils.CreateFileReference(3, 3, 5, 12),
					},
				},
			},
		},
		{
			Filename: "concurrent-jobs.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("concurrent-jobs"),
				Jobs: SortJobs([]*models.Job{
					{
						ID:               utils.GetPtr("job1"),
						Name:             utils.GetPtr("Job 1"),
						ContinueOnError:  utils.GetPtr(false),
						TimeoutMS:        utils.GetPtr(21600000),
						ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("ci")),
						FileReference:    testutils.CreateFileReference(3, 3, 5, 20),
					},
					{
						ID:               utils.GetPtr("job2"),
						Name:             utils.GetPtr("Job 2"),
						ContinueOnError:  utils.GetPtr(false),
						TimeoutMS:        utils.GetPtr(21600000),
						ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("ci")),
						FileReference:    testutils.CreateFileReference(7, 3, 9, 20),
					},
				}),
			},
		},
	}
	executeTestCases(t, testCases, "github", consts.GitHubPlatform)
}

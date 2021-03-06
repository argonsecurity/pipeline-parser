package blackbox

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestAzure(t *testing.T) {
	testCases := []TestCase{
		{
			Filename: "all-triggers.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("all-triggers"),
				Defaults: &models.Defaults{},
				Triggers: &models.Triggers{
					Triggers: []*models.Trigger{
						{
							Event: models.PullRequestEvent,
							Branches: &models.Filter{
								AllowList: []string{"features/*"},
								DenyList:  []string{"features/experimental/*"},
							},
							Paths: &models.Filter{
								AllowList: []string{"path/to/file"},
								DenyList:  []string{"README.md"},
							},
							FileReference: testutils.CreateFileReference(21, 1, 33, 15),
						},
						{
							Event: models.PushEvent,
							Branches: &models.Filter{
								AllowList: []string{"master", "main"},
								DenyList:  []string{"test/*"},
							},
							Paths: &models.Filter{
								AllowList: []string{"path/to/file", "another/path/to/file"},
								DenyList:  []string{"all/*"},
							},
							Tags: &models.Filter{
								AllowList: []string{"v1.0.*"},
								DenyList:  []string{"v2.0.*"},
							},
							FileReference: testutils.CreateFileReference(2, 1, 20, 15),
						},
						{
							Event:         models.ScheduledEvent,
							Schedules:     &[]string{"0 0 * * *", "0 12 * * 0"},
							FileReference: testutils.CreateFileReference(34, -1, 48, 15),
						},
					},
					FileReference: testutils.CreateFileReference(2, 1, 48, 15),
				},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
			},
		},
		{
			Filename: "branch-list-trigger.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("branch-list-trigger"),
				Defaults: &models.Defaults{},
				Triggers: &models.Triggers{
					Triggers: []*models.Trigger{
						{
							Event: models.PullRequestEvent,
							Branches: &models.Filter{
								AllowList: []string{"main", "develop"},
							},
							FileReference: testutils.CreateFileReference(5, -1, 7, 10),
						},
						{
							Event: models.PushEvent,
							Branches: &models.Filter{
								AllowList: []string{"main", "development"},
							},
							FileReference: testutils.CreateFileReference(2, 1, 4, 16),
						},
					},
					FileReference: testutils.CreateFileReference(2, 1, 7, 10),
				},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
			},
		},
		{
			Filename: "no-trigger.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("no-trigger"),
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
				Triggers: nil,
			},
		},
		{
			Filename: "jobs.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("Jobs"),
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						ID:              utils.GetPtr("DeployWeb"),
						Name:            utils.GetPtr("deploy Web App"),
						TimeoutMS:       utils.GetPtr(3600000),
						ContinueOnError: utils.GetPtr(false),
						Dependencies: []*models.JobDependency{
							{JobID: utils.GetPtr("job1")},
							{JobID: utils.GetPtr("job2")},
						},
						FileReference: testutils.CreateFileReference(12, 3, 24, 43),
					},
					{
						ID:              utils.GetPtr("MyJob"),
						Name:            utils.GetPtr("My First Job"),
						TimeoutMS:       utils.GetPtr(3600000),
						ContinueOnError: utils.GetPtr(true),
						Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job")}},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr(""),
									Script: utils.GetPtr("echo My first job"),
								},
								FileReference: testutils.CreateFileReference(11, 5, 11, 30),
							},
						},
						FileReference: testutils.CreateFileReference(4, 3, 11, 30),
					},
				},
			},
		},
		{
			Filename: "parameters.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("parameters"),
				Defaults: &models.Defaults{},
				Parameters: []*models.Parameter{
					{
						Name:          utils.GetPtr("myString"),
						Default:       "a string",
						FileReference: testutils.CreateFileReference(3, 3, 5, 20),
					},
					{
						Name:          utils.GetPtr("myMultiString"),
						Default:       "default",
						Options:       []string{"default", "ubuntu"},
						FileReference: testutils.CreateFileReference(6, 3, 11, 11),
					},
					{
						Name:          utils.GetPtr("myNumber"),
						Default:       2,
						Options:       []string{"1", "2", "4", "8", "16"},
						FileReference: testutils.CreateFileReference(12, 3, 20, 7),
					},
					{
						Name:          utils.GetPtr("myBoolean"),
						Default:       true,
						FileReference: testutils.CreateFileReference(21, 3, 23, 16),
					},
					{
						Name: utils.GetPtr("myObject"),
						Default: map[string]any{
							"foo":    "FOO",
							"bar":    "BAR",
							"things": []any{"one", "two", "three"},
							"nested": map[string]any{
								"one":   "apple",
								"two":   "pear",
								"count": 3,
							},
						},
						FileReference: testutils.CreateFileReference(24, 3, 36, 15),
					},
					{
						Name: utils.GetPtr("myStep"),
						Default: map[string]any{
							"script": "echo my step",
						},
						FileReference: testutils.CreateFileReference(37, 3, 40, 25),
					},
				},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
			},
		},
		{
			Filename: "pool.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("pool"),
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						ID:              utils.GetPtr("jobId"),
						Name:            utils.GetPtr(""),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(3600000),
						FileReference:   testutils.CreateFileReference(9, 3, 13, 18),
					},
				},
			},
		},
		{
			Filename: "resources.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("resources"),
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
			},
		},
		{
			Filename: "stages.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("stages"),
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
			},
		},
		{
			Filename: "steps.yaml",
			Expected: &models.Pipeline{
				Name:     utils.GetPtr("steps"),
				Defaults: &models.Defaults{},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
						Steps: []*models.Step{
							{
								Name: utils.GetPtr("Multiline Bash script"),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"name": "Microsoft",
									},
									FileReference: testutils.CreateFileReference(7, 5, 8, 20),
								},
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("bash"),
									Script: utils.GetPtr("which bash\necho Hello $name\n"),
								},
								FileReference: testutils.CreateFileReference(3, 3, 8, 20),
							},
							{
								Name:          utils.GetPtr(""),
								FileReference: testutils.CreateFileReference(9, 3, 11, 27),
							},
							{
								Name:          utils.GetPtr("Download artifact WebApp"),
								FileReference: testutils.CreateFileReference(12, 3, 15, 40),
							},
							{
								Name:          utils.GetPtr("Download artifact WebApp"),
								FileReference: testutils.CreateFileReference(16, 3, 20, 40),
							},
							{
								Name:          utils.GetPtr(""),
								FileReference: testutils.CreateFileReference(21, 3, 22, 13),
							},
							{
								ID:   utils.GetPtr("firstStep"),
								Name: utils.GetPtr("Say hello"),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"name": "Microsoft",
									},
									FileReference: testutils.CreateFileReference(28, 5, 29, 20),
								},
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("powershell"),
									Script: utils.GetPtr("Write-Host Hello $(name)"),
								},
								WorkingDirectory: utils.GetPtr("$(build.sourcesDirectory)"),
								FileReference:    testutils.CreateFileReference(23, 3, 29, 20),
							},
							{
								Name:          utils.GetPtr("Publish artifact WebApp"),
								FileReference: testutils.CreateFileReference(30, 3, 32, 39),
							},
							{
								ID:   utils.GetPtr("firstStep"),
								Name: utils.GetPtr("Say hello"),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"name": "Microsoft",
									},
									FileReference: testutils.CreateFileReference(38, 5, 39, 20),
								},
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("powershell core"),
									Script: utils.GetPtr("Write-Host Hello $(name)"),
								},
								WorkingDirectory: utils.GetPtr("$(build.sourcesDirectory)"),
								FileReference:    testutils.CreateFileReference(33, 3, 39, 20),
							},
							{
								Name:          utils.GetPtr(""),
								FileReference: testutils.CreateFileReference(40, 3, 40, 20),
							},
							{
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr(""),
									Script: utils.GetPtr("echo This runs in the default shell on any machine"),
								},
								FileReference: testutils.CreateFileReference(41, 3, 41, 61),
							},
							{
								Name: utils.GetPtr("Build"),
								Type: models.TaskStepType,
								Task: &models.Task{
									Name:        utils.GetPtr("VSBuild"),
									Version:     utils.GetPtr("1"),
									VersionType: models.TagVersion,
									Inputs: &[]models.Parameter{
										{
											Name:          utils.GetPtr("solution"),
											Value:         "**\\*.sln",
											FileReference: testutils.CreateFileReference(46, 5, 46, 23),
										},
									},
								},
								Timeout: utils.GetPtr(7200000),
								Metadata: models.Metadata{
									Build: true,
								},
								FileReference: testutils.CreateFileReference(42, 3, 46, 23),
							},
							{
								Name:          utils.GetPtr(""),
								FileReference: testutils.CreateFileReference(47, 3, 49, 15),
							},
						},
						Metadata: models.Metadata{
							Build: true,
						},
					},
				},
			},
		},
		{
			Filename: "variables.yaml",
			Expected: &models.Pipeline{
				Name: utils.GetPtr("variables"),
				Defaults: &models.Defaults{
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"var1": "value1",
							"var2": "value2",
						},
						FileReference: testutils.CreateFileReference(3, 3, 11, 17),
					},
				},
				Jobs: []*models.Job{
					{
						ID:              utils.GetPtr("FirstJob"),
						Name:            utils.GetPtr(""),
						TimeoutMS:       utils.GetPtr(3600000),
						ContinueOnError: utils.GetPtr(false),
						EnvironmentVariables: &models.EnvironmentVariablesRef{
							EnvironmentVariables: models.EnvironmentVariables{
								"JOB_VAR": "a job var",
							},
							FileReference: testutils.CreateFileReference(21, 7, 21, 16),
						},
						Steps: []*models.Step{
							{
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr(""),
									Script: utils.GetPtr("echo $(MY_VAR) $(STAGE_VAR) $(JOB_VAR)"),
								},
								FileReference: testutils.CreateFileReference(23, 7, 23, 53),
							},
						},
						FileReference: testutils.CreateFileReference(19, 5, 23, 53),
					},
				},
			},
		},
	}

	executeTestCases(t, testCases, "azure", consts.AzurePlatform)
}

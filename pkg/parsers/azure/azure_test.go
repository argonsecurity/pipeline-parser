package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name             string
		azurePipeline    *azureModels.Pipeline
		expectedPipeline *models.Pipeline
	}{
		{
			name:             "Pipeline is nil",
			azurePipeline:    nil,
			expectedPipeline: nil,
		},
		{
			name:          "Empty pipeline",
			azurePipeline: &azureModels.Pipeline{},
			expectedPipeline: &models.Pipeline{
				Name: utils.GetPtr(""),
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
					},
				},
			},
		},
		{
			name: "Pipeline with jobs",
			azurePipeline: &azureModels.Pipeline{
				Name:            "pipeline",
				ContinueOnError: utils.GetPtr(true),
				Trigger: &azureModels.TriggerRef{
					Trigger: &azureModels.Trigger{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				PR: &azureModels.PRRef{
					PR: &azureModels.PR{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				Schedules: &azureModels.Schedules{
					Crons: &[]azureModels.Cron{
						{
							Cron:          "1 * * * *",
							FileReference: testutils.CreateFileReference(9, 10, 11, 12),
						},
					},
					FileReference: testutils.CreateFileReference(13, 14, 15, 16),
				},
				Jobs: &azureModels.Jobs{
					CIJobs: []*azureModels.CIJob{
						{
							Job: "job-1",
							BaseJob: azureModels.BaseJob{
								DisplayName:      "job-1",
								DependsOn:        &azureModels.DependsOn{"job-2"},
								Condition:        "job-1-condition",
								ContinueOnError:  true,
								TimeoutInMinutes: 100,
								Steps: &azureModels.Steps{
									{
										Name:      "step-1",
										Bash:      "script",
										Condition: "step-1-condition",
									},
								},
							},
						},
						{
							Job: "job-2",
							BaseJob: azureModels.BaseJob{
								DisplayName:      "job-2",
								DependsOn:        &azureModels.DependsOn{"job-3"},
								Condition:        "job-2-condition",
								ContinueOnError:  true,
								TimeoutInMinutes: 100,
								Steps: &azureModels.Steps{
									{
										Name:      "step-2",
										Bash:      "script",
										Condition: "step-2-condition",
									},
								},
							},
						},
					},
					DeploymentJobs: []*azureModels.DeploymentJob{
						{
							Deployment: "deployment-1",
							BaseJob: azureModels.BaseJob{
								DisplayName:      "job-1",
								DependsOn:        &azureModels.DependsOn{"job-2"},
								Condition:        "job-1-condition",
								ContinueOnError:  true,
								TimeoutInMinutes: 100,
								Steps: &azureModels.Steps{
									{
										Name:      "step-1",
										Bash:      "script",
										Condition: "step-1-condition",
									},
								},
							},
						},
						{
							Deployment: "deployment-2",
							BaseJob: azureModels.BaseJob{
								DisplayName:      "job-2",
								DependsOn:        &azureModels.DependsOn{"job-3"},
								Condition:        "job-2-condition",
								ContinueOnError:  true,
								TimeoutInMinutes: 100,
								Steps: &azureModels.Steps{
									{
										Name:      "step-2",
										Bash:      "script",
										Condition: "step-2-condition",
									},
								},
							},
						},
					},
				},
				Parameters: &azureModels.Parameters{
					{
						Name:          "param1",
						DisplayName:   "Param 1",
						Type:          "string",
						Default:       "default1",
						Values:        []string{"value1", "value2"},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Name:          "param2",
						DisplayName:   "Param 2",
						Type:          "number",
						Default:       2,
						Values:        []string{"1", "2"},
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
				},
			},
			expectedPipeline: &models.Pipeline{
				Name: utils.GetPtr("pipeline"),
				Defaults: &models.Defaults{
					ContinueOnError: utils.GetPtr(true),
				},
				Triggers: &models.Triggers{
					Triggers: []*models.Trigger{
						{
							Event: models.PushEvent,
							Branches: &models.Filter{
								AllowList: []string{"master"},
							},
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
						{
							Event: models.PullRequestEvent,
							Branches: &models.Filter{
								AllowList: []string{"master"},
							},
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
						{
							Event:         models.ScheduledEvent,
							Schedules:     &[]string{"1 * * * *"},
							FileReference: testutils.CreateFileReference(13, 14, 15, 16),
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				Jobs: []*models.Job{
					{
						ID:              utils.GetPtr("job-1"),
						Name:            utils.GetPtr("job-1"),
						ContinueOnError: utils.GetPtr(true),
						TimeoutMS:       utils.GetPtr(6000000),
						Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
						Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
						Steps: []*models.Step{
							{
								ID:   utils.GetPtr("step-1"),
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("bash"),
									Script: utils.GetPtr("script"),
								},
								Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
							},
						},
					},
					{
						ID:              utils.GetPtr("job-2"),
						Name:            utils.GetPtr("job-2"),
						ContinueOnError: utils.GetPtr(true),
						TimeoutMS:       utils.GetPtr(6000000),
						Conditions:      []*models.Condition{{Statement: "job-2-condition"}},
						Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-3")}},
						Steps: []*models.Step{
							{
								ID:   utils.GetPtr("step-2"),
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("bash"),
									Script: utils.GetPtr("script"),
								},
								Conditions: &[]models.Condition{{Statement: "step-2-condition"}},
							},
						},
					},
					{
						ID:              utils.GetPtr("deployment-1"),
						Name:            utils.GetPtr("job-1"),
						ContinueOnError: utils.GetPtr(true),
						TimeoutMS:       utils.GetPtr(6000000),
						Conditions:      []*models.Condition{{Statement: "job-1-condition"}},
						Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-2")}},
						Steps: []*models.Step{
							{
								ID:   utils.GetPtr("step-1"),
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("bash"),
									Script: utils.GetPtr("script"),
								},
								Conditions: &[]models.Condition{{Statement: "step-1-condition"}},
							},
						},
					},
					{
						ID:              utils.GetPtr("deployment-2"),
						Name:            utils.GetPtr("job-2"),
						ContinueOnError: utils.GetPtr(true),
						TimeoutMS:       utils.GetPtr(6000000),
						Conditions:      []*models.Condition{{Statement: "job-2-condition"}},
						Dependencies:    []*models.JobDependency{{JobID: utils.GetPtr("job-3")}},
						Steps: []*models.Step{
							{
								ID:   utils.GetPtr("step-2"),
								Name: utils.GetPtr(""),
								Type: models.ShellStepType,
								Shell: &models.Shell{
									Type:   utils.GetPtr("bash"),
									Script: utils.GetPtr("script"),
								},
								Conditions: &[]models.Condition{{Statement: "step-2-condition"}},
							},
						},
					},
				},
				Parameters: []*models.Parameter{
					{
						Name:          utils.GetPtr("param1"),
						Default:       "default1",
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Name:          utils.GetPtr("param2"),
						Default:       2,
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
				},
			},
		},
		{
			name: "Pipeline with steps",
			azurePipeline: &azureModels.Pipeline{
				Name:            "pipeline",
				ContinueOnError: utils.GetPtr(true),
				Trigger: &azureModels.TriggerRef{
					Trigger: &azureModels.Trigger{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				PR: &azureModels.PRRef{
					PR: &azureModels.PR{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				Schedules: &azureModels.Schedules{
					Crons: &[]azureModels.Cron{
						{
							Cron:          "1 * * * *",
							FileReference: testutils.CreateFileReference(9, 10, 11, 12),
						},
					},
					FileReference: testutils.CreateFileReference(13, 14, 15, 16),
				},
				Steps: &azureModels.Steps{
					{
						Name:        "1",
						DisplayName: "step-name1",
						Env: &azureModels.EnvironmentVariablesRef{
							EnvironmentVariables: models.EnvironmentVariables{
								"key": "value1",
							},
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
						FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
						ContinueOnError:  utils.GetPtr(true),
						Condition:        "condition1",
						TimeoutInMinutes: 1,
						WorkingDirectory: "dir",
						Bash:             "script",
					},
					{
						Name:        "2",
						DisplayName: "step-name2",
						Env: &azureModels.EnvironmentVariablesRef{
							EnvironmentVariables: models.EnvironmentVariables{
								"key": "value2",
							},
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
						FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
						ContinueOnError:  utils.GetPtr(true),
						Condition:        "condition2",
						TimeoutInMinutes: 1,
						WorkingDirectory: "dir",
						Task:             "Task@2",
						Inputs: &azureModels.TaskInputs{
							Inputs:        map[string]any{"key": "value"},
							FileReference: testutils.CreateFileReference(111, 222, 333, 444),
						},
					},
				},
				Parameters: &azureModels.Parameters{
					{
						Name:          "param1",
						DisplayName:   "Param 1",
						Type:          "string",
						Default:       "default1",
						Values:        []string{"value1", "value2"},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Name:          "param2",
						DisplayName:   "Param 2",
						Type:          "number",
						Default:       2,
						Values:        []string{"1", "2"},
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
				},
			},
			expectedPipeline: &models.Pipeline{
				Name: utils.GetPtr("pipeline"),
				Defaults: &models.Defaults{
					ContinueOnError: utils.GetPtr(true),
				},
				Triggers: &models.Triggers{
					Triggers: []*models.Trigger{
						{
							Event: models.PushEvent,
							Branches: &models.Filter{
								AllowList: []string{"master"},
							},
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
						{
							Event: models.PullRequestEvent,
							Branches: &models.Filter{
								AllowList: []string{"master"},
							},
							FileReference: testutils.CreateFileReference(5, 6, 7, 8),
						},
						{
							Event:         models.ScheduledEvent,
							Schedules:     &[]string{"1 * * * *"},
							FileReference: testutils.CreateFileReference(13, 14, 15, 16),
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("default"),
						Steps: []*models.Step{
							{
								ID:   utils.GetPtr("1"),
								Name: utils.GetPtr("step-name1"),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"key": "value1",
									},
									FileReference: testutils.CreateFileReference(1, 2, 3, 4),
								},
								FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
								FailsPipeline:    utils.GetPtr(false),
								Conditions:       &[]models.Condition{{Statement: "condition1"}},
								Timeout:          utils.GetPtr(60000),
								WorkingDirectory: utils.GetPtr("dir"),
								Shell: &models.Shell{
									Script: utils.GetPtr("script"),
									Type:   utils.GetPtr("bash"),
								},
								Type: models.ShellStepType,
							},
							{
								ID:   utils.GetPtr("2"),
								Name: utils.GetPtr("step-name2"),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"key": "value2",
									},
									FileReference: testutils.CreateFileReference(1, 2, 3, 4),
								},
								FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
								FailsPipeline:    utils.GetPtr(false),
								Conditions:       &[]models.Condition{{Statement: "condition2"}},
								Timeout:          utils.GetPtr(60000),
								WorkingDirectory: utils.GetPtr("dir"),
								Task: &models.Task{
									Name:        utils.GetPtr("Task"),
									Version:     utils.GetPtr("2"),
									VersionType: models.TagVersion,
									Inputs: &[]models.Parameter{
										{
											Name:          utils.GetPtr("key"),
											Value:         "value",
											FileReference: testutils.CreateFileReference(112, 224, 112, 234),
										},
									},
								},
								Type: models.TaskStepType,
							},
						},
					},
				},
				Parameters: []*models.Parameter{
					{
						Name:          utils.GetPtr("param1"),
						Default:       "default1",
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Name:          utils.GetPtr("param2"),
						Default:       2,
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := AzureParser{}

			pipeline := parser.Parse(testCase.azurePipeline)

			changelog, err := diff.Diff(testCase.expectedPipeline, pipeline)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

package github

import (
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name             string
		workflow         *githubModels.Workflow
		expectedPipeline *models.Pipeline
	}{
		{
			name:     "Empty workflow",
			workflow: &githubModels.Workflow{},
			expectedPipeline: &models.Pipeline{
				Name: utils.GetPtr(""),
			},
		},
		{
			name: "Workflow with data",
			workflow: &githubModels.Workflow{
				Name: "workflow",
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"job-1": {
							ID:              utils.GetPtr("jobid-1"),
							Name:            "job-1",
							ContinueOnError: true,
							Env: &githubModels.EnvironmentVariablesRef{
								EnvironmentVariables: models.EnvironmentVariables{
									"string": "value",
									"int":    1,
									"bool":   false,
								},
								FileReference: testutils.CreateFileReference(1, 2, 3, 4),
							},
							FileReference:  testutils.CreateFileReference(4, 33, 5, 36),
							TimeoutMinutes: utils.GetPtr(float64(10)),
							If:             "condition",
							Concurrency: &githubModels.Concurrency{
								CancelInProgress: utils.GetPtr(true),
								Group:            utils.GetPtr("group-1"),
							},
							Steps: &githubModels.Steps{
								{
									Id:   "1",
									Name: "step-name",
									Env: &githubModels.EnvironmentVariablesRef{
										EnvironmentVariables: models.EnvironmentVariables{
											"key": "value",
										},
										FileReference: testutils.CreateFileReference(1, 2, 3, 4),
									},
									FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
									ContinueOnError:  utils.GetPtr(true),
									If:               "condition",
									TimeoutMinutes:   1,
									WorkingDirectory: "dir",
									Uses:             "actions/checkout@1.2.3",
									With:             map[string]any{"key": "value"},
								},
							},
							RunsOn: &githubModels.RunsOn{
								OS:            utils.GetPtr("linux"),
								Arch:          utils.GetPtr("amd64"),
								SelfHosted:    false,
								Tags:          []string{"tag-1"},
								FileReference: testutils.CreateFileReference(5, 10, 9, 17),
							},
							Needs: &githubModels.Needs{"job-1", "job-2"},
							Permissions: &githubModels.PermissionsEvent{
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
						},
						"job-2": {
							ID: utils.GetPtr("jobid-1"),
						},
					},
				},
				Permissions: &githubModels.PermissionsEvent{
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
				Env: &githubModels.EnvironmentVariablesRef{
					EnvironmentVariables: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
					FileReference: &models.FileReference{
						StartRef: &models.FileLocation{
							Line:   1,
							Column: 2,
						},
						EndRef: &models.FileLocation{
							Line:   3,
							Column: 4,
						},
					},
				},
			},
			expectedPipeline: &models.Pipeline{
				Name: utils.GetPtr("workflow"),
				Jobs: []*models.Job{
					{
						ID:              utils.GetPtr("jobid-1"),
						Name:            utils.GetPtr("job-1"),
						ContinueOnError: utils.GetPtr(true),
						EnvironmentVariables: &models.EnvironmentVariablesRef{
							EnvironmentVariables: models.EnvironmentVariables{
								"string": "value",
								"int":    1,
								"bool":   false,
							},
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
						FileReference:    testutils.CreateFileReference(4, 33, 5, 36),
						TimeoutMS:        utils.GetPtr(600000),
						Conditions:       []*models.Condition{{Statement: "condition"}},
						ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("group-1")),
						Steps: []*models.Step{
							{
								ID:   utils.GetPtr("1"),
								Name: utils.GetPtr("step-name"),
								EnvironmentVariables: &models.EnvironmentVariablesRef{
									EnvironmentVariables: models.EnvironmentVariables{
										"key": "value",
									},
									FileReference: testutils.CreateFileReference(1, 2, 3, 4),
								},
								FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
								FailsPipeline:    utils.GetPtr(false),
								Conditions:       &[]models.Condition{{Statement: "condition"}},
								Timeout:          utils.GetPtr(60000),
								WorkingDirectory: utils.GetPtr("dir"),
								Task: &models.Task{
									Name:        utils.GetPtr("actions/checkout"),
									Version:     utils.GetPtr("1.2.3"),
									VersionType: models.TagVersion,
									Inputs: &[]models.Parameter{
										{
											Name:  utils.GetPtr("key"),
											Value: "value",
										},
									},
								},
								Type: models.TaskStepType,
							},
						},
						Runner: &models.Runner{
							OS:            utils.GetPtr("linux"),
							Arch:          utils.GetPtr("amd64"),
							Labels:        &[]string{"tag-1"},
							SelfHosted:    utils.GetPtr(false),
							FileReference: testutils.CreateFileReference(5, 10, 9, 17),
						},
						Dependencies: []*models.JobDependency{{JobID: utils.GetPtr("job-1")}, {JobID: utils.GetPtr("job-2")}},
						TokenPermissions: &models.TokenPermissions{
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
					{
						ID:              utils.GetPtr("jobid-1"),
						Name:            utils.GetPtr("jobid-1"),
						ContinueOnError: utils.GetPtr(false),
						TimeoutMS:       utils.GetPtr(21600000),
					},
				},
				Defaults: &models.Defaults{
					TokenPermissions: &models.TokenPermissions{
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
					EnvironmentVariables: &models.EnvironmentVariablesRef{
						EnvironmentVariables: map[string]any{
							"key1": "value1",
							"key2": "value2",
						},
						FileReference: &models.FileReference{
							StartRef: &models.FileLocation{
								Line:   1,
								Column: 2,
							},
							EndRef: &models.FileLocation{
								Line:   3,
								Column: 4,
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := GitHubParser{}

			pipeline, err := parser.Parse(testCase.workflow)
			assert.NoError(t, err)

			changelog, err := diff.Diff(testCase.expectedPipeline, pipeline)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}

}

func TestParseWorkflowDefaults(t *testing.T) {
	testCases := []struct {
		name             string
		workflow         *githubModels.Workflow
		expectedDefaults *models.Defaults
		ShouldError      bool
	}{
		{
			name: "No permissions and no env",
			workflow: &githubModels.Workflow{
				Permissions: nil,
				Env:         nil,
			},
			expectedDefaults: nil,
		},
		{
			name: "Permissions and no env",
			workflow: &githubModels.Workflow{
				Permissions: &githubModels.PermissionsEvent{
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
			},
			expectedDefaults: &models.Defaults{
				TokenPermissions: &models.TokenPermissions{
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
		},
		{
			name: "Env and no permissions",
			workflow: &githubModels.Workflow{
				Env: &githubModels.EnvironmentVariablesRef{
					EnvironmentVariables: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
					FileReference: &models.FileReference{
						StartRef: &models.FileLocation{
							Line:   1,
							Column: 2,
						},
						EndRef: &models.FileLocation{
							Line:   3,
							Column: 4,
						},
					},
				},
			},
			expectedDefaults: &models.Defaults{
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
					FileReference: &models.FileReference{
						StartRef: &models.FileLocation{
							Line:   1,
							Column: 2,
						},
						EndRef: &models.FileLocation{
							Line:   3,
							Column: 4,
						},
					},
				},
			},
		},
		{
			name: "Permissions and env",
			workflow: &githubModels.Workflow{
				Permissions: &githubModels.PermissionsEvent{
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
				Env: &githubModels.EnvironmentVariablesRef{
					EnvironmentVariables: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
					FileReference: &models.FileReference{
						StartRef: &models.FileLocation{
							Line:   1,
							Column: 2,
						},
						EndRef: &models.FileLocation{
							Line:   3,
							Column: 4,
						},
					},
				},
			},
			expectedDefaults: &models.Defaults{
				TokenPermissions: &models.TokenPermissions{
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
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
					FileReference: &models.FileReference{
						StartRef: &models.FileLocation{
							Line:   1,
							Column: 2,
						},
						EndRef: &models.FileLocation{
							Line:   3,
							Column: 4,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseWorkflowDefaults(testCase.workflow)

			if testCase.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				changelog, err := diff.Diff(testCase.expectedDefaults, got)
				assert.NoError(t, err)
				assert.Len(t, changelog, 0)
			}
		})
	}
}

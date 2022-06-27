package github

import (
	"strings"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	commonModels "github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		name             string
		filename         string
		expectedError    error
		expectedWorkflow *models.Workflow
	}{
		{
			name:     "All triggers",
			filename: "../../../test/fixtures/github/all-triggers.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "all-triggers",
				On: &models.On{
					Schedule: &models.Schedule{
						Crons: &[]models.Cron{
							{
								Cron: "30 2 * * *",
							},
						},
						FileReference: testutils.CreateFileReference(3, 3, 4, 23),
					},
					Push: &models.Ref{
						Branches:      []string{"master"},
						FileReference: testutils.CreateFileReference(5, 3, 7, 15),
					},
					PullRequest: &models.Ref{
						PathsIgnore:   []string{"*/test/*"},
						FileReference: testutils.CreateFileReference(8, 3, 10, 17),
					},
					PullRequestTarget: &models.Ref{
						Paths:         []string{"*/test/*"},
						FileReference: testutils.CreateFileReference(11, 3, 13, 17),
					},
					WorkflowDispatch: &models.WorkflowDispatch{
						Inputs: models.Inputs{
							"workflow-input": struct {
								Description string
								Default     interface{}
								Required    bool
								Type        string
								Options     []string
							}{
								Description: "The workflow input",
								Default:     "default-value",
								Required:    true,
							},
						},
						FileReference: testutils.CreateFileReference(14, 3, 19, 23),
					},
					WorkflowCall: &models.WorkflowCall{
						Inputs: models.Inputs{
							"workflow-input": struct {
								Description string
								Default     interface{}
								Required    bool
								Type        string
								Options     []string
							}{
								Description: "The workflow input",
								Default:     "default-value",
								Required:    true,
							},
						},
						FileReference: testutils.CreateFileReference(20, 3, 25, 23),
					},
					WorkflowRun: &models.WorkflowRun{
						Ref: models.Ref{
							BranchesIgnore: []string{"master"},
						},
						FileReference: testutils.CreateFileReference(26, 3, 28, 15),
					},
					Events: models.Events{
						"label": models.Event{
							Types:         []string{"created"},
							FileReference: testutils.CreateFileReference(29, 3, 30, 20),
						},
					},
					FileReference: testutils.CreateFileReference(2, 3, 30, 20),
				},
			},
		},
		{
			name:     "concurrent-jobs",
			filename: "../../../test/fixtures/github/concurrent-jobs.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "concurrent-jobs",
				Jobs: &models.Jobs{
					CIJobs: map[string]*models.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Concurrency: &models.Concurrency{
								Group: utils.GetPtr("ci"),
							},
							FileReference: testutils.CreateFileReference(3, 3, 5, 20),
						},
						"job2": {
							ID:   utils.GetPtr("job2"),
							Name: "Job 2",
							Concurrency: &models.Concurrency{
								Group: utils.GetPtr("ci"),
							},
							FileReference: testutils.CreateFileReference(7, 3, 9, 20),
						},
					},
				},
			},
		},
		{
			name:     "dependant-jobs.yaml",
			filename: "../../../test/fixtures/github/dependant-jobs.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "dependable jobs",
				Jobs: &models.Jobs{
					CIJobs: map[string]*models.Job{
						"dependable-job": {
							ID:            utils.GetPtr("dependable-job"),
							Name:          "Dependable Job",
							FileReference: testutils.CreateFileReference(4, 3, 5, 25),
						},
						"dependant-job": {
							ID:            utils.GetPtr("dependant-job"),
							Name:          "Dependant Job",
							Needs:         &models.Needs{"dependable-job"},
							FileReference: testutils.CreateFileReference(7, 3, 9, 27),
						},
					},
				},
			},
		},
		{
			name:     "environment-variables",
			filename: "../../../test/fixtures/github/environment-variables.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "environment-variables",
				Env: &models.EnvironmentVariablesRef{
					EnvironmentVariables: commonModels.EnvironmentVariables{
						"STRING": "string",
						"NUMBER": 1,
					},
					FileReference: testutils.CreateFileReference(3, 3, 5, 12),
				},
				Jobs: &models.Jobs{
					CIJobs: map[string]*models.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Env: &models.EnvironmentVariablesRef{
								EnvironmentVariables: commonModels.EnvironmentVariables{
									"STRING": "string",
									"NUMBER": 1,
								},
								FileReference: testutils.CreateFileReference(10, 7, 12, 16),
							},
							Steps: &models.Steps{
								{
									Name: "Step 1",
									Run: &models.ShellCommand{
										Script:        "command line",
										FileReference: testutils.CreateFileReference(15, 14, 15, 26),
									},
									Env: &models.EnvironmentVariablesRef{
										EnvironmentVariables: commonModels.EnvironmentVariables{
											"STRING": "string",
											"NUMBER": 1,
										},
										FileReference: testutils.CreateFileReference(16, 11, 18, 20),
									},
									FileReference: testutils.CreateFileReference(14, 9, 18, 20),
								},
							},
							FileReference: testutils.CreateFileReference(8, 3, 18, 20),
						},
					},
				},
			},
		},
		{
			name:     "runners",
			filename: "../../../test/fixtures/github/runners.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "runners",
				Jobs: &models.Jobs{
					CIJobs: map[string]*models.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							RunsOn: &models.RunsOn{
								OS:            utils.GetPtr("linux"),
								SelfHosted:    false,
								Tags:          []string{"ubuntu-latest"},
								FileReference: testutils.CreateFileReference(6, 14, 6, 27),
							},
							FileReference: testutils.CreateFileReference(4, 3, 6, 27),
						},
						"job2": {
							ID:   utils.GetPtr("job2"),
							Name: "Job 2",
							RunsOn: &models.RunsOn{
								OS:            utils.GetPtr("windows"),
								SelfHosted:    true,
								Tags:          []string{"self-hosted", "windows-latest"},
								FileReference: testutils.CreateFileReference(9, 14, 9, 43),
							},
							FileReference: testutils.CreateFileReference(7, 3, 9, 42),
						},
						"job3": {
							ID:   utils.GetPtr("job3"),
							Name: "Job 3",
							RunsOn: &models.RunsOn{
								Arch:          utils.GetPtr("x64"),
								OS:            utils.GetPtr("linux"),
								SelfHosted:    true,
								Tags:          []string{"self-hosted", "linux", "x64"},
								FileReference: testutils.CreateFileReference(12, 14, 12, 39),
							},
							FileReference: testutils.CreateFileReference(10, 3, 12, 38),
						},
					},
				},
			},
		},
		{
			name:     "steps",
			filename: "../../../test/fixtures/github/steps.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "steps",
				Jobs: &models.Jobs{
					CIJobs: map[string]*models.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Steps: &models.Steps{
								{
									Name:          "task without params",
									Uses:          "actions/checkout@v1",
									FileReference: testutils.CreateFileReference(7, 9, 8, 34),
								},
								{
									Name: "task with params",
									Uses: "actions/checkout@v1",
									With: &models.With{
										Inputs:        map[string]any{"repo": "repository"},
										FileReference: testutils.CreateFileReference(12, 9, 13, 27),
									},

									FileReference: testutils.CreateFileReference(10, 9, 13, 27),
								},
								{
									Name: "task with multiline params",
									Uses: "actions/checkout@v1",
									With: &models.With{
										Inputs:        map[string]any{"repos": "repository1\nrepository2\n"},
										FileReference: testutils.CreateFileReference(17, 9, 20, 18),
									},

									FileReference: testutils.CreateFileReference(15, 9, 20, 18),
								},
								{
									Name:          "task with commit ID version",
									Uses:          "actions/checkout@c44948622e1b6bb0eb0cec5b813c1ac561158e1e",
									FileReference: testutils.CreateFileReference(22, 9, 23, 72),
								},
								{
									Name:          "task with branch version",
									Uses:          "actions/checkout@master",
									FileReference: testutils.CreateFileReference(25, 9, 26, 38),
								},
								{
									Name:          "task with tag version",
									Uses:          "actions/checkout@v1.1.1",
									FileReference: testutils.CreateFileReference(28, 9, 29, 38),
								},
								{
									Name: "shell",
									Run: &models.ShellCommand{
										Script:        "command line",
										FileReference: testutils.CreateFileReference(32, 14, 32, 26),
									},
									FileReference: testutils.CreateFileReference(31, 9, 32, 26),
								},
								{
									Name:  "custom shell",
									Shell: "cmd",
									Run: &models.ShellCommand{
										Script:        "command line",
										FileReference: testutils.CreateFileReference(36, 14, 36, 26),
									},
									FileReference: testutils.CreateFileReference(34, 9, 36, 26),
								},
								{
									Name: "shell with break rows",
									Run: &models.ShellCommand{
										Script:        "echo 1\necho 2\necho 3\n",
										FileReference: testutils.CreateFileReference(39, 14, 42, 14),
									},
									FileReference: testutils.CreateFileReference(38, 9, 42, 14),
								},
							},
							FileReference: testutils.CreateFileReference(4, 3, 42, 14),
						},
					},
				},
			},
		},
		{
			name:     "token-permissions",
			filename: "../../../test/fixtures/github/token-permissions.yaml",
			expectedWorkflow: &models.Workflow{
				Name: "permissions",
				Permissions: &models.PermissionsEvent{
					Actions:       "read",
					Statuses:      "write",
					PullRequests:  "read",
					FileReference: testutils.CreateFileReference(2, 3, 5, 22),
				},
				Jobs: &models.Jobs{
					CIJobs: map[string]*models.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Permissions: &models.PermissionsEvent{
								Actions:            "read",
								Checks:             "read",
								Contents:           "read",
								Deployments:        "read",
								Discussions:        "read",
								IdToken:            "read",
								Issues:             "read",
								Packages:           "read",
								PullRequests:       "read",
								Statuses:           "read",
								Pages:              "read",
								RepositoryProjects: "read",
								SecurityEvents:     "read",
							},
							FileReference: testutils.CreateFileReference(8, 3, 10, 26),
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			loader := &GitHubLoader{}
			workflow, err := loader.Load(testutils.ReadFile(testCase.filename))

			if err != testCase.expectedError {
				t.Errorf("Expected error: %v, got: %v", testCase.expectedError, err)
			}

			changelog, _ := diff.Diff(workflow, testCase.expectedWorkflow)

			if len(changelog) > 0 {
				t.Errorf("Loader result is not as expected:")
				for _, change := range changelog {
					t.Errorf("field: %s, got: %v, expected: %v", strings.Join(change.Path, "."), change.From, change.To)
				}
			}
		})
	}
}

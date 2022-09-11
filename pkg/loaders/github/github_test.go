package github

import (
	"strings"
	"testing"

	commonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	pipelineModels "github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
)

func TestLoad(t *testing.T) {
	testCases := []struct {
		name             string
		filename         string
		expectedError    error
		expectedWorkflow *githubModels.Workflow
	}{
		{
			name:     "All triggers",
			filename: "../../../test/fixtures/github/all-triggers.yaml",
			expectedWorkflow: &githubModels.Workflow{
				Name: "all-triggers",
				On: &githubModels.On{
					Schedule: &githubModels.Schedule{
						Crons: &[]githubModels.Cron{
							{
								Cron: "30 2 * * *",
							},
						},
						FileReference: testutils.CreateFileReference(3, 3, 4, 23),
					},
					Push: &githubModels.Ref{
						Branches:      []string{"master"},
						FileReference: testutils.CreateFileReference(5, 3, 7, 15),
					},
					PullRequest: &githubModels.Ref{
						PathsIgnore:   []string{"*/test/*"},
						FileReference: testutils.CreateFileReference(8, 3, 10, 17),
					},
					PullRequestTarget: &githubModels.Ref{
						Paths:         []string{"*/test/*"},
						FileReference: testutils.CreateFileReference(11, 3, 13, 17),
					},
					WorkflowDispatch: &githubModels.WorkflowDispatch{
						Inputs: githubModels.Inputs{
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
					WorkflowCall: &githubModels.WorkflowCall{
						Inputs: githubModels.Inputs{
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
					WorkflowRun: &githubModels.WorkflowRun{
						Ref: githubModels.Ref{
							BranchesIgnore: []string{"master"},
						},
						FileReference: testutils.CreateFileReference(26, 3, 28, 15),
					},
					Events: githubModels.Events{
						"label": githubModels.Event{
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
			expectedWorkflow: &githubModels.Workflow{
				Name: "concurrent-jobs",
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Concurrency: &githubModels.Concurrency{
								Group: utils.GetPtr("ci"),
							},
							FileReference: testutils.CreateFileReference(3, 3, 5, 20),
						},
						"job2": {
							ID:   utils.GetPtr("job2"),
							Name: "Job 2",
							Concurrency: &githubModels.Concurrency{
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
			expectedWorkflow: &githubModels.Workflow{
				Name: "dependable jobs",
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"dependable-job": {
							ID:            utils.GetPtr("dependable-job"),
							Name:          "Dependable Job",
							FileReference: testutils.CreateFileReference(4, 3, 5, 25),
						},
						"dependant-job": {
							ID:            utils.GetPtr("dependant-job"),
							Name:          "Dependant Job",
							Needs:         &githubModels.Needs{"dependable-job"},
							FileReference: testutils.CreateFileReference(7, 3, 9, 27),
						},
					},
				},
			},
		},
		{
			name:     "environment-variables",
			filename: "../../../test/fixtures/github/environment-variables.yaml",
			expectedWorkflow: &githubModels.Workflow{
				Name: "environment-variables",
				Env: &githubModels.EnvironmentVariablesRef{
					EnvironmentVariables: pipelineModels.EnvironmentVariables{
						"STRING": "string",
						"NUMBER": 1,
					},
					FileReference: testutils.CreateFileReference(3, 3, 5, 12),
				},
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Env: &githubModels.EnvironmentVariablesRef{
								EnvironmentVariables: pipelineModels.EnvironmentVariables{
									"STRING": "string",
									"NUMBER": 1,
								},
								FileReference: testutils.CreateFileReference(10, 7, 12, 16),
							},
							Steps: &githubModels.Steps{
								{
									Name: "Step 1",
									Run: &githubModels.ShellCommand{
										Script:        "command line",
										FileReference: testutils.CreateFileReference(15, 14, 15, 26),
									},
									Env: &githubModels.EnvironmentVariablesRef{
										EnvironmentVariables: pipelineModels.EnvironmentVariables{
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
			expectedWorkflow: &githubModels.Workflow{
				Name: "runners",
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							RunsOn: &githubModels.RunsOn{
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
							RunsOn: &githubModels.RunsOn{
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
							RunsOn: &githubModels.RunsOn{
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
			expectedWorkflow: &githubModels.Workflow{
				Name: "steps",
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Steps: &githubModels.Steps{
								{
									Name:          "task without params",
									Uses:          "actions/checkout@v1",
									FileReference: testutils.CreateFileReference(7, 9, 8, 34),
								},
								{
									Name: "task with params",
									Uses: "actions/checkout@v1",
									With: &githubModels.With{
										Values: []*commonModels.MapEntry{
											{
												Key:           "repo",
												Value:         "repository",
												FileReference: testutils.CreateFileReference(13, 11, 13, 21), // End column is supposed to be 27
											},
										},
										FileReference: testutils.CreateFileReference(12, 9, 13, 21), // End column is supposed to be 27
									},
									FileReference: testutils.CreateFileReference(10, 9, 13, 21), // End column is supposed to be 27
								},
								{
									Name: "task with multiline params",
									Uses: "actions/checkout@v1",
									With: &githubModels.With{
										Values: []*commonModels.MapEntry{
											{
												Key:           "repos",
												Value:         "repository1\nrepository2\n",
												FileReference: testutils.CreateFileReference(18, 11, 20, 11), // End column is supposed to be 24
											},
											{
												Key:           "input",
												Value:         "value",
												FileReference: testutils.CreateFileReference(21, 11, 21, 16), // End column is supposed to be 23
											},
										},
										FileReference: testutils.CreateFileReference(17, 9, 21, 16), // End column is supposed to be 23
									},

									FileReference: testutils.CreateFileReference(15, 9, 21, 16), // End column is supposed to be 23
								},
								{
									Name:          "task with commit ID version",
									Uses:          "actions/checkout@c44948622e1b6bb0eb0cec5b813c1ac561158e1e",
									FileReference: testutils.CreateFileReference(23, 9, 24, 72),
								},
								{
									Name:          "task with branch version",
									Uses:          "actions/checkout@master",
									FileReference: testutils.CreateFileReference(26, 9, 27, 38),
								},
								{
									Name:          "task with tag version",
									Uses:          "actions/checkout@v1.1.1",
									FileReference: testutils.CreateFileReference(29, 9, 30, 38),
								},
								{
									Name: "shell",
									Run: &githubModels.ShellCommand{
										Script:        "command line",
										FileReference: testutils.CreateFileReference(33, 14, 33, 26),
									},
									FileReference: testutils.CreateFileReference(32, 9, 33, 26),
								},
								{
									Name:  "custom shell",
									Shell: "cmd",
									Run: &githubModels.ShellCommand{
										Script:        "command line",
										FileReference: testutils.CreateFileReference(37, 14, 37, 26),
									},
									FileReference: testutils.CreateFileReference(35, 9, 37, 26),
								},
								{
									Name: "shell with break rows",
									Run: &githubModels.ShellCommand{
										Script:        "echo 1\necho 2\necho 3\n",
										FileReference: testutils.CreateFileReference(40, 14, 43, 14),
									},
									FileReference: testutils.CreateFileReference(39, 9, 43, 14),
								},
							},
							FileReference: testutils.CreateFileReference(4, 3, 43, 14),
						},
					},
				},
			},
		},
		{
			name:     "token-permissions",
			filename: "../../../test/fixtures/github/token-permissions.yaml",
			expectedWorkflow: &githubModels.Workflow{
				Name: "permissions",
				Permissions: &githubModels.PermissionsEvent{
					Actions:       "read",
					Statuses:      "write",
					PullRequests:  "read",
					FileReference: testutils.CreateFileReference(2, 3, 5, 22),
				},
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"job1": {
							ID:   utils.GetPtr("job1"),
							Name: "Job 1",
							Permissions: &githubModels.PermissionsEvent{
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
		{
			name:     "matrix",
			filename: "../../../test/fixtures/github/matrix.yaml",
			expectedWorkflow: &githubModels.Workflow{
				Name: "matrix",
				Jobs: &githubModels.Jobs{
					CIJobs: map[string]*githubModels.Job{
						"matrix-job": {
							ID: utils.GetPtr("matrix-job"),
							Strategy: &githubModels.Strategy{
								Matrix: &githubModels.Matrix{
									Values: map[string][]any{
										"artifact": {"docker/image", "docker/tar", "go", "java", "node", "php", "python/tar", "python/wheel", "ruby/gemspec"},
										"os":       {"ubuntu-latest", "macos-latest", "windows-latest"},
									},
									Include: []map[string]any{
										{
											"os":       "ubuntu-latest",
											"artifact": "docker/image",
										},
									},
									Exclude: []map[string]any{
										{
											"os":       "ubuntu-latest",
											"artifact": "docker/tar",
										},
									},
									FileReference: testutils.CreateFileReference(6, 9, 25, 33),
								},
							},
							Steps: &githubModels.Steps{
								{
									Name:          "task without params",
									Uses:          "actions/checkout@v1",
									FileReference: testutils.CreateFileReference(28, 9, 29, 34),
								},
								{
									Name: "task with params",
									Uses: "actions/checkout@v1",
									With: &githubModels.With{
										Values: []*commonModels.MapEntry{
											{
												Key:           "repo",
												Value:         "${{ matrix.artifact }}",
												FileReference: testutils.CreateFileReference(34, 11, 34, 33), // End column is supposed to be 27
											},
										},
										FileReference: testutils.CreateFileReference(33, 9, 34, 33), // End column is supposed to be 27
									},
									FileReference: testutils.CreateFileReference(31, 9, 34, 33), // End column is supposed to be 27
								},
							},
							FileReference: testutils.CreateFileReference(4, 3, 34, 33),
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

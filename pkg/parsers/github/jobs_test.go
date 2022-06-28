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

func TestParseWorkflowJobs(t *testing.T) {
	testCases := []struct {
		name         string
		workflow     *githubModels.Workflow
		expectedJobs []*models.Job
	}{
		{
			name:         "Empty workflow",
			workflow:     &githubModels.Workflow{},
			expectedJobs: nil,
		},
		{
			name: "Workflow with one job",
			workflow: &githubModels.Workflow{
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
									With: &githubModels.With{
										Inputs:        map[string]any{"key": "value"},
										FileReference: testutils.CreateFileReference(111, 222, 333, 444),
									},
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
					},
				},
			},
			expectedJobs: []*models.Job{
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
										Name:          utils.GetPtr("key"),
										Value:         "value",
										FileReference: testutils.CreateFileReference(112, 224, 112, 234),
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
			},
		},
		{
			name: "Workflow with multiple jobs",
			workflow: &githubModels.Workflow{
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
									With: &githubModels.With{
										Inputs:        map[string]any{"key": "value"},
										FileReference: testutils.CreateFileReference(111, 222, 333, 444),
									},
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
			},
			expectedJobs: []*models.Job{
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
										Name:          utils.GetPtr("key"),
										Value:         "value",
										FileReference: testutils.CreateFileReference(112, 224, 112, 234),
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
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseWorkflowJobs(testCase.workflow)
			assert.NoError(t, err, testCase.name)

			changelog, err := diff.Diff(testCase.expectedJobs, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

func TestParseJob(t *testing.T) {
	testCases := []struct {
		name        string
		job         *githubModels.Job
		expectedJob *models.Job
	}{
		{
			name: "Empty job",
			job:  &githubModels.Job{},
			expectedJob: &models.Job{
				ContinueOnError: utils.GetPtr(false),
				TimeoutMS:       &defaultTimeoutMS,
			},
		},
		{
			name: "Job with name",
			job: &githubModels.Job{
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
						With: &githubModels.With{
							Inputs:        map[string]any{"key": "value"},
							FileReference: testutils.CreateFileReference(111, 222, 333, 444),
						},
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
			expectedJob: &models.Job{
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
									Name:          utils.GetPtr("key"),
									Value:         "value",
									FileReference: testutils.CreateFileReference(112, 224, 112, 234),
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
		},
		{
			name: "Job without name",
			job: &githubModels.Job{
				ID: utils.GetPtr("jobid-1"),
			},
			expectedJob: &models.Job{
				ID:              utils.GetPtr("jobid-1"),
				Name:            utils.GetPtr("jobid-1"),
				ContinueOnError: utils.GetPtr(false),
				TimeoutMS:       utils.GetPtr(21600000),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseJob("", testCase.job)

			assert.NoError(t, err)

			changelog, err := diff.Diff(testCase.expectedJob, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)

		})
	}
}

func TestParseDependencies(t *testing.T) {
	testCases := []struct {
		name                 string
		needs                *githubModels.Needs
		expectedDependencies []*models.JobDependency
	}{
		{
			name:                 "Empty needs",
			needs:                &githubModels.Needs{},
			expectedDependencies: []*models.JobDependency{},
		},
		{
			name: "Needs with one dependency",
			needs: &githubModels.Needs{
				"job-1",
			},
			expectedDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
			},
		},
		{
			name: "Needs with some dependencies",
			needs: &githubModels.Needs{
				"job-1",
				"job-2",
				"job-3",
			},
			expectedDependencies: []*models.JobDependency{
				{
					JobID: utils.GetPtr("job-1"),
				},
				{
					JobID: utils.GetPtr("job-2"),
				},
				{
					JobID: utils.GetPtr("job-3"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseDependencies(testCase.needs)

			changelog, err := diff.Diff(testCase.expectedDependencies, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

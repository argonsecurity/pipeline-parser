package github

import (
	"reflect"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersCommonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_parseWorkflowJobs(t *testing.T) {
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
							ContinueOnError: utils.GetPtr("true"),
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
									ContinueOnError:  utils.GetPtr("true"),
									If:               "condition",
									TimeoutMinutes:   1,
									WorkingDirectory: "dir",
									Uses:             "actions/checkout@1.2.3",
									With: &githubModels.With{
										Values: []*loadersCommonModels.MapEntry{
											{
												Key:           "key",
												Value:         "value",
												FileReference: testutils.CreateFileReference(112, 224, 112, 234),
											},
										},
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
					ContinueOnError: utils.GetPtr("true"),
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
								Inputs: []*models.Parameter{
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
							ContinueOnError: utils.GetPtr("true"),
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
									ContinueOnError:  utils.GetPtr("true"),
									If:               "condition",
									TimeoutMinutes:   1,
									WorkingDirectory: "dir",
									Uses:             "actions/checkout@1.2.3",
									With: &githubModels.With{
										Values: []*loadersCommonModels.MapEntry{
											{
												Key:           "key",
												Value:         "value",
												FileReference: testutils.CreateFileReference(112, 224, 112, 234),
											},
										},
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
							ID: utils.GetPtr("jobid-2"),
						},
					},
				},
			},
			expectedJobs: []*models.Job{
				{
					ID:              utils.GetPtr("jobid-1"),
					Name:            utils.GetPtr("job-1"),
					ContinueOnError: utils.GetPtr("true"),
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
								Inputs: []*models.Parameter{
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
					ID:        utils.GetPtr("jobid-2"),
					Name:      utils.GetPtr("jobid-2"),
					TimeoutMS: utils.GetPtr(21600000),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseWorkflowJobs(testCase.workflow)
			assert.NoError(t, err, testCase.name)

			testutils.SortJobs(got)
			testutils.SortJobs(testCase.expectedJobs)

			testutils.DeepCompare(t, testCase.expectedJobs, got)
		})
	}
}

func Test_parseCIJob(t *testing.T) {
	testCases := []struct {
		name        string
		job         *githubModels.Job
		expectedJob *models.Job
	}{
		{
			name: "Empty job",
			job:  &githubModels.Job{},
			expectedJob: &models.Job{
				TimeoutMS: &defaultTimeoutMS,
			},
		},
		{
			name: "Job with name",
			job: &githubModels.Job{
				ID:              utils.GetPtr("jobid-1"),
				Name:            "job-1",
				ContinueOnError: utils.GetPtr("true"),
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
						ContinueOnError:  utils.GetPtr("true"),
						If:               "condition",
						TimeoutMinutes:   1,
						WorkingDirectory: "dir",
						Uses:             "actions/checkout@1.2.3",
						With: &githubModels.With{
							Values: []*loadersCommonModels.MapEntry{
								{
									Key:           "key",
									Value:         "value",
									FileReference: testutils.CreateFileReference(112, 224, 112, 234),
								},
							},
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
				ContinueOnError: utils.GetPtr("true"),
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
							Inputs: []*models.Parameter{
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
				ID:        utils.GetPtr("jobid-1"),
				Name:      utils.GetPtr("jobid-1"),
				TimeoutMS: utils.GetPtr(21600000),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseCIJob("", testCase.job)

			assert.NoError(t, err)

			testutils.DeepCompare(t, testCase.expectedJob, got)
		})
	}
}

func Test_parseReusableWorkflowCallJob(t *testing.T) {
	type args struct {
		jobName string
		job     *githubModels.ReusableWorkflowCallJob
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Job
		wantErr bool
	}{
		{
			name: "Job with all fields",
			args: args{
				jobName: "job-1",
				job: &githubModels.ReusableWorkflowCallJob{
					ID:            utils.GetPtr("jobid-1"),
					Name:          "job-1",
					FileReference: testutils.CreateFileReference(4, 33, 5, 36),
					If:            "condition",
					Needs:         &githubModels.Needs{"job-1", "job-2"},
					Uses:          "org/repo/.github/workflows/workflow.yml@main",
					With: map[string]any{
						"key1": "value1",
						"key2": 2,
					},
					Secrets: "inherit",
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
			want: &models.Job{
				ID:            utils.GetPtr("jobid-1"),
				Name:          utils.GetPtr("job-1"),
				FileReference: testutils.CreateFileReference(4, 33, 5, 36),
				Conditions:    []*models.Condition{{Statement: "condition"}},
				Dependencies:  []*models.JobDependency{{JobID: utils.GetPtr("job-1")}, {JobID: utils.GetPtr("job-2")}},
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
				Import: &models.Import{
					Source: &models.ImportSource{
						SCM:          consts.GitHubPlatform,
						Organization: utils.GetPtr("org"),
						Repository:   utils.GetPtr("repo"),
						Path:         utils.GetPtr(".github/workflows/workflow.yml"),
						Type:         models.SourceTypeRemote,
					},
					Version:     utils.GetPtr("main"),
					VersionType: models.BranchVersion,
					Parameters: map[string]any{
						"key1": "value1",
						"key2": 2,
					},
					Secrets: &models.SecretsRef{
						Inherit: true,
					},
				},
			},
		},
		{
			name: "Job without name",
			args: args{
				job: &githubModels.ReusableWorkflowCallJob{
					ID: utils.GetPtr("jobid-1"),
				},
			},
			want: &models.Job{
				ID:   utils.GetPtr("jobid-1"),
				Name: utils.GetPtr("jobid-1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReusableWorkflowCallJob(tt.args.jobName, tt.args.job)

			assert.NoError(t, err)

			testutils.DeepCompare(t, tt.want, got)
		})
	}
}

func Test_parseDependencies(t *testing.T) {
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

			testutils.DeepCompare(t, testCase.expectedDependencies, got)
		})
	}
}

func Test_convertMatrixMap(t *testing.T) {
	type args struct {
		matrix map[string][]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "empty matrix",
			args: args{
				matrix: map[string][]any{},
			},
			want: map[string]any{},
		},
		{
			name: "matrix with one key",
			args: args{
				matrix: map[string][]any{
					"key": {"value"},
				},
			},
			want: map[string]any{
				"key": []any{"value"},
			},
		},
		{
			name: "matrix with list",
			args: args{
				matrix: map[string][]any{
					"key1": {"value1", "value2"},
					"key2": {"value2"},
				},
			},
			want: map[string]any{
				"key1": []any{"value1", "value2"},
				"key2": []any{"value2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertMatrixMap(tt.args.matrix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertMatrixMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMatrix(t *testing.T) {
	type args struct {
		matrix *githubModels.Matrix
	}
	tests := []struct {
		name string
		args args
		want *models.Matrix
	}{
		{
			name: "empty matrix",
			args: args{
				matrix: nil,
			},
			want: nil,
		},
		{
			name: "matrix with one key",
			args: args{
				matrix: &githubModels.Matrix{
					Values: map[string][]any{
						"key": {"value"},
					},
				},
			},
			want: &models.Matrix{
				Matrix: map[string]any{
					"key": []any{"value"},
				},
			},
		},
		{
			name: "matrix with include",
			args: args{
				matrix: &githubModels.Matrix{
					Include: []map[string]any{
						{
							"key1": "value1",
							"key2": "value2",
						},
					},
					Values: map[string][]any{
						"key1": {"value1"},
						"key2": {"value2"},
					},
				},
			},
			want: &models.Matrix{
				Matrix: map[string]any{
					"key1": []any{"value1"},
					"key2": []any{"value2"},
				},
				Include: []map[string]any{
					{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
		},
		{
			name: "matrix with exclude",
			args: args{
				matrix: &githubModels.Matrix{
					Exclude: []map[string]any{
						{
							"key1": "value1",
							"key2": "value2",
						},
					},
					Values: map[string][]any{
						"key1": {"value1"},
						"key2": {"value2"},
					},
				},
			},
			want: &models.Matrix{
				Matrix: map[string]any{
					"key1": []any{"value1"},
					"key2": []any{"value2"},
				},
				Exclude: []map[string]any{
					{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseMatrix(tt.args.matrix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseJobUses(t *testing.T) {
	type args struct {
		uses string
	}
	tests := []struct {
		name            string
		args            args
		wantOrg         string
		wantRepo        string
		wantPath        string
		wantVersion     string
		wantVersionType models.VersionType
		wantSourceType  models.SourceType
	}{
		{
			name: "empty uses",
			args: args{
				uses: "",
			},
			wantOrg:         "",
			wantRepo:        "",
			wantPath:        "",
			wantVersion:     "",
			wantVersionType: "",
			wantSourceType:  "",
		},
		{
			name: "uses remote workflow",
			args: args{
				uses: "organization/repository/.github/workflows/test.yml@version",
			},
			wantOrg:         "organization",
			wantRepo:        "repository",
			wantPath:        ".github/workflows/test.yml",
			wantVersion:     "version",
			wantVersionType: models.BranchVersion,
			wantSourceType:  models.SourceTypeRemote,
		},
		{
			name: "uses local workflow",
			args: args{
				uses: "./.github/workflows/ci.yml",
			},
			wantOrg:         "",
			wantRepo:        "",
			wantPath:        "./.github/workflows/ci.yml",
			wantVersion:     "",
			wantVersionType: models.None,
			wantSourceType:  models.SourceTypeLocal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrg, gotRepo, gotPath, gotVersion, gotVersionType, gotSourceType := parseJobUses(tt.args.uses)
			if gotOrg != tt.wantOrg {
				t.Errorf("parseJobUses() gotOrg = %v, want %v", gotOrg, tt.wantOrg)
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("parseJobUses() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseJobUses() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("parseJobUses() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
			if !reflect.DeepEqual(gotVersionType, tt.wantVersionType) {
				t.Errorf("parseJobUses() gotVersionType = %v, want %v", gotVersionType, tt.wantVersionType)
			}
			if !reflect.DeepEqual(gotSourceType, tt.wantSourceType) {
				t.Errorf("parseJobUses() gotSourceType = %v, want %v", gotSourceType, tt.wantSourceType)
			}
		})
	}
}

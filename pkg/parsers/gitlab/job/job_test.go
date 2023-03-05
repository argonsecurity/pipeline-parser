package job

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/job"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseJob(t *testing.T) {
	testCases := []struct {
		name        string
		jobID       string
		job         *gitlabModels.Job
		expectedJob *models.Job
	}{
		{
			name:  "Job is empty",
			jobID: "1",
			job:   &gitlabModels.Job{},
			expectedJob: &models.Job{
				ID:   utils.GetPtr("1"),
				Name: utils.GetPtr("1"),
			},
		},
		{
			name:  "Job with data",
			jobID: "1",
			job: &gitlabModels.Job{
				AllowFailure: &job.AllowFailure{
					Enabled: utils.GetPtr(true),
				},
				Stage:        "stage",
				Tags:         []string{"1", "2"},
				Image:        &common.Image{Name: "image:tag"},
				BeforeScript: &common.Script{Commands: []string{"before"}},
				AfterScript:  &common.Script{Commands: []string{"after"}},
				Script:       &common.Script{Commands: []string{"script"}},
				Variables: &common.EnvironmentVariablesRef{
					Variables: &common.Variables{
						"key": "value",
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedJob: &models.Job{
				ID:               utils.GetPtr("1"),
				Name:             utils.GetPtr("1"),
				ContinueOnError:  utils.GetPtr("true"),
				ConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("stage")),
				Tags:             []string{"1", "2"},
				PreSteps: []*models.Step{
					{Type: models.ShellStepType,
						Shell: &models.Shell{
							Script: utils.GetPtr("before"),
						},
					},
				},
				PostSteps: []*models.Step{
					{Type: models.ShellStepType,
						Shell: &models.Shell{
							Script: utils.GetPtr("after"),
						},
					},
				},
				Steps: []*models.Step{
					{Type: models.ShellStepType,
						Shell: &models.Shell{
							Script: utils.GetPtr("script"),
						},
					},
				},
				EnvironmentVariables: &models.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
				},
				Runner: &models.Runner{
					DockerMetadata: &models.DockerMetadata{
						Image: utils.GetPtr("image"),
						Label: utils.GetPtr("tag"),
					},
				},

				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseJob(testCase.jobID, testCase.job)

			assert.NoError(t, err)

			testutils.DeepCompare(t, testCase.expectedJob, got)

		})
	}
}

func TestGetJobContinueOnError(t *testing.T) {
	testCases := []struct {
		name             string
		job              *gitlabModels.Job
		expectedContinue *string
	}{
		{
			name:             "Job is empty",
			job:              &gitlabModels.Job{},
			expectedContinue: nil,
		},
		{
			name: "Job with AllowFailure is enabled",
			job: &gitlabModels.Job{
				AllowFailure: &job.AllowFailure{
					Enabled: utils.GetPtr(true),
				},
			},
			expectedContinue: utils.GetPtr("true"),
		},
		{
			name: "Job with AllowFailure is disabled",
			job: &gitlabModels.Job{
				AllowFailure: &job.AllowFailure{
					Enabled: utils.GetPtr(false),
				},
			},
			expectedContinue: utils.GetPtr("false"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getJobContinueOnError(testCase.job)

			testutils.DeepCompare(t, testCase.expectedContinue, got)
		})
	}
}

func TestGetJobConcurrencyGroup(t *testing.T) {
	testCases := []struct {
		name                     string
		job                      *gitlabModels.Job
		expectedConcurrencyGroup *models.ConcurrencyGroup
	}{
		{
			name:                     "Job is empty",
			job:                      &gitlabModels.Job{},
			expectedConcurrencyGroup: nil,
		},
		{
			name: "Job with stage",
			job: &gitlabModels.Job{
				Stage: "stage",
			},
			expectedConcurrencyGroup: utils.GetPtr(models.ConcurrencyGroup("stage")),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getJobConcurrencyGroup(testCase.job)

			testutils.DeepCompare(t, testCase.expectedConcurrencyGroup, got)
		})
	}
}

func TestGetJobConditions(t *testing.T) {
	testCases := []struct {
		name               string
		job                *gitlabModels.Job
		expectedConditions []*models.Condition
	}{
		{
			name:               "Job is empty",
			job:                &gitlabModels.Job{},
			expectedConditions: nil,
		},
		{
			name: "Job with rule",
			job: &gitlabModels.Job{
				Rules: &common.Rules{
					RulesList: []*common.Rule{
						{
							If:      "condition",
							When:    "never",
							Changes: []string{"a", "b", "c"},
							Exists:  []string{"d", "e", "f"},
							Variables: &common.EnvironmentVariablesRef{
								Variables: &common.Variables{
									"string": "string",
									"number": 1,
									"bool":   true,
								},
							},
						},
					},
				},
			},
			expectedConditions: []*models.Condition{
				{
					Statement: "condition",
					Allow:     utils.GetPtr(false),
					Paths: &models.Filter{
						DenyList: []string{"a", "b", "c"},
					},
					Exists: &models.Filter{
						DenyList: []string{"d", "e", "f"},
					},
					Variables: map[string]string{
						"string": "string",
						"number": "1",
						"bool":   "true",
					},
				},
			},
		},
		{
			name: "Job with except",
			job: &gitlabModels.Job{
				Rules: &common.Rules{
					RulesList: []*common.Rule{
						{
							If:      "condition",
							When:    "never",
							Changes: []string{"a", "b", "c"},
							Exists:  []string{"d", "e", "f"},
							Variables: &common.EnvironmentVariablesRef{
								Variables: &common.Variables{
									"string": "string",
									"number": 1,
									"bool":   true,
								},
							},
						},
					},
				},
				Except: &job.Controls{
					Refs: []string{
						"master",
						"test",
					},
					Changes: []string{
						"/path/to/file",
					},
				},
			},
			expectedConditions: []*models.Condition{
				{
					Statement: "condition",
					Allow:     utils.GetPtr(false),
					Paths: &models.Filter{
						DenyList: []string{"a", "b", "c"},
					},
					Exists: &models.Filter{
						DenyList: []string{"d", "e", "f"},
					},
					Variables: map[string]string{
						"string": "string",
						"number": "1",
						"bool":   "true",
					},
				},
				{
					Allow: utils.GetPtr(false),
					Branches: &models.Filter{
						DenyList: []string{
							"master",
							"test",
						},
					},
					Paths: &models.Filter{
						DenyList: []string{
							"/path/to/file",
						},
					},
					Events:    []models.EventType{},
					Variables: nil,
				},
			},
		},
		{
			name: "Job with except and only",
			job: &gitlabModels.Job{
				Rules: &common.Rules{
					RulesList: []*common.Rule{
						{
							If:      "condition",
							When:    "never",
							Changes: []string{"a", "b", "c"},
							Exists:  []string{"d", "e", "f"},
							Variables: &common.EnvironmentVariablesRef{
								Variables: &common.Variables{
									"string": "string",
									"number": 1,
									"bool":   true,
								},
							},
						},
					},
				},
				Except: &job.Controls{
					Refs: []string{
						"master",
						"test",
					},
					Changes: []string{
						"/path/to/file",
					},
				},
				Only: &job.Controls{
					Refs: []string{
						"master",
						"test",
					},
					Changes: []string{
						"/path/to/file",
					},
				},
			},
			expectedConditions: []*models.Condition{
				{
					Statement: "condition",
					Allow:     utils.GetPtr(false),
					Paths: &models.Filter{
						DenyList: []string{"a", "b", "c"},
					},
					Exists: &models.Filter{
						DenyList: []string{"d", "e", "f"},
					},
					Variables: map[string]string{
						"string": "string",
						"number": "1",
						"bool":   "true",
					},
				},
				{
					Allow: utils.GetPtr(false),
					Branches: &models.Filter{
						DenyList: []string{
							"master",
							"test",
						},
					},
					Paths: &models.Filter{
						DenyList: []string{
							"/path/to/file",
						},
					},
					Events:    []models.EventType{},
					Variables: nil,
				},
				{
					Allow: utils.GetPtr(true),
					Branches: &models.Filter{
						AllowList: []string{
							"master",
							"test",
						},
					},
					Paths: &models.Filter{
						AllowList: []string{
							"/path/to/file",
						},
					},
					Events:    []models.EventType{},
					Variables: nil,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getJobConditions(testCase.job)

			testutils.DeepCompare(t, testCase.expectedConditions, got)
		})
	}
}

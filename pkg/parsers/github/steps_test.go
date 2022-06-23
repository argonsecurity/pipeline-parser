package github

import (
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseJobSteps(t *testing.T) {
	testCases := []struct {
		name          string
		steps         *githubModels.Steps
		expectedSteps []*models.Step
	}{
		{
			name:          "nil steps",
			steps:         nil,
			expectedSteps: nil,
		},
		{
			name: "steps with one step",
			steps: &githubModels.Steps{
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
					Run: &githubModels.ShellCommand{
						Script:        "script",
						FileReference: testutils.CreateFileReference(111, 222, 333, 444),
					},
					Shell: "ubuntu",
				},
			},
			expectedSteps: []*models.Step{
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
					Shell: &models.Shell{
						Script:        utils.GetPtr("script"),
						Type:          utils.GetPtr("ubuntu"),
						FileReference: testutils.CreateFileReference(111, 222, 333, 444),
					},
					Type: models.ShellStepType,
				},
			},
		},
		{
			name: "steps with some steps",
			steps: &githubModels.Steps{
				{
					Id:   "1",
					Name: "step-name1",
					Env: &githubModels.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"key": "value1",
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
					ContinueOnError:  utils.GetPtr(true),
					If:               "condition1",
					TimeoutMinutes:   1,
					WorkingDirectory: "dir",
					Run: &githubModels.ShellCommand{
						Script:        "script",
						FileReference: testutils.CreateFileReference(111, 222, 333, 444),
					},
					Shell: "ubuntu",
				},
				{
					Id:   "2",
					Name: "step-name2",
					Env: &githubModels.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"key": "value2",
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
					ContinueOnError:  utils.GetPtr(true),
					If:               "condition2",
					TimeoutMinutes:   1,
					WorkingDirectory: "dir",
					Uses:             "actions/checkout@1.2.3",
					With:             map[string]any{"key": "value"},
				},
			},
			expectedSteps: []*models.Step{
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
						Script:        utils.GetPtr("script"),
						Type:          utils.GetPtr("ubuntu"),
						FileReference: testutils.CreateFileReference(111, 222, 333, 444),
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
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseJobSteps(testCase.steps)
			assert.ElementsMatch(t, testCase.expectedSteps, got, testCase.name)
		})
	}
}

func TestParseJobStep(t *testing.T) {
	testCases := []struct {
		name         string
		step         githubModels.Step
		expectedStep *models.Step
	}{
		{
			name: "Empty step",
			step: githubModels.Step{},
			expectedStep: &models.Step{
				Name: utils.GetPtr(""),
			},
		},
		{
			name: "Shell step",
			step: githubModels.Step{
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
				Run: &githubModels.ShellCommand{
					Script:        "script",
					FileReference: testutils.CreateFileReference(111, 222, 333, 444),
				},
				Shell: "ubuntu",
			},
			expectedStep: &models.Step{
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
				Shell: &models.Shell{
					Script:        utils.GetPtr("script"),
					Type:          utils.GetPtr("ubuntu"),
					FileReference: testutils.CreateFileReference(111, 222, 333, 444),
				},
				Type: models.ShellStepType,
			},
		},
		{
			name: "Task step",
			step: githubModels.Step{
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
			expectedStep: &models.Step{
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseJobStep(testCase.step)
			assert.Equal(t, testCase.expectedStep, got, testCase.name)
		})
	}
}

func TestParseActionHeader(t *testing.T) {
	testCases := []struct {
		name                string
		header              string
		expectedActionName  string
		expectedVersion     string
		expectedVersionType models.VersionType
	}{
		{
			name:                "Header doesn't fit regex",
			header:              "action",
			expectedActionName:  "action",
			expectedVersion:     "",
			expectedVersionType: models.None,
		},
		{
			name:                "Header with no version",
			header:              "actions/checkout",
			expectedActionName:  "actions/checkout",
			expectedVersion:     "",
			expectedVersionType: models.None,
		},
		{
			name:                "Header semver version",
			header:              "actions/checkout@1.2.3",
			expectedActionName:  "actions/checkout",
			expectedVersion:     "1.2.3",
			expectedVersionType: models.TagVersion,
		},
		{
			name:                "Header semver version",
			header:              "actions/checkout@1e204e9a9253d643386038d443f96446fa156a97",
			expectedActionName:  "actions/checkout",
			expectedVersion:     "1e204e9a9253d643386038d443f96446fa156a97",
			expectedVersionType: models.CommitSHA,
		},
		{
			name:                "Header semver version",
			header:              "actions/checkout@branch-name",
			expectedActionName:  "actions/checkout",
			expectedVersion:     "branch-name",
			expectedVersionType: models.BranchVersion,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actionName, version, versionType := parseActionHeader(testCase.header)
			assert.Equal(t, testCase.expectedActionName, actionName, testCase.name)
			assert.Equal(t, testCase.expectedVersion, version, testCase.name)
			assert.Equal(t, testCase.expectedVersionType, versionType, testCase.name)
		})
	}
}

func TestParseActionInput(t *testing.T) {
	testCases := []struct {
		name              string
		with              map[string]any
		expectedParameter *[]models.Parameter
	}{
		{
			name:              "with nil",
			with:              nil,
			expectedParameter: nil,
		},
		{
			name: "with values",
			with: map[string]any{
				"string": "string",
				"int":    1,
				"bool":   true,
			},
			expectedParameter: &[]models.Parameter{
				{
					Name:  utils.GetPtr("string"),
					Value: "string",
				},
				{
					Name:  utils.GetPtr("int"),
					Value: 1,
				},
				{
					Name:  utils.GetPtr("bool"),
					Value: true,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseActionInput(testCase.with)
			assert.Equal(t, testCase.expectedParameter, got, testCase.name)
		})
	}
}

func TestParseActionInputItem(t *testing.T) {
	testCases := []struct {
		name              string
		k                 string
		val               any
		expectedParameter models.Parameter
	}{
		{
			name: "string value",
			k:    "string",
			val:  "value",
			expectedParameter: models.Parameter{
				Name:  utils.GetPtr("string"),
				Value: "value",
			},
		},
		{
			name: "int value",
			k:    "int",
			val:  1,
			expectedParameter: models.Parameter{
				Name:  utils.GetPtr("int"),
				Value: 1,
			},
		},
		{
			name: "boolean value",
			k:    "boolean",
			val:  true,
			expectedParameter: models.Parameter{
				Name:  utils.GetPtr("boolean"),
				Value: true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseActionInputItem(testCase.k, testCase.val)
			assert.Equal(t, testCase.expectedParameter, got, testCase.name)
		})
	}
}

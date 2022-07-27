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

func TestParseSteps(t *testing.T) {
	testCases := []struct {
		name          string
		steps         *azureModels.Steps
		expectedSteps []*models.Step
	}{
		{
			name:          "nil steps",
			steps:         nil,
			expectedSteps: nil,
		},
		{
			name: "steps with one step",
			steps: &azureModels.Steps{
				{
					Name:        "1",
					DisplayName: "step-name",
					Env: &azureModels.EnvironmentVariablesRef{
						EnvironmentVariables: models.EnvironmentVariables{
							"key": "value",
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
					ContinueOnError:  utils.GetPtr(true),
					Condition:        "condition",
					TimeoutInMinutes: 1,
					WorkingDirectory: "dir",
					Bash:             "script",
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
						Script: utils.GetPtr("script"),
						Type:   utils.GetPtr("bash"),
					},
					Type: models.ShellStepType,
				},
			},
		},
		{
			name: "steps with some steps",
			steps: &azureModels.Steps{
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseSteps(testCase.steps)
			assert.ElementsMatch(t, testCase.expectedSteps, got, testCase.name)
		})
	}
}

func TestParseStep(t *testing.T) {
	testCases := []struct {
		name         string
		step         azureModels.Step
		expectedStep *models.Step
	}{
		{
			name: "Empty step",
			step: azureModels.Step{},
			expectedStep: &models.Step{
				Name: utils.GetPtr(""),
			},
		},
		{
			name: "Script step",
			step: azureModels.Step{
				Name:        "1",
				DisplayName: "step-name",
				Env: &azureModels.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
				ContinueOnError:  utils.GetPtr(true),
				Condition:        "condition",
				TimeoutInMinutes: 1,
				WorkingDirectory: "dir",
				Script:           "script",
				Enabled:          utils.GetPtr(true),
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
					Script: utils.GetPtr("script"),
					Type:   utils.GetPtr(""),
				},
				Disabled: utils.GetPtr(false),
				Type:     models.ShellStepType,
			},
		},
		{
			name: "Bash step",
			step: azureModels.Step{
				Name:        "1",
				DisplayName: "step-name",
				Env: &azureModels.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
				ContinueOnError:  utils.GetPtr(true),
				Condition:        "condition",
				TimeoutInMinutes: 1,
				WorkingDirectory: "dir",
				Bash:             "script",
				Enabled:          utils.GetPtr(true),
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
					Script: utils.GetPtr("script"),
					Type:   utils.GetPtr("bash"),
				},
				Disabled: utils.GetPtr(false),
				Type:     models.ShellStepType,
			},
		},
		{
			name: "Powershell step",
			step: azureModels.Step{
				Name:        "1",
				DisplayName: "step-name",
				Env: &azureModels.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
				ContinueOnError:  utils.GetPtr(true),
				Condition:        "condition",
				TimeoutInMinutes: 1,
				WorkingDirectory: "dir",
				Powershell:       "script",
				Enabled:          utils.GetPtr(true),
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
					Script: utils.GetPtr("script"),
					Type:   utils.GetPtr("powershell"),
				},
				Disabled: utils.GetPtr(false),
				Type:     models.ShellStepType,
			},
		},
		{
			name: "Script step",
			step: azureModels.Step{
				Name:        "1",
				DisplayName: "step-name",
				Env: &azureModels.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
				ContinueOnError:  utils.GetPtr(true),
				Condition:        "condition",
				TimeoutInMinutes: 1,
				WorkingDirectory: "dir",
				Pwsh:             "script",
				Enabled:          utils.GetPtr(true),
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
					Script: utils.GetPtr("script"),
					Type:   utils.GetPtr("powershell core"),
				},
				Disabled: utils.GetPtr(false),
				Type:     models.ShellStepType,
			},
		},
		{
			name: "Task step",
			step: azureModels.Step{
				Name:        "1",
				DisplayName: "step-name",
				Env: &azureModels.EnvironmentVariablesRef{
					EnvironmentVariables: models.EnvironmentVariables{
						"key": "value",
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				FileReference:    testutils.CreateFileReference(11, 21, 31, 41),
				ContinueOnError:  utils.GetPtr(true),
				Condition:        "condition",
				TimeoutInMinutes: 1,
				WorkingDirectory: "dir",
				Task:             "Task@2",
				Inputs: &azureModels.TaskInputs{
					Inputs:        map[string]any{"key": "value"},
					FileReference: testutils.CreateFileReference(111, 222, 333, 444),
				},
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseStep(testCase.step)

			changelog, err := diff.Diff(testCase.expectedStep, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseStepScript(t *testing.T) {
	testCases := []struct {
		name          string
		step          azureModels.Step
		expectedShell *models.Shell
	}{
		{
			name:          "Empty step",
			step:          azureModels.Step{},
			expectedShell: nil,
		},
		{
			name: "Script step",
			step: azureModels.Step{
				Script: "script",
			},
			expectedShell: &models.Shell{
				Script: utils.GetPtr("script"),
				Type:   utils.GetPtr(""),
			},
		},
		{
			name: "Bash step",
			step: azureModels.Step{
				Bash: "script",
			},
			expectedShell: &models.Shell{
				Script: utils.GetPtr("script"),
				Type:   utils.GetPtr("bash"),
			},
		},
		{
			name: "Powershell step",
			step: azureModels.Step{
				Powershell: "script",
			},
			expectedShell: &models.Shell{
				Script: utils.GetPtr("script"),
				Type:   utils.GetPtr("powershell"),
			},
		},
		{
			name: "Pwsh step",
			step: azureModels.Step{
				Pwsh: "script",
			},
			expectedShell: &models.Shell{
				Script: utils.GetPtr("script"),
				Type:   utils.GetPtr("powershell core"),
			},
		},
		{
			name: "Task step",
			step: azureModels.Step{
				Task: "Task",
			},
			expectedShell: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseStepScript(testCase.step)

			changelog, err := diff.Diff(testCase.expectedShell, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
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
			name:                "Header doesn't have @",
			header:              "action",
			expectedActionName:  "action",
			expectedVersion:     "",
			expectedVersionType: models.None,
		},
		{
			name:                "Header has @ with only major",
			header:              "Task@1",
			expectedActionName:  "Task",
			expectedVersion:     "1",
			expectedVersionType: models.TagVersion,
		},
		{
			name:                "Header has @ with semver version",
			header:              "Task@1.2.3",
			expectedActionName:  "Task",
			expectedVersion:     "1.2.3",
			expectedVersionType: models.TagVersion,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actionName, version, versionType := parseTaskHeader(testCase.header)
			assert.Equal(t, testCase.expectedActionName, actionName, testCase.name)
			assert.Equal(t, testCase.expectedVersion, version, testCase.name)
			assert.Equal(t, testCase.expectedVersionType, versionType, testCase.name)
		})
	}
}

func TestParseTaskInput(t *testing.T) {
	testCases := []struct {
		name               string
		taskInputs         *azureModels.TaskInputs
		expectedParameters *[]models.Parameter
	}{
		{
			name:               "Task inputs are nil",
			taskInputs:         nil,
			expectedParameters: nil,
		},
		// {
		// 	name: "Task inputs with values",
		// 	taskInputs: &azureModels.TaskInputs{
		// 		Inputs: map[string]any{
		// 			"string": "string",
		// 			"int":    1,
		// 			"bool":   true,
		// 		},
		// 		FileReference: testutils.CreateFileReference(111, 222, 333, 444),
		// 	},
		// 	expectedParameters: &[]models.Parameter{
		// 		{
		// 			Name:          utils.GetPtr("string"),
		// 			Value:         "string",
		// 			FileReference: testutils.CreateFileReference(112, 224, 112, 238),
		// 		},
		// 		{
		// 			Name:          utils.GetPtr("int"),
		// 			Value:         1,
		// 			FileReference: testutils.CreateFileReference(113, 224, 113, 230),
		// 		},
		// 		{
		// 			Name:          utils.GetPtr("bool"),
		// 			Value:         true,
		// 			FileReference: testutils.CreateFileReference(114, 224, 114, 234),
		// 		},
		// 	},
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseTaskInput(testCase.taskInputs)

			// assert.ElementsMatch(t, testCase.expectedParameters, got, testCase.name)
			changelog, err := diff.Diff(testCase.expectedParameters, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

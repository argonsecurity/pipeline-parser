package common

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestParseScript(t *testing.T) {
	testCases := []struct {
		name          string
		script        *common.Script
		expectedSteps []*models.Step
	}{
		{
			name:          "Script is nil",
			script:        nil,
			expectedSteps: nil,
		},
		{
			name:          "Script is empty",
			script:        &common.Script{},
			expectedSteps: []*models.Step{},
		},
		{
			name: "Script with one command",
			script: &common.Script{
				Commands:      []string{"command1"},
				FileReference: testutils.CreateFileReference(1, 2, 1, 10),
			},
			expectedSteps: []*models.Step{
				{
					Type: models.ShellStepType,
					Shell: &models.Shell{
						Script: utils.GetPtr("command1"),
					},
					FileReference: testutils.CreateFileReference(1, 2, 1, 10),
				},
			},
		},
		{
			name: "Script with some commands",
			script: &common.Script{
				Commands:      []string{"command1", "command2", "echo"},
				FileReference: testutils.CreateFileReference(1, 2, 1, 10),
			},
			expectedSteps: []*models.Step{
				{
					Type: models.ShellStepType,
					Shell: &models.Shell{
						Script: utils.GetPtr("command1"),
					},
					FileReference: testutils.CreateFileReference(2, 4, 2, 18),
				},
				{
					Type: models.ShellStepType,
					Shell: &models.Shell{
						Script: utils.GetPtr("command2"),
					},
					FileReference: testutils.CreateFileReference(3, 4, 3, 18),
				},
				{
					Type: models.ShellStepType,
					Shell: &models.Shell{
						Script: utils.GetPtr("echo"),
					},
					FileReference: testutils.CreateFileReference(4, 4, 4, 14),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseScript(testCase.script)

			testutils.DeepCompare(t, testCase.expectedSteps, got)
		})
	}
}

func TestParseCommandFileReference(t *testing.T) {
	testCases := []struct {
		name                  string
		script                *common.Script
		commandIndex          int
		expectedFileReference *models.FileReference
	}{
		{
			name: "Script with one command",
			script: &common.Script{
				Commands:      []string{"command1"},
				FileReference: testutils.CreateFileReference(1, 2, 1, 10),
			},
			commandIndex:          0,
			expectedFileReference: testutils.CreateFileReference(2, 4, 2, 18),
		},
		{
			name: "Script with some commands",
			script: &common.Script{
				Commands:      []string{"command1", "command2", "command3"},
				FileReference: testutils.CreateFileReference(1, 2, 1, 10),
			},
			commandIndex:          2,
			expectedFileReference: testutils.CreateFileReference(4, 4, 4, 18),
		},
		{
			name: "Script with multiline commands",
			script: &common.Script{
				Commands:      []string{"command1", "command2\ncommand3"},
				FileReference: testutils.CreateFileReference(1, 2, 1, 10),
			},
			commandIndex:          1,
			expectedFileReference: testutils.CreateFileReference(3, 4, 3, 27),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseCommandFileReference(testCase.script, testCase.commandIndex)

			testutils.DeepCompare(t, testCase.expectedFileReference, got)
		})
	}
}

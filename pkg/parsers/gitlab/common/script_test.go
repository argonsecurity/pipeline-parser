package common

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
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
			name: "Script with some command",
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
			actualSteps := ParseScript(testCase.script)

			assert.ElementsMatch(t, testCase.expectedSteps, actualSteps, testCase.name)
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseCommandFileReference(testCase.script, testCase.commandIndex)

			changelog, err := diff.Diff(testCase.expectedFileReference, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

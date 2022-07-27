package common

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseEnvironmentVariables(t *testing.T) {
	testCases := []struct {
		name                 string
		environmentVariables *gitlabModels.EnvironmentVariablesRef
		expectedEnvs         *models.EnvironmentVariablesRef
	}{
		{
			name:                 "EnvironmentVariables is nil",
			environmentVariables: nil,
			expectedEnvs:         nil,
		},
		{
			name: "EnvironmentVariables with data",
			environmentVariables: &gitlabModels.EnvironmentVariablesRef{
				Variables: &gitlabModels.Variables{
					"string": "string",
					"number": 1,
					"bool":   true,
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedEnvs: &models.EnvironmentVariablesRef{
				EnvironmentVariables: models.EnvironmentVariables{
					"string": "string",
					"number": 1,
					"bool":   true,
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseEnvironmentVariables(testCase.environmentVariables)

			changelog, err := diff.Diff(testCase.expectedEnvs, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

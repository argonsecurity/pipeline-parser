package github

import (
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestParseEnvironmentVariablesRef(t *testing.T) {
	testCases := []struct {
		name        string
		envRef      *githubModels.EnvironmentVariablesRef
		expectedEnv *models.EnvironmentVariablesRef
	}{
		{
			name:        "Input is nil",
			envRef:      nil,
			expectedEnv: nil,
		},
		{
			name: "Input is not nil",
			envRef: &githubModels.EnvironmentVariablesRef{
				EnvironmentVariables: map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
				FileReference: &models.FileReference{
					StartRef: &models.FileLocation{
						Line:   1,
						Column: 2,
					},
					EndRef: &models.FileLocation{
						Line:   3,
						Column: 4,
					},
				},
			},
			expectedEnv: &models.EnvironmentVariablesRef{
				EnvironmentVariables: map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
				FileReference: &models.FileReference{
					StartRef: &models.FileLocation{
						Line:   1,
						Column: 2,
					},
					EndRef: &models.FileLocation{
						Line:   3,
						Column: 4,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseEnvironmentVariablesRef(testCase.envRef)

			testutils.DeepCompare(t, testCase.expectedEnv, got)
		})
	}
}

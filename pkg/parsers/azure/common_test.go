package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestParseEnvironmentVariablesRef(t *testing.T) {
	testCases := []struct {
		name        string
		envRef      *azureModels.EnvironmentVariablesRef
		expectedEnv *models.EnvironmentVariablesRef
	}{
		{
			name:        "Input is nil",
			envRef:      nil,
			expectedEnv: nil,
		},
		{
			name: "Input is not nil",
			envRef: &azureModels.EnvironmentVariablesRef{
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

func TestParseExtends(t *testing.T) {
	testCases := []struct {
		name            string
		extends         *azureModels.Extends
		expectedImports []string
	}{
		{
			name:            "Extends is nil",
			extends:         nil,
			expectedImports: nil,
		},
		{
			name: "Extends with data",
			extends: &azureModels.Extends{
				Template: azureModels.Template{
					Template: "template1",
				},
			},
			expectedImports: []string{"template1"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseExtends(testCase.extends)
			testutils.DeepCompare(t, testCase.expectedImports, got)
		})
	}
}

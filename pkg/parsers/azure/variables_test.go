package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestParseVariables(t *testing.T) {
	testCases := []struct {
		name      string
		variables *azureModels.Variables
		expected  *models.EnvironmentVariablesRef
	}{
		{
			name:      "variables are nil",
			variables: nil,
			expected:  nil,
		},
		{
			name:      "variables are empty",
			variables: &azureModels.Variables{},
			expected:  nil,
		},
		{
			name: "variables with data",
			variables: &azureModels.Variables{
				{
					Name:          "var1",
					Value:         "value1",
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					Group:         "group1",
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				{
					Name:          "var2",
					Value:         "value2",
					FileReference: testutils.CreateFileReference(9, 10, 11, 12),
				},
			},
			expected: &models.EnvironmentVariablesRef{
				EnvironmentVariables: models.EnvironmentVariables{
					"var1": "value1",
					"var2": "value2",
				},
				FileReference: testutils.CreateFileReference(1, 2, 11, 12),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseVariables(testCase.variables)

			testutils.DeepCompare(t, testCase.expected, got)
		})
	}
}

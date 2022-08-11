package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseParameters(t *testing.T) {
	testCases := []struct {
		name               string
		parameters         *azureModels.Parameters
		expectedParameters []*models.Parameter
	}{
		{
			name:               "parameters are nil",
			parameters:         nil,
			expectedParameters: nil,
		},
		{
			name:               "parameters are empty",
			parameters:         &azureModels.Parameters{},
			expectedParameters: nil,
		},
		{
			name: "parameters with data",
			parameters: &azureModels.Parameters{
				{
					Name:          "param1",
					DisplayName:   "Param 1",
					Type:          "string",
					Default:       "default1",
					Values:        []string{"value1", "value2"},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					Name:          "param2",
					DisplayName:   "Param 2",
					Type:          "number",
					Default:       2,
					Values:        []string{"1", "2"},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
			expectedParameters: []*models.Parameter{
				{
					Name:          utils.GetPtr("param1"),
					Default:       "default1",
					Options:       []string{"value1", "value2"},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					Name:          utils.GetPtr("param2"),
					Default:       2,
					Options:       []string{"1", "2"},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseParameters(testCase.parameters)
			testutils.DeepCompare(t, testCase.expectedParameters, got)
		})
	}
}

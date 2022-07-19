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

func TestParseParameters(t *testing.T) {
	testCases := []struct {
		name               string
		pipeline           *azureModels.Pipeline
		expectedParameters []*models.Parameter
	}{
		{
			name:               "pipeline is nil",
			pipeline:           nil,
			expectedParameters: nil,
		},
		{
			name:               "pipeline with no parameters",
			pipeline:           &azureModels.Pipeline{},
			expectedParameters: nil,
		},
		{
			name: "pipeline with parameters",
			pipeline: &azureModels.Pipeline{
				Parameters: &azureModels.Parameters{
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
			},
			expectedParameters: []*models.Parameter{
				{
					Name:          utils.GetPtr("param1"),
					Default:       "default1",
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					Name:          utils.GetPtr("param2"),
					Default:       2,
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseParameters(testCase.pipeline)

			changelog, err := diff.Diff(testCase.expectedParameters, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

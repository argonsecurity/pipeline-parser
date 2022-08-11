package utils

import (
	"testing"

	loadersCommonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	testUtils "github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseMapToParameters(t *testing.T) {
	testCases := []struct {
		name               string
		with               loadersCommonModels.Map
		expectedParameters []*models.Parameter
	}{
		{
			name:               "with nil",
			with:               loadersCommonModels.Map{},
			expectedParameters: []*models.Parameter{},
		},
		{
			name: "with values",
			with: loadersCommonModels.Map{
				Values: []*loadersCommonModels.MapEntry{
					{
						Key:           "string",
						Value:         "string",
						FileReference: testUtils.CreateFileReference(112, 224, 112, 238),
					},
					{
						Key:           "int",
						Value:         1,
						FileReference: testUtils.CreateFileReference(113, 224, 113, 230),
					},
					{
						Key:           "bool",
						Value:         true,
						FileReference: testUtils.CreateFileReference(114, 224, 114, 234),
					},
				},
				FileReference: testUtils.CreateFileReference(111, 222, 333, 444),
			},
			expectedParameters: []*models.Parameter{
				{
					Name:          utils.GetPtr("string"),
					Value:         "string",
					FileReference: testUtils.CreateFileReference(112, 224, 112, 238),
				},
				{
					Name:          utils.GetPtr("int"),
					Value:         1,
					FileReference: testUtils.CreateFileReference(113, 224, 113, 230),
				},
				{
					Name:          utils.GetPtr("bool"),
					Value:         true,
					FileReference: testUtils.CreateFileReference(114, 224, 114, 234),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseMapToParameters(testCase.with)
			testUtils.DeepCompare(t, testCase.expectedParameters, got)
		})
	}
}

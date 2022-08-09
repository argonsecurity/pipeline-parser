package utils

import (
	"testing"

	loadersCommonModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/common/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParseMapToParameters(t *testing.T) {
	testCases := []struct {
		name               string
		with               *loadersCommonModels.Map
		expectedParameters []*models.Parameter
	}{
		{
			name:               "with nil",
			with:               nil,
			expectedParameters: nil,
		},
		{
			name: "with values",
			with: &loadersCommonModels.Map{
				Values: []*loadersCommonModels.MapEntry{
					{
						Key:           "string",
						Value:         "string",
						FileReference: testutils.CreateFileReference(112, 224, 112, 238),
					},
					{
						Key:           "int",
						Value:         1,
						FileReference: testutils.CreateFileReference(113, 224, 113, 230),
					},
					{
						Key:           "bool",
						Value:         true,
						FileReference: testutils.CreateFileReference(114, 224, 114, 234),
					},
				},
				FileReference: testutils.CreateFileReference(111, 222, 333, 444),
			},
			expectedParameters: []*models.Parameter{
				{
					Name:          utils.GetPtr("string"),
					Value:         "string",
					FileReference: testutils.CreateFileReference(112, 224, 112, 238),
				},
				{
					Name:          utils.GetPtr("int"),
					Value:         1,
					FileReference: testutils.CreateFileReference(113, 224, 113, 230),
				},
				{
					Name:          utils.GetPtr("bool"),
					Value:         true,
					FileReference: testutils.CreateFileReference(114, 224, 114, 234),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var m loadersCommonModels.Map
			if testCase.with != nil {
				m = loadersCommonModels.Map(*testCase.with)
			}

			got := ParseMapToParameters(m)
			changeLog, err := diff.Diff(testCase.expectedParameters, got)
			assert.NoError(t, err)
			assert.Len(t, changeLog, 0)
		})
	}
}

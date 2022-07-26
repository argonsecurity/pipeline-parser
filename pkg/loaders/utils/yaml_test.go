package utils

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCalcParameterFileReference(t *testing.T) {
	testCases := []struct {
		name                  string
		startLine             int
		startColumn           int
		key                   string
		val                   any
		expectedFileReference *models.FileReference
	}{
		{
			name:                  "Start line is -1",
			startLine:             -1,
			expectedFileReference: nil,
		},
		{
			name:                  "Start column is -1",
			startColumn:           -1,
			expectedFileReference: nil,
		},
		{
			name:                  "Value without \\n",
			startLine:             2,
			startColumn:           4,
			key:                   "key",
			val:                   "value",
			expectedFileReference: testutils.CreateFileReference(2, 4, 2, 14),
		},
		{
			name: "Value with \\n",

			startLine:             2,
			startColumn:           4,
			key:                   "key",
			val:                   "value\nvalue\n",
			expectedFileReference: testutils.CreateFileReference(2, 4, 4, 14),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := CalculateParameterFileReference(testCase.startLine, testCase.startColumn, testCase.key, testCase.val)
			assert.Equal(t, testCase.expectedFileReference, got, testCase.name)
		})
	}
}

package triggers

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestIsPositive(t *testing.T) {
	testCases := []struct {
		name       string
		comparison *Comparison
		expected   bool
	}{
		{
			name:       "Comparison is empty",
			comparison: &Comparison{},
			expected:   false,
		},
		{
			name:       "Comparison is equals",
			comparison: &Comparison{Operator: equals},
			expected:   true,
		},
		// {
		// 	name:       "Comparison is notEquals",
		// 	comparison: &Comparison{Operator: notEquals},
		// 	expected:   false,
		// },
		{
			name:       "Comparison is match",
			comparison: &Comparison{Operator: match},
			expected:   true,
		},
		// {
		// 	name:       "Comparison is notMatch",
		// 	comparison: &Comparison{Operator: notMatch},
		// 	expected:   false,
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.comparison.IsPositive()

			testutils.DeepCompare(t, testCase.expected, got)

		})
	}
}

// func TestGetComparisons(t *testing.T) {
// 	testCases := []struct {
// 		name                string
// 		expression          string
// 		expectedComparisons []*Comparison
// 	}{
// 		{
// 			name:                "Expressions is empty",
// 			expression:          "",
// 			expectedComparisons: nil,
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			got := getComparisons(testCase.expression)

// 			changelog, err := diff.Diff(testCase.expectedComparisons, got)
// 			assert.NoError(t, err)
// 			assert.Len(t, changelog, 0)
// 		})
// 	}
// }

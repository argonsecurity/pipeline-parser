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
		{
			name:       "Comparison is match",
			comparison: &Comparison{Operator: match},
			expected:   true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.comparison.IsPositive()

			testutils.DeepCompare(t, testCase.expected, got)

		})
	}
}

func TestGetComparisons(t *testing.T) {
	testCases := []struct {
		name                string
		expression          string
		expectedComparisons []*Comparison
	}{
		{
			name:                "Expressions is empty",
			expression:          "",
			expectedComparisons: []*Comparison{},
		},
		{
			name:       "Expressions with == operator and parentheses",
			expression: `$var == "value"`,
			expectedComparisons: []*Comparison{
				{
					Variable: "$var",
					Value:    `"value"`,
					Operator: equals,
				},
			},
		},
		{
			name:       "Expressions with == operator and regex",
			expression: `$var == /value/`,
			expectedComparisons: []*Comparison{
				{
					Variable: "$var",
					Value:    `/value/`,
					Operator: equals,
				},
			},
		},
		{
			name:       "Expressions with =~ operator and parentheses",
			expression: `$var =~ "value"`,
			expectedComparisons: []*Comparison{
				{
					Variable: "$var",
					Value:    `"value"`,
					Operator: match,
				},
			},
		},
		{
			name:       "Expressions with =~ operator and regex",
			expression: `$var =~ /value/`,
			expectedComparisons: []*Comparison{
				{
					Variable: "$var",
					Value:    `/value/`,
					Operator: match,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := getComparisons(testCase.expression)

			testutils.DeepCompare(t, testCase.expectedComparisons, got)
		})
	}
}

package testutils

import (
	"testing"

	"github.com/go-test/deep"
)

func DeepCompare(t *testing.T, struct1 any, struct2 any) bool {
	if !validateNil(struct1, struct2) {
		return false
	}

	if diffs := deep.Equal(struct1, struct2); diffs != nil {
		for _, diff := range diffs {
			t.Errorf(diff)
		}

		return false
	}

	return true
}

func validateNil(struct1, struct2 any) bool {
	if struct1 != struct2 && (struct1 == nil || struct2 == nil) {
		return false
	}

	return true
}

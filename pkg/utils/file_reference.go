package utils

import (
	"reflect"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func CompareFileReferences(a, b *models.FileReference) bool {
	if a == nil || b == nil {
		return false
	}
	return reflect.DeepEqual(a, b)
}

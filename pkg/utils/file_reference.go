package utils

import "github.com/argonsecurity/pipeline-parser/pkg/models"

func CompareFileReferences(a, b *models.FileReference) bool {
	if a == nil || b == nil {
		return false
	}
	if a.StartRef.Line == b.StartRef.Line &&
		a.EndRef.Line == b.EndRef.Line &&
		a.StartRef.Column == b.StartRef.Column &&
		a.EndRef.Column == b.EndRef.Column &&
		a.IsAlias == b.IsAlias {
		return true
	}
	return false
}

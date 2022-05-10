package testutils

import "github.com/argonsecurity/pipeline-parser/pkg/models"

func CreateFileReference(l1, c1, l2, c2 int) *models.FileReference {
	return &models.FileReference{
		StartRef: &models.FileLocation{
			Line:   l1,
			Column: c1,
		},
		EndRef: &models.FileLocation{
			Line:   l2,
			Column: c2,
		},
	}
}

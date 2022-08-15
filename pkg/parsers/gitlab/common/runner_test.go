package common

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func TestParseRunner(t *testing.T) {
	testCases := []struct {
		name           string
		image          *gitlabModels.Image
		expectedRunner *models.Runner
	}{
		{
			name:           "Image is nil",
			image:          nil,
			expectedRunner: nil,
		},
		{
			name:  "Image is empty",
			image: &gitlabModels.Image{},
			expectedRunner: &models.Runner{
				DockerMetadata: &models.DockerMetadata{
					Image:       utils.GetPtrOrNil(""),
					Label:       utils.GetPtrOrNil(""),
					RegistryURL: utils.GetPtrOrNil(""),
				},
			},
		},
		{
			name: "Image with data",
			image: &gitlabModels.Image{
				Name:          "registry/namespace/image:tag",
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedRunner: &models.Runner{
				DockerMetadata: &models.DockerMetadata{
					Image:       utils.GetPtr("namespace/image"),
					Label:       utils.GetPtr("tag"),
					RegistryURL: utils.GetPtr("registry"),
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := ParseRunner(testCase.image)

			testutils.DeepCompare(t, testCase.expectedRunner, got)
		})
	}
}

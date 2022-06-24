package common

import (
	"testing"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
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

			changelog, err := diff.Diff(testCase.expectedRunner, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0)
		})
	}
}

func TestParseImageName(t *testing.T) {
	testCases := []struct {
		name              string
		imageName         string
		expectedRegistry  string
		expectedNamespace string
		expectedImageName string
		expectedTag       string
	}{
		{
			name:              "Empty image name",
			imageName:         "",
			expectedRegistry:  "",
			expectedNamespace: "",
			expectedImageName: "",
			expectedTag:       "",
		},
		{
			name:              "Image name with tag",
			imageName:         "image:tag",
			expectedRegistry:  "",
			expectedNamespace: "",
			expectedImageName: "image",
			expectedTag:       "tag",
		},
		{
			name:              "Image name without registry and tag",
			imageName:         "repository/image",
			expectedRegistry:  "",
			expectedNamespace: "repository",
			expectedImageName: "image",
			expectedTag:       "",
		},
		{
			name:              "Image name with registry and namespace without tag",
			imageName:         "registry/namespace/image",
			expectedRegistry:  "registry",
			expectedNamespace: "namespace",
			expectedImageName: "image",
			expectedTag:       "",
		},

		{
			name:              "Image name with tag and namespace",
			imageName:         "namespace/image:tag",
			expectedRegistry:  "",
			expectedNamespace: "namespace",
			expectedImageName: "image",
			expectedTag:       "tag",
		},
		{
			name:              "Image name with tag and registry and namespace",
			imageName:         "registry/namespace/image:tag",
			expectedRegistry:  "registry",
			expectedNamespace: "namespace",
			expectedImageName: "image",
			expectedTag:       "tag",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registry, namespace, imageName, tag := parseImageName(testCase.imageName)
			assert.Equal(t, testCase.expectedRegistry, registry)
			assert.Equal(t, testCase.expectedNamespace, namespace)
			assert.Equal(t, testCase.expectedImageName, imageName)
			assert.Equal(t, testCase.expectedTag, tag)
		})
	}
}

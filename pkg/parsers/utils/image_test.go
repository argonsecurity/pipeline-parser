package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			registry, namespace, imageName, tag := ParseImageName(testCase.imageName)
			assert.Equal(t, testCase.expectedRegistry, registry)
			assert.Equal(t, testCase.expectedNamespace, namespace)
			assert.Equal(t, testCase.expectedImageName, imageName)
			assert.Equal(t, testCase.expectedTag, tag)
		})
	}
}

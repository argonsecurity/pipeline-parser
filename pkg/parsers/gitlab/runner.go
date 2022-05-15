package gitlab

import (
	"strings"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseRunner(image *gitlabModels.Image) *models.Runner {
	registry, namespace, imageName, tag := parseImageName(image.Name)
	if namespace != "" {
		imageName = namespace + "/" + imageName
	}
	return &models.Runner{
		DockerMetadata: &models.DockerMetadata{
			Image:       &imageName,
			Label:       &tag,
			RegistryURL: &registry,
		},
		FileReference: image.FileReference,
	}
}

func parseImageName(imageName string) (string, string, string, string) {
	var registry, namespace, tag string
	image := imageName

	if split := strings.Split(imageName, "/"); len(split) == 3 {
		registry = split[0]
		namespace = split[1]
		image = split[2]
	} else if len(split) == 2 {
		namespace = split[0]
		image = split[1]
	}

	if split := strings.Split(image, ":"); len(split) == 2 {
		image = split[0]
		tag = split[1]
	}

	return registry, namespace, image, tag
}

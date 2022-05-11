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
	var registry, namespace, image, tag string
	var split = strings.Split(imageName, "/")
	if len(split) == 3 {
		registry = split[0]
		namespace = split[1]
		image = split[2]
	} else if len(split) == 2 {
		registry = ""
		namespace = split[0]
		image = split[1]
	} else {
		registry = ""
		namespace = ""
		image = imageName
	}

	split = strings.Split(image, ":")
	if len(split) == 2 {
		image = split[0]
		tag = split[1]
	}
	return registry, namespace, image, tag
}

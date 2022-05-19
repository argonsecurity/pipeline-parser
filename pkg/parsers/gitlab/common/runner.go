package common

import (
	"strings"

	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func ParseRunner(image *gitlabModels.Image) *models.Runner {
	if image == nil {
		return nil
	}
	registry, namespace, imageName, tag := parseImageName(image.Name)
	if namespace != "" {
		imageName = namespace + "/" + imageName
	}

	return &models.Runner{
		DockerMetadata: &models.DockerMetadata{
			Image:       utils.GetPtrOrNil(imageName),
			Label:       utils.GetPtrOrNil(tag),
			RegistryURL: utils.GetPtrOrNil(registry),
		},
		FileReference: image.FileReference,
	}
}

func parseImageName(imageName string) (string, string, string, string) {
	var registry, namespace, tag string
	image := imageName

	if split := strings.Split(imageName, "/"); len(split) == 3 { // imageName contains registry/repository/image
		registry = split[0]
		namespace = split[1]
		image = split[2]
	} else if len(split) == 2 { // imageName contains repository/image
		namespace = split[0]
		image = split[1]
	}

	if split := strings.Split(image, ":"); len(split) == 2 { // image contains image:tag
		image = split[0]
		tag = split[1]
	}

	return registry, namespace, image, tag
}

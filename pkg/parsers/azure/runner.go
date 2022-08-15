package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	parsersUtils "github.com/argonsecurity/pipeline-parser/pkg/parsers/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parseRunner(job *azureModels.BaseJob) *models.Runner {
	if job == nil || job.Pool == nil && job.Container == nil {
		return nil
	}

	runner := &models.Runner{}

	if job.Pool != nil {
		runner = parsePool(job.Pool, runner)
	}

	if job.Container != nil {
		runner = parseContainer(job.Container, runner)
	}
	return runner
}

func parsePool(pool *azureModels.Pool, runner *models.Runner) *models.Runner {
	if runner == nil || pool == nil || pool.VmImage == "" {
		return runner
	}

	return parsersUtils.ParseRunnerTag(pool.VmImage, runner)
}

func parseContainer(container *azureModels.JobContainer, runner *models.Runner) *models.Runner {
	if runner == nil || container == nil || container.Image == "" {
		return runner
	}

	registry, namespace, imageName, tag := parsersUtils.ParseImageName(container.Image)
	if namespace != "" {
		imageName = namespace + "/" + imageName
	}

	runner.DockerMetadata = &models.DockerMetadata{
		Image:       utils.GetPtrOrNil(imageName),
		Label:       utils.GetPtrOrNil(tag),
		RegistryURL: utils.GetPtrOrNil(registry),
	}
	return runner
}

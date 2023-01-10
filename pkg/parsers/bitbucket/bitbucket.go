package bitbucket

import (
	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type BitbucketParser struct{}

func (g *BitbucketParser) Parse(bitbucketPipeline *bitbucketModels.Pipeline) (*models.Pipeline, error) {
	if bitbucketPipeline == nil {
		return nil, nil
	}

	// pipeline := &models.Pipeline{
	// 	Name: &bitbucketPipeline.Name,
	// }

	// pipeline.Defaults = parsePipelineDefaults(azurePipeline)
	// pipeline.Triggers = parsePipelineTriggers(azurePipeline)
	// pipeline.Parameters = parseParameters(azurePipeline.Parameters)
	// pipeline.Imports = parseExtends(azurePipeline.Extends)

	// var jobs []*models.Job

	// if azurePipeline.Stages != nil {
	// 	jobs = append(jobs, parseStages(azurePipeline.Stages)...)
	// }

	// if azurePipeline.Jobs != nil {
	// 	jobs = append(pipeline.Jobs, parseJobs(azurePipeline.Jobs)...)
	// }

	// if len(jobs) == 0 {
	// 	jobs = []*models.Job{generateDefaultJob()}
	// 	if azurePipeline.Pool != nil {
	// 		jobs[0].Runner = parsePool(azurePipeline.Pool, jobs[0].Runner)
	// 	}

	// 	if azurePipeline.Container != nil {
	// 		jobs[0].Runner = parseContainer(azurePipeline.Container, jobs[0].Runner)
	// 	}
	// }

	// pipeline.Jobs = jobs

	// if azurePipeline.Steps != nil {
	// 	pipeline.Jobs[0].Steps = parseSteps(azurePipeline.Steps)
	// }

	// return pipeline, nil
	return nil, nil
}

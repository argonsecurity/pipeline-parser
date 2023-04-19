package azure

import (
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

type AzureEnhancer struct{}

func (a *AzureEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials, organization string) ([]*enhancers.ImportedPipeline, error) {
	importedPipelines, err := getTemplates(data, credentials, organization)
	if err != nil {
		return importedPipelines, err
	}

	return importedPipelines, nil
}

func (a *AzureEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline) (*models.Pipeline, error) {
	if len(importedPipelines) == 0 {
		return data, nil
	}

	for _, importedPipeline := range importedPipelines {
		data = mergePipelines(data, importedPipeline)
	}

	return data, nil
}

func mergePipelines(pipeline *models.Pipeline, importedPipeline *enhancers.ImportedPipeline) *models.Pipeline {
	if pipeline == nil || importedPipeline == nil {
		return pipeline
	}

	if pipeline.Imports != nil || len(pipeline.Imports) > 0 {
		for _, imported := range pipeline.Imports {
			if *imported.Source.Path == importedPipeline.JobName {
				imported.Pipeline = importedPipeline.Pipeline
				return pipeline
			}
		}
	}

	if pipeline.Jobs != nil && len(pipeline.Jobs) > 0 {
		for _, job := range pipeline.Jobs {
			if job != nil {
				if isDirectImport(job, importedPipeline) {
					job.Imports.Pipeline = importedPipeline.Pipeline
					return pipeline
				}

				if isJobVariableImport(job, importedPipeline) {
					job.EnvironmentVariables.Imports.Pipeline = importedPipeline.Pipeline
					return pipeline
				}

				if found := job.Steps != nil && fillStepsImports(job.Steps, importedPipeline); found {
					return pipeline
				}
				if found := job.PreSteps != nil && fillStepsImports(job.PreSteps, importedPipeline); found {
					return pipeline
				}
				if found := job.PostSteps != nil && fillStepsImports(job.PostSteps, importedPipeline); found {
					return pipeline
				}
			}
		}
	}

	if isDefaultsImport(pipeline.Defaults, importedPipeline) {
		for _, imported := range pipeline.Imports {
			if *imported.Source.Path == importedPipeline.JobName {
				imported.Pipeline = importedPipeline.Pipeline
			}
		}
	}
	return pipeline
}

func isDirectImport(job *models.Job, importedPipeline *enhancers.ImportedPipeline) bool {
	return job.Imports != nil &&
		job.Imports.Source != nil &&
		job.Imports.FileReference != nil &&
		importedPipeline.OriginFileReference != nil &&
		*job.Imports.Source.Path == importedPipeline.JobName &&
		utils.CompareFileReferences(job.Imports.FileReference, importedPipeline.OriginFileReference)
}

func isJobVariableImport(job *models.Job, importedPipeline *enhancers.ImportedPipeline) bool {
	return job.EnvironmentVariables != nil &&
		job.EnvironmentVariables.Imports != nil &&
		job.EnvironmentVariables.Imports.Source != nil &&
		job.EnvironmentVariables.Imports.FileReference != nil &&
		*job.EnvironmentVariables.Imports.Source.Path == importedPipeline.JobName &&
		utils.CompareFileReferences(job.EnvironmentVariables.Imports.FileReference, importedPipeline.OriginFileReference)
}

func isDirectStepImport(step *models.Step, importedPipeline *enhancers.ImportedPipeline) bool {
	return step.Imports != nil &&
		step.Imports.Source != nil &&
		step.Imports.FileReference != nil &&
		importedPipeline.OriginFileReference != nil &&
		*step.Imports.Source.Path == importedPipeline.JobName &&
		utils.CompareFileReferences(step.Imports.FileReference, importedPipeline.OriginFileReference)
}

func isStepVariableImport(step *models.Step, importedPipeline *enhancers.ImportedPipeline) bool {
	return step.EnvironmentVariables != nil &&
		step.EnvironmentVariables.Imports != nil &&
		step.EnvironmentVariables.Imports.Source != nil &&
		step.EnvironmentVariables.Imports.FileReference != nil &&
		*step.EnvironmentVariables.Imports.Source.Path == importedPipeline.JobName &&
		utils.CompareFileReferences(step.EnvironmentVariables.Imports.FileReference, importedPipeline.OriginFileReference)
}

func isDefaultsImport(defaults *models.Defaults, importedPipeline *enhancers.ImportedPipeline) bool {
	return defaults != nil &&
		defaults.EnvironmentVariables != nil &&
		defaults.EnvironmentVariables.Imports != nil &&
		defaults.EnvironmentVariables.Imports.Source != nil &&
		defaults.EnvironmentVariables.Imports.FileReference != nil &&
		*defaults.EnvironmentVariables.Imports.Source.Path == importedPipeline.JobName &&
		utils.CompareFileReferences(defaults.EnvironmentVariables.Imports.FileReference, importedPipeline.OriginFileReference)
}

func fillStepsImports(steps []*models.Step, importedPipeline *enhancers.ImportedPipeline) bool {
	if len(steps) == 0 {
		return false
	}
	for _, step := range steps {
		if isDirectStepImport(step, importedPipeline) {
			step.Imports.Pipeline = importedPipeline.Pipeline
			return true
		}

		if isStepVariableImport(step, importedPipeline) {
			step.EnvironmentVariables.Imports.Pipeline = importedPipeline.Pipeline
			return true
		}
	}
	return false
}

package bitbucket

import (
	"fmt"

	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parseJobs(pipeline *bitbucketModels.Pipeline) []*models.Job {
	if pipeline == nil {
		return nil
	}

	var jobs []*models.Job

	if pipeline.Pipelines != nil {
		if pipeline.Pipelines.Default != nil {
			defaultJob := createJob("default")
			defaultJob.Steps = parseStepArray(pipeline.Pipelines.Default, defaultJob)
			jobs = append(jobs, defaultJob)
		}

		if pipeline.Pipelines.PullRequests != nil {
			jobs = append(jobs, parseStepMapToJob(pipeline.Pipelines.PullRequests)...)
		}

		if pipeline.Pipelines.Branches != nil {
			jobs = append(jobs, parseStepMapToJob(pipeline.Pipelines.Branches)...)
		}

		if pipeline.Pipelines.Tags != nil {
			jobs = append(jobs, parseStepMapToJob(pipeline.Pipelines.Tags)...)
		}

		if pipeline.Pipelines.Bookmarks != nil {
			jobs = append(jobs, parseStepMapToJob(pipeline.Pipelines.Bookmarks)...)
		}

		if pipeline.Pipelines.Custom != nil {
			jobs = append(jobs, parseStepMapToJob(pipeline.Pipelines.Custom)...)
		}
	}

	return jobs
}

func parseStepMapToJob(jobMap *bitbucketModels.StepMap) []*models.Job {
	var jobs []*models.Job
	for jobName, steps := range *jobMap {
		job := createJob(jobName)
		job.Steps = parseStepArray(steps, job)
		jobs = append(jobs, job)
	}
	return jobs
}

func parseStepArray(jobSteps []*bitbucketModels.Step, job *models.Job) []*models.Step {
	var steps []*models.Step
	for _, step := range jobSteps {
		steps = append(steps, parseStep(step)...)
	}
	return steps
}

func createJob(jobName string) *models.Job {
	var job models.Job
	id := fmt.Sprintf("job-%s", jobName)
	job.ID = &id
	job.Name = &jobName
	return &job
}

func parseStep(step *bitbucketModels.Step) []*models.Step {
	if step == nil {
		return nil
	}

	var steps []*models.Step
	if step.Step != nil {
		steps = append(steps, parseExecutionUnitToStep(step.Step))
	}

	if step.Parallel != nil {
		for _, parallelStep := range step.Parallel {
			steps = append(steps, parseExecutionUnitToStep(parallelStep.Step))
		}
	}

	return steps
}

func parseExecutionUnitToStep(executionUnitRef *bitbucketModels.ExecutionUnitRef) *models.Step {
	if executionUnitRef == nil {
		return nil
	}

	var step models.Step
	step.Name = executionUnitRef.ExecutionUnit.Name
	step.FileReference = executionUnitRef.FileReference
	if executionUnitRef.ExecutionUnit.MaxTime != nil {
		var timeout int = int(*executionUnitRef.ExecutionUnit.MaxTime)
		step.Timeout = &timeout
	}
	step.Shell = parseScriptToShell(executionUnitRef.ExecutionUnit.Script)
	if step.Shell != nil && step.Shell.Type != nil {
		step.Type = models.StepType(*step.Shell.Type)
	}
	step.AfterScript = parseScriptToShell(executionUnitRef.ExecutionUnit.AfterScript)
	var scripts = executionUnitRef.ExecutionUnit.Script
	if step.Shell != nil { // script env vars
		for _, script := range scripts {
			if script.PipeToExecute != nil {
				step.EnvironmentVariables = parseEnvironmentVariables(script.PipeToExecute.Variables)
			}
		}
	}
	return &step
}

func parseEnvironmentVariables(srcEnvVars *bitbucketModels.EnvironmentVariablesRef) *models.EnvironmentVariablesRef {
	if srcEnvVars == nil {
		return nil
	}
	envVars := models.EnvironmentVariablesRef{
		EnvironmentVariables: make(map[string]any),
	}
	for key, env := range srcEnvVars.EnvironmentVariables {
		envVars.EnvironmentVariables[key] = env
	}
	envVars.FileReference = &models.FileReference{
		StartRef: srcEnvVars.FileReference.StartRef,
		EndRef:   srcEnvVars.FileReference.EndRef,
	}
	return &envVars
}

func parseScriptToShell(scripts []*bitbucketModels.Script) *models.Shell {
	if scripts == nil {
		return nil
	}

	var shell models.Shell
	var scriptString string
	var pipeFileReference *models.FileReference
	for _, script := range scripts {
		if script != nil {
			if script.String != nil {
				scriptString += addScriptLine(*script.String)
			}
			if (script.PipeToExecute) != nil {
				scriptString += addScriptLine(*script.PipeToExecute.Pipe.String)
				if pipeFileReference == nil {
					pipeFileReference = script.PipeToExecute.Pipe.FileReference
					continue
				}
				pipeFileReference.EndRef = script.PipeToExecute.Pipe.FileReference.EndRef
			}
		}
	}

	shell.Script = &scriptString
	if pipeFileReference != nil {
		shell.FileReference = pipeFileReference
		return &shell
	}

	shell.FileReference = &models.FileReference{
		StartRef: scripts[0].FileReference.StartRef,
		EndRef:   scripts[len(scripts)-1].FileReference.EndRef,
	}
	shell.Type = utils.GetPtr(string(models.ShellStepType))
	return &shell
}

func addScriptLine(script string) string {
	return fmt.Sprintf("%s\n", script)
}

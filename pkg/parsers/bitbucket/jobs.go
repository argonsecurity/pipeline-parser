package bitbucket

import (
	"fmt"

	bitbucketModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/bitbucket/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func parseJobs(pipeline *bitbucketModels.Pipeline) []*models.Job {
	if pipeline == nil {
		return nil
	}

	var jobs []*models.Job

	if pipeline.Pipelines.Default != nil {
		defaultJob := createJob("default")
		defaultJob.Steps = parseStepArray(pipeline.Pipelines.Default, defaultJob)
		jobs = append(jobs, defaultJob)
	}

	if pipeline.Pipelines.PullRequests != nil {
		jobs = append(jobs, parseJobMap(pipeline.Pipelines.PullRequests)...)
	}

	if pipeline.Pipelines.Branches != nil {
		jobs = append(jobs, parseJobMap(pipeline.Pipelines.Branches)...)
	}

	if pipeline.Pipelines.Tags != nil {
		jobs = append(jobs, parseJobMap(pipeline.Pipelines.Tags)...)
	}

	if pipeline.Pipelines.Bookmarks != nil {
		jobs = append(jobs, parseJobMap(pipeline.Pipelines.Bookmarks)...)
	}

	if pipeline.Pipelines.Custom != nil {
		jobs = append(jobs, parseJobMap(pipeline.Pipelines.Custom)...)
	}

	return jobs
}

func parseJobMap(jobMap *bitbucketModels.StepMap) []*models.Job {
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
	job.Name = &jobName
	return &job
}

func parseStep(step *bitbucketModels.Step) []*models.Step {
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
	var step models.Step
	step.Name = &executionUnitRef.ExecutionUnit.Name
	step.FileReference = executionUnitRef.FileReference
	if executionUnitRef.ExecutionUnit.MaxTime != nil {
		var timeout int = int(*executionUnitRef.ExecutionUnit.MaxTime)
		step.Timeout = &timeout
	}
	step.Shell = parseScript(executionUnitRef.ExecutionUnit.Script)
	var scripts = executionUnitRef.ExecutionUnit.Script
	if step.Shell != nil { // script env vars
		for _, script := range scripts {
			if script.PipeToExecute != nil {
				step.EnvironmentVariables = &models.EnvironmentVariablesRef{
					EnvironmentVariables: make(map[string]any),
				}
				for key, env := range script.PipeToExecute.Variables.EnvironmentVariables {
					step.EnvironmentVariables.EnvironmentVariables[key] = env
				}
			}
		}
		if step.EnvironmentVariables != nil {
			step.EnvironmentVariables.FileReference = &models.FileReference{
				StartRef: scripts[0].FileReference.StartRef,
				EndRef:   scripts[len(scripts)-1].FileReference.EndRef,
			}
		}
	}
	return &step
}

func parseScript(scripts []bitbucketModels.Script) *models.Shell {
	if scripts == nil {
		return nil
	}

	var shell models.Shell
	var scriptString string
	for _, script := range scripts {
		if script.String != "" {
			scriptString += addScriptLine(script.String)
		}
		if (script.PipeToExecute) != nil {
			scriptString += addScriptLine(script.PipeToExecute.Pipe)
		}
	}
	shell.Script = &scriptString
	shell.FileReference = &models.FileReference{
		StartRef: scripts[0].FileReference.StartRef,
		EndRef:   scripts[len(scripts)-1].FileReference.EndRef,
	}
	return &shell
}

func addScriptLine(script string) string {
	return fmt.Sprintf("- %s \n", script)
}

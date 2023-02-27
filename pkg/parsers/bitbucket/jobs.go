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
			defaultJob := parseJob("default", pipeline.Pipelines.Default)
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
		job := parseJob(jobName, steps)
		jobs = append(jobs, job)
	}
	return jobs
}

func parseJob(jobName string, steps []*bitbucketModels.Step) *models.Job {
	job := createJob(jobName)
	job.Steps = parseStepArray(steps, job)
	job.FileReference = generateJobFileReference(job)
	return job
}

func generateJobFileReference(job *models.Job) *models.FileReference {
	if job.Steps != nil || len(job.Steps) > 0 {
		return &models.FileReference{
			StartRef: job.Steps[0].FileReference.StartRef,
			EndRef:   job.Steps[len(job.Steps)-1].FileReference.EndRef,
		}
	}
	return nil
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
	step.Task = parseScriptToTask(executionUnitRef.ExecutionUnit.Script)
	step.Type = getStepType(&step)
	var scripts = executionUnitRef.ExecutionUnit.Script
	if step.Task != nil { // script env vars belong to tasks
		for _, script := range scripts {
			if script.PipeToExecute != nil {
				step.EnvironmentVariables = parseEnvironmentVariables(step.EnvironmentVariables, script.PipeToExecute.Variables)
			}
		}
	}
	var afterScripts = executionUnitRef.ExecutionUnit.AfterScript
	step.AfterScript = parseScriptToShell(afterScripts)
	if step.AfterScript != nil {
		for _, script := range afterScripts {
			if script.PipeToExecute != nil {
				step.EnvironmentVariables = parseEnvironmentVariables(step.EnvironmentVariables, script.PipeToExecute.Variables)
			}
		}
	}
	return &step
}

<<<<<<< HEAD
func setStepType(step *models.Step) models.StepType {
	if step.Shell != nil && step.Shell.Type != nil {
		return models.StepType(*step.Shell.Type)
=======
func getStepType(step *models.Step) models.StepType {
	if step.Shell != nil {
		return models.ShellStepType
>>>>>>> 03a07db89e12ba3e60363c836346f15fc25933bd
	}
	if step.Task != nil {
		return models.TaskStepType
	}

	return ""
}

func parseEnvironmentVariables(existing *models.EnvironmentVariablesRef, srcEnvVars *bitbucketModels.EnvironmentVariablesRef) *models.EnvironmentVariablesRef {
	if srcEnvVars == nil {
		return existing
	}
	if existing == nil {
		existing = &models.EnvironmentVariablesRef{
			EnvironmentVariables: make(map[string]any),
			FileReference:        srcEnvVars.FileReference,
		}
	}
	for key, env := range srcEnvVars.EnvironmentVariables {
		_, ok := existing.EnvironmentVariables[key]
		if !ok {
			existing.EnvironmentVariables[key] = env
		}
	}
	existing.FileReference = &models.FileReference{
		StartRef: existing.FileReference.StartRef,
		EndRef:   srcEnvVars.FileReference.EndRef,
	}
	return existing
}

func parseScriptToShell(scripts []*bitbucketModels.Script) *models.Shell {
	if scripts == nil {
		return nil
	}

	var shell models.Shell
	var scriptString string
	var fileReference *models.FileReference
	for _, script := range scripts {
		if script != nil {
			if script.String != nil {
				scriptString += addScriptLine(*script.String)
				if fileReference == nil {
					fileReference = script.FileReference
					continue
				}
				fileReference.EndRef = script.FileReference.EndRef
			}
		}
	}

	if scriptString == "" {
		return nil
	}
	shell.Script = &scriptString
	if fileReference != nil {
		shell.FileReference = fileReference
	}
	shell.Type = utils.GetPtr(string(models.ShellStepType))
	return &shell
}

func parseScriptToTask(scripts []*bitbucketModels.Script) *models.Task {
	if scripts == nil {
		return nil
	}

	var task models.Task
	var scriptString string
	for _, script := range scripts {
		if script != nil {
			if (script.PipeToExecute) != nil {
				scriptString += addScriptLine(*script.PipeToExecute.Pipe.String)
			}
		}
	}

	if scriptString == "" {
		return nil
	}
	task.Name = &scriptString
	task.VersionType = models.None // not supported
	return &task
}

func addScriptLine(script string) string {
	return fmt.Sprintf("%s\n", script)
}

package github

import (
	"regexp"
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	defaultTimeoutMS int = 360 * 60 * 1000

	githubWorkflowCallRegex = regexp.MustCompile(`(?P<org>[^/]+)/(?P<repo>[^/]+)/(?P<path>[^@]+)(?:@(?P<version>.+)|$)`)
)

func parseWorkflowJobs(workflow *githubModels.Workflow) ([]*models.Job, error) {
	if workflow == nil || workflow.Jobs == nil {
		return nil, nil
	}
	ciJobs, err := utils.MapToSliceErr(workflow.Jobs.CIJobs, parseCIJob)
	if err != nil {
		return nil, err
	}

	reusableWorkflowCallJobs, err := utils.MapToSliceErr(workflow.Jobs.ReusableWorkflowCallJobs, parseReusableWorkflowCallJob)
	if err != nil {
		return nil, err
	}

	return append(ciJobs, reusableWorkflowCallJobs...), nil
}

func parseCIJob(jobName string, job *githubModels.Job) (*models.Job, error) {
	parsedJob := &models.Job{
		ID:                   job.ID,
		Name:                 &job.Name,
		ContinueOnError:      job.ContinueOnError,
		EnvironmentVariables: parseEnvironmentVariablesRef(job.Env),
		FileReference:        job.FileReference,
	}

	if job.Name == "" {
		parsedJob.Name = job.ID
	}

	if job.TimeoutMinutes != nil && *job.TimeoutMinutes != 0 {
		timeout := int(*job.TimeoutMinutes) * 60 * 1000
		parsedJob.TimeoutMS = &timeout
	} else {
		parsedJob.TimeoutMS = &defaultTimeoutMS
	}

	if job.If != "" {
		parsedJob.Conditions = []*models.Condition{{Statement: job.If}}
	}

	if job.Concurrency != nil {
		parsedJob.ConcurrencyGroup = (*models.ConcurrencyGroup)(job.Concurrency.Group)
	}

	if job.Steps != nil {
		parsedJob.Steps = parseJobSteps(job.Steps)
	}

	if job.RunsOn != nil {
		parsedJob.Runner = parseRunsOnToRunner(job.RunsOn)
	}

	if job.Needs != nil {
		parsedJob.Dependencies = parseDependencies(job.Needs)
	}

	if job.Permissions != nil {
		permissions, err := parseTokenPermissions(job.Permissions)
		if err != nil {
			return nil, err
		}
		parsedJob.TokenPermissions = permissions
	}

	if job.Strategy != nil && job.Strategy.Matrix != nil {
		parsedJob.Matrix = parseMatrix(job.Strategy.Matrix)
	}

	return parsedJob, nil
}

func parseReusableWorkflowCallJob(jobName string, job *githubModels.ReusableWorkflowCallJob) (*models.Job, error) {
	parsedJob := &models.Job{
		ID:            job.ID,
		Name:          &job.Name,
		FileReference: job.FileReference,
	}

	if job.Name == "" {
		parsedJob.Name = job.ID
	}

	if job.If != "" {
		parsedJob.Conditions = []*models.Condition{{Statement: job.If}}
	}

	if job.Needs != nil {
		parsedJob.Dependencies = parseDependencies(job.Needs)
	}

	if job.Permissions != nil {
		permissions, err := parseTokenPermissions(job.Permissions)
		if err != nil {
			return nil, err
		}
		parsedJob.TokenPermissions = permissions
	}

	if job.Strategy != nil && job.Strategy.Matrix != nil {
		parsedJob.Matrix = parseMatrix(job.Strategy.Matrix)
	}

	if job.Uses != "" {
		org, repo, path, version, versionType, sourceType := parseJobUses(job.Uses)
		secretsMap, inherit := parseSecrets(job.Secrets)
		parsedJob.Imports = &models.Import{
			Source: &models.ImportSource{
				SCM:          consts.GitHubPlatform,
				Organization: &org,
				Repository:   &repo,
				Path:         &path,
				Type:         sourceType,
			},
			Version:     &version,
			VersionType: versionType,
			Parameters:  job.With,
			Secrets: &models.SecretsRef{
				Secrets: secretsMap,
				Inherit: inherit,
			},
		}
	}

	return parsedJob, nil
}

func parseMatrix(matrix *githubModels.Matrix) *models.Matrix {
	if matrix == nil {
		return nil
	}

	return &models.Matrix{
		Matrix:        convertMatrixMap(matrix.Values),
		Include:       matrix.Include,
		Exclude:       matrix.Exclude,
		FileReference: matrix.FileReference,
	}
}

func convertMatrixMap(matrix map[string][]any) map[string]any {
	convertedMatrix := map[string]any{}
	for key, value := range matrix {
		convertedMatrix[key] = value
	}
	return convertedMatrix
}

func parseDependencies(needs *githubModels.Needs) []*models.JobDependency {
	return utils.Map(([]string)(*needs), func(dependency string) *models.JobDependency {
		return &models.JobDependency{
			JobID: &dependency,
		}
	})
}

func parseSecrets(secrets any) (map[string]any, bool) {
	if secrets == nil {
		return nil, false
	}

	if secretsString, ok := secrets.(string); ok && secretsString == "inherit" {
		return nil, true
	}

	if secretsMap, ok := secrets.(map[string]any); ok {
		return secretsMap, false
	}

	return nil, false
}

func parseJobUses(uses string) (org string, repo string, path string, version string, versionType models.VersionType, sourceType models.SourceType) {
	if uses == "" {
		return
	}

	if strings.HasPrefix(uses, "./") {
		path = uses
		versionType = models.None
		sourceType = models.SourceTypeLocal
		return
	}

	result := githubWorkflowCallRegex.FindStringSubmatch(uses)
	if len(result) == 0 {
		return
	}

	org = result[1]
	repo = result[2]
	path = result[3]
	sourceType = models.SourceTypeRemote

	if len(result) == 5 {
		version = result[4]
	}

	versionType = detectVersionType(version)
	return
}

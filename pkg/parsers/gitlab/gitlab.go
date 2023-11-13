package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/common"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/job"
	"github.com/argonsecurity/pipeline-parser/pkg/parsers/gitlab/triggers"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

type GitLabParser struct{}

func (g *GitLabParser) Parse(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) (*models.Pipeline, error) {
	var err error
	pipeline := &models.Pipeline{
		Imports: ParseImports(gitlabCIConfiguration.Include),
	}

	pipeline.Defaults = parseDefaults(gitlabCIConfiguration)

	if gitlabCIConfiguration.Workflow != nil {
		pipeline.Triggers, pipeline.Defaults.Conditions = triggers.ParseRules(gitlabCIConfiguration.Workflow.Rules)
	}
	pipeline.Jobs, err = job.ParseJobs(gitlabCIConfiguration)
	if err != nil {
		return nil, err
	}

	if len(gitlabCIConfiguration.Jobs) > 0 {
		_, err := utils.MapToSliceErr(gitlabCIConfiguration.Jobs, func(jobID string, job *gitlabModels.Job) ([]*models.Import, error) {
			return appendJobTriggerIncludes(job, &pipeline.Imports)
		})
		if err != nil {
			return nil, err
		}
	}
	return pipeline, nil
}

func parseScans(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) *models.Scans {
	if gitlabCIConfiguration.Default == nil || gitlabCIConfiguration.Default.Artifacts == nil {
		return nil
	}
	reports := gitlabCIConfiguration.Default.Artifacts.Reports

	return &models.Scans{
		Secrets:      utils.GetPtr(reports.SecretDetection != nil),
		SAST:         utils.GetPtr(reports.Sast != nil),
		Dependencies: utils.GetPtr(reports.DependencyScanning != nil),
		Iac:          utils.GetPtr(reports.Terraform != nil),
		License:      utils.GetPtr(reports.LicenseScanning != nil),
	}
}

func parseDefaults(gitlabCIConfiguration *gitlabModels.GitlabCIConfiguration) *models.Defaults {
	defaults := &models.Defaults{
		EnvironmentVariables: common.ParseEnvironmentVariables(gitlabCIConfiguration.Variables),
		Runner:               common.ParseRunner(gitlabCIConfiguration.Image),
		PostSteps:            common.ParseScript(gitlabCIConfiguration.AfterScript),
		PreSteps:             common.ParseScript(gitlabCIConfiguration.BeforeScript),
		Scans:                parseScans(gitlabCIConfiguration),
	}
	return defaults
}

func appendJobTriggerIncludes(job *gitlabModels.Job, imports *[]*models.Import) ([]*models.Import, error) {
	if job.Trigger != nil && job.Trigger.Include != nil {
		if jobImport := ParseImports(job.Trigger.Include); jobImport != nil {
			if imports == nil {
				*imports = []*models.Import{}
			}
			*imports = append(*imports, jobImport...)
		}
	}
	return *imports, nil
}

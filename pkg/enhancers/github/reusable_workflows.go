package github

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	GITHUB_BASE_URL = "https://raw.githubusercontent.com"
)

func getReusableWorkflows(pipeline *models.Pipeline, credentials *models.Credentials, baseUrl *string) ([]*enhancers.ImportedPipeline, error) {
	var errs error
	importedPipelines := []*enhancers.ImportedPipeline{}
	for _, job := range pipeline.Jobs {
		if job.Imports != nil {
			importedPipelineBuf, err := handleImport(job.Imports, credentials, baseUrl)
			if err != nil {
				if errs == nil {
					errs = errors.New("got error(s) importing pipeline(s):")
				}
				errs = errors.Wrap(errs, fmt.Sprintf("error importing pipeline for job %s: %s", *job.Name, err.Error()))
			}
			importedPipelines = append(importedPipelines, &enhancers.ImportedPipeline{
				JobName: *job.Name,
				Data:    importedPipelineBuf,
			})
		}

	}
	return importedPipelines, errs
}

func handleImport(jobImport *models.Import, credentials *models.Credentials, baseUrl *string) ([]byte, error) {
	if jobImport == nil || jobImport.Source == nil {
		return nil, nil
	}

	if jobImport.Source.Type == models.SourceTypeRemote && jobImport.Source.Organization != nil && jobImport.Source.Repository != nil && jobImport.Source.Path != nil && jobImport.Version != nil {
		return loadRemoteFile(*jobImport.Source.Organization, *jobImport.Source.Repository, *jobImport.Version, *jobImport.Source.Path, credentials, baseUrl)
	}

	if jobImport.Source.Type == models.SourceTypeLocal && jobImport.Source.Path != nil {
		return loadLocalFile(*jobImport.Source.Path)
	}

	return nil, nil
}

func loadRemoteFile(org, repo, version, path string, credentials *models.Credentials, baseUrl *string) ([]byte, error) {
	if org == "" || repo == "" || path == "" {
		return nil, nil
	}

	if version == "" {
		version = "main"
	}

	url := fmt.Sprintf("%s/%s/%s/%s/%s", GITHUB_BASE_URL, org, repo, version, path)
	if baseUrl != nil && *baseUrl != "" {
		url = fmt.Sprintf("%s/raw/%s/%s/%s/%s", *baseUrl, org, repo, version, path)
	}

	client := utils.GetHttpClient(credentials)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, errors.New(resp.Response.Status)
	}

	buf := resp.Bytes()
	return buf, nil
}

func loadLocalFile(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

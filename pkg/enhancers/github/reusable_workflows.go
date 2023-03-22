package github

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/http"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const githubRemoteFileUrl = "https://raw.githubusercontent.com/%s/%s/%s/%s"

func getReusableWorkflows(pipeline *models.Pipeline, credentials *models.Credentials) ([]*enhancers.ImportedPipeline, error) {
	var errs error
	importedPipelines := []*enhancers.ImportedPipeline{}
	for _, job := range pipeline.Jobs {
		if job.Imports != nil {
			importedPipelineBuf, err := handleImport(job.Imports, credentials)
			if err != nil {
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

func handleImport(imports *models.Import, credentials *models.Credentials) ([]byte, error) {
	if imports == nil || imports.Source == nil {
		return nil, nil
	}

	if imports.Source.Type == models.SourceTypeRemote && imports.Source.Organization != nil && imports.Source.Repository != nil && imports.Source.Path != nil && imports.Version != nil {
		return loadRemotePipeline(*imports.Source.Organization, *imports.Source.Repository, *imports.Version, *imports.Source.Path, credentials)
	}

	if imports.Source.Type == models.SourceTypeLocal && imports.Source.Path != nil {
		return loadLocalPipeline(*imports.Source.Path)
	}

	return nil, nil
}

func loadRemotePipeline(org, repo, version, path string, credentials *models.Credentials) ([]byte, error) {
	if org == "" || repo == "" || path == "" {
		return nil, nil
	}

	if version == "" {
		version = "main"
	}

	url := fmt.Sprintf(githubRemoteFileUrl, org, repo, version, path)
	buf, err := getHttpClient(credentials).Get(url, nil, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func loadLocalPipeline(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	buf, err := utils.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func getHttpClient(credentials *models.Credentials) *http.HTTPClient {
	if credentials == nil {
		return http.NewHTTPClient(nil)
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", credentials.Token),
	}

	return http.NewHTTPClient(headers)
}

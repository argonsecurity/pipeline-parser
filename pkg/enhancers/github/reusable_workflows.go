package github

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/imroc/req/v3"
)

var (
	GithubBaseURL = "https://raw.githubusercontent.com"
)

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
		return loadRemoteFile(*imports.Source.Organization, *imports.Source.Repository, *imports.Version, *imports.Source.Path, credentials)
	}

	if imports.Source.Type == models.SourceTypeLocal && imports.Source.Path != nil {
		return loadLocalFile(*imports.Source.Path)
	}

	return nil, nil
}

func loadRemoteFile(org, repo, version, path string, credentials *models.Credentials) ([]byte, error) {
	if org == "" || repo == "" || path == "" {
		return nil, nil
	}

	if version == "" {
		version = "main"
	}

	url := fmt.Sprintf("%s/%s/%s/%s/%s", GithubBaseURL, org, repo, version, path)
	client := getHttpClient(credentials)
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

func getHttpClient(credentials *models.Credentials) *req.Client {
	client := req.C()
	if credentials == nil {
		return client
	}

	return client.SetCommonHeader("Authorization", fmt.Sprintf("token %s", credentials.Token))

}

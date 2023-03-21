package github

import (
	"errors"
	"fmt"

	"github.com/argonsecurity/pipeline-parser/pkg/http"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const githubRemoteFileUrl = "https://raw.githubusercontent.com/%s/%s/%s/%s"

func enhanceReusableWorkflows(pipeline *models.Pipeline, credentials *models.Credentials) (*models.Pipeline, error) {
	var errs []error
	for _, job := range pipeline.Jobs {
		if job.Imports != nil {
			importedPipeline, err := handleImport(job.Imports, credentials)
			if err != nil {
				errs = append(errs, fmt.Errorf("error importing pipeline for job %s: %w", *job.Name, err))
			}
			job.Imports.Pipeline = importedPipeline
		}

	}
	return pipeline, errors.Join(errs...)
}

func handleImport(imports *models.Import, credentials *models.Credentials) (*models.Pipeline, error) {
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

func loadRemotePipeline(org, repo, version, path string, credentials *models.Credentials) (*models.Pipeline, error) {
	if org == "" || repo == "" || path == "" {
		return nil, nil
	}

	if version == "" {
		version = "main"
	}

	url := fmt.Sprintf(githubRemoteFileUrl, org, repo, version, path)
	_, err := getHttpClient(credentials).Get(url, nil, nil)
	if err != nil {
		return nil, err
	}

	return &models.Pipeline{
		Name: utils.GetPtr("remote"),
	}, nil
	// return handler.Handle(buf, consts.GitHubPlatform, credentials)
}

func loadLocalPipeline(path string) (*models.Pipeline, error) {
	if path == "" {
		return nil, nil
	}

	_, err := utils.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &models.Pipeline{
		Name: utils.GetPtr("local"),
	}, nil
	// return handler.Handle(buf, consts.GitHubPlatform, models.Credentials{})
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

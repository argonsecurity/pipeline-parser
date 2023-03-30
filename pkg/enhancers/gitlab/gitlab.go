package gitlab

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	GITLAB_BASE_URL = "https://gitlab.com"
)

type GitLabEnhancer struct{}

func (g *GitLabEnhancer) LoadImportedPipelines(data *models.Pipeline, credentials *models.Credentials) ([]*enhancers.ImportedPipeline, error) {
	var errs error
	importedPipelines := []*enhancers.ImportedPipeline{}
	if data.Imports != nil {
		for _, importData := range data.Imports {
			importedPipeline, err := handleImport(importData, credentials)
			if err != nil {
				if errs == nil {
					errs = errors.New("got error(s) importing pipeline(s):")
				}
				errs = errors.Wrap(errs, fmt.Sprintf("error importing pipeline: %s", err.Error()))
			}

			// We append nil imported pipelines to maintain the order of the imported pipelines
			importedPipelines = append(importedPipelines, importedPipeline)
		}
	}
	return importedPipelines, errs
}

func handleImport(importData *models.Import, credentials *models.Credentials) (*enhancers.ImportedPipeline, error) {
	if importData == nil || importData.Source == nil {
		return nil, nil
	}

	if importData.Source.Type == models.SourceTypeRemote {
		return handleRemoteImport(importData, credentials)
	}

	if importData.Source.Type == models.SourceTypeLocal {
		return handleLocalImport(importData)
	}

	return nil, nil
}

func handleRemoteImport(importData *models.Import, credentials *models.Credentials) (*enhancers.ImportedPipeline, error) {
	if importData.Source.Type != models.SourceTypeRemote {
		return nil, errors.New("invalid source type for remote import")
	}

	if importData.Source.Organization == nil ||
		importData.Source.Repository == nil ||
		importData.Source.Path == nil ||
		importData.Version == nil {
		return nil, errors.New("missing required fields for remote import")
	}

	url := fmt.Sprintf("%s/%s/%s/-/raw/%s/%s",
		GITLAB_BASE_URL,
		*importData.Source.Organization,
		*importData.Source.Repository,
		*importData.Version,
		*importData.Source.Path,
	)
	client := utils.GetHttpClient(credentials)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, errors.New(resp.Response.Status)
	}

	buf := resp.Bytes()
	return &enhancers.ImportedPipeline{Data: buf}, nil
}

func handleLocalImport(importData *models.Import) (*enhancers.ImportedPipeline, error) {
	buf, err := os.ReadFile(strings.TrimPrefix(*importData.Source.Path, "/"))
	if err != nil {
		return nil, err
	}

	return &enhancers.ImportedPipeline{
		Data: buf,
	}, nil
}

func (g *GitLabEnhancer) Enhance(data *models.Pipeline, importedPipelines []*enhancers.ImportedPipeline) (*models.Pipeline, error) {
	for i, importData := range data.Imports {
		importedPipeline := importedPipelines[i]
		if importedPipeline != nil {
			importData.Pipeline = importedPipeline.Pipeline
		}
	}

	return data, nil
}

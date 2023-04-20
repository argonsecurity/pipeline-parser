package azure

import (
	"fmt"
	"os"
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/pkg/errors"
)

var (
	AZURE_BASE_URL = "https://dev.azure.com"
	ITEMS_API      = "/{ORGANIZATION}/{PROJECT}/_apis/git/repositories/{REPOSITORY}/items?path={PATH}"
	VERSION_QUERY  = "&versionDescriptor.versionType=tag&version="
)

func getTemplates(pipeline *models.Pipeline, credentials *models.Credentials, organization *string) ([]*enhancers.ImportedPipeline, error) {
	var errs error
	importedPipelines := []*enhancers.ImportedPipeline{}
	var resources *models.Resources
	if pipeline.Defaults != nil && pipeline.Defaults.Resources != nil {
		resources = pipeline.Defaults.Resources
	}

	// import from default variables
	if pipeline.Defaults != nil &&
		pipeline.Defaults.EnvironmentVariables != nil &&
		pipeline.Defaults.EnvironmentVariables.Imports != nil {
		importedPipelines, errs = appendImports(importedPipelines, pipeline.Defaults.EnvironmentVariables.Imports, resources, credentials, organization, errs)
	}

	// main imports (extends field)
	for _, imported := range pipeline.Imports {
		importedPipelines, errs = appendImports(importedPipelines, imported, resources, credentials, organization, errs)
	}

	// job imports (job, step, variable imports)
	for _, job := range pipeline.Jobs {
		if job != nil {
			if job.Imports != nil {
				importedPipelines, errs = appendImports(importedPipelines, job.Imports, resources, credentials, organization, errs)
			}
		}

		if job.EnvironmentVariables != nil && job.EnvironmentVariables.Imports != nil {
			importedPipelines, errs = appendImports(importedPipelines, job.EnvironmentVariables.Imports, resources, credentials, organization, errs)
		}

		if len(job.PreSteps) > 0 {
			importedPipelines, errs = iterateSteps(job.PreSteps, importedPipelines, resources, credentials, organization, errs)
		}

		if len(job.Steps) > 0 {
			importedPipelines, errs = iterateSteps(job.Steps, importedPipelines, resources, credentials, organization, errs)

		}

		if len(job.PostSteps) > 0 {
			importedPipelines, errs = iterateSteps(job.PostSteps, importedPipelines, resources, credentials, organization, errs)

		}
	}

	return importedPipelines, errs
}

func handleImport(jobImport *models.Import, resources *models.Resources, credentials *models.Credentials, organization *string) ([]byte, error) {
	if jobImport == nil || jobImport.Source == nil {
		return nil, nil
	}

	if jobImport.Source.Type == models.SourceTypeRemote && resources != nil && len(resources.Repositories) > 0 {
		return loadRemoteFile(jobImport, resources, credentials, organization)
	}

	if jobImport.Source.Type == models.SourceTypeLocal && jobImport.Source.Path != nil {
		return loadLocalFile(*jobImport.Source.Path)
	}

	return nil, nil
}

func loadRemoteFile(jobImport *models.Import, resources *models.Resources, credentials *models.Credentials, organization *string) ([]byte, error) {
	project, repo, path, version, _ := extractRemoteParams(jobImport, resources)
	if project == "" || repo == "" || path == "" || organization == nil || *organization == "" {
		return nil, nil
	}
	url := generateRequestUrl(project, repo, path, version, *organization)
	client := utils.GetHttpClientWithBasicAuth(credentials)
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

func extractRemoteParams(jobImport *models.Import, resources *models.Resources) (string, string, string, string, models.Platform) {
	var proj, repo, path, version string
	var platform models.Platform
	if jobImport.Source.Path != nil {
		path = *jobImport.Source.Path
	}
	if jobImport.Source.SCM != "" {
		platform = jobImport.Source.SCM
	}

	if resources != nil && len(resources.Repositories) > 0 {
		for _, item := range resources.Repositories {
			// atm we only support current azure organization import (no github import or import from another azure org)
			if item != nil && item.SCM == consts.AzurePlatform && *item.RepositoryAlias == *jobImport.Source.RepositoryAlias && item.Repository != nil {
				parts := strings.Split(*item.Repository, "/")
				if len(parts) == 1 {
					repo = parts[0]
				} else {
					proj = parts[0]
					repo = parts[1]
				}
				if item.Reference != nil {
					version = *item.Reference
				}
				break
			}
		}
	}

	return proj, repo, path, version, platform
}

func generateRequestUrl(proj, repo, path, version, organization string) string {
	url := AZURE_BASE_URL + ITEMS_API
	url = strings.Replace(url, "{ORGANIZATION}", organization, 1)
	url = strings.Replace(url, "{PROJECT}", proj, 1)
	url = strings.Replace(url, "{REPOSITORY}", repo, 1)
	url = strings.Replace(url, "{PATH}", path, 1)
	if version != "" {
		url = url + VERSION_QUERY + version
	}
	return url
}

func getImportedData(
	imported *models.Import,
	resources *models.Resources,
	credentials *models.Credentials,
	organization *string) (*enhancers.ImportedPipeline, error) {
	if imported != nil && imported.Source != nil && imported.Source.Path != nil {
		importedPipelineBuf, err := handleImport(imported, resources, credentials, organization)
		if err != nil {
			return nil, err
		}
		return &enhancers.ImportedPipeline{
			JobName:             *imported.Source.Path,
			OriginFileReference: imported.FileReference,
			Data:                importedPipelineBuf,
		}, nil
	}
	return nil, nil
}

func appendImports(
	list []*enhancers.ImportedPipeline,
	imports *models.Import,
	resources *models.Resources,
	credentials *models.Credentials,
	organization *string,
	errs error) ([]*enhancers.ImportedPipeline, error) {
	importedPipeline, err := getImportedData(imports, resources, credentials, organization)
	if err != nil {
		if errs == nil {
			errs = errors.New("got error(s) importing pipeline(s):")
		}
		errs = errors.Wrap(errs, fmt.Sprintf("error importing pipeline %s: %s", *imports.Source.Path, err.Error()))
	}
	if importedPipeline != nil {
		list = append(list, importedPipeline)
	}

	return list, errs
}

func iterateSteps(
	steps []*models.Step,
	list []*enhancers.ImportedPipeline,
	resources *models.Resources,
	credentials *models.Credentials,
	organization *string,
	errs error) ([]*enhancers.ImportedPipeline, error) {
	for _, step := range steps {
		if step != nil && step.Imports != nil {
			list, errs = appendImports(list, step.Imports, resources, credentials, organization, errs)
		}
		if step != nil && step.EnvironmentVariables != nil && step.EnvironmentVariables.Imports != nil {
			list, errs = appendImports(list, step.EnvironmentVariables.Imports, resources, credentials, organization, errs)
		}
	}

	return list, errs
}

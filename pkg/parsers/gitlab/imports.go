package gitlab

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	parserUtils "github.com/argonsecurity/pipeline-parser/pkg/parsers/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

const (
	TEMPLATE_URL_FORMAT = "https://gitlab.com/gitlab-org/gitlab/-/raw/master/lib/gitlab/ci/templates/%s"
)

var (
	gitlabRemotePipelineRegex = regexp.MustCompile(`https://gitlab\.com/(?P<group>[\w-_]+)/(?P<project>.*?)/(?:-/)?raw/(?P<ref>.*?)/(?P<filePath>.*\.ya?ml)`)
)

func parseImports(include *gitlabModels.Include) []*models.Import {
	if include == nil {
		return nil
	}

	imports := []*models.Import{}
	for _, item := range *include {
		if importItem := parseIncludeItem(item); importItem != nil {
			imports = append(imports, importItem)
		}
	}
	return imports
}

func parseIncludeItem(item gitlabModels.IncludeItem) *models.Import {
	if item.Local != "" {
		return parseLocalImport(&item)
	}

	if item.Remote != "" {
		return parseRemoteImport(&item)
	}

	if item.File != "" {
		return parseFileImport(&item)
	}

	if item.Template != "" {
		return parseTemplateImport(&item)
	}

	return nil
}

func parseLocalImport(item *gitlabModels.IncludeItem) *models.Import {
	if item.Local == "" {
		return nil
	}

	return &models.Import{
		Source: &models.ImportSource{
			SCM:  consts.GitLabPlatform,
			Type: models.SourceTypeLocal,
			Path: &item.Local,
		},
		FileReference: item.FileReference,
	}
}

func parseRemoteImport(item *gitlabModels.IncludeItem) *models.Import {
	if item.Remote == "" {
		return nil
	}

	group, project, ref, filePath := extractRemotePipelineInfo(item.Remote)
	versionType := parserUtils.DetectVersionType(ref)
	return &models.Import{
		Source: &models.ImportSource{
			SCM:          consts.GitLabPlatform,
			Type:         models.SourceTypeRemote,
			Organization: utils.GetPtr(group),
			Repository:   utils.GetPtr(project),
			Path:         utils.GetPtr(filePath),
		},
		Version:       utils.GetPtr(ref),
		VersionType:   versionType,
		FileReference: item.FileReference,
	}
}

func extractRemotePipelineInfo(url string) (group string, project string, ref string, filePath string) {
	match := gitlabRemotePipelineRegex.FindStringSubmatch(url)
	if len(match) == 5 {
		return match[1], match[2], match[3], match[4]
	}
	return "", "", "", ""
}

func parseFileImport(item *gitlabModels.IncludeItem) *models.Import {
	if item.File == "" {
		return nil
	}

	splitProject := strings.SplitN(item.Project, "/", 2)
	if len(splitProject) != 2 {
		return nil
	}

	importData := &models.Import{
		Source: &models.ImportSource{
			SCM:          consts.GitLabPlatform,
			Type:         models.SourceTypeRemote,
			Path:         &item.File,
			Repository:   utils.GetPtr(splitProject[1]),
			Organization: utils.GetPtr(splitProject[0]),
		},
		FileReference: item.FileReference,
	}

	if item.Ref != "" {
		importData.Version = &item.Ref
		importData.VersionType = parserUtils.DetectVersionType(item.Ref)
	}

	return importData
}

func parseTemplateImport(item *gitlabModels.IncludeItem) *models.Import {
	if item.Template == "" {
		return nil
	}

	fullTemplateUrl := fmt.Sprintf(TEMPLATE_URL_FORMAT, item.Template)
	item.Remote = fullTemplateUrl
	return parseRemoteImport(item)
}

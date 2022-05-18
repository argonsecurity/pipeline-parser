package gitlab

import (
	gitlabModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models"
)

// TODO: add ref and project
func parseImports(include *gitlabModels.Include) []string {
	if include == nil {
		return nil
	}

	imports := []string{}
	for _, item := range *include {
		if item.Template != "" {
			imports = append(imports, item.Template)
		}
		if item.Local != "" {
			imports = append(imports, item.Local)
		}
		if item.Remote != "" {
			imports = append(imports, item.Remote)
		}

	}
	return imports
}

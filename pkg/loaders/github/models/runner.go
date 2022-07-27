package models

import (
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"gopkg.in/yaml.v3"
)

type RunsOn struct {
	OS            *string
	Arch          *string
	SelfHosted    bool
	Tags          []string
	FileReference *models.FileReference
}

func (r *RunsOn) UnmarshalYAML(node *yaml.Node) error {
	var tags []string
	var err error
	if node.Tag == consts.StringTag {
		tags = []string{node.Value}
	} else if node.Tag == consts.SequenceTag {
		if tags, err = loadersUtils.ParseYamlStringSequenceToSlice(node); err != nil {
			return err
		}
	} else {
		return consts.NewErrInvalidYamlTag(node.Tag, "RunsOn")
	}

	*r = *generateRunsOnFromTags(tags)
	r.FileReference = loadersUtils.GetFileReference(node)
	return nil
}

func generateRunsOnFromTags(tags []string) *RunsOn {
	r := &RunsOn{}
	r.Tags = tags
	for _, tag := range tags {
		r = parseTag(r, tag)
	}
	return r
}

func parseTag(r *RunsOn, tag string) *RunsOn {
	if tag == consts.SelfHosted {
		r.SelfHosted = true
	}

	for os, keywords := range consts.OsToKeywords {
		didFind := false
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(tag), keyword) {
				r.OS = utils.GetPtr(string(os))
				didFind = true
				break
			}
		}
		if didFind {
			break
		}
	}

	for _, arch := range consts.ArchKeywords {
		if strings.Contains(strings.ToLower(tag), arch) {
			r.Arch = utils.GetPtr(arch)
			break
		}
	}

	return r
}

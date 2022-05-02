package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

const (
	macos1015     = "macos-10.15"
	macos11       = "macos-11"
	macosLatest   = "macos-latest"
	selfHosted    = "self-hosted"
	ubuntu1804    = "ubuntu-18.04"
	ubuntu2004    = "ubuntu-20.04"
	ubuntuLatest  = "ubuntu-latest"
	windows2016   = "windows-2016"
	windows2019   = "windows-2019"
	windows2022   = "windows-2022"
	windowsLatest = "windows-latest"

	armArch = "arm32"
	x64Arch = "x64"
	x32Arch = "x32"
)

var (
	windowsKeywords = []string{"windows", windows2016, windows2019, windows2022, windowsLatest}
	linuxKeywords   = []string{"linux", "ubuntu", "debian", ubuntu1804, ubuntu2004, ubuntuLatest}
	macKeywords     = []string{"macos", "darwin", "osx", macos1015, macos11, macosLatest}

	archKeywords = []string{armArch, x64Arch, x32Arch}

	osToKeywords = map[models.OS][]string{
		models.WindowsOS: windowsKeywords,
		models.LinuxOS:   linuxKeywords,
		models.MacOS:     macKeywords,
	}
)

type RunsOn struct {
	OS         *string
	Arch       *string
	SelfHosted bool
	Tags       []string
}

func (r *RunsOn) UnmarshalYAML(node *yaml.Node) error {
	var tags []string
	if node.Tag == "!!str" {
		tags = []string{node.Value}
	} else if node.Tag == "!!seq" {
		tags = utils.Map(node.Content, func(n *yaml.Node) string {
			return n.Value
		})
	} else {
		return errors.New("invalid RunsOn tags")
	}

	r = generateRunsOnFromTags(tags)
	return nil
}

func DecodeRunsOnHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any) (any, error) {
		if t != reflect.TypeOf(RunsOn{}) {
			return data, nil
		}

		var tags []string
		if f.Kind() == reflect.String {
			tags = []string{data.(string)}
		} else {
			if err := mapstructure.Decode(data, &tags); err != nil {
				return data, err
			}
		}

		return generateRunsOnFromTags(tags), nil
	}
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
	if tag == selfHosted {
		r.SelfHosted = true
	}

	for os, keywords := range osToKeywords {
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

	for _, arch := range archKeywords {
		if strings.Contains(strings.ToLower(tag), arch) {
			r.Arch = utils.GetPtr(arch)
			break
		}
	}

	return r
}

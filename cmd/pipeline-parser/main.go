package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/handler"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	platform              string
	platformFlagName      = "platform"
	platformShortFlagName = "p"
	platformDefaultValue  = string(consts.GitHubPlatform)
	platformUsage         = fmt.Sprintf("CI platform to parse - %v", consts.Platforms)

	output              string
	outputFlagName      = "output"
	outputShortFlagName = "o"
	outputDefaultValue  = string(consts.Stdout)
	outputUsage         = fmt.Sprintf("Output target - %v", consts.OutputTargets)

	fileSuffix             string
	fileSuffixFlagName     = "file-suffix"
	fileSuffixDefaultValue = "parsed"
	fileSuffixUsage        = "File suffix for output file. This flag is useless if 'output' flag is not set to 'file'"

	token             string
	tokenFlagName     = "token"
	tokenDefaultValue = ""
	tokenUsage        = "SCM token to use for fetching remote files if necessary"

	organization             string
	organizationFlagName     = "organization"
	organizationDefaultValue = ""
	organizationUsage        = "The target organization when fetching remote files (used for Azure Pipelines)"

	version string
)

func main() {
	c := GetCommand(version)
	c.Execute()
}

func GetCommand(version string) *cobra.Command {
	command := &cobra.Command{
		Use:   "pipeline-parser",
		Short: "Parses a pipeline file",
		Long:  "Parses a pipeline file",
		Example: `pipeline-parser --platform github workflow.yml
pipeline-parser --platform gitlab .gitlab-ci.yml
pipeline-parser --platform azure azure-pipelines.yml`,
		SilenceUsage: true,
		Version:      version,
		PreRunE:      preRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, pipelinePath := range args {
				if fi, err := os.Stat(pipelinePath); !os.IsNotExist(err) && !fi.IsDir() {
					buf, err := ioutil.ReadFile(pipelinePath)
					if err != nil {
						return nil
					}
					pipeline, err := handler.Handle(buf, models.Platform(platform), &models.Credentials{Token: token}, organization)
					if err != nil {
						return err
					}
					if err := writePipelineToOutput(pipeline, consts.OutputTarget(output), pipelinePath); err != nil {
						return err
					}
				} else {
					return err
				}
			}
			return nil
		},
	}

	command.PersistentFlags().StringVarP(&platform, platformFlagName, platformShortFlagName, platformDefaultValue, platformUsage)
	command.PersistentFlags().StringVarP(&output, outputFlagName, outputShortFlagName, outputDefaultValue, outputUsage)
	command.PersistentFlags().StringVar(&fileSuffix, fileSuffixFlagName, fileSuffixDefaultValue, fileSuffixUsage)
	command.PersistentFlags().StringVar(&token, tokenFlagName, tokenDefaultValue, tokenUsage)
	command.PersistentFlags().StringVar(&token, tokenFlagName, tokenDefaultValue, tokenUsage)

	return command
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return consts.NewErrInvalidArgumentsCount(len(args))
	}

	if !slices.Contains(consts.Platforms, models.Platform(platform)) {
		return consts.NewErrInvalidPlatform(models.Platform(platform))
	}

	if !slices.Contains(consts.OutputTargets, consts.OutputTarget(output)) {
		return consts.NewErrInvalidOutputTarget(consts.OutputTarget(output))
	}

	return nil
}

func writePipelineToOutput(pipeline *models.Pipeline, outputTarget consts.OutputTarget, pipelinePath string) error {
	jsonPipeline, err := json.MarshalIndent(pipeline, "", " ")
	if err != nil {
		return err
	}

	switch outputTarget {
	case consts.Stdout:
		fmt.Printf("%s:\n", pipelinePath)
		fmt.Println(string(jsonPipeline))
	case consts.File:
		outputFilePath := getOutputFilePath(pipelinePath, fileSuffix)
		if err = ioutil.WriteFile(outputFilePath, jsonPipeline, 0644); err != nil {
			return err
		}
	}

	return nil
}

func getOutputFilePath(pipelinePath string, fileSuffix string) string {
	ext := filepath.Ext(pipelinePath)
	base := filepath.Base(pipelinePath)

	return filepath.Join(filepath.Dir(pipelinePath), fmt.Sprintf("%s_%s.json", base[0:len(base)-len(ext)], fileSuffix))
}

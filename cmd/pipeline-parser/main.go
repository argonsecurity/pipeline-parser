package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/handler"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	platform              string
	platformFlagName      = "platform"
	platformShortFlagName = "p"
	platformDefaultValue  = string(consts.GitHubPlatform)
	platformUsage         = "CI platform to parse"

	version string
)

func main() {
	c := GetCommand(version)
	c.Execute()
}

func GetCommand(version string) *cobra.Command {
	command := &cobra.Command{
		Use:          "pipeline-parser",
		Short:        "Parses a pipeline file",
		Long:         "Parses a pipeline file",
		SilenceUsage: true,
		Version:      version,
		PreRunE:      preRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, workflowPath := range args {
				if fi, err := os.Stat(workflowPath); !os.IsNotExist(err) && !fi.IsDir() {
					buf, err := ioutil.ReadFile(workflowPath)
					if err != nil {
						return nil
					}
					pipeline, err := handler.Handle(buf, consts.GitHubPlatform)
					if err != nil {
						return err
					}
					fmt.Println(pipeline)
				} else {
					return err
				}
			}
			return nil
		},
	}

	command.PersistentFlags().StringVarP(&platform, platformFlagName, platformShortFlagName, platformDefaultValue, platformUsage)

	return command
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Invalid arguments length")
	}

	if !slices.Contains(consts.Platforms, consts.Platform(platform)) {
		return fmt.Errorf("Invalid platform: %s. Supported platforms: %v", platform, consts.Platforms)
	}

	return nil
}

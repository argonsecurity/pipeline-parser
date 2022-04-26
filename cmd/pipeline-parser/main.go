package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/handler"
	"github.com/spf13/cobra"
)

var (
	platform     string
	platformFlag = "f"
)

func GetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "pipeline-parser",
		Short: "Parses a pipeline file",
		Long:  "Parses a pipeline file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}

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
	command.PersistentFlags().StringVarP(&platform, "platform", platformFlag, string(consts.GitHubPlatform), "Platform to parse")
	return command
}

func main() {
	c := GetCommand()
	c.Execute()
}

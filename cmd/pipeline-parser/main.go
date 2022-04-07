package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/argonsecurity/pipeline-parser/pkg/parsers/github"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pipeline-parser",
		Short: "Parses a pipeline file",
		Long:  "Parses a pipeline file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}

			for _, workflowPath := range args {
				if fi, err := os.Stat(workflowPath); os.IsNotExist(err) || !fi.IsDir() {
					buf, err := ioutil.ReadFile(workflowPath)
					if err != nil {
						return nil
					}
					root, err := github.Parse(buf)
					if err != nil {
						return err
					}

					fmt.Println(root)
				}
			}
			return nil
		},
	}
}

func main() {
	c := GetCommand()
	c.Execute()
}

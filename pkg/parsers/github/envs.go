package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/mitchellh/mapstructure"
)

func parseEnvironmentVariables(env interface{}) *models.EnvironmentVariables {
	if env == nil {
		return nil
	}
	var envVars models.EnvironmentVariables
	mapstructure.Decode(env, &envVars)
	return &envVars
}

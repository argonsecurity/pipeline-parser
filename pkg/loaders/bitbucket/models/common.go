package models

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

type EnvironmentVariablesRef struct {
	models.EnvironmentVariables
	FileReference *models.FileReference
}

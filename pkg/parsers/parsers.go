package parsers

import "github.com/argonsecurity/pipeline-parser/pkg/models"

type Parser[T any] interface {
	Parse(*T) (*models.Pipeline, error)
}

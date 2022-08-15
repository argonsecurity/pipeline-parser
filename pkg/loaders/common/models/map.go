package models

import (
	loadersUtils "github.com/argonsecurity/pipeline-parser/pkg/loaders/utils"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"gopkg.in/yaml.v3"
)

/*
	Map is a struct that is used to describe a yaml tag that is a map of key value pairs.
*/
type Map struct {
	FileReference *models.FileReference
	Values        []*MapEntry
}

type MapEntry struct {
	Key           string
	Value         any
	FileReference *models.FileReference
}

func (m *Map) UnmarshalYAML(node *yaml.Node) error {
	keyValues := []*MapEntry{}
	if err := loadersUtils.IterateOnMap(node,
		func(key string, value *yaml.Node) error {
			keyValues = append(keyValues, &MapEntry{
				Key:           key,
				Value:         loadersUtils.GetNodeValue(value),
				FileReference: loadersUtils.GetFileReference(value),
			})
			return nil
		},
		"Map"); err != nil {
		return err
	}

	m.Values = keyValues

	m.FileReference = loadersUtils.GetFileReference(node)

	// node.Line is the line number of the first key-value pair of the map
	// but we want to set the file reference from the map declaration
	// e.g `
	// map:
	//   key: value`
	// set the file reference from "map:"
	m.FileReference.StartRef.Line--
	m.FileReference.StartRef.Column -= 2

	return nil
}

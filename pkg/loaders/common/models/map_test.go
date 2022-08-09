package models

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type testMap struct {
	Map *Map `yaml:"map"`
}

func TestMapLoader(t *testing.T) {
	testCases := []struct {
		name        string
		yamlBuffer  []byte
		expectedMap *testMap
	}{
		{
			name:        "Map is nil",
			yamlBuffer:  nil,
			expectedMap: nil,
		},
		{
			name:        "Map is empty",
			yamlBuffer:  []byte{},
			expectedMap: nil,
		},
		{
			name: "Map with data",
			yamlBuffer: []byte(`map:
  string: "string"
  int: 1
  bool: true`,
			),
			expectedMap: &testMap{
				Map: &Map{
					Values: []*MapEntry{
						{
							Key:           "string",
							Value:         "string",
							FileReference: testutils.CreateFileReference(2, 3, 2, 9),
						},
						{
							Key:           "int",
							Value:         1,
							FileReference: testutils.CreateFileReference(3, 3, 3, 4),
						},
						{
							Key:           "bool",
							Value:         true,
							FileReference: testutils.CreateFileReference(4, 3, 4, 7),
						},
					},
					FileReference: testutils.CreateFileReference(1, 1, 4, 7),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var m *testMap
			err := yaml.Unmarshal(testCase.yamlBuffer, &m)
			assert.NoError(t, err)

			changelog, err := diff.Diff(testCase.expectedMap, m)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, "test case failed with %d modifications", len(changelog))
			for _, change := range changelog {
				t.Log(change)
			}
		})
	}
}

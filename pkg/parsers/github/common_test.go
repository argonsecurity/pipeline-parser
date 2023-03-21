package github

import (
	"reflect"
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
)

func TestParseEnvironmentVariablesRef(t *testing.T) {
	testCases := []struct {
		name        string
		envRef      *githubModels.EnvironmentVariablesRef
		expectedEnv *models.EnvironmentVariablesRef
	}{
		{
			name:        "Input is nil",
			envRef:      nil,
			expectedEnv: nil,
		},
		{
			name: "Input is not nil",
			envRef: &githubModels.EnvironmentVariablesRef{
				EnvironmentVariables: map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
				FileReference: &models.FileReference{
					StartRef: &models.FileLocation{
						Line:   1,
						Column: 2,
					},
					EndRef: &models.FileLocation{
						Line:   3,
						Column: 4,
					},
				},
			},
			expectedEnv: &models.EnvironmentVariablesRef{
				EnvironmentVariables: map[string]any{
					"key1": "value1",
					"key2": "value2",
				},
				FileReference: &models.FileReference{
					StartRef: &models.FileLocation{
						Line:   1,
						Column: 2,
					},
					EndRef: &models.FileLocation{
						Line:   3,
						Column: 4,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseEnvironmentVariablesRef(testCase.envRef)

			testutils.DeepCompare(t, testCase.expectedEnv, got)
		})
	}
}

func Test_detectVersionType(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want models.VersionType
	}{
		{
			name: "Empty version",
			args: args{
				version: "",
			},
			want: models.None,
		},
		{
			name: "Commit SHA",
			args: args{
				version: "1234567890123456789012345678901234567890",
			},
			want: models.CommitSHA,
		},
		{
			name: "Tag version",
			args: args{
				version: "v1.2.3",
			},
			want: models.TagVersion,
		},
		{
			name: "Branch version",
			args: args{
				version: "master",
			},
			want: models.BranchVersion,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectVersionType(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("detectVersionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

package utils

import (
	"reflect"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func TestDetectVersionType(t *testing.T) {
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
		{
			name: "Latest version",
			args: args{
				version: "latest",
			},
			want: models.Latest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectVersionType(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectVersionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

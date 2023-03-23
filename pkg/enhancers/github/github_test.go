package github

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_mergePipelines(t *testing.T) {
	type args struct {
		pipeline         *models.Pipeline
		importedPipeline *enhancers.ImportedPipeline
	}
	tests := []struct {
		name string
		args args
		want *models.Pipeline
	}{
		{
			name: "happy flow",
			args: args{
				pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							Name:    utils.GetPtr("job1"),
							Imports: &models.Import{},
						},
					},
				},
				importedPipeline: &enhancers.ImportedPipeline{
					JobName: "job1",
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("imported"),
					},
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name: utils.GetPtr("job1"),
						Imports: &models.Import{
							Pipeline: &models.Pipeline{
								Name: utils.GetPtr("imported"),
							},
						},
					},
				},
			},
		},
		{
			name: "jobs names don't match",
			args: args{
				pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							Name:    utils.GetPtr("job1"),
							Imports: &models.Import{},
						},
					},
				},
				importedPipeline: &enhancers.ImportedPipeline{
					JobName: "job2",
					Pipeline: &models.Pipeline{
						Name: utils.GetPtr("imported"),
					},
				},
			},
			want: &models.Pipeline{
				Jobs: []*models.Job{
					{
						Name:    utils.GetPtr("job1"),
						Imports: &models.Import{},
					},
				},
			},
		},
		{
			name: "pipeline is nil",
			args: args{
				pipeline: nil,
			},
			want: nil,
		},
		{
			name: "pipeline jobs are nil",
			args: args{
				pipeline: &models.Pipeline{
					Name: utils.GetPtr("pipeline"),
				},
			},
			want: &models.Pipeline{
				Name: utils.GetPtr("pipeline"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergePipelines(tt.args.pipeline, tt.args.importedPipeline)
			assert.Equal(t, tt.want, got)
		})
	}
}

package azure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_getTemplates(t *testing.T) {
	type args struct {
		pipeline    *models.Pipeline
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    []*enhancers.ImportedPipeline
		wantErr bool
	}{
		{
			name: "happy flow",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									RepositoryAlias: utils.GetPtr("templates"),
									Organization:    utils.GetPtr("proj"),
									Repository:      utils.GetPtr("repo"),
									SCM:             consts.AzurePlatform,
								},
							},
						},
					},
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								RepositoryAlias: utils.GetPtr("templates"),
								Path:            utils.GetPtr("file"),
								Type:            models.SourceTypeRemote,
							},
						},
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/file"),
								Type: models.SourceTypeLocal,
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "file",
					Data:    []byte("file content"),
				},
				{
					JobName: "testdata/file",
					Data:    []byte("file content"),
				},
			},
		},
		{
			name: "local does not exist",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									RepositoryAlias: utils.GetPtr("templates"),
									Organization:    utils.GetPtr("proj"),
									Repository:      utils.GetPtr("repo"),
									SCM:             consts.AzurePlatform,
								},
							},
						},
					},
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								RepositoryAlias: utils.GetPtr("templates"),
								Path:            utils.GetPtr("file"),
								Type:            models.SourceTypeRemote,
							},
						},
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/other"),
								Type: models.SourceTypeLocal,
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "remote",
					Data:    []byte("file content"),
				},
				{
					JobName: "local",
					Data:    nil,
				},
			},
			wantErr: true,
		},
		{
			name: "remote does not exist",
			args: args{
				pipeline: &models.Pipeline{
					Defaults: &models.Defaults{
						Resources: &models.Resources{
							Repositories: []*models.ImportSource{
								{
									RepositoryAlias: utils.GetPtr("templates"),
									Organization:    utils.GetPtr("proj"),
									Repository:      utils.GetPtr("other"),
									SCM:             consts.AzurePlatform,
								},
							},
						},
					},
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								RepositoryAlias: utils.GetPtr("templates"),
								Path:            utils.GetPtr("file"),
								Type:            models.SourceTypeRemote,
							},
						},
						{
							Source: &models.ImportSource{
								Path: utils.GetPtr("testdata/file"),
								Type: models.SourceTypeLocal,
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "remote",
					Data:    nil,
				},
				{
					JobName: "local",
					Data:    []byte("file content"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.FileServer(http.Dir("testdata"))
			ts := httptest.NewServer(h)
			defer ts.Close()
			AZURE_BASE_URL = ts.URL

			got, err := getTemplates(tt.args.pipeline, tt.args.credentials, "azure-org")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

package gitlab

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func Test_handleLocalImport(t *testing.T) {
	type args struct {
		importData *models.Import
	}
	tests := []struct {
		name    string
		args    args
		want    *enhancers.ImportedPipeline
		wantErr bool
	}{
		{
			name: "valid local import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type: models.SourceTypeLocal,
						SCM:  consts.GitLabPlatform,
						Path: utils.GetPtr("testdata/pipeline.yaml"),
					},
				},
			},
			want: &enhancers.ImportedPipeline{
				Data: []byte("test data\n"),
			},
		},
		{
			name: "invalid local import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type: models.SourceTypeLocal,
						SCM:  consts.GitLabPlatform,
						Path: utils.GetPtr("testdata/invalid.yaml"),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleLocalImport(tt.args.importData)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleLocalImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleLocalImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleRemoteImport(t *testing.T) {
	type args struct {
		importData  *models.Import
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    *enhancers.ImportedPipeline
		wantErr bool
	}{
		{
			name: "valid remote import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						SCM:          consts.GitLabPlatform,
						Path:         utils.GetPtr("pipeline.yaml"),
						Organization: utils.GetPtr("group"),
						Repository:   utils.GetPtr("subgroup/project"),
					},
					Version:     utils.GetPtr("master"),
					VersionType: models.BranchVersion,
				},
			},
			want: &enhancers.ImportedPipeline{
				Data: []byte("test data\n"),
			},
		},
		{
			name: "invalid remote import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						SCM:          consts.GitLabPlatform,
						Path:         utils.GetPtr("pipeline.yaml"),
						Organization: utils.GetPtr("group"),
						Repository:   utils.GetPtr("subgroup/project"),
					},
					Version:     utils.GetPtr("invalidbranch"),
					VersionType: models.BranchVersion,
				},
			},
			wantErr: true,
		},
		{
			name: "nil values",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type: models.SourceTypeRemote,
						SCM:  consts.GitLabPlatform,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		h := http.FileServer(http.Dir("testdata"))
		ts := httptest.NewServer(h)
		defer ts.Close()
		GITLAB_BASE_URL = ts.URL

		t.Run(tt.name, func(t *testing.T) {
			got, err := handleRemoteImport(tt.args.importData, tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleRemoteImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleRemoteImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleImport(t *testing.T) {
	type args struct {
		importData  *models.Import
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    *enhancers.ImportedPipeline
		wantErr bool
	}{
		{
			name: "valid local import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type: models.SourceTypeLocal,
						SCM:  consts.GitLabPlatform,
						Path: utils.GetPtr("testdata/pipeline.yaml"),
					},
				},
			},
			want: &enhancers.ImportedPipeline{
				Data: []byte("test data\n"),
			},
		},
		{
			name: "valid remote import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						SCM:          consts.GitLabPlatform,
						Path:         utils.GetPtr("pipeline.yaml"),
						Organization: utils.GetPtr("group"),
						Repository:   utils.GetPtr("subgroup/project"),
					},
					Version:     utils.GetPtr("master"),
					VersionType: models.BranchVersion,
				},
			},
			want: &enhancers.ImportedPipeline{
				Data: []byte("test data\n"),
			},
		},
		{
			name: "invalid remote import",
			args: args{
				importData: &models.Import{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						SCM:          consts.GitLabPlatform,
						Path:         utils.GetPtr("pipeline.yaml"),
						Organization: utils.GetPtr("group"),
						Repository:   utils.GetPtr("subgroup/project"),
					},
					Version:     utils.GetPtr("invalidbranch"),
					VersionType: models.BranchVersion,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		h := http.FileServer(http.Dir("testdata"))
		ts := httptest.NewServer(h)
		defer ts.Close()
		GITLAB_BASE_URL = ts.URL

		t.Run(tt.name, func(t *testing.T) {
			got, err := handleImport(tt.args.importData, tt.args.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitLabEnhancer_LoadImportedPipelines(t *testing.T) {
	type args struct {
		data        *models.Pipeline
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		g       *GitLabEnhancer
		args    args
		want    []*enhancers.ImportedPipeline
		wantErr bool
	}{
		{
			name: "valid local import",
			args: args{
				data: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Type: models.SourceTypeLocal,
								SCM:  consts.GitLabPlatform,
								Path: utils.GetPtr("testdata/pipeline.yaml"),
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					Data: []byte("test data\n"),
				},
			},
		},
		{
			name: "valid remote import",
			args: args{
				data: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Type:         models.SourceTypeRemote,
								SCM:          consts.GitLabPlatform,
								Path:         utils.GetPtr("pipeline.yaml"),
								Organization: utils.GetPtr("group"),
								Repository:   utils.GetPtr("subgroup/project"),
							},
							Version:     utils.GetPtr("master"),
							VersionType: models.BranchVersion,
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					Data: []byte("test data\n"),
				},
			},
		},
		{
			name: "multiple imports - some invalid",
			args: args{
				data: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Type: models.SourceTypeLocal,
								SCM:  consts.GitLabPlatform,
								Path: utils.GetPtr("testdata/pipeline.yaml"),
							},
						},
						{
							Source: &models.ImportSource{
								Type:         models.SourceTypeRemote,
								SCM:          consts.GitLabPlatform,
								Path:         utils.GetPtr("pipeline.yaml"),
								Organization: utils.GetPtr("group"),
								Repository:   utils.GetPtr("subgroup/project"),
							},
							Version:     utils.GetPtr("master"),
							VersionType: models.BranchVersion,
						},
						{
							Source: &models.ImportSource{
								Type:         models.SourceTypeRemote,
								SCM:          consts.GitLabPlatform,
								Path:         utils.GetPtr("pipeline.yaml"),
								Organization: utils.GetPtr("group"),
								Repository:   utils.GetPtr("subgroup/project"),
							},
							Version:     utils.GetPtr("invalidbranch"),
							VersionType: models.BranchVersion,
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					Data: []byte("test data\n"),
				},
				{
					Data: []byte("test data\n"),
				},
				nil,
			},
			wantErr: true,
		},
		{
			name: "invalid remote import",
			args: args{
				data: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Type:         models.SourceTypeRemote,
								SCM:          consts.GitLabPlatform,
								Path:         utils.GetPtr("pipeline.yaml"),
								Organization: utils.GetPtr("group"),
								Repository:   utils.GetPtr("subgroup/project"),
							},
							Version:     utils.GetPtr("invalidbranch"),
							VersionType: models.BranchVersion,
						},
					},
				},
			},
			want:    []*enhancers.ImportedPipeline{nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		h := http.FileServer(http.Dir("testdata"))
		ts := httptest.NewServer(h)
		defer ts.Close()
		GITLAB_BASE_URL = ts.URL
		t.Run(tt.name, func(t *testing.T) {
			g := &GitLabEnhancer{}
			got, err := g.LoadImportedPipelines(tt.args.data, tt.args.credentials, utils.GetPtr(""), utils.GetPtr(""))
			if (err != nil) != tt.wantErr {
				t.Errorf("GitLabEnhancer.LoadImportedPipelines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitLabEnhancer.LoadImportedPipelines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitLabEnhancer_Enhance(t *testing.T) {
	type args struct {
		data              *models.Pipeline
		importedPipelines []*enhancers.ImportedPipeline
	}
	tests := []struct {
		name    string
		g       *GitLabEnhancer
		args    args
		want    *models.Pipeline
		wantErr bool
	}{
		{
			name: "valid local import",
			args: args{
				data: &models.Pipeline{
					Imports: []*models.Import{
						{
							Source: &models.ImportSource{
								Type: models.SourceTypeLocal,
								SCM:  consts.GitLabPlatform,
								Path: utils.GetPtr("testdata/pipeline.yaml"),
							},
						},
						{
							Source: &models.ImportSource{
								Type:         models.SourceTypeRemote,
								SCM:          consts.GitLabPlatform,
								Path:         utils.GetPtr("pipeline.yaml"),
								Organization: utils.GetPtr("group"),
								Repository:   utils.GetPtr("subgroup/project"),
							},
							Version:     utils.GetPtr("invalidbranch"),
							VersionType: models.BranchVersion,
						},
						{
							Source: &models.ImportSource{
								Type: models.SourceTypeLocal,
								SCM:  consts.GitLabPlatform,
								Path: utils.GetPtr("testdata/pipeline.yaml"),
							},
						},
					},
				},
				importedPipelines: []*enhancers.ImportedPipeline{
					{
						Data:     []byte("test data\n"),
						Pipeline: &models.Pipeline{Name: utils.GetPtr("test")},
					},
					nil,
					{
						Data:     []byte("test data\n"),
						Pipeline: &models.Pipeline{Name: utils.GetPtr("test")},
					},
				},
			},
			want: &models.Pipeline{
				Imports: []*models.Import{
					{
						Source: &models.ImportSource{
							Type: models.SourceTypeLocal,
							SCM:  consts.GitLabPlatform,
							Path: utils.GetPtr("testdata/pipeline.yaml"),
						},
						Pipeline: &models.Pipeline{Name: utils.GetPtr("test")},
					},
					{
						Source: &models.ImportSource{
							Type:         models.SourceTypeRemote,
							SCM:          consts.GitLabPlatform,
							Path:         utils.GetPtr("pipeline.yaml"),
							Organization: utils.GetPtr("group"),
							Repository:   utils.GetPtr("subgroup/project"),
						},
						Version:     utils.GetPtr("invalidbranch"),
						VersionType: models.BranchVersion,
					},
					{
						Source: &models.ImportSource{
							Type: models.SourceTypeLocal,
							SCM:  consts.GitLabPlatform,
							Path: utils.GetPtr("testdata/pipeline.yaml"),
						},
						Pipeline: &models.Pipeline{Name: utils.GetPtr("test")},
					},
				},
			},
		},
		{
			name: "nil pipeline",
			args: args{
				data:              nil,
				importedPipelines: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitLabEnhancer{}
			got, err := g.Enhance(tt.args.data, tt.args.importedPipelines)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitLabEnhancer.Enhance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitLabEnhancer.Enhance() = %v, want %v", got, tt.want)
			}
		})
	}
}

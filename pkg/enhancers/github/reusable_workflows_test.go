package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/enhancers"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_getReusableWorkflows(t *testing.T) {
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
					Jobs: []*models.Job{
						{
							Name: utils.GetPtr("remote"),
							Imports: &models.Import{
								Source: &models.ImportSource{
									Organization: utils.GetPtr("org"),
									Repository:   utils.GetPtr("repo"),
									Path:         utils.GetPtr("file"),
									Type:         models.SourceTypeRemote,
								},
								Version: utils.GetPtr("version"),
							},
						},
						{
							Name: utils.GetPtr("local"),
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("testdata/org/repo/version/file"),
									Type: models.SourceTypeLocal,
								},
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "remote",
					Data:    []byte("file data"),
				},
				{
					JobName: "local",
					Data:    []byte("file data"),
				},
			},
		},
		{
			name: "local does not exist",
			args: args{
				pipeline: &models.Pipeline{
					Jobs: []*models.Job{
						{
							Name: utils.GetPtr("remote"),
							Imports: &models.Import{
								Source: &models.ImportSource{
									Organization: utils.GetPtr("org"),
									Repository:   utils.GetPtr("repo"),
									Path:         utils.GetPtr("file"),
									Type:         models.SourceTypeRemote,
								},
								Version: utils.GetPtr("version"),
							},
						},
						{
							Name: utils.GetPtr("local"),
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("testdata/org/repo/version/does-not-exist"),
									Type: models.SourceTypeLocal,
								},
							},
						},
					},
				},
			},
			want: []*enhancers.ImportedPipeline{
				{
					JobName: "remote",
					Data:    []byte("file data"),
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
					Jobs: []*models.Job{
						{
							Name: utils.GetPtr("remote"),
							Imports: &models.Import{
								Source: &models.ImportSource{
									Organization: utils.GetPtr("org"),
									Repository:   utils.GetPtr("repo"),
									Path:         utils.GetPtr("does-not-exist"),
									Type:         models.SourceTypeRemote,
								},
								Version: utils.GetPtr("version"),
							},
						},
						{
							Name: utils.GetPtr("local"),
							Imports: &models.Import{
								Source: &models.ImportSource{
									Path: utils.GetPtr("testdata/org/repo/version/file"),
									Type: models.SourceTypeLocal,
								},
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
					Data:    []byte("file data"),
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
			GITHUB_BASE_URL = ts.URL

			got, err := getReusableWorkflows(tt.args.pipeline, tt.args.credentials)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_handleImport(t *testing.T) {
	type args struct {
		imports     *models.Import
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "remote file - happy flow",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Organization: utils.GetPtr("org"),
						Repository:   utils.GetPtr("repo"),
						Path:         utils.GetPtr("file"),
						Type:         models.SourceTypeRemote,
					},
					Version: utils.GetPtr("version"),
				},
			},
			want: []byte("file data"),
		},
		{
			name: "remote file - file not found",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Organization: utils.GetPtr("org"),
						Repository:   utils.GetPtr("repo"),
						Path:         utils.GetPtr("does-not-exist"),
						Type:         models.SourceTypeRemote,
					},
					Version: utils.GetPtr("version"),
				},
			},
			wantErr: true,
		},
		{
			name: "imports is nil",
			args: args{
				imports: nil,
			},
			want: nil,
		},
		{
			name: "imports.Source is nil",
			args: args{
				imports: &models.Import{
					Source:  nil,
					Version: utils.GetPtr("version"),
				},
			},
			want: nil,
		},
		{
			name: "remote file - imports.Version is nil",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Organization: utils.GetPtr("org"),
						Repository:   utils.GetPtr("repo"),
						Path:         utils.GetPtr("file"),
						Type:         models.SourceTypeRemote,
					},
					Version: nil,
				},
			},
			want: nil,
		},
		{
			name: "remote file - imports with empty data",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Organization: utils.GetPtr(""),
						Repository:   utils.GetPtr(""),
						Path:         utils.GetPtr(""),
						Type:         models.SourceTypeRemote,
					},
					Version: utils.GetPtr(""),
				},
			},
			want: nil,
		},
		{
			name: "local file - happy flow",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Organization: utils.GetPtr(""),
						Repository:   utils.GetPtr(""),
						Path:         utils.GetPtr("testdata/org/repo/version/file"),
						Type:         models.SourceTypeLocal,
					},
					Version: utils.GetPtr(""),
				},
			},
			want: []byte("file data"),
		},
		{
			name: "local file - imports with empty data",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Path: utils.GetPtr(""),
						Type: models.SourceTypeLocal,
					},
					Version: utils.GetPtr(""),
				},
			},
			want: nil,
		},
		{
			name: "local file - file not found",
			args: args{
				imports: &models.Import{
					Source: &models.ImportSource{
						Organization: utils.GetPtr(""),
						Repository:   utils.GetPtr(""),
						Path:         utils.GetPtr("does-not-exist"),
						Type:         models.SourceTypeLocal,
					},
					Version: utils.GetPtr(""),
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
			GITHUB_BASE_URL = ts.URL

			got, err := handleImport(tt.args.imports, tt.args.credentials)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_loadRemoteFile(t *testing.T) {
	type args struct {
		org         string
		repo        string
		version     string
		path        string
		credentials *models.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "happy flow",
			args: args{
				org:     "org",
				repo:    "repo",
				version: "version",
				path:    "file",
			},
			want: []byte("file data"),
		},
		{
			name: "happy flow with credentials",
			args: args{
				org:     "org",
				repo:    "repo",
				version: "version",
				path:    "file",
				credentials: &models.Credentials{
					Token: "token",
				},
			},
			want: []byte("file data"),
		},
		{
			name: "file does not exist",
			args: args{
				org:     "org",
				repo:    "repo",
				version: "version",
				path:    "does-not-exist",
			},
			wantErr: true,
		},
		{
			name: "empty args",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.FileServer(http.Dir("testdata"))
			ts := httptest.NewServer(h)
			defer ts.Close()
			GITHUB_BASE_URL = ts.URL

			got, err := loadRemoteFile(tt.args.org, tt.args.repo, tt.args.version, tt.args.path, tt.args.credentials)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_loadLocalFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "happy flow",
			args: args{
				path: "testdata/org/repo/version/file",
			},
			want: []byte("file data"),
		},
		{
			name: "path does not exist",
			args: args{
				path: "testdata/org/repo/version/does-not-exist",
			},
			wantErr: true,
		},
		{
			name: "empty args",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadLocalFile(tt.args.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

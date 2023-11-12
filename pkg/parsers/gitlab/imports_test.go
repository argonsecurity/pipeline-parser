package gitlab

import (
	"testing"

	"github.com/argonsecurity/pipeline-parser/pkg/consts"
	gitlabCommon "github.com/argonsecurity/pipeline-parser/pkg/loaders/gitlab/models/common"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_parseLocalImport(t *testing.T) {
	type args struct {
		item *gitlabCommon.IncludeItem
	}
	tests := []struct {
		name string
		args args
		want *models.Import
	}{
		{
			name: "Local is empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Local:         "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Local is not empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Local:         "local",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type: models.SourceTypeLocal,
					Path: utils.GetPtr("local"),
					SCM:  consts.GitLabPlatform,
				},
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLocalImport(tt.args.item)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseRemoteImport(t *testing.T) {
	type args struct {
		item *gitlabCommon.IncludeItem
	}
	tests := []struct {
		name string
		args args
		want *models.Import
	}{
		{
			name: "Remote is empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Remote:        "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Remote is not empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Remote:        "https://gitlab.com/group/subgroup/project/-/raw/master/.gitlab-ci.yml",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(".gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("subgroup/project"),
					Organization: utils.GetPtr("group"),
				},
				Version:       utils.GetPtr("master"),
				VersionType:   models.BranchVersion,
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "Remote is an invalid url",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Remote:        "https://gitlab.com/group/subgroup/project/-/raw/master",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(""),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr(""),
					Organization: utils.GetPtr(""),
				},
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				Version:       utils.GetPtr(""),
				VersionType:   models.None,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRemoteImport(tt.args.item)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_extractRemotePipelineInfo(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name         string
		args         args
		wantGroup    string
		wantProject  string
		wantRef      string
		wantFilePath string
	}{
		{
			name: "Valid url",
			args: args{
				url: "https://gitlab.com/group/subgroup/project/-/raw/master/.gitlab-ci.yml",
			},
			wantGroup:    "group",
			wantProject:  "subgroup/project",
			wantRef:      "master",
			wantFilePath: ".gitlab-ci.yml",
		},
		{
			name: "Invalid url",
			args: args{
				url: "https://gitlab.com/group/subgroup/project/-/raw/master",
			},
			wantGroup:    "",
			wantProject:  "",
			wantRef:      "",
			wantFilePath: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroup, gotProject, gotRef, gotFilePath := extractRemotePipelineInfo(tt.args.url)
			if gotGroup != tt.wantGroup {
				t.Errorf("extractRemotePipelineInfo() gotGroup = %v, want %v", gotGroup, tt.wantGroup)
			}
			if gotProject != tt.wantProject {
				t.Errorf("extractRemotePipelineInfo() gotProject = %v, want %v", gotProject, tt.wantProject)
			}
			if gotRef != tt.wantRef {
				t.Errorf("extractRemotePipelineInfo() gotRef = %v, want %v", gotRef, tt.wantRef)
			}
			if gotFilePath != tt.wantFilePath {
				t.Errorf("extractRemotePipelineInfo() gotFilePath = %v, want %v", gotFilePath, tt.wantFilePath)
			}
		})
	}
}

func Test_parseFileImport(t *testing.T) {
	type args struct {
		item *gitlabCommon.IncludeItem
	}
	tests := []struct {
		name string
		args args
		want *models.Import
	}{
		{
			name: "File is empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					File:          "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "File is not empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					File:          ".gitlab-ci.yml",
					Project:       "group/subgroup/project",
					Ref:           "master",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(".gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("subgroup/project"),
					Organization: utils.GetPtr("group"),
				},
				Version:       utils.GetPtr("master"),
				VersionType:   models.BranchVersion,
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "ref is empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					File:          ".gitlab-ci.yml",
					Project:       "group/subgroup/project",
					Ref:           "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(".gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("subgroup/project"),
					Organization: utils.GetPtr("group"),
				},
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "File is not empty and project is invalid",
			args: args{
				item: &gitlabCommon.IncludeItem{
					File:          ".gitlab-ci.yml",
					Project:       "invalid",
					Ref:           "master",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFileImport(tt.args.item)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseTemplateImport(t *testing.T) {
	type args struct {
		item *gitlabCommon.IncludeItem
	}
	tests := []struct {
		name string
		args args
		want *models.Import
	}{
		{
			name: "Template is empty",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Template:      "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Regular template",
			args: args{
				item: &gitlabCommon.IncludeItem{
					Template:      "Android.gitlab-ci.yml",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr("lib/gitlab/ci/templates/Android.gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("gitlab"),
					Organization: utils.GetPtr("gitlab-org"),
				},
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				Version:       utils.GetPtr("master"),
				VersionType:   models.BranchVersion,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTemplateImport(tt.args.item)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseIncludeItem(t *testing.T) {
	type args struct {
		item gitlabCommon.IncludeItem
	}
	tests := []struct {
		name string
		args args
		want *models.Import
	}{
		{
			name: "Local is empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					Local:         "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Local is not empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					Local:         "local",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type: models.SourceTypeLocal,
					Path: utils.GetPtr("local"),
					SCM:  consts.GitLabPlatform,
				},
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "Remote is empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					Remote:        "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Remote is not empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					Remote:        "https://gitlab.com/group/subgroup/project/-/raw/master/.gitlab-ci.yml",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(".gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("subgroup/project"),
					Organization: utils.GetPtr("group"),
				},
				Version:       utils.GetPtr("master"),
				VersionType:   models.BranchVersion,
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "Remote is an invalid url",
			args: args{
				item: gitlabCommon.IncludeItem{
					Remote:        "https://gitlab.com/group/subgroup/project/-/raw/master",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(""),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr(""),
					Organization: utils.GetPtr(""),
				},
				Version:       utils.GetPtr(""),
				VersionType:   models.None,
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "File is empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					File:          "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "File is not empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					File:          ".gitlab-ci.yml",
					Project:       "group/subgroup/project",
					Ref:           "master",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(".gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("subgroup/project"),
					Organization: utils.GetPtr("group"),
				},
				Version:       utils.GetPtr("master"),
				VersionType:   models.BranchVersion,
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "ref is empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					File:          ".gitlab-ci.yml",
					Project:       "group/subgroup/project",
					Ref:           "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr(".gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("subgroup/project"),
					Organization: utils.GetPtr("group"),
				},
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
		{
			name: "File is not empty and project is invalid",
			args: args{
				item: gitlabCommon.IncludeItem{
					File:          ".gitlab-ci.yml",
					Project:       "invalid",
					Ref:           "master",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Template is empty",
			args: args{
				item: gitlabCommon.IncludeItem{
					Template:      "",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: nil,
		},
		{
			name: "Regular template",
			args: args{
				item: gitlabCommon.IncludeItem{
					Template:      "Android.gitlab-ci.yml",
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
			},
			want: &models.Import{
				Source: &models.ImportSource{
					Type:         models.SourceTypeRemote,
					Path:         utils.GetPtr("lib/gitlab/ci/templates/Android.gitlab-ci.yml"),
					SCM:          consts.GitLabPlatform,
					Repository:   utils.GetPtr("gitlab"),
					Organization: utils.GetPtr("gitlab-org"),
				},
				Version:       utils.GetPtr("master"),
				VersionType:   models.BranchVersion,
				FileReference: testutils.CreateFileReference(1, 1, 1, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseIncludeItem(tt.args.item)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseImports(t *testing.T) {
	type args struct {
		include *gitlabCommon.Include
	}
	tests := []struct {
		name string
		args args
		want []*models.Import
	}{
		{
			name: "Include is nil",
			args: args{
				include: nil,
			},
			want: nil,
		},
		{
			name: "Include is empty",
			args: args{
				include: &gitlabCommon.Include{},
			},
			want: []*models.Import{},
		},
		{
			name: "Include all type of imports",
			args: args{
				include: &gitlabCommon.Include{
					{
						Local:         ".gitlab-ci.yml",
						FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					},
					{
						Remote:        "https://gitlab.com/group/subgroup/project/-/raw/f7b05602eb1bf22dc808838622c5d8d00b39d473/.gitlab-ci.yml",
						FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					},
					{
						File:          ".gitlab-ci.yml",
						Project:       "group/subgroup/project",
						Ref:           "v1",
						FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					},
					{
						Template:      "Android.gitlab-ci.yml",
						FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					},
				},
			},
			want: []*models.Import{
				{
					Source: &models.ImportSource{
						Type: models.SourceTypeLocal,
						Path: utils.GetPtr(".gitlab-ci.yml"),
						SCM:  consts.GitLabPlatform,
					},
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
				{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						Path:         utils.GetPtr(".gitlab-ci.yml"),
						SCM:          consts.GitLabPlatform,
						Repository:   utils.GetPtr("subgroup/project"),
						Organization: utils.GetPtr("group"),
					},
					Version:       utils.GetPtr("f7b05602eb1bf22dc808838622c5d8d00b39d473"),
					VersionType:   models.CommitSHA,
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
				},
				{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						Path:         utils.GetPtr(".gitlab-ci.yml"),
						SCM:          consts.GitLabPlatform,
						Repository:   utils.GetPtr("subgroup/project"),
						Organization: utils.GetPtr("group"),
					},
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					Version:       utils.GetPtr("v1"),
					VersionType:   models.TagVersion,
				},
				{
					Source: &models.ImportSource{
						Type:         models.SourceTypeRemote,
						Path:         utils.GetPtr("lib/gitlab/ci/templates/Android.gitlab-ci.yml"),
						SCM:          consts.GitLabPlatform,
						Repository:   utils.GetPtr("gitlab"),
						Organization: utils.GetPtr("gitlab-org"),
					},
					FileReference: testutils.CreateFileReference(1, 1, 1, 1),
					Version:       utils.GetPtr("master"),
					VersionType:   models.BranchVersion,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseImports(tt.args.include)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

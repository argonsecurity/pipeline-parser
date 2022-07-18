package azure

import (
	"testing"

	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/r3labs/diff/v3"
	"github.com/stretchr/testify/assert"
)

func TestParsePipelineTriggers(t *testing.T) {
	testCases := []struct {
		name            string
		pipeline        *azureModels.Pipeline
		expectedTrigger *models.Triggers
	}{
		{
			name:            "Pipeline is nil",
			pipeline:        nil,
			expectedTrigger: nil,
		},
		{
			name:            "Pipeline with no triggers",
			pipeline:        &azureModels.Pipeline{},
			expectedTrigger: nil,
		},
		{
			name: "Pipeline with trigger ref",
			pipeline: &azureModels.Pipeline{
				Trigger: &azureModels.TriggerRef{
					Trigger: &azureModels.Trigger{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
			expectedTrigger: &models.Triggers{
				Triggers: []*models.Trigger{
					{
						Event: models.PushEvent,
						Branches: &models.Filter{
							AllowList: []string{"master"},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Pipeline with pr ref",
			pipeline: &azureModels.Pipeline{
				PR: &azureModels.PRRef{
					PR: &azureModels.PR{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
			expectedTrigger: &models.Triggers{
				Triggers: []*models.Trigger{
					{
						Event: models.PullRequestEvent,
						Branches: &models.Filter{
							AllowList: []string{"master"},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Pipeline with schedules",
			pipeline: &azureModels.Pipeline{
				Schedules: &azureModels.Schedules{
					Crons: &[]azureModels.Cron{
						{
							Cron:          "1 * * * *",
							FileReference: testutils.CreateFileReference(1, 2, 3, 4),
						},
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
			expectedTrigger: &models.Triggers{
				Triggers: []*models.Trigger{
					{
						Event:         models.ScheduledEvent,
						Schedules:     &[]string{"1 * * * *"},
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
				},
				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
			},
		},
		{
			name: "Pipeline with all triggers",
			pipeline: &azureModels.Pipeline{
				Trigger: &azureModels.TriggerRef{
					Trigger: &azureModels.Trigger{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				PR: &azureModels.PRRef{
					PR: &azureModels.PR{
						Branches: azureModels.Filter{
							Include: []string{"master"},
						},
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
				Schedules: &azureModels.Schedules{
					Crons: &[]azureModels.Cron{
						{
							Cron:          "1 * * * *",
							FileReference: testutils.CreateFileReference(9, 10, 11, 12),
						},
					},
					FileReference: testutils.CreateFileReference(13, 14, 15, 16),
				},
			},
			expectedTrigger: &models.Triggers{
				Triggers: []*models.Trigger{
					{
						Event: models.PushEvent,
						Branches: &models.Filter{
							AllowList: []string{"master"},
						},
						FileReference: testutils.CreateFileReference(1, 2, 3, 4),
					},
					{
						Event: models.PullRequestEvent,
						Branches: &models.Filter{
							AllowList: []string{"master"},
						},
						FileReference: testutils.CreateFileReference(5, 6, 7, 8),
					},
					{
						Event:         models.ScheduledEvent,
						Schedules:     &[]string{"1 * * * *"},
						FileReference: testutils.CreateFileReference(13, 14, 15, 16),
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parsePipelineTriggers(testCase.pipeline)
			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

func TestParseTrigger(t *testing.T) {
	testCases := []struct {
		name            string
		ref             *azureModels.TriggerRef
		expectedTrigger *models.Trigger
	}{
		{
			name:            "Trigger ref is nil",
			ref:             nil,
			expectedTrigger: nil,
		},
		{
			name:            "Empty trigger ref",
			ref:             &azureModels.TriggerRef{},
			expectedTrigger: nil,
		},
		{
			name: "Trigger ref with branches include only",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Branches: azureModels.Filter{
						Include: []string{"master"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Branches: &models.Filter{
					AllowList: []string{"master"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Trigger ref with branches exclude only",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Branches: azureModels.Filter{
						Exclude: []string{"master", "develop"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Branches: &models.Filter{
					DenyList: []string{"master", "develop"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Trigger ref with paths include only",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Paths: azureModels.Filter{
						Include: []string{"/path1", "/path2"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Paths: &models.Filter{
					AllowList: []string{"/path1", "/path2"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Trigger ref with paths exclude only",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Paths: azureModels.Filter{
						Exclude: []string{"/path1", "/path2"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Paths: &models.Filter{
					DenyList: []string{"/path1", "/path2"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Trigger ref with tags include only",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Tags: azureModels.Filter{
						Include: []string{"tag1", "tag2"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Tags: &models.Filter{
					AllowList: []string{"tag1", "tag2"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Trigger ref with tags exclude only",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Tags: azureModels.Filter{
						Exclude: []string{"tag1", "tag2"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Tags: &models.Filter{
					DenyList: []string{"tag1", "tag2"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "Trigger ref with full data",
			ref: &azureModels.TriggerRef{
				Trigger: &azureModels.Trigger{
					Branches: azureModels.Filter{
						Include: []string{"master"},
						Exclude: []string{"develop"},
					},
					Paths: azureModels.Filter{
						Include: []string{"/path1", "/path2"},
						Exclude: []string{"/path3", "/path4"},
					},
					Tags: azureModels.Filter{
						Include: []string{"tag1", "tag2"},
						Exclude: []string{"tag3", "tag4"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
				Branches: &models.Filter{
					AllowList: []string{"master"},
					DenyList:  []string{"develop"},
				},
				Paths: &models.Filter{
					AllowList: []string{"/path1", "/path2"},
					DenyList:  []string{"/path3", "/path4"},
				},
				Tags: &models.Filter{
					AllowList: []string{"tag1", "tag2"},
					DenyList:  []string{"tag3", "tag4"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseTrigger(testCase.ref)

			changelog, err := diff.Diff(testCase.expectedTrigger, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParsePR(t *testing.T) {
	tests := []struct {
		name            string
		ref             *azureModels.PRRef
		expectedTrigger *models.Trigger
	}{
		{
			name:            "PR ref is nil",
			ref:             nil,
			expectedTrigger: nil,
		},
		{
			name:            "Empty PR ref",
			ref:             &azureModels.PRRef{},
			expectedTrigger: nil,
		},
		{
			name: "PR ref branches include only",
			ref: &azureModels.PRRef{
				PR: &azureModels.PR{
					Branches: azureModels.Filter{
						Include: []string{"master", "dev"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PullRequestEvent,
				Branches: &models.Filter{
					AllowList: []string{"master", "dev"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "PR ref branches exclude only",
			ref: &azureModels.PRRef{
				PR: &azureModels.PR{
					Branches: azureModels.Filter{
						Exclude: []string{"master", "dev"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PullRequestEvent,
				Branches: &models.Filter{
					DenyList: []string{"master", "dev"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "PR ref paths include only",
			ref: &azureModels.PRRef{
				PR: &azureModels.PR{
					Paths: azureModels.Filter{
						Include: []string{"/path/to/file", "/path/to/other/file"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PullRequestEvent,
				Paths: &models.Filter{
					AllowList: []string{"/path/to/file", "/path/to/other/file"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "PR ref paths exclude only",
			ref: &azureModels.PRRef{
				PR: &azureModels.PR{
					Paths: azureModels.Filter{
						Exclude: []string{"/path/to/file", "/path/to/other/file"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PullRequestEvent,
				Paths: &models.Filter{
					DenyList: []string{"/path/to/file", "/path/to/other/file"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		{
			name: "PR ref with full data",
			ref: &azureModels.PRRef{
				PR: &azureModels.PR{
					Branches: azureModels.Filter{
						Include: []string{"master", "dev"},
						Exclude: []string{"release"},
					},
					Paths: azureModels.Filter{
						Include: []string{"/path/to/file", "/path/to/other/file"},
						Exclude: []string{"/tmp/path"},
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event: models.PullRequestEvent,
				Branches: &models.Filter{
					AllowList: []string{"master", "dev"},
					DenyList:  []string{"release"},
				},
				Paths: &models.Filter{
					AllowList: []string{"/path/to/file", "/path/to/other/file"},
					DenyList:  []string{"/tmp/path"},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
		// {
		// 	name: "Full PR",
		// 	ref:  &azureModels.PRRef{},
		// 	expectedTrigger: &models.Trigger{
		// 		Event:         models.PullRequestEvent,
		// 		FileReference: testutils.CreateFileReference(1, 2, 3, 4),
		// 	},
		// },
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := parsePR(testCase.ref)

			changelog, err := diff.Diff(testCase.expectedTrigger, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

func TestParseSchedules(t *testing.T) {
	tests := []struct {
		name            string
		schedule        *azureModels.Schedules
		expectedTrigger *models.Trigger
	}{
		{
			name:            "Schedules is nil",
			schedule:        nil,
			expectedTrigger: nil,
		},
		{
			name:     "Empty schedule",
			schedule: &azureModels.Schedules{},
			expectedTrigger: &models.Trigger{
				Event:     models.ScheduledEvent,
				Schedules: utils.GetPtr([]string{}),
			},
		},
		{
			name: "Schedule with empty crons",
			schedule: &azureModels.Schedules{
				Crons: &[]azureModels.Cron{},
			},
			expectedTrigger: &models.Trigger{
				Event:     models.ScheduledEvent,
				Schedules: utils.GetPtr([]string{}),
			},
		},
		{
			name: "Full Schedule with crons",
			schedule: &azureModels.Schedules{
				Crons: &[]azureModels.Cron{
					{
						Cron: "1 * * *",
					},
					{
						Cron: "* * * 1",
					},
				},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event:         models.ScheduledEvent,
				Schedules:     utils.GetPtr([]string{"1 * * *", "* * * 1"}),
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseSchedules(testCase.schedule)

			changelog, err := diff.Diff(testCase.expectedTrigger, got)
			assert.NoError(t, err)
			assert.Len(t, changelog, 0, testCase.name)
		})
	}
}

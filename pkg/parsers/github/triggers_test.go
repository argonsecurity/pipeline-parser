package github

import (
	"testing"

	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/testutils"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var (
	mockInputs = githubModels.Inputs{
		"input1": struct {
			Description string
			Default     interface{}
			Required    bool
			Type        string
			Options     []string
		}{
			Description: "Description1",
			Default:     "Default Val string",
		},
		"input2": struct {
			Description string
			Default     interface{}
			Required    bool
			Type        string
			Options     []string
		}{
			Description: "Description2",
			Default:     1,
		},
		"input3": struct {
			Description string
			Default     interface{}
			Required    bool
			Type        string
			Options     []string
		}{
			Description: "Description3",
			Default:     true,
		},
	}

	mockExpectedParameters = []models.Parameter{
		{
			Name:        utils.GetPtr("input1"),
			Description: utils.GetPtr("Description1"),
			Default:     "Default Val string",
		},
		{
			Name:        utils.GetPtr("input2"),
			Description: utils.GetPtr("Description2"),
			Default:     1,
		},
		{
			Name:        utils.GetPtr("input3"),
			Description: utils.GetPtr("Description3"),
			Default:     true,
		},
	}
)

func TestParseWorkflowTriggers(t *testing.T) {
	testCases := []struct {
		name             string
		workflow         *githubModels.Workflow
		expectedTriggers *models.Triggers
	}{
		{
			name:             "Empty workflow",
			workflow:         &githubModels.Workflow{},
			expectedTriggers: nil,
		},
		// {
		// 	name: "Workflow with workflow.on with full data",
		// 	workflow: &githubModels.Workflow{
		// 		On: &githubModels.On{
		// 			Push: &githubModels.Ref{
		// 				Branches:       []string{"branch1", "branch2"},
		// 				BranchesIgnore: []string{"branch3", "branch4"},
		// 				Paths:          []string{"path1", "path2"},
		// 				PathsIgnore:    []string{"path3", "path4"},
		// 				FileReference:  testutils.CreateFileReference(1, 2, 3, 4),
		// 			},
		// 			PullRequest: &githubModels.Ref{
		// 				Branches:       []string{"branch1", "branch2"},
		// 				BranchesIgnore: []string{"branch3", "branch4"},
		// 				Paths:          []string{"path1", "path2"},
		// 				PathsIgnore:    []string{"path3", "path4"},
		// 				FileReference:  testutils.CreateFileReference(11, 21, 31, 41),
		// 			},
		// 			PullRequestTarget: &githubModels.Ref{
		// 				Branches:       []string{"branch1", "branch2"},
		// 				BranchesIgnore: []string{"branch3", "branch4"},
		// 				Paths:          []string{"path1", "path2"},
		// 				PathsIgnore:    []string{"path3", "path4"},
		// 				FileReference:  testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 			WorkflowDispatch: &githubModels.WorkflowDispatch{
		// 				Inputs:        mockInputs,
		// 				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
		// 			},
		// 			WorkflowCall: &githubModels.WorkflowCall{
		// 				Inputs:        mockInputs,
		// 				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 			WorkflowRun: &githubModels.WorkflowRun{
		// 				Workflows:     []string{"workflow1", "workflow2"},
		// 				FileReference: testutils.CreateFileReference(9, 10, 11, 12),
		// 				Ref: githubModels.Ref{
		// 					Branches:       []string{"branch1", "branch2"},
		// 					BranchesIgnore: []string{"branch3", "branch4"},
		// 				},
		// 			},
		// 			Schedule: &githubModels.Schedule{
		// 				Crons: &[]githubModels.Cron{
		// 					{
		// 						Cron: "1 * * *",
		// 					},
		// 				},
		// 				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 			Events: githubModels.Events{
		// 				"event1": githubModels.Event{
		// 					Types:         []string{"type1", "type2"},
		// 					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
		// 				},
		// 				"event2": githubModels.Event{
		// 					Types:         []string{"type3", "type4"},
		// 					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 				},
		// 			},
		// 			FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 		},
		// 	},
		// 	expectedTriggers: &models.Triggers{
		// 		Triggers: []*models.Trigger{
		// 			{
		// 				Event: models.PushEvent,
		// 				Branches: &models.Filter{
		// 					AllowList: []string{"branch1", "branch2"},
		// 					DenyList:  []string{"branch3", "branch4"},
		// 				},
		// 				Paths: &models.Filter{
		// 					AllowList: []string{"path1", "path2"},
		// 					DenyList:  []string{"path3", "path4"},
		// 				},
		// 				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
		// 			},
		// 			{
		// 				Event: models.PullRequestEvent,
		// 				Branches: &models.Filter{
		// 					AllowList: []string{"branch1", "branch2"},
		// 					DenyList:  []string{"branch3", "branch4"},
		// 				},
		// 				Paths: &models.Filter{
		// 					AllowList: []string{"path1", "path2"},
		// 					DenyList:  []string{"path3", "path4"},
		// 				},
		// 				FileReference: testutils.CreateFileReference(11, 21, 31, 41),
		// 			},
		// 			{
		// 				Event: models.EventType(pullRequestTargetEvent),
		// 				Branches: &models.Filter{
		// 					AllowList: []string{"branch1", "branch2"},
		// 					DenyList:  []string{"branch3", "branch4"},
		// 				},
		// 				Paths: &models.Filter{
		// 					AllowList: []string{"path1", "path2"},
		// 					DenyList:  []string{"path3", "path4"},
		// 				},
		// 				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 			{
		// 				Event:         models.ManualEvent,
		// 				Parameters:    mockExpectedParameters,
		// 				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
		// 			},
		// 			{
		// 				Event:         models.PipelineTriggerEvent,
		// 				Parameters:    mockExpectedParameters,
		// 				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 			{
		// 				Event:         models.PipelineRunEvent,
		// 				Pipelines:     []string{"workflow1", "workflow2"},
		// 				FileReference: testutils.CreateFileReference(9, 10, 11, 12),
		// 				Branches: &models.Filter{
		// 					AllowList: []string{"branch1", "branch2"},
		// 					DenyList:  []string{"branch3", "branch4"},
		// 				},
		// 			},
		// 			{
		// 				Event:         models.ScheduledEvent,
		// 				Schedules:     &[]string{"1 * * *"},
		// 				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 			{
		// 				Event: models.EventType("event1"),
		// 				Filters: map[string]any{
		// 					"types": []string{"type1", "type2"},
		// 				},
		// 				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
		// 			},
		// 			{
		// 				Event: models.EventType("event2"),
		// 				Filters: map[string]any{
		// 					"types": []string{"type3", "type4"},
		// 				},
		// 				FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 			},
		// 		},
		// 		FileReference: testutils.CreateFileReference(5, 6, 7, 8),
		// 	},
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			triggers := parseWorkflowTriggers(testCase.workflow)
			assert.Equal(t, testCase.expectedTriggers, triggers)
		})
	}
}

func TestParseEvents(t *testing.T) {
	testCases := []struct {
		name             string
		events           githubModels.Events
		expectedTriggers []*models.Trigger
	}{
		{
			name:             "Empty events map",
			events:           githubModels.Events{},
			expectedTriggers: []*models.Trigger{},
		},
		{
			name: "Events map with one event",
			events: githubModels.Events{
				"event1": githubModels.Event{
					Types:         []string{"type1", "type2"},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
			expectedTriggers: []*models.Trigger{
				{
					Event: models.EventType("event1"),
					Filters: map[string]any{
						"types": []string{"type1", "type2"},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
			},
		},
		{
			name: "Events map with some events",
			events: githubModels.Events{
				"event1": githubModels.Event{
					Types:         []string{"type1", "type2"},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				"event2": githubModels.Event{
					Types:         []string{"type3", "type4"},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
			expectedTriggers: []*models.Trigger{
				{
					Event: models.EventType("event1"),
					Filters: map[string]any{
						"types": []string{"type1", "type2"},
					},
					FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				},
				{
					Event: models.EventType("event2"),
					Filters: map[string]any{
						"types": []string{"type3", "type4"},
					},
					FileReference: testutils.CreateFileReference(5, 6, 7, 8),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseEvents(testCase.events)

			assert.Equal(t, testCase.expectedTriggers, got, testCase.name)
		})
	}
}

func TestParseWorkflowRun(t *testing.T) {
	testCases := []struct {
		name            string
		workflowRun     *githubModels.WorkflowRun
		expectedTrigger *models.Trigger
	}{
		{
			name:        "Empty workflow run",
			workflowRun: &githubModels.WorkflowRun{},
			expectedTrigger: &models.Trigger{
				Event:    models.PipelineRunEvent,
				Branches: &models.Filter{},
			},
		},
		{
			name: "Workflow run with data",
			workflowRun: &githubModels.WorkflowRun{
				Workflows:     []string{"workflow1", "workflow2"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				Ref: githubModels.Ref{
					Branches:       []string{"branch 1", "branch 2"},
					BranchesIgnore: []string{"branch 3", "branch 4"},
				},
			},
			expectedTrigger: &models.Trigger{
				Event:         models.PipelineRunEvent,
				Pipelines:     []string{"workflow1", "workflow2"},
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
				Branches: &models.Filter{
					AllowList: []string{"branch 1", "branch 2"},
					DenyList:  []string{"branch 3", "branch 4"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseWorkflowRun(testCase.workflowRun)

			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

func TestParseWorkflowCall(t *testing.T) {
	testCases := []struct {
		name            string
		workflowCall    *githubModels.WorkflowCall
		expectedTrigger *models.Trigger
	}{
		{
			name:         "Empty workflow call",
			workflowCall: &githubModels.WorkflowCall{},
			expectedTrigger: &models.Trigger{
				Event:      models.PipelineTriggerEvent,
				Parameters: []models.Parameter{},
			},
		},
		{
			name: "Workflow call with data",
			workflowCall: &githubModels.WorkflowCall{
				Inputs:        mockInputs,
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event:         models.PipelineTriggerEvent,
				Parameters:    mockExpectedParameters,
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseWorkflowCall(testCase.workflowCall)

			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

func TestParseWorkflowDispatch(t *testing.T) {
	testCases := []struct {
		name             string
		workflowDispatch *githubModels.WorkflowDispatch
		expectedTrigger  *models.Trigger
	}{
		{
			name:             "Empty workflow dispatch",
			workflowDispatch: &githubModels.WorkflowDispatch{},
			expectedTrigger: &models.Trigger{
				Event: models.ManualEvent,
			},
		},
		{
			name: "Workflow dispatch with data",
			workflowDispatch: &githubModels.WorkflowDispatch{
				Inputs:        mockInputs,
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
			expectedTrigger: &models.Trigger{
				Event:         models.ManualEvent,
				Parameters:    mockExpectedParameters,
				FileReference: testutils.CreateFileReference(1, 2, 3, 4),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseWorkflowDispatch(testCase.workflowDispatch)

			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

func TestParseInputs(t *testing.T) {
	testCases := []struct {
		name               string
		inputs             githubModels.Inputs
		expectedParameters []models.Parameter
	}{
		{
			name:               "nil inputs",
			inputs:             nil,
			expectedParameters: []models.Parameter{},
		},
		{
			name:               "Empty Inputs",
			inputs:             githubModels.Inputs{},
			expectedParameters: []models.Parameter{},
		},
		{
			name:               "Inputs with data",
			inputs:             mockInputs,
			expectedParameters: mockExpectedParameters,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseInputs(testCase.inputs)

			assert.ElementsMatch(t, testCase.expectedParameters, got, testCase.name)
		})
	}
}

func TestParseRef(t *testing.T) {
	testCases := []struct {
		name            string
		ref             *githubModels.Ref
		event           models.EventType
		expectedTrigger *models.Trigger
	}{
		{
			name:            "Empty ref with no event",
			ref:             &githubModels.Ref{},
			event:           models.EventType(""),
			expectedTrigger: &models.Trigger{},
		},
		{
			name:  "Empty ref with event",
			ref:   &githubModels.Ref{},
			event: models.EventType("event"),
			expectedTrigger: &models.Trigger{
				Event: models.EventType("event"),
			},
		},
		{
			name: "Ref with branches only",
			ref: &githubModels.Ref{
				Branches: []string{"branch 1", "branch 2"},
			},
			event: models.EventType("event"),
			expectedTrigger: &models.Trigger{
				Event: models.EventType("event"),
				Branches: &models.Filter{
					AllowList: []string{"branch 1", "branch 2"},
				},
			},
		},
		{
			name: "Ref with branches ignore only",
			ref: &githubModels.Ref{
				BranchesIgnore: []string{"branch 1", "branch 2"},
			},
			event: models.EventType("event"),
			expectedTrigger: &models.Trigger{
				Event: models.EventType("event"),
				Branches: &models.Filter{
					DenyList: []string{"branch 1", "branch 2"},
				},
			},
		},
		{
			name: "Ref with paths only",
			ref: &githubModels.Ref{
				Paths: []string{"path 1", "path 2"},
			},
			event: models.EventType("event"),
			expectedTrigger: &models.Trigger{
				Event: models.EventType("event"),
				Paths: &models.Filter{
					AllowList: []string{"path 1", "path 2"},
				},
			},
		},
		{
			name: "Ref with paths ignore only",
			ref: &githubModels.Ref{
				PathsIgnore: []string{"path 1", "path 2"},
			},
			event: models.EventType("event"),
			expectedTrigger: &models.Trigger{
				Event: models.EventType("event"),
				Paths: &models.Filter{
					DenyList: []string{"path 1", "path 2"},
				},
			},
		},
		{
			name: "Ref with full data",
			ref: &githubModels.Ref{
				Branches:       []string{"branch 1", "branch 2"},
				BranchesIgnore: []string{"branch 3", "branch 4"},
				Paths:          []string{"path 1", "path 2"},
				PathsIgnore:    []string{"path 3", "path 4"},
			},
			event: models.EventType("event"),
			expectedTrigger: &models.Trigger{
				Event: models.EventType("event"),
				Branches: &models.Filter{
					AllowList: []string{"branch 1", "branch 2"},
					DenyList:  []string{"branch 3", "branch 4"},
				},
				Paths: &models.Filter{
					AllowList: []string{"path 1", "path 2"},
					DenyList:  []string{"path 3", "path 4"},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseRef(testCase.ref, testCase.event)

			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

func TestParseSchedule(t *testing.T) {
	testCases := []struct {
		name            string
		schedule        *githubModels.Schedule
		expectedTrigger *models.Trigger
	}{
		{
			name:     "Empty schedule",
			schedule: &githubModels.Schedule{},
			expectedTrigger: &models.Trigger{
				Event:     models.ScheduledEvent,
				Schedules: utils.GetPtr([]string{}),
			},
		},
		{
			name: "Schedule with empty crons",
			schedule: &githubModels.Schedule{
				Crons: &[]githubModels.Cron{},
			},
			expectedTrigger: &models.Trigger{
				Event:     models.ScheduledEvent,
				Schedules: utils.GetPtr([]string{}),
			},
		},
		{
			name: "Full Schedule with crons",
			schedule: &githubModels.Schedule{
				Crons: &[]githubModels.Cron{
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseSchedule(testCase.schedule)

			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

func TestGenerateTriggersFromEvents(t *testing.T) {
	testCases := []struct {
		name             string
		events           []string
		expectedTriggers []*models.Trigger
	}{
		{
			name:   "Existing events in event map",
			events: []string{pushEvent, pullRequestEvent},
			expectedTriggers: []*models.Trigger{
				{
					Event: models.PushEvent,
				},
				{
					Event: models.PullRequestEvent,
				},
			},
		}, {
			name:   "Non existing events in event map",
			events: []string{"event1", "event2"},
			expectedTriggers: []*models.Trigger{
				{
					Event: models.EventType("event1"),
				},
				{
					Event: models.EventType("event2"),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateTriggersFromEvents(testCase.events)

			assert.Equal(t, testCase.expectedTriggers, got, testCase.name)
		})
	}
}

func TestGenerateTriggerFromEvent(t *testing.T) {
	testCases := []struct {
		name            string
		event           string
		expectedTrigger *models.Trigger
	}{
		{
			name:  "Existing event in event map",
			event: pushEvent,
			expectedTrigger: &models.Trigger{
				Event: models.PushEvent,
			},
		},
		{
			name:  "Non existing event in event map",
			event: "some_event",
			expectedTrigger: &models.Trigger{
				Event: models.EventType("some_event"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := generateTriggerFromEvent(testCase.event)

			assert.Equal(t, testCase.expectedTrigger, got, testCase.name)
		})
	}
}

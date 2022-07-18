package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parsePipelineTriggers(pipeline *azureModels.Pipeline) *models.Triggers {
	var triggers []*models.Trigger
	if pipeline.Trigger != nil {
		triggers = append(triggers, parseTrigger(pipeline.Trigger))
	}

	if pipeline.PR != nil {
		triggers = append(triggers, parsePR(pipeline.PR))
	}

	if pipeline.Schedules != nil {
		triggers = append(triggers, parseSchedules(pipeline.Schedules))
	}

	if len(triggers) == 0 {
		return nil
	}

	return &models.Triggers{
		Triggers:      triggers,
		FileReference: triggers[0].FileReference,
	}
}

func parseTrigger(ref *azureModels.TriggerRef) *models.Trigger {
	trigger := &models.Trigger{
		Event:         models.PushEvent,
		FileReference: ref.FileReference,
	}

	if len(ref.Trigger.Branches.Include)+len(ref.Trigger.Branches.Exclude) > 0 {
		trigger.Branches = &models.Filter{}
		trigger.Branches.AllowList = append(trigger.Branches.AllowList, ref.Trigger.Branches.Include...)
		trigger.Branches.DenyList = append(trigger.Branches.DenyList, ref.Trigger.Branches.Exclude...)
	}

	if len(ref.Trigger.Paths.Include)+len(ref.Trigger.Paths.Exclude) > 0 {
		trigger.Branches = &models.Filter{}
		trigger.Paths.AllowList = append(trigger.Paths.AllowList, ref.Trigger.Paths.Include...)
		trigger.Paths.DenyList = append(trigger.Paths.AllowList, ref.Trigger.Paths.Exclude...)
	}

	if len(ref.Trigger.Tags.Include)+len(ref.Trigger.Tags.Exclude) > 0 {
		trigger.Branches = &models.Filter{}
		trigger.Tags.AllowList = append(trigger.Tags.AllowList, ref.Trigger.Tags.Include...)
		trigger.Tags.DenyList = append(trigger.Tags.DenyList, ref.Trigger.Tags.Exclude...)
	}

	return trigger
}

func parsePR(ref *azureModels.PRRef) *models.Trigger {
	trigger := &models.Trigger{
		Event:         models.PullRequestEvent,
		FileReference: ref.FileReference,
	}

	if len(ref.PR.Branches.Include)+len(ref.PR.Branches.Exclude) > 0 {
		trigger.Branches = &models.Filter{}
		trigger.Branches.AllowList = append(trigger.Branches.AllowList, ref.PR.Branches.Include...)
		trigger.Branches.DenyList = append(trigger.Branches.DenyList, ref.PR.Branches.Exclude...)
	}

	if len(ref.PR.Paths.Include)+len(ref.PR.Paths.Exclude) > 0 {
		trigger.Branches = &models.Filter{}
		trigger.Paths.AllowList = append(trigger.Paths.AllowList, ref.PR.Paths.Include...)
		trigger.Paths.DenyList = append(trigger.Paths.AllowList, ref.PR.Paths.Exclude...)
	}

	return trigger
}

func parseSchedules(schedule *azureModels.Schedules) *models.Trigger {
	schedules := []string{}
	if schedule.Crons != nil {
		schedules = utils.Map(*schedule.Crons, func(cron azureModels.Cron) string {
			return cron.Cron
		})
	}

	return &models.Trigger{
		Event:         models.ScheduledEvent,
		Schedules:     &schedules,
		FileReference: schedule.FileReference,
	}
}

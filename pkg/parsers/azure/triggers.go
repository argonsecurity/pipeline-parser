package azure

import (
	azureModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/azure/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

func parsePipelineTriggers(pipeline *azureModels.Pipeline) *models.Triggers {
	if pipeline == nil {
		return nil
	}

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

	triggers = utils.Filter(triggers, func(t *models.Trigger) bool { return t != nil })

	if len(triggers) == 0 {
		return nil
	}

	return &models.Triggers{
		Triggers: triggers,
		FileReference: &models.FileReference{
			StartRef: triggers[0].FileReference.StartRef,
			EndRef:   triggers[len(triggers)-1].FileReference.EndRef,
		},
	}
}

func parseTrigger(ref *azureModels.TriggerRef) *models.Trigger {
	if ref == nil || ref.Trigger == nil {
		return nil
	}

	trigger := &models.Trigger{
		Event:         models.PushEvent,
		FileReference: ref.FileReference,
	}

	if len(ref.Trigger.Branches.Include)+len(ref.Trigger.Branches.Exclude) > 0 {
		trigger.Branches = &models.Filter{
			AllowList: ref.Trigger.Branches.Include,
			DenyList:  ref.Trigger.Branches.Exclude,
		}
	}

	if len(ref.Trigger.Paths.Include)+len(ref.Trigger.Paths.Exclude) > 0 {
		trigger.Paths = &models.Filter{
			AllowList: ref.Trigger.Paths.Include,
			DenyList:  ref.Trigger.Paths.Exclude,
		}
	}

	if len(ref.Trigger.Tags.Include)+len(ref.Trigger.Tags.Exclude) > 0 {
		trigger.Tags = &models.Filter{
			AllowList: ref.Trigger.Tags.Include,
			DenyList:  ref.Trigger.Tags.Exclude,
		}
	}

	return trigger
}

func parsePR(ref *azureModels.PRRef) *models.Trigger {
	if ref == nil || ref.PR == nil {
		return nil
	}

	trigger := &models.Trigger{
		Event:         models.PullRequestEvent,
		FileReference: ref.FileReference,
	}

	if len(ref.PR.Branches.Include)+len(ref.PR.Branches.Exclude) > 0 {
		trigger.Branches = &models.Filter{
			AllowList: ref.PR.Branches.Include,
			DenyList:  ref.PR.Branches.Exclude,
		}
	}

	if len(ref.PR.Paths.Include)+len(ref.PR.Paths.Exclude) > 0 {
		trigger.Paths = &models.Filter{
			AllowList: ref.PR.Paths.Include,
			DenyList:  ref.PR.Paths.Exclude,
		}
	}

	return trigger
}

func parseSchedules(schedule *azureModels.Schedules) *models.Trigger {
	if schedule == nil {
		return nil
	}
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

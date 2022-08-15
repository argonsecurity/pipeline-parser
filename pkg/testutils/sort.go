package testutils

import (
	"sort"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
)

func SortParameters(params []models.Parameter) {
	sort.Slice(params, func(i, j int) bool {
		return *params[i].Name < *params[j].Name
	})
}

func SortTrigger(trigger *models.Trigger) {
	sort.Slice(trigger.Parameters, func(i, j int) bool {
		return *trigger.Parameters[i].Name < *trigger.Parameters[j].Name
	})
}

func SortTriggers(triggers []*models.Trigger) {
	sort.Slice(triggers, func(i, j int) bool {
		return triggers[i].Event < triggers[j].Event
	})

	SortMany(triggers, SortTrigger)
}

func SortJobs(jobs []*models.Job) {
	sort.Slice(jobs, func(i, j int) bool {
		return *jobs[i].ID < *jobs[j].ID
	})
}

func SortMany[T any](s []T, sortFunc func(T)) {
	for _, item := range s {
		sortFunc(item)
	}
}

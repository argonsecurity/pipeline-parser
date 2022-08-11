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
	SortMany(triggers, SortTrigger)
}

func SortMany[T any](s []T, sortFunc func(T)) {
	for _, item := range s {
		sortFunc(item)
	}
}

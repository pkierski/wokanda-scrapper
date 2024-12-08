package trialdownloader

import (
	"cmp"
	"slices"
	"strings"
)

func SortTrials(trials []Trial) []Trial {
	slices.SortFunc(trials, func(a, b Trial) int {
		return cmp.Or(
			strings.Compare(a.CaseID, b.CaseID),
			a.Date.Compare(b.Date),
		)
	})
	return trials
}

package trialdownloader

import (
	"cmp"
	"slices"
	"strings"
)

func (t Trial) Compare(other Trial) int {
	return cmp.Or(
		strings.Compare(t.CaseID, other.CaseID),
		t.Date.Compare(other.Date),
		slices.Compare(t.Judges, other.Judges),
		strings.Compare(t.Room, other.Room),
		strings.Compare(t.Department, other.Department),
	)
}

func SortTrials(trials []Trial) []Trial {
	slices.SortFunc(trials, func(a, b Trial) int {
		return a.Compare(b)
	})
	return trials
}

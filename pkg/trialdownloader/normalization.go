package trialdownloader

import (
	"regexp"
	"slices"
	"strings"
)

func normalizeTrials(trials []Trial) []Trial {
	for t := range MutableValues(trials) {
		normalizeTrial(t)
	}

	// sort and compact trials
	// some courts add the same entry many times
	trials = slices.CompactFunc(SortTrials(trials), func(a, b Trial) bool {
		return a.Compare(b) == 0
	})

	return trials
}

func normalizeTrial(t *Trial) {
	splitJudges(t)
	for j := range MutableValues(t.Judges) {
		*j = normalizeJudgeName(*j)
	}
}

func splitJudges(t *Trial) {
	replacement := []string{}
	for _, j := range t.Judges {
		replacement = append(replacement, strings.Split(j, ",")...)
	}
	t.Judges = replacement
}

var (
	reHypenWithSpaces = regexp.MustCompile(` ?- ?`)
	reRemovePrefixes  = regexp.MustCompile("^ *(([Ss]ędzia)|(([Aa]sesor|[Rr]eferendarz)( +sądowy)?)|(SS[AOR]))")
	reRemoveSuffixes  = regexp.MustCompile("- *(([Pp]rzewodniczący)|([Ss]ędzia)([Łł]awnik))$")
)

func normalizeJudgeName(j string) string {
	j = reHypenWithSpaces.ReplaceAllString(j, "-")
	j = reRemovePrefixes.ReplaceAllString(j, "")
	j = reRemoveSuffixes.ReplaceAllString(j, "")
	return strings.TrimSpace(j)
}

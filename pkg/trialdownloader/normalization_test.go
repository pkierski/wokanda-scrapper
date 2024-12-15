package trialdownloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalization(t *testing.T) {
	trials := []Trial{
		{
			Judges: []string{
				"sędzia first  ",
				"Asesor second  ",
			},
		},
		{
			Judges: []string{
				" second, thrid with -fancy- name - here  ",
			},
		},
	}
	normalizeTrials(trials)

	assert.Equal(t, trials, []Trial{
		{
			Judges: []string{
				"first",
				"second",
			},
		},
		{
			Judges: []string{
				"second",
				"thrid with-fancy-name-here",
			},
		},
	})
}

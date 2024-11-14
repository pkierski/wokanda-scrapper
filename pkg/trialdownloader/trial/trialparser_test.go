package trial

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:embed "V1.html"
var v1 []byte

//go:embed "V2-bialystok.html"
var v2 []byte

func TestParseV1(t *testing.T) {
	trial, err := ParseV1(v1)

	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02 15:04:05", "2024-11-13 09:00:00")
	assert.NoError(t, err)
	expectedTime = localizeTime(expectedTime)

	assert.Equal(t,
		Trial{
			CaseID:     "I C 268/24",
			Department: "I Wydział Cywilny",
			Judge:      "SSR Krzysztof Świeczkowski",
			Date:       expectedTime,
			Room:       "4",
		},
		trial,
	)
}

func TestParseV2(t *testing.T) {
	trials, err := ParseV2(v2)
	assert.NoError(t, err)
	assert.Len(t, trials, 388)

	expectedTime, err := time.Parse("2006-01-02 15:04:05", "2024-11-28 14:30:00")
	assert.NoError(t, err)
	expectedTime = localizeTime(expectedTime)

	assert.Equal(t,
		Trial{
			CaseID:     "I ACa 998/23",
			Department: "I Wydział Cywilny",
			Judge:      "SSA Alicja Dubij (sprawozdawca)",
			Date:       expectedTime,
			Room:       "Sala II",
		},
		trials[387], // the last one
	)

}

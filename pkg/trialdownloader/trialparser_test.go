package trialdownloader

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

//go:embed "V2-bialystok-small.html"
var v2small []byte

func TestParseV1(t *testing.T) {
	trial, err := parseV1DetailPage(v1)

	assert.NoError(t, err)

	expectedTime, err := time.Parse("2006-01-02 15:04:05", "2024-11-13 09:00:00")
	assert.NoError(t, err)
	expectedTime = localizeTime(expectedTime)

	assert.Equal(t,
		Trial{
			CaseID:     "I C 268/24",
			Department: "I Wydział Cywilny",
			Judges:     []string{"SSR Krzysztof Świeczkowski"},
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
			Date:       expectedTime,
			Room:       "Sala II",
			Judges:     []string{"SSA Alicja Dubij (sprawozdawca)"},
		},
		trials[387], // the last one
	)

}

func timeMustParse(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return localizeTime(t)
}

func TestParseV2small(t *testing.T) {
	trials, err := ParseV2(v2small)
	assert.NoError(t, err)

	assert.NoError(t, err)

	assert.Equal(t,
		Trial{
			CaseID:     "II AKa 92/24",
			Department: "II Wydział Karny",
			Date:       timeMustParse("2024-11-13 11:00:00"),
			Room:       "Sala IV",
			Judges:     []string{"SSA Brandeta Hryniewicka (sprawozdawca)", "SSA Tomasz Uściłko", "SSA Jacek Dunikowski"},
		},
		trials[0], // the last one
	)

}

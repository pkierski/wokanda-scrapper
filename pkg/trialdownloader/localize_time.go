package trialdownloader

import (
	"fmt"
	"time"
)

var (
	warsawTime = func() *time.Location {
		l, err := time.LoadLocation("Europe/Warsaw")
		if err != nil {
			panic(err)
		}
		return l
	}()

	warsawTimeOffset = func() time.Duration {
		t := time.Now().UTC()
		_, offset := t.In(warsawTime).Zone()
		return time.Duration(offset) * time.Second
	}()
)

func parseAndLocalizeTime(dateStr, timeStr string) (time.Time, error) {
	if len(timeStr) == 5 {
		return parseAndLocalizeTimeWithTimeLayout(dateStr, timeStr, "15:04")
	}
	if len(timeStr) == 8 {
		return parseAndLocalizeTimeWithTimeLayout(dateStr, timeStr, "15:04:05")
	}
	return time.Time{}, fmt.Errorf("can't parse date and time: '%v', '%v'", dateStr, timeStr)
}

func parseAndLocalizeTimeWithTimeLayout(dateStr, timeStr string, timeLayout string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 "+timeLayout, dateStr+" "+timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return localizeTime(t), nil
}

func localizeTime(t time.Time) time.Time {
	return t.In(warsawTime).Add(-warsawTimeOffset)
}

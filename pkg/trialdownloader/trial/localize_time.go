package trial

import "time"

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

func parseAndLocalizeTime(dateStr, timeStr string, timeLayout string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 "+timeLayout, dateStr+" "+timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return localizeTime(t), nil
}

func localizeTime(t time.Time) time.Time {
	return t.In(warsawTime).Add(-warsawTimeOffset)
}

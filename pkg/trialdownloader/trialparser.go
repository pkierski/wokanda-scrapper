package trialdownloader

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var ErrNoDataOnPage = errors.New("can't find trial data")

// ParseV1 parses one page from type "<url>/wokanda,N".
func ParseV1(data []byte) (trial Trial, err error) {
	if !bytes.Contains(data, []byte(`<dl class="dl-horizontal case-description-list">`)) {
		err = fmt.Errorf("parsing trial page: %w", ErrNoDataOnPage)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}

	var (
		dateStr string
		timeStr string
	)

	s := doc.Selection.Find("dl[class='dl-horizontal case-description-list']")
	dts := s.Find("dt")
	dds := s.Find("dd")
	for i, dt := range dts.EachIter() {
		header := dt.Text()
		val := dds.Eq(i).Text()

		switch {
		case strings.Contains(header, "Sygnatura"):
			trial.CaseID = val

		case strings.Contains(header, "Wydział"):
			trial.Department = val

		case strings.Contains(header, "Godzina"):
			timeStr = val

		case strings.Contains(header, "Data"):
			dateStr = val

		case strings.Contains(header, "Sala"):
			trial.Room = val

		case strings.Contains(header, "Przewodniczący"):
			trial.Judges = append(trial.Judges, val)
		}

	}

	trial.Date, err = parseAndLocalizeTime(dateStr, timeStr, "15:04:05")

	return
}

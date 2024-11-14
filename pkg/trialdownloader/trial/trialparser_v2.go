package trial

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ParseV1 parses one page from type pages like
// "https://bialystok.sa.gov.pl/zalatw-sprawe/e-wokanda".
func ParseV2(data []byte) (trials []Trial, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}

	// base information (revealed part of table cell)
	baseInfos := doc.Selection.Find("tr[class='category table table-striped table-bordered table-hover']")
	baseInfos.Each(func(i int, s *goquery.Selection) {
		baseVals := s.Find("span.strong")
		bv := baseVals.Map(func(i int, s *goquery.Selection) string {
			return s.Text()
		})
		var dateTime time.Time
		dateTime, err = parseAndLocalizeTime(bv[2], bv[4], "15:04")
		trial := Trial{
			CaseID:     bv[0],
			Department: bv[1],
			Date:       dateTime,
			Room:       bv[3],
		}

		trials = append(trials, trial)
	})

	if err != nil {
		return nil, err
	}

	// information hidded under "więcej..." button
	additionalInfo := doc.Selection.Find("tr[class='row_sklad category table table-striped table-bordered table-hover']").Filter(":contains('Przewod')")
	if additionalInfo.Length() != baseInfos.Length() {
		return nil, fmt.Errorf("parsing trial page: base info and additional info length mismatch")
	}

	// update trials with additional info: the first elemnent is the judge name
	additionalInfo.Each(func(i int, s *goquery.Selection) {
		trial := trials[i]
		trial.Judge = s.First().Find("td.strong").Text()
		trials[i] = trial
	})

	return
}

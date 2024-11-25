package trialdownloader

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/pageparser"
	"golang.org/x/net/html"
)

// ParseV1 parses one page from type pages like
// "https://bialystok.sa.gov.pl/zalatw-sprawe/e-wokanda".
func ParseV2(data []byte) (trials []Trial, err error) {
	if !bytes.Contains(data, []byte(`<form action="/zalatw-sprawe/e-wokanda" method="post">`)) {
		return nil, fmt.Errorf("parsing trial page: %w", ErrNoDataOnPage)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parsing trial page: %w", err)
	}

	// base information (revealed part of table cell)
	baseInfos := doc.Selection.Find("tr[class='category table table-striped table-bordered table-hover']")
	for i, s := range baseInfos.EachIter() {
		var (
			strDate, strTime string
			trial            Trial
		)
		baseVals := s.Find("span.strong")
		for _, s := range baseVals.EachIter() {
			keyNode := s.Nodes[0].PrevSibling
			if keyNode == nil {
				continue
			}

			switch {
			case strings.Contains(keyNode.Data, "Sygnatura"):
				trial.CaseID = s.Text()
			case strings.Contains(keyNode.Data, "Wydział"):
				trial.Department = s.Text()
			case strings.Contains(keyNode.Data, "Data"):
				strDate = s.Text()
			case strings.Contains(keyNode.Data, "Sala"):
				trial.Room = s.Text()
			case strings.Contains(keyNode.Data, "Godzina"):
				strTime = s.Text()
			}
		}

		trial.Date, err = parseAndLocalizeTime(strDate, strTime, "15:04")
		if err != nil {
			return nil, fmt.Errorf("parsing trial page (case: %v): %w", i, err)
		}

		trials = append(trials, trial)
	}

	// information hidded under "więcej..." button
	additionalInfo := doc.Selection.Find("tr[class='row_sklad category table table-striped table-bordered table-hover']").Filter(":contains('Przewod')")
	if additionalInfo.Length() != baseInfos.Length() {
		return nil, fmt.Errorf("parsing trial page: base info and additional info length mismatch")
	}

	// update trials with additional info: the judge names
	for i, s := range additionalInfo.EachIter() {
		trial := trials[i]
		trial.Judges = make([]string, 0)

		judgeTable := s.Find("table[class='sklad_sedziowski']").Find("td.strong")
		for _, s := range judgeTable.EachIter() {
			judge := strings.TrimSpace(s.Text())
			if !strings.Contains(judge, "nie ustalono") {
				trial.Judges = append(trial.Judges, judge)
			}
		}

		judgesTable := s.Find("table[class='sklad_sedziowski category table table-striped table-bordered table-hover']").Find("td.strong")
		if judgesTable.Length() > 0 {
			pageparser.WalkNodes(judgesTable.Nodes[0], func(node *html.Node) {
				if node.DataAtom == 0 && !strings.Contains(node.Data, "nie ustalono") {
					trial.Judges = append(trial.Judges, strings.TrimSpace(node.Data))
				}
			})
		}

		trials[i] = trial
	}

	return trials, nil
}

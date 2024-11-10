package trial

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/pageparser"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var ErrNoDataOnPage = errors.New("can't find trial data")

var warsawTime = func() *time.Location {
	l, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		panic(err)
	}
	return l
}()

func Parse(data []byte) (trial Trial, err error) {
	if !bytes.Contains(data, []byte(`<dl class="dl-horizontal case-description-list">`)) {
		err = fmt.Errorf("parsing trial page: %w", ErrNoDataOnPage)
		return
	}

	node, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}

	var (
		dateStr string
		timeStr string
	)

	pageparser.WalkNodes(node, func(node *html.Node) {
		if node.DataAtom != atom.Dl || !strings.Contains(pageparser.FindAttrValue(node, "class"), "case-description-list") {
			return
		}

		var lastTerm string

		for d := node.FirstChild; d != nil; d = d.NextSibling {
			if d.FirstChild == nil {
				continue
			}

			if d.DataAtom == atom.Dt {
				lastTerm = d.FirstChild.Data
			}

			if d.DataAtom == atom.Dd {
				switch {
				case strings.Contains(lastTerm, "Sygnatura"):
					trial.CaseID = d.FirstChild.Data

				case strings.Contains(lastTerm, "Wydział"):
					trial.Department = d.FirstChild.Data

				case strings.Contains(lastTerm, "Godzina"):
					timeStr = d.FirstChild.Data

				case strings.Contains(lastTerm, "Data"):
					dateStr = d.FirstChild.Data

				case strings.Contains(lastTerm, "Sala"):
					trial.Room = d.FirstChild.Data

				case strings.Contains(lastTerm, "Przewodniczący"):
					trial.Judge = d.FirstChild.Data
				}

			}
		}
	})

	trial.Date, err = time.Parse("2006-01-02 15:04:05", dateStr+" "+timeStr)
	trial.Date = trial.Date.In(warsawTime)
	_, offset := trial.Date.Zone()
	trial.Date = trial.Date.Add(-time.Duration(offset) * time.Second)

	return
}

package trialdownloader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	maxPerPageV3 = 10
)

type V3Wokanda commonDownloader

// check if V3Wokanda implements Downloader
var _ Downloader = (*V3Wokanda)(nil)

func NewV3Wokanda(client *http.Client, baseUrl string) V3Wokanda {
	return V3Wokanda{
		client:  client,
		baseUrl: "https://" + baseUrl,
	}
}

func (d V3Wokanda) Download(ctx context.Context, date string) ([]Trial, error) {
	// get first page
	trials, pages, err := getOnePageV3(ctx, d.client, d.baseUrl, date, 1, maxPerPageV3)
	if err != nil {
		return nil, err
	}
	if pages == 0 {
		return trials, nil
	}

	// TODO: make it concurrent
	for page := 2; page <= pages; page++ {
		trialsPage, _, err := getOnePageV3(ctx, d.client, d.baseUrl, date, page, maxPerPageV3)
		if err != nil {
			return nil, err
		}
		trials = append(trials, trialsPage...)
	}

	return trials, nil
}

func getOnePageV3(ctx context.Context, client *http.Client, url string, date string, page int, perPage int) (trials []Trial, pages int, err error) {
	query := fmt.Sprintf("/e-wokanda/szukaj?start_time=%v&perPage=%v&page=%v", date, perPage, page)
	pageContent, err := postOne(ctx, client, url+query, "")
	if err != nil {
		return nil, 0, fmt.Errorf("downloading trial page: %w", err)
	}

	os.WriteFile("page.html", pageContent, 0o644)

	// TODO: check if format is correct for V3

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageContent))
	if err != nil {
		return nil, 0, fmt.Errorf("parsing trial page: %w", err)
	}

	pages, err = parseV3Pagination(doc)
	if err != nil {
		return nil, 0, fmt.Errorf("parsing trial page: %w", err)
	}

	trials, err = parseV3PageTrials(doc)
	if err != nil {
		return nil, 0, fmt.Errorf("parsing trial page: %w", err)
	}

	trials = normalizeTrials(trials)
	return trials, pages, nil
}

func parseV3PageTrials(doc *goquery.Document) (trials []Trial, err error) {
	tables := doc.Find("table[class='table table-borderless']")
	trials = make([]Trial, 0, tables.Length())
	for _, table := range tables.EachIter() {

		var (
			trial            Trial
			dateStr, timeStr string
			judgesStr        string
		)
		for _, row := range table.Find("tr").EachIter() {
			header := row.Find("th")
			headerStr := strings.TrimSpace(header.Text())
			val := row.Find("td")
			valStr := strings.TrimSpace(val.Text())

			// fmt.Println(header.Text(), val.Text())
			fmt.Printf("'%v': '%v'\n", headerStr, valStr)

			switch {
			case strings.Contains(headerStr, "Sygnatura"):
				trial.CaseID = valStr
			case strings.Contains(headerStr, "Wydział"):
				trial.Department = valStr
			case strings.Contains(headerStr, "Miejsce"):
				trial.Room = valStr
			// requires parsing
			case strings.Contains(headerStr, "Data"):
				dateStr = valStr
			case strings.Contains(headerStr, "Czas trwania"):
				timeStr = valStr
			case strings.Contains(headerStr, "Skład"):
				judgesStr = valStr
			}
		}

		trial.Date, err = parseAndLocalizeTimeV3(dateStr, timeStr)
		if err != nil {
			return nil, fmt.Errorf("parsing trial page: %w", err)
		}

		trial.Judges = parseJudges(judgesStr)

		trials = append(trials, trial)
	}

	return
}

func parseJudges(judgesStr string) []string {
	judges := strings.Split(judgesStr, "\n")
	for i, j := range judges {
		judges[i] = strings.TrimSpace(j)
	}
	return judges
}

func parseAndLocalizeTimeV3(dateStr, timeStr string) (time.Time, error) {
	timeStrSplit := strings.Split(timeStr, "-")
	if len(timeStr) < 1 {
		return time.Time{}, fmt.Errorf("parsing time")
	}
	timeStr = strings.TrimSpace(timeStrSplit[0])
	date, err := time.Parse("02.01.2006 15:04", dateStr+" "+timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return localizeTime(date), nil
}

func parseV3Pagination(doc *goquery.Document) (maxPage int, err error) {
	navBar := doc.Find("nav.pagination")
	if navBar.Length() == 0 {
		return 1, nil
	}

	maxPage = 1
	navBar.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		u, err1 := url.Parse(href)
		if err1 != nil {
			err = err1
			return
		}
		pageStr := u.Query().Get("page")
		if pageStr == "" {
			return
		}
		page, err1 := strconv.Atoi(pageStr)
		if err1 != nil {
			err = err1
			return
		}

		maxPage = max(maxPage, page)
	})

	return maxPage, nil
}

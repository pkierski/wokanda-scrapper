package trialdownloader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

const (
	maxConcurrentIndexPageRequests   = 2
	maxConcurrentDetailsPageRequests = 4
)

type V1Wokanda commonDownloader

// check if V1Wokanda implements Downloader
var _ Downloader = (*V1Wokanda)(nil)

func NewV1Wokanda(client *http.Client, baseUrl string) V1Wokanda {
	return V1Wokanda{
		client:  client,
		baseUrl: "https://" + baseUrl,
	}
}

// Downloads all trials.
// date is string in format YYYY-MM-DD.
func (d V1Wokanda) Download(ctx context.Context, date string) ([]Trial, error) {
	di, pages, err := d.getListPage(ctx, date, 0)
	if err != nil {
		return nil, err
	}

	// download other pages
	var diMu sync.Mutex
	egPages, taskCtx := errgroup.WithContext(ctx)
	egPages.SetLimit(maxConcurrentIndexPageRequests)
	for page := range pages {
		if page == 0 { // first page already downloaded and parsed
			continue
		}
		egPages.Go(func() error {
			diPage, _, errPage := d.getListPage(taskCtx, date, page)
			diMu.Lock()
			di = append(di, diPage...)
			diMu.Unlock()
			return errPage
		})
	}

	err = egPages.Wait()
	if err != nil {
		return nil, err
	}

	// download details pages
	var trialMu sync.Mutex
	trials := make([]Trial, 0, len(di))
	egDetails, taskCtx := errgroup.WithContext(ctx)
	egDetails.SetLimit(maxConcurrentDetailsPageRequests)
	for _, pageIndex := range di {
		egDetails.Go(func() error {
			trial, err := d.getDetailPage(ctx, d.client, pageIndex)
			if err == nil {
				trialMu.Lock()
				trials = append(trials, trial)
				trialMu.Unlock()
			}
			return err
		})
	}

	err = egDetails.Wait()

	return SortTrials(trials), err
}

// getListPage downloads list of cases (selected page)
// returns:
//   - indices to cases on this pages (https://<baseUrl>/wokanda,I)
//   - number of list pages,
//   - error
func (d V1Wokanda) getListPage(ctx context.Context, date string, page int) (detailIndices []int, pages int, err error) {
	// get first page, extract indices to detail page and number of next pages
	pageContent, err := getOne(ctx, d.client, fmt.Sprintf("%v/index.php?p=cases&action=search&data=%v&s=%v", d.baseUrl, date, page))
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}
	if !bytes.Contains(pageContent, []byte(`<form action="index.php" method="GET" class="cases-form">`)) {
		err = fmt.Errorf("parsing trial page: %w", ErrNoDataOnPage)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageContent))
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}

	// list of pages:
	// <ul class="main-news-pagination list-unstyled list-inline text-center">
	// get the last element
	// check if any (avoid error on parsing empty string)
	lastPageList := doc.Find(`ul[class="main-news-pagination list-unstyled list-inline text-center"]`)
	if lastPageList.Length() == 0 {
		return
	}
	lastPage := lastPageList.First().Find("span.title").Last().Text()
	lp, err := strconv.ParseUint(lastPage, 10, 64)
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}

	// pages are indexed from 1, so the last number is the number of pages
	pages = int(lp)

	// get and parse index pages (extend details page indices)
	moreLinks := doc.Find("a.more-link")
	for i, s := range moreLinks.EachIter() {
		href, exists := s.Attr("href")
		if !exists {
			err = fmt.Errorf("parsing trial page: missing link %v on page %v", i, page)
			return
		}
		detailIndexStr, exists := strings.CutPrefix(href, "wokanda,")
		if !exists {
			err = fmt.Errorf("parsing trial page: bad format of link %v on page %v (%v)", i, page, href)
			return
		}
		di, err1 := strconv.ParseUint(detailIndexStr, 10, 64)
		if err != nil {
			err = fmt.Errorf("parsing trial page: bad format of link %v on page %v (%v): %w", i, page, href, err1)
			return
		}
		detailIndices = append(detailIndices, int(di))
	}
	return
}

func (d V1Wokanda) getDetailPage(ctx context.Context, client *http.Client, index int) (Trial, error) {
	data, err := getOne(ctx, client, fmt.Sprintf("%v/wokanda,%v", d.baseUrl, index))
	if err != nil {
		return Trial{}, fmt.Errorf("downloading trial detail page: %w", err)
	}
	return parseV1DetailPage(data)
}

// parseV1DetailPage parses one page from type "<url>/wokanda,N".
func parseV1DetailPage(data []byte) (trial Trial, err error) {
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

	caseIdTitle := doc.Selection.Find("h2.main-header").Find("span.title").Text()
	caseIdTitle = strings.ReplaceAll(caseIdTitle, "Sprawa ", "")
	if caseIdTitle != "" {
		trial.CaseID = caseIdTitle
	}

	s := doc.Selection.Find("dl[class='dl-horizontal case-description-list']")
	dts := s.Find("dt")
	dds := s.Find("dd")
	for i, dt := range dts.EachIter() {
		header := dt.Text()
		val := dds.Eq(i).Text()

		switch {
		case strings.Contains(header, "Sygnatura"):
			if val != "" {
				trial.CaseID = val
			}

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

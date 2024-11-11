package trialdownloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/trial"
)

// GetV1 parses all pages from type "<url>/wokanda,N".
//
// Url is expected in form of bare domain, ex.: https://poznan.so.gov.pl
func GetV1(ctx context.Context, client *http.Client, url string) ([]trial.Trial, error) {
	trialNo := 0
	var done atomic.Bool
	requestCh := make(chan int)

	// request generator -> requestCh
	go func() {
		for !done.Load() {
			trialNo++
			requestCh <- trialNo
		}
		close(requestCh)
	}()

	// requestCh -> workers -> results
	resultsCh := make(chan trial.Trial)
	errorsCh := make(chan error, 1)

	wg := sync.WaitGroup{}
	for range 16 {
		wg.Add(1)
		go func() { // worker
			defer wg.Done()
			for trialNo := range requestCh {
				t, err := getOneAndParseV1(ctx, client, fmt.Sprintf("%v/wokanda,%v", url, trialNo))
				if err != nil {
					// ignore ErrNoDataOnPage because it's the page out of range
					// except the first page has no data (no data at all, not in proper format)
					if !errors.Is(err, trial.ErrNoDataOnPage) || trialNo == 1 {
						errorsCh <- err
					}
					done.Store(true)
					<-requestCh // force generator to check done
					break
				}
				resultsCh <- t
			}
		}()
	}

	errs := make([]error, 0)
	go collect(errorsCh, errs)
	results := make([]trial.Trial, 0)
	go collect(resultsCh, results)

	wg.Wait()
	close(errorsCh)
	close(resultsCh)

	return results, errors.Join(errs...)
}

func getOneAndParseV1(ctx context.Context, client *http.Client, url string) (trial.Trial, error) {
	data, err := getOne(ctx, client, url)
	if err != nil {
		return trial.Trial{}, err
	}

	return trial.ParseV1(data)
}

func collect[E any](c <-chan E, s []E) {
	for e := range c {
		s = append(s, e)
	}
}

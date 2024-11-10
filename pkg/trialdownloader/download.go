package trialdownloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/trial"
)

func Get(ctx context.Context, client *http.Client, url string) ([]trial.Trial, error) {
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
	results := make([]trial.Trial, 0)

	wg := sync.WaitGroup{}
	for range 16 {
		wg.Add(1)
		go func() { // worker
			defer wg.Done()
			for trialNo := range requestCh {
				t, err := getOneAndParse(ctx, client, fmt.Sprintf("%v/wokanda,%v", url, trialNo))
				if err != nil {
					if !errors.Is(err, trial.ErrNoDataOnPage) {
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

	// errorsCh -> errors collector
	errs := make([]error, 0)
	go func() {
		for err := range errorsCh {
			errs = append(errs, err)
		}
	}()

	go func() {
		for t := range resultsCh {
			results = append(results, t)
		}
	}()

	wg.Wait()

	return results, errors.Join(errs...)
}

func getOneAndParse(ctx context.Context, client *http.Client, url string) (trial.Trial, error) {
	data, err := getOne(ctx, client, url)
	if err != nil {
		return trial.Trial{}, err
	}

	return trial.Parse(data)
}

func getOne(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch page: building request: %w", err)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch page: request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch page: unexpected status: %v (%v)", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch page body: %w", err)
	}

	return data, nil
}

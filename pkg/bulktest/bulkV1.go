package bulktest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
	"golang.org/x/sync/errgroup"
)

type resultType struct {
	Url         string                  `json:"url"`
	Trials      []trialdownloader.Trial `json:"trials"`
	Err         error                   `json:"err"`
	DateAquired time.Time               `json:"date_aquired"`
}

func BulkV1Test(ctx context.Context, client *http.Client) {
	eg, taskCtx := errgroup.WithContext(ctx)
	eg.SetLimit(16)

	results := make(map[string]resultType)
	resultsCh := make(chan resultType)

	go func() {
		for result := range resultsCh {
			results[result.Url] = result
			writeResultV1(result)
		}
	}()

	// domains := slices.Clone(data.Domains)
	// TODO: remove already checked
	domains := []string{"legnica.so.gov.pl"}

	for _, url := range domains {
		eg.Go(func() error {
			log.Printf("starting %v", url)
			defer log.Printf("finished %v", url)

			var result resultType
			result.Url = url
			downloader := trialdownloader.NewV1Wokanda(client, fmt.Sprintf("https://%v", url))
			result.Trials, result.Err = downloader.Download(taskCtx, "2006-02-01")
			result.DateAquired = time.Now().UTC()
			resultsCh <- result
			return nil
		})
	}

	eg.Wait()
	close(resultsCh)
}

func writeResultV1(result resultType) {
	var filename string
	if result.Err == nil {
		// log.Println("***", result.Url, result.Trials)
		filename = filepath.Join("results", "v1", result.Url)
	} else {
		filename = filepath.Join("results", "v1", "err", result.Url)
	}
	data, _ := json.MarshalIndent(result, "", "    ")
	os.WriteFile(filename, data, 0o644)

	if result.Err != nil {
		return
	}

	total, err := os.OpenFile(filepath.Join("results", "v1", "results.csv"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}
	defer total.Close()

	var firstDate, lastDate time.Time
	if result.Err == nil && len(result.Trials) > 0 {
		firstTrial := slices.MinFunc(result.Trials, func(a, b trialdownloader.Trial) int {
			return a.Date.Compare(b.Date)
		})
		firstDate = firstTrial.Date

		lastTrial := slices.MaxFunc(result.Trials, func(a, b trialdownloader.Trial) int {
			return a.Date.Compare(b.Date)
		})
		lastDate = lastTrial.Date
	}
	fmt.Fprintf(total, "%v,%v,%v,%v,%v\n", result.Url, len(result.Trials), result.Err, firstDate, lastDate)
}

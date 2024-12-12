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
	"sync"
	"time"

	"github.com/pkierski/wokanda-scrapper/pkg/data"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
	"golang.org/x/sync/errgroup"
)

type resultType struct {
	Url         string                  `json:"url"`
	Trials      []trialdownloader.Trial `json:"trials"`
	Err         error                   `json:"err"`
	DateAquired time.Time               `json:"date_aquired"`
}

func BulkV1Test(ctx context.Context, client *http.Client) ([]string, []bool) {
	eg, taskCtx := errgroup.WithContext(ctx)
	//eg.SetLimit(128)

	domains := slices.Clone(data.Domains)
	// TODO: remove already checked
	// domains := []string{"legnica.so.gov.pl"}
	resultsV1 := make([]bool, len(domains))
	var resultsV1Mu sync.Mutex

	for _, url := range domains {
		eg.Go(func() error {
			url = "https://" + url
			log.Printf("starting %v", url)
			defer log.Printf("finished %v", url)

			r := trialdownloader.Detect(taskCtx, client, url)
			if len(r) == 1 && r[0] == trialdownloader.AppTypeV1 {
				i := slices.Index(domains, url)
				resultsV1Mu.Lock()
				resultsV1[i] = true
				resultsV1Mu.Unlock()
			}
			return nil
		})
	}

	eg.Wait()

	return domains, resultsV1
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

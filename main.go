package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"maps"
	"net"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
)

func main() {
	// judges, err := extractJudges(`dr_2024-12-13.json`)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(strings.Join(judges, "\n"))
	// return
	insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "Don't check server certificate (default false)")
	const todayPlaceholder = "today"
	scrapDate := flag.String("scrap-date", todayPlaceholder, "Scrap trials from specified date: format YYYY-MM-DD")
	relativeScrapDate := flag.Int("relative-scrap-date", 0, "Scrap date as relative day (-1 means yesterday, 1 means tomorrow). Doesn't apply if -scrap-date used.")

	flag.Parse()

	if *scrapDate == todayPlaceholder {
		*scrapDate = time.Now().AddDate(0, 0, *relativeScrapDate).Format("2006-01-02")
	}

	// fmt.Println(*insecureSkipVerify, *scrapDate)
	// return

	transport := cleanhttp.DefaultPooledTransport()
	transport.MaxConnsPerHost = 1000
	transport.MaxIdleConns = 0
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: *insecureSkipVerify}
	client := retryablehttp.NewClient()
	client.RetryMax = 7
	client.RetryWaitMax = 10 * time.Second
	client.HTTPClient = &http.Client{
		Transport: transport,
	}

	cd, err := trialdownloader.LoadCourtsData("courts.json")
	if err != nil {
		panic(err)
	}

	date := *scrapDate

	start := time.Now().Format("2006-01-02T15-04-05")
	dr := trialdownloader.BulkDownload(context.Background(), client.StandardClient(), date, cd)

	trialdownloader.SaveJson(fmt.Sprintf("trials_%v_fetched-%v.json", date, start), dr)
}

func extractJudges(filename string) (result []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	var data []trialdownloader.DownloadResult
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return
	}

	m := map[string]struct{}{}
	for _, dr := range data {
		for _, t := range dr.Trials {
			for _, judge := range t.Judges {
				m[judge] = struct{}{}
			}
		}
	}

	result = slices.Collect(maps.Keys(m))
	slices.Sort(result)

	return
}

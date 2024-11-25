package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
)

func main() {
	client := retryablehttp.NewClient()
	// bulktest.BulkV1Test(context.Background(), client.StandardClient())
	// return

	if len(os.Args) < 2 {
		fmt.Printf("Podaj adres sądu do sprawdzenia wokandy (np.: https://poznan.so.gov.pl)")
		return
	}

	// TODO: use constructor based on url
	downloader := trialdownloader.NewV2Wokanda(client.StandardClient(), os.Args[1])

	trials, err := downloader.Download(context.Background(), "2006-01-02")
	if err != nil {
		panic(err)
	}

	slices.SortFunc(trials, func(a, b trialdownloader.Trial) int {
		return a.Date.Compare(b.Date)
	})
	j, _ := json.MarshalIndent(trials, "", "  ")
	fmt.Println(string(j))
}

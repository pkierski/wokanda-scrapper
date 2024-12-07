package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
)

func main() {
	transport := cleanhttp.DefaultPooledTransport()
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConns = 100
	transport.DialContext = (&net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext
	client := retryablehttp.NewClient()
	client.RetryMax = 7
	client.HTTPClient = &http.Client{
		Transport: transport,
	}

	cd, err := trialdownloader.LoadCourtsData("courts.json")
	if err != nil {
		panic(err)
	}

	const date = "2024-12-06"
	dr := trialdownloader.BulkDownload(context.Background(), client.StandardClient(), date, cd)

	trialdownloader.SaveJson(fmt.Sprintf("dr_%v.json", date), dr)

	// courtsData := trialdownloader.DetectBulk(context.Background(), client.StandardClient(), data.Domains)
	// f1, err := os.Create("courts.json")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f1.Close()
	// encoder := json.NewEncoder(f1)
	// encoder.SetIndent("", "  ")
	// encoder.Encode(courtsData)

	// b := &bytes.Buffer{}
	// for _, cd := range courtsData {
	// 	isV1 := slices.Contains(cd.AppTypes, trialdownloader.AppTypeV1)
	// 	fmt.Fprintf(b, "%v,%v\n", cd.Domain, isV1)
	// }
	// os.WriteFile("courts.txt", b.Bytes(), 0o644)

	return

	// domains, v1Results := bulktest.BulkV1Test(context.Background(), client.StandardClient())

	// f1, err := os.Create("v1_with_false.csv")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f1.Close()
	// for i, d := range domains {
	// 	fmt.Fprintf(f1, "%v,%v\n", d, v1Results[i])
	// }

	// f2, err := os.Create("v1_with_empty.csv")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f2.Close()
	// for i, d := range domains {
	// 	if !v1Results[i] {
	// 		d = ""
	// 	}
	// 	fmt.Fprintln(f2, d)
	// 	// fmt.Fprintf(f2, "%v: %v\n", d, v1Results[i])
	// }
	// return

	if len(os.Args) < 2 {
		fmt.Printf("Podaj adres sądu do sprawdzenia wokandy (np.: https://poznan.so.gov.pl)")
		return
	}

	// TODO: use constructor based on url
	downloader := trialdownloader.NewV1Wokanda(client.StandardClient(), os.Args[1])

	trials, err := downloader.Download(context.Background(), "2024-11-27")
	if err != nil {
		panic(err)
	}

	trialdownloader.SortTrials(trials)
	j, _ := json.MarshalIndent(trials, "", "  ")
	fmt.Println(string(j))
}

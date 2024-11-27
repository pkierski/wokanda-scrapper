package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkierski/wokanda-scrapper/pkg/data"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
)

func main() {
	client := retryablehttp.NewClient()

	courtsData := trialdownloader.DetectBulk(context.Background(), client.StandardClient(), data.Domains)
	f1, err := os.Create("courts.json")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	encoder := json.NewEncoder(f1)
	encoder.SetIndent("", "  ")
	encoder.Encode(courtsData)

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

	slices.SortFunc(trials, func(a, b trialdownloader.Trial) int {
		if c := strings.Compare(a.CaseID, b.CaseID); c != 0 {
			return c
		}
		return a.Date.Compare(b.Date)
	})
	j, _ := json.MarshalIndent(trials, "", "  ")
	fmt.Println(string(j))
}

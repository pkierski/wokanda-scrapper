package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/trial"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Podaj adres sądu do sprawdzenia wokandy (np.: https://poznan.so.gov.pl)")
		return
	}

	trials, err := trialdownloader.GetV2(context.Background(), http.DefaultClient, os.Args[1])
	if err != nil {
		panic(err)
	}

	slices.SortFunc(trials, func(a, b trial.Trial) int {
		return a.Date.Compare(b.Date)
	})
	j, _ := json.MarshalIndent(trials, "", "  ")
	fmt.Println(string(j))
}

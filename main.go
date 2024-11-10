package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Podaj adres sądu do sprawdzenia wokandy (np.: https://poznan.so.gov.pl)")
		return
	}

	trials, err := trialdownloader.Get(context.Background(), http.DefaultClient, os.Args[1])
	if err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(trials, "", "  ")
	fmt.Println(string(j))
}

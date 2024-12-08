package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
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
}

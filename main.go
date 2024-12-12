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
	transport.MaxConnsPerHost = 1000
	transport.MaxIdleConns = 0
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext
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

	const date = "2024-12-11"
	dr := trialdownloader.BulkDownload(context.Background(), client.StandardClient(), date, cd)

	trialdownloader.SaveJson(fmt.Sprintf("dr_%v.json", date), dr)
}

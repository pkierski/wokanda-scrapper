package trialdownloader

import (
	"context"
	"net/http"
)

type V2Wokanda commonDownloader

// check if V2Wokanda implements Downloader
var _ Downloader = (*V1Wokanda)(nil)

func NewV2Wokanda(client *http.Client, baseUrl string) V2Wokanda {
	return V2Wokanda{
		client:  client,
		baseUrl: baseUrl,
	}
}

// Downloads all trials.
// date is string in format YYYY-MM-DD.
func (d V2Wokanda) Download(ctx context.Context, date string) ([]Trial, error) {
	// TODO: filter by date
	return getV2(ctx, d.client, d.baseUrl)
}

func getV2(ctx context.Context, client *http.Client, url string) ([]Trial, error) {
	data, err := postOne(ctx, client, url, "akcja=szukaj&wydzial=wszystko&data_s=wszystko&sygnatura=")
	if err != nil {
		return nil, err
	}

	return ParseV2(data)
}

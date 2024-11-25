package trialdownloader

import (
	"context"
	"net/http"
)

func GetV2(ctx context.Context, client *http.Client, url string) ([]Trial, error) {
	data, err := postOne(ctx, client, url, "akcja=szukaj&wydzial=wszystko&data_s=wszystko&sygnatura=")
	if err != nil {
		return nil, err
	}

	return ParseV2(data)
}

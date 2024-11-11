package trialdownloader

import (
	"context"
	"net/http"
	"os"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/trial"
)

func GetV2(ctx context.Context, client *http.Client, url string) ([]trial.Trial, error) {
	data, err := postOne(ctx, client, url, "akcja=szukaj&wydzial=wszystko&data_s=wszystko&sygnatura=")
	if err != nil {
		return nil, err
	}

	os.WriteFile("foo.html", data, 0o666)

	return trial.ParseV2(data)
}

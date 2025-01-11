package trialdownloader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

type CourtData struct {
	Domain   string    `json:"domain"`
	AppTypes []AppType `json:"app_types"`
}

const (
	AppTypeV1                 AppType = "V1:url/wokanda,I"
	AppTypeV3LogonetBydgoszcz AppType = "V3:logonet-bydgoszcz-okreg"
)

func Detect(ctx context.Context, client *http.Client, baseUrl string) (result []AppType) {
	detectors := [](func(ctx context.Context, client *http.Client, baseUrl string) (bool, AppType)){
		detectV1,
		detectV3Bydgoszcz,
	}

	var resultMu sync.Mutex

	eg, ctx := errgroup.WithContext(ctx)
	for _, d := range detectors {
		eg.Go(func() error {
			if found, typ := d(ctx, client, baseUrl); found {
				resultMu.Lock()
				result = append(result, typ)
				resultMu.Unlock()
			}
			return nil
		})
	}

	eg.Wait()

	return
}

func DetectBulk(ctx context.Context, client *http.Client, domains []string) []CourtData {
	courts := make([]CourtData, len(domains))
	eg, taskCtx := errgroup.WithContext(ctx)

	for i, domain := range domains {
		eg.Go(func() error {
			court := CourtData{
				Domain:   domain,
				AppTypes: append(make([]AppType, 0), Detect(taskCtx, client, domain)...),
			}
			courts[i] = court
			return nil
		})
	}

	eg.Wait()
	return courts
}

func detectV1(ctx context.Context, client *http.Client, baseUrl string) (found bool, typ AppType) {
	typ = AppTypeV1
	page, err := getOne(ctx, client, fmt.Sprintf("https://%v/wokanda", baseUrl))
	if err != nil {
		return
	}

	found = bytes.Contains(page, []byte(`<form action="index.php" method="GET" class="cases-form">`)) &&
		bytes.Contains(page, []byte(`<input name="p" type="hidden" value="cases"`)) &&
		bytes.Contains(page, []byte(`<input name="action" type="hidden" value="search"`))

	return
}

func detectV3Bydgoszcz(ctx context.Context, client *http.Client, baseUrl string) (found bool, typ AppType) {
	typ = AppTypeV3LogonetBydgoszcz
	page, err := getOne(ctx, client, fmt.Sprintf("https://%v", baseUrl))
	if err != nil {
		return
	}

	found = bytes.Contains(page, []byte(`CMS i hosting: Logonet Sp. z o.o. w Bydgoszczy`))

	return
}

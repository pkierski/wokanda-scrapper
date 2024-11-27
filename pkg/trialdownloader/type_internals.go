package trialdownloader

import "net/http"

type (
	commonDownloader struct {
		client  *http.Client
		baseUrl string
	}
)

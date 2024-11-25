package trialdownloader

import (
	"context"
	"time"
)

type (
	Trial struct {
		CaseID     string    `json:"case_id"`
		Department string    `json:"department"`
		Judges     []string  `json:"judges"`
		Date       time.Time `json:"date"`
		Room       string    `json:"room"`
	}

	Downloader interface {
		// Downloads all trials.
		//
		// date is string in format YYYY-MM-DD.
		Download(ctx context.Context, date string) ([]Trial, error)
	}
)

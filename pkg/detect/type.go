package detect

import (
	"context"
	"net/http"
)

type Detector interface {
	Detect(ctx context.Context, client *http.Client, url string) bool
}

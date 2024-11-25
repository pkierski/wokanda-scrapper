package detect

import (
	"context"
	"net/http"
)

type V1 struct{}

func (V1) Detect(ctx context.Context, client *http.Client, url string) bool {
	return false
}

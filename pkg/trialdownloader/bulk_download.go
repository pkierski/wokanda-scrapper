package trialdownloader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

type DownloadResult struct {
	CourtID string
	Err     string
	Trials  []Trial
	Date    string
}

func BulkDownload(ctx context.Context, client *http.Client, date string, courtData []CourtData) []DownloadResult {
	result := make([]DownloadResult, 0, len(courtData))
	var resultMu sync.Mutex

	eg, taskCtx := errgroup.WithContext(ctx)

	for _, cd := range courtData {
		eg.Go(func() error {
			dr := DownloadResult{
				CourtID: cd.Domain,
				Date:    date,
			}

			if len(cd.AppTypes) == 1 {
				var downloader Downloader
				switch cd.AppTypes[0] {
				case AppTypeV1:
					downloader = NewV1Wokanda(client, cd.Domain)
				}
				if downloader == nil {
					dr.Err = fmt.Errorf("unknown or app type: %v", cd.AppTypes).Error()
				} else {
					trials, err := downloader.Download(taskCtx, date)
					var errStr string
					if err != nil {
						errStr = err.Error()
					}
					dr = DownloadResult{
						CourtID: cd.Domain,
						Err:     errStr,
						Trials:  trials,
						Date:    date,
					}
				}
			} else {
				dr.Err = fmt.Errorf("unknown or ambiguous app type: %v", cd.AppTypes).Error()
			}

			resultMu.Lock()
			result = append(result, dr)
			resultMu.Unlock()
			SaveJson(fmt.Sprintf("dr_%v_%v.json", dr.Date, dr.CourtID), dr)
			return nil
		})
	}

	// all errors are enclosed in DownloadResult for each court
	// download task always return nil
	eg.Wait()

	slices.SortFunc(result, func(a, b DownloadResult) int {
		return strings.Compare(a.CourtID, b.CourtID)
	})
	return result
}

func LoadCourtsData(filename string) ([]CourtData, error) {
	return loadJson[[]CourtData](filename)
}

func loadJson[T any](filename string) (res T, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&res)
	return
}

func SaveJson[T any](filename string, data T) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

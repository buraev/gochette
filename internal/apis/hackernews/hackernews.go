package hackernews

import (
	"lightweight-cache-proxy-service/internal/cache"
	"lightweight-cache-proxy-service/pgc/lcp"
	"net/http"
	"time"

	"github.com/buraev/barelog"
)

const cacheInstance = cache.HackerNews

func Setup(mux *http.ServeMux) {
	initialData, err := FetchTop30()
	if err != nil {
		barelog.Error(err, "fetching initial HackerNews failed")
	}

	hnCache := cache.New(cacheInstance, initialData, err == nil)

	mux.Handle("/api/hn/top", hnCache)

	go cache.UpdatePeriodically(
		hnCache,
		struct{}{},
		func(_ struct{}) ([]lcp.HackerNews, error) {
			return FetchTop30()
		},
		15*time.Minute,
	)

	barelog.Info(cacheInstance, "setup cache and endpoint")
}

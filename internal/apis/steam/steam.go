package steam

import (
	"lightweight-cache-proxy-service/internal/cache"
	"lightweight-cache-proxy-service/pgc/lcp"
	"net/http"
	"time"

	"github.com/buraev/barelog"
)

const cacheInstance = cache.Steam

func Setup(mux *http.ServeMux, client *http.Client) {
	games, err := fetchRecentlyPlayedGames(client)
	if err != nil {
		barelog.Error(err, "initial fetch of steam games failed")
	}

	steamCache := cache.New(cacheInstance, games, err == nil)
	mux.Handle("/api/steam", steamCache)
	go cache.UpdatePeriodically(
		steamCache,
		client,
		func(client *http.Client) ([]lcp.SteamGame, error) {
			return fetchRecentlyPlayedGames(client)
		},
		15*time.Minute,
	)
	barelog.Info(cacheInstance.LogPrefix(), "setup cache and endpoint")
}

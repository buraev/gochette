package github

import (
	"context"
	"lightweight-cache-proxy-service/internal/apis/secrets"
	"lightweight-cache-proxy-service/internal/cache"
	"net/http"
	"time"

	"github.com/buraev/barelog"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const cacheInstance = cache.GitHub

func Setup(mux *http.ServeMux) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: withHeaders(oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: secrets.ENV.GitHubAccessToken}),
		).Transport, map[string]string{
			"User-Agent": "lightweight-cache-proxy/1.0",
			"Accept":     "application/json",
		}),
	}

	githubClient := githubv4.NewClient(httpClient)

	pinnedRepos, err := fetchPinnedRepos(githubClient)
	if err != nil {
		barelog.Error(err, "fetching initial pinned repos failed")
	}

	githubCache := cache.New(cacheInstance, pinnedRepos, err == nil)
	mux.Handle("/api/github", githubCache)
	go cache.UpdatePeriodically(githubCache, githubClient, fetchPinnedRepos, 1*time.Minute)

	barelog.Info(cacheInstance, "setup cache and endpoint")
}

func withHeaders(rt http.RoundTripper, headers map[string]string) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return rt.RoundTrip(req)
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

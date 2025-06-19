package hackernews

import (
	"encoding/json"
	"fmt"
	"lightweight-cache-proxy-service/pgc/lcp"
	"net/http"
	"time"
)

const (
	topStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	itemURL       = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

func FetchTop30() ([]lcp.HackerNews, error) {
	resp, err := http.Get(topStoriesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}
	if len(ids) > 30 {
		ids = ids[:30]
	}

	var items []lcp.HackerNews
	for _, id := range ids {
		item, err := fetchItem(id)
		if err == nil {
			items = append(items, item)
		}
		time.Sleep(30 * time.Millisecond) // избежать rate limiting
	}

	return items, nil
}

func fetchItem(id int) (lcp.HackerNews, error) {
	url := fmt.Sprintf(itemURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return lcp.HackerNews{}, err
	}
	defer resp.Body.Close()

	var item lcp.HackerNews
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return lcp.HackerNews{}, err
	}
	return item, nil
}

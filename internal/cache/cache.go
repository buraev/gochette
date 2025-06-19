package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"lightweight-cache-proxy-service/internal/apis"
	"lightweight-cache-proxy-service/internal/apis/secrets"
	"lightweight-cache-proxy-service/pgc/lcp"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/buraev/barelog"
)

type CacheInstance int

const (
	AppleMusic CacheInstance = iota
	GitHub
	Steam
	HackerNews
)

func (c CacheInstance) String() string {
	switch c {
	case AppleMusic:
		return "applemusic"
	case HackerNews:
		return "hackernews"
	case GitHub:
		return "github"
	case Steam:
		return "steam"
	}
	return "unknown"
}

func (c CacheInstance) LogPrefix() string {
	return fmt.Sprintf("[%s]", c.String())
}

type Cache[T lcp.CacheData] struct {
	instance CacheInstance
	filePath string
	Mutex    sync.RWMutex
	Data     T
	Updated  time.Time
}

func New[T lcp.CacheData](instance CacheInstance, data T, update bool) *Cache[T] {
	cache := Cache[T]{
		instance: instance,
		Updated:  time.Now().UTC(),
		filePath: filepath.Join(secrets.ENV.CacheFolder, fmt.Sprintf("%s.json", instance.String())),
	}
	cache.loadFromFile()
	if update {
		cache.Update(data)
	}
	return &cache
}

type CacheResponse[T any] struct {
	Data    T         `json:"data"`
	Updated time.Time `json:"updated"`
}

func (c *Cache[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: вернуть авторизацию
	// if !auth.IsAuthorized(w, r) {
	//     return
	// }

	w.Header().Set("Content-Type", "application/json")

	c.Mutex.RLock()
	data := CacheResponse[T]{
		Data:    c.Data,
		Updated: c.Updated,
	}
	c.Mutex.RUnlock()

	if err := json.NewEncoder(w).Encode(data); err != nil {
		wrappedErr := fmt.Errorf("failed to write JSON response: %w", err)
		barelog.Error(fmt.Sprintf("ServeHTTP error [%s]: %v", c.instance, wrappedErr))
		http.Error(w, wrappedErr.Error(), http.StatusInternalServerError)
	}
}

func (c *Cache[T]) Update(data T) {
	c.Mutex.RLock()
	currentData := c.Data
	c.Mutex.RUnlock()

	oldBin, err := json.Marshal(currentData)
	if err != nil {
		barelog.Error(fmt.Sprintf("failed to marshal old data [%s]: %v", c.instance, err))
		return
	}

	newBin, err := json.Marshal(data)
	if err != nil {
		barelog.Error(fmt.Sprintf("failed to marshal new data [%s]: %v", c.instance, err))
		return
	}

	newStr := strings.TrimSpace(string(newBin))

	if !bytes.Equal(oldBin, newBin) && newStr != "" && newStr != "null" {
		c.Mutex.Lock()
		c.Data = data
		c.Updated = time.Now().UTC()
		c.Mutex.Unlock()

		c.persistToFile()
		barelog.Info(fmt.Sprintf("cache updated [%s]", c.instance))
	} else {
		barelog.Debug(fmt.Sprintf("no cache update needed [%s]", c.instance))
	}
}

func UpdatePeriodically[T lcp.CacheData, C any](
	cache *Cache[T],
	client C,
	update func(C) (T, error),
	interval time.Duration,
) {
	for {
		time.Sleep(interval)
		data, err := update(client)
		if err != nil {
			if !errors.Is(err, apis.WarningError) {
				barelog.Error("updating", err, "cache failed", cache.instance)
			}
		} else {
			cache.Update(data)
		}
	}
}

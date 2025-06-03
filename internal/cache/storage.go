package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buraev/barelog"
)

func (c *Cache[T]) persistToFile() {
	folder := filepath.Dir(c.filePath)
	if err := os.MkdirAll(folder, 0o700); err != nil {
		barelog.Error(fmt.Sprintf("failed to create directory %q: %v", folder, err))
		return
	}

	file, err := os.OpenFile(c.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		barelog.Error(fmt.Sprintf("failed to open or create cache file %q: %v", c.filePath, err))
		return
	}
	defer file.Close()

	data := CacheResponse[T]{
		Data:    c.Data,
		Updated: c.Updated,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		barelog.Error(fmt.Sprintf("failed to encode cache data to %q: %v", c.filePath, err))
		return
	}

	barelog.Info(fmt.Sprintf("cache successfully persisted to file: %q", c.filePath))
}

func (c *Cache[T]) loadFromFile() {
	const maxFileSize = 5 * 1024 * 1024

	info, err := os.Stat(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		barelog.Error(fmt.Sprintf("failed to stat cache file %q: %v", c.filePath, err))
		return
	}

	if info.IsDir() {
		barelog.Error(fmt.Sprintf("cache path %q is a directory, expected a file", c.filePath))
		return
	}

	if info.Size() > maxFileSize {
		barelog.Error(fmt.Sprintf("cache file %q is too large (%d bytes), max allowed is %d bytes", c.filePath, info.Size(), maxFileSize))
		return
	}

	b, err := os.ReadFile(c.filePath)
	if err != nil {
		barelog.Error(fmt.Sprintf("failed to read cache file %q: %v", c.filePath, err))
		return
	}

	if len(b) == 0 {
		return
	}

	var data CacheResponse[T]
	if err := json.Unmarshal(b, &data); err != nil {
		barelog.Error(fmt.Sprintf("failed to unmarshal cache file %q: %v. Raw data: %q", c.filePath, err, b))
		return
	}

	c.Data = data.Data
	c.Updated = data.Updated
}

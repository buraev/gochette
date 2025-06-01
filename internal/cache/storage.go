package cache

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/buraev/barelog"
)

func (c *Cache[T]) persistToFile() {
	var file *os.File
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		folder := filepath.Dir(c.filePath)
		err := os.MkdirAll(folder, 0700)
		if err != nil {
			barelog.Error(err.Error(), "failed to create folder at path:", folder)
			return
		}
		file, err = os.Create(c.filePath)
		if err != nil {
			barelog.Error(err.Error(), "failed to create file at path:", c.filePath)
			return
		}
	} else {
		file, err = os.OpenFile(c.filePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			barelog.Error(err.Error(), "failed to read file at path:", c.filePath)
			return
		}
	}
	defer file.Close()

	c.Mutex.RLock()
	bin, err := json.Marshal(CacheResponse[T]{
		Data:    c.Data,
		Updated: c.Updated,
	})
	c.Mutex.RUnlock()
	if err != nil {
		barelog.Error(err.Error(), "encoding data to json failed")
		return
	}
	_, err = file.Write(bin)
	if err != nil {
		barelog.Error(err.Error(), "writing data to json failed")
	}
}

func (c *Cache[T]) loadFromFile() {
	if _, err := os.Stat(c.filePath); !os.IsNotExist(err) {
		b, err := os.ReadFile(c.filePath)
		if err != nil {
			barelog.Error(err.Error(), "reading from cache file from", c.filePath, "failed")
		}

		var data CacheResponse[T]
		err = json.Unmarshal(b, &data)
		if err != nil {
			barelog.Error(err.Error(), "unmarshaling json data from", c.filePath, "failed:", string(b))
		}

		c.Data = data.Data
		c.Updated = data.Updated
	}
}

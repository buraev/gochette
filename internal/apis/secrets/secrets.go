package secrets

import (
	"errors"
	"io/fs"
	"os"

	"github.com/buraev/barelog"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ENV Secrets

type Secrets struct {
	ValidTokens   string `env:"VALID_TOKENS"`
	CacheFolder   string `env:"CACHE_FOLDER"`
	AllowFrontend string `env:"ALLOW_FRONTEND"`
	CorePath      string `env:"CORE_PATH"`
	// github
	GitHubAccessToken string `env:"GITHUB_ACCESS_TOKEN"`
}

func Load() {
	if _, err := os.Stat(".env"); !errors.Is(err, fs.ErrNotExist) {
		err := godotenv.Load()
		if err != nil {
			barelog.Error(err, "loading .env file failed")
		}
	}

	secrets, err := env.ParseAs[Secrets]()
	if err != nil {
		barelog.Error(err, "parsing required env vars failed")
	}
	ENV = secrets
	barelog.Info("loaded secrets")
}

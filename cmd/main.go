package main

import (
	"fmt"
	"lightweight-cache-proxy-service/internal/apis/github"
	"lightweight-cache-proxy-service/internal/apis/hackernews"
	"lightweight-cache-proxy-service/internal/apis/secrets"
	"lightweight-cache-proxy-service/internal/apis/steam"
	"lightweight-cache-proxy-service/internal/middleware"
	"time"

	"net/http"
)

func main() {
	middleware.SetupLogger()
	fmt.Println("booted")

	secrets.Load()

	var (
		mux    = http.NewServeMux()
		client = http.Client{Timeout: 20 * time.Second}
	)

	// mux.HandleFunc("/", middleware.RootRedirect)

	github.Setup(mux)
	hackernews.Setup(mux)
	steam.Setup(mux, &client)

	// Обернём mux в CORS middleware
	handler := middleware.WithCORS(mux)

	fmt.Println("starting server")
	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		fmt.Println(err)
	}
}

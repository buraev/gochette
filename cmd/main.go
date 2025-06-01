package main

import (
	"fmt"
	"lightweight-cache-proxy-service/internal/apis/github"
	"lightweight-cache-proxy-service/internal/apis/secrets"
	"net/http"
	"time"
)

func main() {
	setupLogger()
	fmt.Println("booted")

	secrets.Load()

	var (
		client = http.Client{}
		mux    = http.NewServeMux()
	)
	fmt.Println(client)

	mux.HandleFunc("/", rootRedirect)
	github.Setup(mux)

	// Обернём mux в CORS middleware
	handler := withCORS(mux)

	fmt.Println("starting server")
	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		fmt.Println(err)
	}
}

func setupLogger() {
	ny, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("failed to load timezone")
	}
	fmt.Println(ny)
}

func rootRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://buraev.com", http.StatusPermanentRedirect)
}

// CORS middleware
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с frontend'а
		w.Header().Set("Access-Control-Allow-Origin", secrets.ENV.AllowFrontend)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обработка preflight-запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

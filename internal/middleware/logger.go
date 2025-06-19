package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func SetupLogger() {
	ny, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("failed to load timezone")
	}
	fmt.Println(ny)
}

func RootRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	target := "https://www.buraev.com"
	http.Redirect(w, r, target, http.StatusPermanentRedirect)
}

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
	http.Redirect(w, r, "https://buraev.com", http.StatusPermanentRedirect)
}

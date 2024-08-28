package middleware

import (
	"log"
	"net/http"
)

type Cena struct{}

func LogRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request made to %s %s by %s\n", r.Method, r.URL.Path, r.RemoteAddr)
}

func LogResponse(w http.ResponseWriter, r *http.Request) {
	log.Println("logging response")
}

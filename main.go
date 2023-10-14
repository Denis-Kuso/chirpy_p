package main

import (
	"fmt"
	"net/http"
    "log"
)

func checkStatus(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
    const portNum = "8080"
    fmt.Println("Spinning up")

    rootPath := "." // home for now
    readinessEndpoint := "/healthz"

    mux := http.NewServeMux()
    corsMux := middlewareCors(mux)

    mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath))))
    mux.Handle(readinessEndpoint, http.HandlerFunc(checkStatus))
    server := &http.Server{
        Addr:   ":" + portNum,
        Handler: corsMux,
        MaxHeaderBytes: 1 << 20,
        }
    log.Printf("Serving on port: %s\n",portNum)
    server.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

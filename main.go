package main

import (
	"fmt"
	"net/http"
    "log"
)

//build custom handler
/*
Your handler should do the following:

    Write the Content-Type header
    Write the status code using w.WriteHeader
    Write the body text using w.Write
*/

func main() {
    const portNum = "8080"
    fmt.Println("Spinning up")

    rootPath := "." // home for now
    readinessEndpoint := "/healthz"

    mux := http.NewServeMux()
    corsMux := middlewareCors(mux)

    mux.Handle("/",http.FileServer(http.Dir(rootPath)))
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

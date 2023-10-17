package main

import (
	"fmt"
	"net/http"
    "log"
)



func main() {
    const portNum = "8080"

    state := apiState {
        ViewCount: 0,
    }
    rootPath := "." // home for now
    readinessEndpoint := "/healthz"
    metrics := "/metrics"
    resetCount := "/reset"

    mux := http.NewServeMux()
    corsMux := middlewareCors(mux)

    mux.Handle("/app/",state.middlewareMetrics(http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))))
    mux.HandleFunc(readinessEndpoint, http.HandlerFunc(checkStatus))
    mux.HandleFunc(metrics,state.showPageViews)
    mux.HandleFunc(resetCount, state.resetViews)

    server := &http.Server{
        Addr:   ":" + portNum,
        Handler: corsMux,
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

type apiState struct {
    ViewCount int
}


// Show number of views
func (s *apiState) showPageViews(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Hits: %d",s.ViewCount)))
}


// increment number of page views
func (s *apiState) middlewareMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        s.ViewCount++
        next.ServeHTTP(w, req)
})
}

// Reset page view count to zero
func (s *apiState) resetViews(w http.ResponseWriter, req *http.Request) {
    s.ViewCount = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}


func checkStatus(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(http.StatusText(http.StatusOK)))
    }

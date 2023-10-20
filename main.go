package main

import (
	"fmt"
	"net/http"
    "log"
    "github.com/go-chi/chi/v5"
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

    r := chi.NewRouter()

    fsHandler := state.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
    r.Handle("/app", fsHandler)
    r.Handle("/app/*", fsHandler)
    r.Get(metrics, state.showPageViews)
    r.Get(readinessEndpoint,http.HandlerFunc(checkStatus))
    r.Get(resetCount, state.resetViews)

    corsMux := middlewareCors(r)
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

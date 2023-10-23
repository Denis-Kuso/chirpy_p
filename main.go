package main

import (
	"fmt"
	"net/http"
    "log"
    "github.com/go-chi/chi/v5"
    "github.com/Denis-Kuso/chirpy_p/handlers"
)



func main() {
    const portNum = "8080"

    state := apiState {
        ViewCount: 0,
    }
    rootPath := "." // home for now
    readinessEndpoint := "/healthz"
    metrics := "/metrics"
    reset := "/reset"
    valid := "/validate_chirp"

    r := chi.NewRouter()

    fsHandler := state.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
    r.Handle("/app", fsHandler)
    r.Handle("/app/*", fsHandler)

    apiRouter := chi.NewRouter()
	apiRouter.Get(readinessEndpoint, checkStatus)
	apiRouter.Get(reset, state.resetViews)
    apiRouter.Post(valid, handlers.ValidateChirp)

    adminRouter := chi.NewRouter()
    adminRouter.Get(metrics,state.showPageViews)
    r.Mount("/admin", adminRouter)
	r.Mount("/api", apiRouter)

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
    w.Header().Add("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    displayInfo := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",s.ViewCount)
    w.Write([]byte(displayInfo))
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

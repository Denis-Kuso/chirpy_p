package main

import (
	"net/http"
    "log"
    "github.com/go-chi/chi/v5"
    "github.com/Denis-Kuso/chirpy_p/handlers"
)



func main() {
    const portNum = "8080"

    state := handlers.ApiState {
        ViewCount: 0,
    }
    rootPath := "." // home for now
    readinessEndpoint := "/healthz"
    metrics := "/metrics"
    reset := "/reset"
    valid := "/chirps"

    r := chi.NewRouter()

    fsHandler := state.MiddlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
    r.Handle("/app", fsHandler)
    r.Handle("/app/*", fsHandler)

    apiRouter := chi.NewRouter()
	apiRouter.Get(readinessEndpoint, handlers.CheckStatus)
	apiRouter.Get(reset, state.ResetViews)
    apiRouter.Post(valid, state.ValidateChirp)
    apiRouter.Get(valid, state.GetChirps)

    adminRouter := chi.NewRouter()
    adminRouter.Get(metrics,state.ShowPageViews)
    r.Mount("/admin", adminRouter)
	r.Mount("/api", apiRouter)

    corsMux := handlers.MiddlewareCors(r)
    server := &http.Server{
        Addr:   ":" + portNum,
        Handler: corsMux,
        }
    log.Printf("Serving on port: %s\n",portNum)
    server.ListenAndServe()
}


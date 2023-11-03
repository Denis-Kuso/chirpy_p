package main

import (
	"log"
	"net/http"

//	"github.com/Denis-Kuso/chirpy_p/handlers"
	"github.com/Denis-Kuso/chirpy_p/internal/database"
	"github.com/go-chi/chi/v5"
)

type ApiState struct {
    ViewCount int
    DB *database.DB
}


func main() {
    const portNum = "8080"

    loc := "database.json"
    db, err := database.NewDB(loc)
    if err != nil {
        log.Fatalf("Error with database: %v\n",err)
    }
    state := ApiState {
        ViewCount: 0,
        DB: db,
    }
    rootPath := "." // home for now
    readinessEndpoint := "/healthz"
    metrics := "/metrics"
    reset := "/reset"
    valid := "/chirps"
    users := "/users"
    login := "/login"

    r := chi.NewRouter()

    fsHandler := state.MiddlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
    r.Handle("/app", fsHandler)
    r.Handle("/app/*", fsHandler)

    apiRouter := chi.NewRouter()
	apiRouter.Get(readinessEndpoint, CheckStatus)
	apiRouter.Get(reset, state.ResetViews)
    apiRouter.Post(valid, state.ValidateChirp)
    apiRouter.Get(valid, state.GetChirps)
    apiRouter.Get("/chirps/{chirpID}", state.GetChirp)
    apiRouter.Post(users,state.CreateUser)
    apiRouter.Post(login, state.LoginUser)

    adminRouter := chi.NewRouter()
    adminRouter.Get(metrics,state.ShowPageViews)
    r.Mount("/admin", adminRouter)
	r.Mount("/api", apiRouter)

    corsMux := MiddlewareCors(r)
    server := &http.Server{
        Addr:   ":" + portNum,
        Handler: corsMux,
        }
    log.Printf("Serving on port: %s\n",portNum)
    server.ListenAndServe()
}


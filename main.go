package main

import (
	"log"
	"net/http"
	"os"
	//	"github.com/Denis-Kuso/chirpy_p/handlers"
	"github.com/Denis-Kuso/chirpy_p/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type ApiState struct {
    ViewCount int
    DB *database.DB
    Token string
    WebhookKey string
}


func main() {
    err := godotenv.Load()
    if err != nil {
	log.Fatal("Error loading .env file")
    }
    const portNum = "8080"

    loc := os.Getenv("DB_FILE")
    db, err := database.NewDB(loc)
    if err != nil {
        log.Fatalf("Failed loading database: %v\n",err)
    }

    token := os.Getenv("JWT_SECRET")
    webhookKey := os.Getenv("POLKA_KEY")
    state := ApiState {
        ViewCount: 0,
        DB: db,
	Token: token,
	WebhookKey: webhookKey,
    }
    rootPath := "." // home for now
    readinessEndpoint := "/healthz"
    metrics := "/metrics"
    reset := "/reset"
    valid := "/chirps"
    users := "/users"
    login := "/login"
    refresh := "/refresh"
    revoke := "/revoke"
    polka_webhook := "/polka/webhooks"

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
    apiRouter.Delete("/chirps/{chirpID}", state.RemoveChirp)

    apiRouter.Post(users,state.CreateUser)
    apiRouter.Put(users,state.UpdateUser)

    apiRouter.Post(login, state.LoginUser)
    apiRouter.Post(refresh,state.RefreshToken)
    apiRouter.Post(revoke,state.RevokeToken)
    apiRouter.Post(polka_webhook, state.UpgradeMembership)

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


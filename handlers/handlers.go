package handlers

import (
    "net/http"
    "log"
    "encoding/json"
    "fmt"
    "io"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {

    const charLimit = 140
    var statusCode int
    type requestBody struct {
        Message string `json:"body"`
    }
    type responseBody struct {
        Valid bool `json:"valid"`
    }

    data, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Println("Oh well, mistake reading response")
        return 
    }
    params := requestBody{}
    err = json.Unmarshal(data, &params)
    if err != nil {
        fmt.Println("Oh well, mistake unmarshaling")
        return 
    }

    resp := responseBody{}
    if len(params.Message) > charLimit {
        statusCode = 400
        resp.Valid = false
    }else {
        statusCode = 200
        resp.Valid = true 
    }
    jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
    w.WriteHeader(statusCode)
    w.Header().Set("Content-type", "application/json")
    w.Write(jsonResp)
}


func MiddlewareCors(next http.Handler) http.Handler {
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

type ApiState struct {
    ViewCount int
}


// Show number of views
func (s *ApiState) ShowPageViews(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    displayInfo := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",s.ViewCount)
    w.Write([]byte(displayInfo))
}


// increment number of page views
func (s *ApiState) MiddlewareMetrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        s.ViewCount++
        next.ServeHTTP(w, req)
})
}

// Reset page view count to zero
func (s *ApiState) ResetViews(w http.ResponseWriter, req *http.Request) {
    s.ViewCount = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}


func CheckStatus(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(http.StatusText(http.StatusOK)))
    }


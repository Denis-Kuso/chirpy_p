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

package handlers

import (
    "net/http"
    "log"
    "encoding/json"
    "fmt"
    //"io"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {

/*Add a new endpoint to the Chirpy API that accepts a POST request at /api/validate_chirp. It should expect a JSON body of this shape:
{
    "body": "This is an opinion I need to share with the world"
}*/
// example
//func handler(w http.ResponseWriter, r *http.Request){
    // ...
    /* decode incoming body data into a struct containing only the text
    I ll ignore other data for now
    give feedback irespective of status code in form of json*/
    const charLimit = 14
    var statusCode int
    resp := make(map[string]string)
    req := make(map[string]string)
    req["message"] = "HHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendHello My friendello My friend"
    var feedback string
    if len(req["message"]) > charLimit {
        statusCode = 400
        feedback = "Too long my friend"
    }else {
        statusCode = 200
        feedback = "all good"
    }
    resp["message"] = feedback
    jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
    w.WriteHeader(statusCode)
    w.Header().Set("Content-type", "application/json")
    w.Write(jsonResp)
}

package handlers

import (
    "net/http"
    "log"
    "encoding/json"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {

//At Chirpy, we have a silly rule that says all Chirps must be 140 characters long or less.

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
    const charLimit = 140
    var statusCode int
    type clientData struct {
        text string
    }
    type userResponse struct {
        status bool `json:"valid"`
    }
    dec := json.NewDecoder(r.Body)
    cData := clientData{}
    err := dec.Decode(&cData)
    if err != nil {
        errMessage := "err: something went wrong"
        log.Printf("Error decoding parameters: %s", err)
        w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
        w.Write([]byte(string(errMessage)))
		return
    }
    resp := userResponse{}
    err = json.NewEncoder(w).Encode(resp)
    if err != nil {
                http.Error(w, err.Error(), 500)
                return
    }
    // exceed limit
    if len(cData.text) > charLimit {
        resp.status = false
        statusCode = 400
    } else {
        statusCode = 200
        resp.status = true
    }
    data, _ := json.Marshal(resp)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(data)
}
    


//func handler(w http.ResponseWriter, r *http.Request){
//    // ...
//
//    type returnVals struct {
//        // the key will be the name of struct field unless you give it an explicit JSON tag
//        CreatedAt time.Time `json:"created_at"`
//        ID int `json:"id"`
//    }
//    respBody := returnVals{
//        CreatedAt: time.Now(),
//        ID: 123,
//    }
//    dat, err := json.Marshal(respBody)
//	if err != nil {
//		log.Printf("Error marshalling JSON: %s", err)
//		w.WriteHeader(500)
//		return
//	}
//    w.Header().Set("Content-Type", "application/json")
//    w.WriteHeader(200)
//	w.Write(dat)
//}

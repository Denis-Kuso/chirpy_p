package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"github.com/Denis-Kuso/chirpy_p/handlers"
	"github.com/Denis-Kuso/chirpy_p/internal/auth"
	"github.com/Denis-Kuso/chirpy_p/internal/database"
	"github.com/go-chi/chi/v5"
)
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func (s *ApiState) UpdateUser(w http.ResponseWriter, r *http.Request){
    // parse request
    type loginRequest struct {
     Email string `json:"email"`
     Password string `json:"password"`
    }
    data, reqErr := io.ReadAll(r.Body)
    if reqErr != nil {
	respondWithError(w, http.StatusBadRequest,"Bad request my friend")
	return
    }

    reqData := loginRequest{}
    err := json.Unmarshal(data,&reqData)
    if err!= nil {
	respondWithError(w, http.StatusInternalServerError, "Our bad man")
	return
    }

    // is JWT present?
    // TODO YOU can do better than this(What if there is no prefix????)
//    token := r.Header.Get("Authorization")
//    // strip Bearer from key
//    prefix := "Bearer "
//    token = strings.TrimPrefix(token, prefix)
    token, tokErr := auth.GetBearerToken(r.Header)
    if tokErr != nil {
	fmt.Printf("ERROR:%v\n",tokErr)
	return
    }
    //fmt.Printf("Actual authorization data: %s\n",token)

    // is JWT valid/in date?
    id, err := ValidateJWT(token,s.Token)
    if err!= nil {
	respondWithError(w, http.StatusUnauthorized,"Sorry mate, I don't believe you")
	return 
    }
    // else proceed with update
    type response struct {
	Email string `json:"email"`
	Id int `json:"id"`
    }
    intID, convErr := strconv.Atoi(id)
    if convErr != nil {
	return}
    user, usrErr := s.DB.GetUser(intID)

    if usrErr != nil{
	respondWithError(w,http.StatusInternalServerError,"We messed up")// IS THIS THE RIGHT ERROR TO USE
	return 
    }else{
	user,err = s.DB.UpdateUser(intID, reqData.Email, reqData.Password)
	if err != nil{
	    fmt.Printf("could not update user:%v\n",err)
	    return
	}
	log.Printf("Created user: %s\n",user.Email)
	respondWithJSON(w,http.StatusOK, response{
	    Email: user.Email,
	    Id : user.Id})
    return
    }
}

func (s *ApiState) LoginUser(w http.ResponseWriter, r *http.Request){
type loginRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
    ExpTime int `json:"expires_in_seconds"`
    }
    type response struct {
	Email string `json:"email"`
	Id int `json:"id"`
	Token string `json:"token"`
	Rtoken string `json:"refresh_token"`
    }
    data, err := io.ReadAll(r.Body)
    if err != nil {
	respondWithError(w, http.StatusInternalServerError,"Can't login")
	return
    }
    
    //does user exist?
    users,err := s.DB.GetUsers()
    if err != nil {
	respondWithError(w, http.StatusInternalServerError, "Sorry mate")
	return
    }
    reqData := loginRequest{}
    userData := database.User{}
    err = json.Unmarshal(data,&reqData)
    if err!= nil{
	respondWithError(w, http.StatusInternalServerError, "Sorry mate")
	return
    }

    // TODO abstract to getUSER
    for _, user := range users{
	if user.Email == reqData.Email{
	    userData = user
	    break
	}
    }
    // user does not exist
    if userData.Email == "" {
	respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
	return
    }

    // do the passwords match?
    err = auth.ValidatePswd(userData.Password,userData.Salt+reqData.Password)
    if err != nil {
	respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
	return
    }else{
	stoken := CreateAccessToken(userData.Id,s.Token)
	rtoken := CreateRefreshToken(userData.Id,s.Token)
	log.Printf("Logged in user: %s\n",userData.Email)
	respondWithJSON(w, http.StatusOK,response{
	    Email: userData.Email,
	    Id: userData.Id,
	    Token: stoken,
	    Rtoken: rtoken,
	})
    }  
}


func (s *ApiState) CreateUser(w http.ResponseWriter, r *http.Request){

    type response struct{
	Email string `json:"email"`
	Id int `json:"id"`
    }
    data, err := io.ReadAll(r.Body)
    if err != nil {
	respondWithError(w, http.StatusInternalServerError,"Can't create user")
	return
    }
    // TODO refactor instantiation of request data
    reqData := database.User{}
    err = json.Unmarshal(data, &reqData)
    if err != nil {
	respondWithError(w, http.StatusInternalServerError, "Err during json processing")
	return
    }
    users, DBerr := s.DB.GetUsers()
    if DBerr != nil {
	respondWithError(w, http.StatusInternalServerError,"Error at our end")
	return
    }
    for _, usr := range users {
	if usr.Email == reqData.Email {
	    respondWithError(w, http.StatusBadRequest, "User already exists")
	    return
	}
    }
    newUser,dberr := s.DB.CreateUser(reqData.Email, reqData.Password)
    if dberr != nil {
	respondWithError(w, http.StatusInternalServerError,"Cannot create user")
	return
    }
    respondWithJSON(w, http.StatusCreated,response{
		    Email: newUser.Email,
		    Id: newUser.Id})
    return
}

func (s *ApiState) GetChirp(w http.ResponseWriter, r *http.Request){
    // read id
    desiredID := r.URL.Path
    desiredID = chi.URLParam(r,"chirpID")
    id, err := strconv.Atoi(desiredID)

    // if valid
    if err != nil {
	respondWithError(w, http.StatusBadRequest, "Can't recognize id")
	return
    }
    // check if chirp exists
    chirps,DBerr := s.DB.GetChirps()
    if DBerr != nil {
	respondWithError(w, http.StatusInternalServerError, "Cant read chirps")
	return
    }
    // TODO sort by ID to reduce redundant iteration
    if id <= len(chirps){
	for idx, chirp := range chirps{
	    if chirp.Id == id {
		respondWithJSON(w, http.StatusOK,chirps[idx])
		break
	    }
	}
	return
    }else {
	respondWithError(w, http.StatusNotFound, "Can't find chirp")
	return
    }
}


func (s *ApiState) GetChirps(w http.ResponseWriter, r *http.Request){
    // read from db and return []Chirp
    chirps, err := s.DB.GetChirps()
    if err!=nil{
	respondWithError(w, http.StatusInternalServerError,"Error retrieving chirps")
    }
    respondWithJSON(w, http.StatusOK, chirps)
}

func (s *ApiState) ValidateChirp(w http.ResponseWriter, r *http.Request) {

    const charLimit = 140
    type requestBody struct {
        Message string `json:"body"`
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

    if len(params.Message) > charLimit {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }else {
        params.Message = handlers.FilterText(params.Message)
        // Increment chirp number
	chirp, chiErr := s.DB.CreateChirp(params.Message)
	if chiErr != nil {
	    respondWithError(w, http.StatusInternalServerError,"Can't create chirp")
	    return
	}
        respondWithJSON(w, http.StatusCreated, Chirp{
            Body: chirp.Body, 
            Id: chirp.Id},
        )
    }
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


type Chirp struct {
    Body string `json:"body"`
    Id int `json:"id"`
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


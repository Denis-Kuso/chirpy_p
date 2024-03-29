package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
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


func(s *ApiState) UpgradeMembership(w http.ResponseWriter, r *http.Request){
    type request struct{
	Event string `json:"event"`
	Data map[string]int `json:"data"`
    }
    desiredEvent := "user.upgraded"
    data, reqErr := io.ReadAll(r.Body)
    if reqErr != nil {
	respondWithError(w, http.StatusUnauthorized,"Bad request my friend")
	return
    }

    reqData := request{}
    err := json.Unmarshal(data,&reqData)
    if err!= nil {
	respondWithError(w, http.StatusInternalServerError, "Our bad man")
	return
    }
    // check event
    // if not user.upgraded return 200
    apiKey, err := auth.GetApiKey(r.Header)
    if err != nil {
	log.Printf("ERROR:%v\n",err)
	respondWithError(w,http.StatusUnauthorized, "No key, no entry")
	return
    }

    if apiKey != s.WebhookKey{
	log.Printf("Webhook attempted: %v\n",apiKey)
	respondWithError(w, http.StatusUnauthorized,"Invalid API key")
	return
    }
    if reqData.Event == desiredEvent {
	userID:= reqData.Data["user_id"]// TODO consider adding "ok" for mroe robustness
	_, err := s.DB.GetUser(userID)
	if err != nil{
	    log.Printf("Could not find user %v\n",userID)
	    respondWithError(w, http.StatusNotFound, "Cannot find user")
	    return
	}
	err = s.DB.MakeUserRed(userID)
	if err != nil {
	respondWithError(w, http.StatusInternalServerError, "Our bad man")
	return
	}
    }
    respondWithJSON(w,http.StatusOK,"")
    return
}


func (s *ApiState) RemoveChirp(w http.ResponseWriter, r *http.Request){
    // is the person authenticated
    token, tokErr := auth.GetBearerToken(r.Header)
    if tokErr != nil {
	log.Printf("ERROR:%v\n",tokErr)
	respondWithError(w,http.StatusBadRequest,"No token, no entry")
	return
    }
    // read id
    desiredChirp := r.URL.Path
    desiredChirp = chi.URLParam(r,"chirpID")
    userID, err := ValidateJWT(token,s.Token,"chirpy-access")
    if err!=nil{
	log.Printf("TOKEN ERR: %v for user:%v\n", err, userID)
	respondWithError(w,http.StatusUnauthorized,"invalid token")
	return
    }
    // does tweet exists
    desiredChirpINT, err := strconv.Atoi(desiredChirp)
    chirp, err := s.DB.GetChirpByID(desiredChirpINT)
    if err != nil {
	respondWithError(w,http.StatusNotFound,"Cannot find chirp")
	return
    }
    // is user author of chirp
    if strconv.Itoa(chirp.Author) == userID {
	log.Printf("About to delete chirp\n")
	s.DB.RemoveChirp(desiredChirpINT)
	respondWithJSON(w, http.StatusOK,"")
	return
    }
    respondWithError(w,http.StatusForbidden,"Not allowed to delete this.")
    return
}

func (s *ApiState) RevokeToken(w http.ResponseWriter, r *http.Request){
    token, tokErr := auth.GetBearerToken(r.Header)
    if tokErr != nil {
	log.Printf("ERROR:%v\n",tokErr)
	respondWithError(w,http.StatusBadRequest,"No token, no entry")
	return
    }
    err := s.DB.StoreRevokedToken(token)
    if err != nil {
	log.Printf("Err storing revoked token: %v\n",err)
	respondWithError(w, http.StatusInternalServerError,"")
	return
    }
    respondWithJSON(w, http.StatusOK,"")
    return
}


func (s *ApiState) RefreshToken(w http.ResponseWriter, r *http.Request){
    type response struct {
	Token string `json:"token"`
    }

    // check request data (prologoue)
    token, tokErr := auth.GetBearerToken(r.Header)
    if tokErr != nil {
	log.Printf("ERROR:%v\n",tokErr)
	respondWithError(w,http.StatusBadRequest,"No token, no entry")
	return
    }
    issuer := "chirpy-refresh"
    // is JWT valid
    userID, err := ValidateJWT(token, s.Token, issuer)
    if err != nil {
	log.Printf("TOKEN ERR: %v for user:%v\n",err,userID)
	respondWithError(w,http.StatusUnauthorized,"invalid token")
	return
    }
    //  are there no revokations
    ok, err := s.DB.IsRevoked(token)
    if err!= nil {
	log.Print(err)
	respondWithError(w,http.StatusInternalServerError, "our bad man")
	return
    }
    // only then return 200 and new access token
    if !ok{
	ID, serr := strconv.Atoi(userID)
	if serr!= nil {
	    log.Printf("ERR: %v during conversion of: %v\n",serr,userID)
	}
        newToken := CreateAccessToken(ID, s.Token)
        respondWithJSON(w, http.StatusOK, response{
    	Token: newToken,
    	})
	return
    }else {
	respondWithError(w,http.StatusUnauthorized,"")
	return
    }

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
    token, tokErr := auth.GetBearerToken(r.Header)
    if tokErr != nil {
	fmt.Printf("ERROR:%v\n",tokErr)
	return
    }

    log.Printf("Recevied update request with :%v\n",reqData)
    updateIssuer := "chirpy-access"
    // is JWT valid/in date?
    id, err := ValidateJWT(token,s.Token, updateIssuer)
    if err!= nil {
	log.Printf("User %v provided invalid token: %v\n", reqData, err)
	respondWithError(w, http.StatusUnauthorized,"Sorry mate, I don't believe you")
	return 
    }
    log.Printf("From user: %v, and token: %v, got id: %v\n",reqData,token,id)
    // else proceed with update
    type response struct {
	Email string `json:"email"`
	Id int `json:"id"`
    }
    intID, convErr := strconv.Atoi(id)
    if convErr != nil {
	log.Printf("ERR during converting str: %s to int\n",id)
	return
    }
    user, usrErr := s.DB.GetUser(intID)

    if usrErr != nil{
	log.Printf("ERR during fetching user:%v from DB: %v\n",user,usrErr)
	respondWithError(w,http.StatusInternalServerError,"We messed up")// IS THIS THE RIGHT ERROR TO USE
	return 
    }else{
	user,err = s.DB.UpdateUser(intID, reqData.Email, reqData.Password)
	if err != nil{
	    fmt.Printf("could not update user:%v\n",err)
	    return
	}
	log.Printf("Updated user: %s\n",user.Email)
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
	IsRed bool `json:"is_chirpy_red"`
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
	    IsRed: userData.IsRed,
	    Token: stoken,
	    Rtoken: rtoken,
	})
    }  
}


func (s *ApiState) CreateUser(w http.ResponseWriter, r *http.Request){

    type response struct{
	Email string `json:"email"`
	Id int `json:"id"`
	IsRed bool `json:"is_chirpy_red"`
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
		    Id: newUser.Id,
		    IsRed: false})
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
    // is there a query

    query := r.URL.Query().Get("author_id")
    sortFlag := r.URL.Query().Get("sort")
    // on empty query return all chirps (sorted)
    log.Printf("Query provided author_id:%v and sort_flag: %s\n", query, sortFlag)
    chirps, err := s.DB.GetChirps()

    if err!=nil{
    	respondWithError(w, http.StatusInternalServerError,"Error retrieving chirps")
    }
    slices.SortFunc(chirps, func(a, b database.Chirp) int {
		if sortFlag == "desc" {
		    a,b = b, a
		}
		if n := cmp.Compare(a.Id, b.Id); n != 0 {
			return n
		}
		// If chirp_IDs are equal, order by author_id
		return cmp.Compare(a.Author, b.Author)
	})
    if len(query) == 0 {
        respondWithJSON(w, http.StatusOK, chirps)
	return
    }
    var users_chirps []database.Chirp
    queryID, err := strconv.Atoi(query)
    if err != nil {
	log.Printf("ERR: %v, user tried to to find author_id with query: %v\n",err,query)
	respondWithError(w,http.StatusBadRequest, "Invalid query")
	return
    }
    for _, chirp := range chirps {
	if chirp.Author == queryID {
	    users_chirps = append(users_chirps,chirp)
	}
    }
    if len(users_chirps) > 0 {
	respondWithJSON(w, http.StatusOK, users_chirps)
	return
    }else {
	respondWithError(w, http.StatusNotFound, "No chirps founds for author_ID")
	return
    }
}

func (s *ApiState) ValidateChirp(w http.ResponseWriter, r *http.Request) {

    const charLimit = 140
    type requestBody struct {
        Message string `json:"body"`
    }
    token, tokErr := auth.GetBearerToken(r.Header)
    if tokErr != nil {
	fmt.Printf("ERROR:%v\n",tokErr)
	return
    }
    log.Printf("Token after extraction:%v\n", token)

    // is user allowed to post
    userID, err := ValidateJWT(token, s.Token, "chirpy-access")
    if err != nil {
	log.Printf("Token: %v\n",token)
	log.Printf("User: %v, unauthorised attempt: %v\n",userID, err)
	respondWithError(w, http.StatusUnauthorized,"")
	return
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
	intID, _ := strconv.Atoi(userID)//TODO HANDLE ERROR
	chirp, chiErr := s.DB.CreateChirp(params.Message, intID)
	if chiErr != nil {
	    respondWithError(w, http.StatusInternalServerError,"Can't create chirp")
	    return
	}
        respondWithJSON(w, http.StatusCreated, Chirp{
	    AuthorId: intID,
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
    AuthorId int `json:"author_id"`
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


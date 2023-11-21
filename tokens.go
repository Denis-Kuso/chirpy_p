package main

import (
    "github.com/golang-jwt/jwt/v5"
    "time"
    "errors"
    "fmt"
    "strconv"
)


func CreateRefreshToken(userID int, key string) string {
    // calls create userToken with different issuer
    issuer := "chirpy-refresh"
    expiresAt := 60*24*60 //60 days in minutes
    token := createUserToken(userID, expiresAt, issuer, key)
    return token
}



// Create access token derived from 
func CreateAccessToken(userID int, key string) string {
    issuer := "chirpy-access"
    expiresAt := 60//minutes
    token := createUserToken(userID, expiresAt, issuer, key)
    return token
}

// Create custom token for user, expTime in minutes
func createUserToken(userID int, expTime int, issuer string, key string) string {
    // TODO handle errors from mutating methods
    var t *jwt.Token
    c := jwt.RegisteredClaims{
	    ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expTime)* time.Hour)),
	    IssuedAt: jwt.NewNumericDate(time.Now()),
	    Issuer: issuer,
	    Subject: strconv.Itoa(userID),
	    }
    

    t = jwt.NewWithClaims(jwt.SigningMethodHS256,c)

    // sign token
    signedKey, err := t.SignedString([]byte(key))
    if err != nil {
	fmt.Printf("Some error during signing your token man: %v\n",err)
	return ""
    }
    return signedKey 
}

func ValidateJWT(userToken, tokenSecret, issuer string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		userToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	ClaimIssuer, err := token.Claims.GetIssuer()
	if ClaimIssuer != issuer {
	    return "",errors.New("Invalid issuer")
	}
	return userIDString, nil
}


// Modifies expirationTime if provided and smaller than 24h
// should time be of type int or time.Second
func (c *myCustomClaims) setExpirationTime(expTime int) error {
    const S_DAY time.Duration = 24 * 60 * 60
    if (expTime > 0) && (time.Duration(expTime) * time.Second < S_DAY * time.Second) {
	c.Claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(expTime) * time.Second))}
    return nil // TODO modify
}

// customise subject
func (c *myCustomClaims) setSubject(id int) error {
    c.Claims.Subject = fmt.Sprint(id)
    return nil//TODO modify
}


type myCustomClaims struct{
    Claims jwt.RegisteredClaims
}


// abstract into function
// what should the return type be?a pointer to a myCustomClaims instance?
//func getDefaultClaims() myCustomClaims{
//    defaultClaims := myCustomClaims{
//	jwt.RegisteredClaims{
//	    ExpiresAt: jwt.NewNumericDate(time.Now().Add(24* time.Hour)),
//	    IssuedAt: jwt.NewNumericDate(time.Now()),
//	    Issuer: "chirpy",
//	    Subject: "userID",},
//    }
//return defaultClaims
//}

package main

import (
    "github.com/golang-jwt/jwt/v5"
    "time"
    "fmt"
)

// Create custom token for user
func CreateUserToken(userID int, expTime int, key string) string {
    // TODO handle errors from mutating methods
    var t *jwt.Token
//    // customise expectation field if provided
//     c.setExpirationTime(expTime)
//
    c := getDefaultClaims()
    // set subject id
    c.setSubject(userID)
    // modify exp time if needed 
    c.setExpirationTime(expTime)
    t = jwt.NewWithClaims(jwt.SigningMethodHS256,c.Claims)

    // sign token
    signedKey, err := t.SignedString([]byte(key))// how should the key be accessed?Parameter to function?
    if err != nil {
	fmt.Printf("Some error during signing your token man: %v\n",err)
	return ""
    }
    return signedKey 
}

// Modifies expirationTime if provided and smaller than 24h
// should time be of type int or time.Second
func (c *myCustomClaims) setExpirationTime(expTime int) error {
    const S_DAY time.Duration = 24 * 60 * 60
    if time.Duration(expTime) * time.Second < S_DAY * time.Second {
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
func getDefaultClaims() myCustomClaims{
    defaultClaims := myCustomClaims{
	jwt.RegisteredClaims{
	    ExpiresAt: jwt.NewNumericDate(time.Now().Add(24* time.Hour)),
	    IssuedAt: jwt.NewNumericDate(time.Now()),
	    Issuer: "chirpy",
	    Subject: "userID",},
    }
return defaultClaims
}


// Create default claims settings
//func newClaims() myCustomClaims {
//var claims myCustomClaims; 
//    claims.Claims
//    ExpiresAt: jwt.NewNumericDate(time.Now().Add(24* time.Hour)),
//    IssuedAt: jwt.NewNumericDate(time.Now()),
//    Issuer: "chirpy",
//    Subject: "userID",}

// Create claims with multiple fields populated
//t := jwt.RegisteredClaims{
//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
//		IssuedAt:  jwt.NewNumericDate(time.Now()),
//		Issuer:    "chirpy",
//		Subject:   "userId",
//	},
//}


// Create claims while leaving out some of the optional fields
//claims = MyCustomClaims{
//	"bar",
//	jwt.RegisteredClaims{
//		// Also fixed dates can be used for the NumericDate
//		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
//		Issuer:    "test",
//	},
//}

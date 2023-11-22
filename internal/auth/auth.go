package auth

import (
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "errors"
    "strings"
)

// Create hash of password
func GeneratePswd(p string) (string, error) {
bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
    if err != nil {
        return string(bs), err
    }
    return string(bs),nil
}


// Parse authorisation header and extract token
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No auth header")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// Parse authorisation header and extract key 
func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No auth header")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}


// validate a password
// Compare the hash of the given password with the hash from the database. If they match, the password is correct. Otherwise, the password is incorrect.
func ValidatePswd(hashedPswd string, plainPswd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPswd), []byte(plainPswd))
	if err != nil {
		return err
	}
    	return nil
}

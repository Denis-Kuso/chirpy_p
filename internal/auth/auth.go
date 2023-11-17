package auth

import (
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "errors"
    "strings"
)

// Create hash of password
//Generate a long random salt using a CSPRNG.
//Prepend the salt to the password and hash it with a standard password hashing function like bcrypt
//Save both the salt and the hash in the user's database record
func GeneratePswd(p string) (string, error) {
bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
    if err != nil {
        return string(bs), err
    }
    return string(bs),nil
}


// create salt for password
// USE SOMETHING ELSE, perhaps crypto/ran
//func createSalt(u string) (string, error) {
//    err, salt := generatePswd(u + "chirp")
//    if err!=nil{
//	return string(salt), nil
//    }
//    return string(salt),nil
//}


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

// validate a password
// Retrieve the user's salt and hash from the database.
// Prepend the salt to the given password and hash it using the same hash function.
// Compare the hash of the given password with the hash from the database. If they match, the password is correct. Otherwise, the password is incorrect.
func ValidatePswd(hashedPswd string, plainPswd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPswd), []byte(plainPswd))
	if err != nil {
		return err
	}
    	return nil
}

//use in handler as auth.ValidatePswd(user.pswd, user.salt+req.paswd)

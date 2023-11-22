package database

import (
    "log"
    "time"
)

// Saves token to database
func (db *DB) StoreRevokedToken(token string) error {
    dbStructure, err := db.loadDB()
	if err != nil {
	    log.Printf("Err during db.loadDB(): %v\n",err)	
	    return ErrReadingDB 
	}

	dbStructure.RevokedTokens[token] = time.Now()
	err = db.writeDB(dbStructure)
	if err != nil {
		log.Printf("ERR during writting token to db:%v\n",err)
		return ErrReadingDB
	}
	return  nil
}


// Search database to see if token has been previously revoked.
// Return true iff token is found, false otherwise
func (db *DB) IsRevoked(token string) (bool, error) {
    tokens,err := db.getRevokedTokens()
    if err!= nil{
	log.Println(err)
	return false, ErrReadingDB  
    }
    for tok := range tokens{
	if tok == token{
	    return true, nil
	}
    }
    return false,nil
}


func (db *DB) getRevokedTokens() (map[string]time.Time, error) {
    var tokens map[string]time.Time

    dbStructure, err := db.loadDB()
    if err != nil {
	return tokens, err
    }
    tokens = dbStructure.RevokedTokens
    return tokens,nil 
}

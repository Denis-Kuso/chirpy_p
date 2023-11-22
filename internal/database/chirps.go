package database

import (
    "log"
)


type Chirp struct {
    Author int `json:"author_id"`
    Body string `json:"body"`
    Id int `json:"id"`
}


// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error){
    dbStructure, err := db.loadDB()
	if err != nil {
	    log.Printf("Err during db.loadDB(): %v\n",err)	
	    return Chirp{}, ErrReadingDB
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Author: authorId,
		Id:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}



func (db *DB) RemoveChirp(id int)  error {
    dbStructure, err := db.loadDB()
    if err != nil {
	return ErrReadingDB
    }
    delete(dbStructure.Chirps, id)
    err = db.writeDB(dbStructure)
    if err != nil {
	return err
    }
    return nil
}


func (db *DB) GetChirpByID(id int) (Chirp, error) {
    data, loadErr := db.loadDB()
    if loadErr != nil {
	//debug log
	log.Printf("ERR during loading DB:%v\n",loadErr)
        return Chirp{}, ErrReadingDB
    }
    if id > len(data.Chirps){
	return Chirp{}, ErrNotExist
    }
    chirp, found:= data.Chirps[id]
    if !found {
	return Chirp{},ErrNotExist
    }
    return chirp, nil
}


// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
    data, loadErr := db.loadDB()
    if loadErr != nil {
        return nil, loadErr
    }
    var chirps []Chirp
    for _, chirp := range data.Chirps {
        chirps = append(chirps, chirp)
    }

    return chirps, nil
}

package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
	"log"
)

var ErrNotExist = errors.New("does not exist")
var ErrReadingDB = errors.New("database issues")

type DB struct {
	path string
	mu  *sync.RWMutex
}


type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
	RevokedTokens map[string]time.Time`json:"revoked_tokens"`
}


func NewDB(path string) (*DB, error){
    db := DB{
        path: path,
        mu:   &sync.RWMutex{},
        // initialise with unlocked RWmutex
    }
    err := db.ensureDB()
	return &db, err
}


func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
		RevokedTokens: map[string]time.Time{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
    // create a file if it does not exist
    _, err := os.Stat(db.path)
    if errors.Is(err, os.ErrNotExist) {
        return db.createDB()
    }
    return err 
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
    db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		log.Printf("Err during unmarshaling: %v\n",err)
		return dbStructure, err
	}

	return dbStructure, nil
}


// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
    db.mu.Lock()
    defer db.mu.Unlock()
    data, err := json.Marshal(dbStructure)
    if err != nil {
        return err
    }
    err = os.WriteFile(db.path, data, 0600)
    if err != nil {
        return err
    }

    return nil
}

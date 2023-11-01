package database

import (
    "sync"
    "fmt"
    "os"
    "errors"
    "encoding/json"
)


type DB struct {
	path string
	mu  *sync.RWMutex
}

type Chirp struct {
    Body string `json:"body"`
    Id int `json:"id"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){
    dbStructure, err := db.loadDB()
	if err != nil {
	    fmt.Printf("Err during db.loadDB(): %v\n",err)	
	    return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
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

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
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
		fmt.Printf("Err during unmarshaling: %v\n",err)
		return dbStructure, err
	}
	fmt.Println(dbStructure)

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


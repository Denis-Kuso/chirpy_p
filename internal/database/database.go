package database

import (
    "sync"
    "fmt"
    "os"
    "errors"
)


type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
    Body string
    Id int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}


func NewDB(path string) (*DB, error){
    db := DB{
        path: path,
        // initialise with unlocked RWmutex
    }
    return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error){
    // how to determine id? BY length of DBStructure.Chirps?
    return Chirp{}, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
    data, loadErr := db.loadDB()
    if loadErr != nil {
        fmt.Printf("Error during loading: %v\n", loadErr)
        return nil, loadErr
    }
    var chirps []Chirp
    for _, chirp := range data.Chirps {
        chirps = append(chirps, chirp)
    }

    return chirps, nil
}


func (db *DB) ensureDB() error {
    // create a file if it does not exist
    if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist) {
        fmt.Println(err)
        file,_ := os.Create(db.path)
        fmt.Println("created",file)
        return nil
    }
        //if _, createErr := os.Create(db.path){
            //return createErr
        //}
    return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
//    err := db.ensureDB()
//    if err != nil {
//        fmt.Printf("Err loading database: %v\n",err)
//        return DBStructure{}, err
//    }
//    data, readErr := os.ReadFile(db.path)
//	if err != nil {
//		log.Fatal(err)
//        return DBStructure{}, readErr
//	}
//    return DBStructure{string(data)}, nil
return DBStructure{make(map[int]Chirp)},nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
    return nil
}


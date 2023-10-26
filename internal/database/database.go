package database

import (
    "sync"
    "fmt"
    "os"
)


type DB struct {
	path string
	mux  *sync.RWMutex
}



type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}


type Chirp struct {
    Body string `json:"body"`
    id int `json:"id"`
}


func (db *DB) CreateChirp(body string) (Chirp, error){

}


func (db *DB) GetChirps() ([]Chirp, error) {

}


func (db *DB) ensureDB() error {
}


func (db *DB) loadDB() (DBStructure, error) {

}


func (db *DB) writeDB(dbStructure DBStructure) error {
}


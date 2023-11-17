package database

import (
    "sync"
    "fmt"
    "os"
    "errors"
    "encoding/json"
    "github.com/Denis-Kuso/chirpy_p/internal/auth"
)
var ErrNotExist = errors.New("does not exist")

type DB struct {
	path string
	mu  *sync.RWMutex
}

type Chirp struct {
    Body string `json:"body"`
    Id int `json:"id"`
}

// TODO consider making user a private type
type User struct {
    Email string `json:"email"`
    Password string `json:"password"`
    Id int `json:"id"`
    Salt string
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
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


// Create user creates a new user and saves to disk
func (db *DB) CreateUser(body string, pswd string) (User, error){
    dbStructure, err := db.loadDB()
	if err != nil {
	    fmt.Printf("Err during db.loadDB(): %v\n",err)	
	    return User{}, err
	}

	id := len(dbStructure.Users) + 1
	userSalt,sErr := auth.GeneratePswd(body)
	if sErr != nil {
	    return User{},err
	}
	hashedPswd, err := auth.GeneratePswd(userSalt + pswd)
	if err != nil {
	    return User{},err
	}
	user := User{
		Id:   id,
		Email: body,
		Password: string(hashedPswd),
		Salt: userSalt,
	}
	dbStructure.Users[id] = user 
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}


func (db *DB) UpdateUser(ID int, email string, pswd string) (User, error){
    dbStructure, err := db.loadDB()
    if err != nil {
	return User{}, err
    }
    // does user exist?
    //FIND BY ID user, err := db.GetUserByEmail(email)
    user, err := db.GetUser(ID)
    // find user
    if err != nil {
	fmt.Printf("Cannot update user:%v. ERR:%v\n",email,err)
	return User{}, err
    }
    //id := user.Id// extract ID from email

    // updated desired credentials
    // TODO duplicated code - try abstracting into routine
    userSalt,sErr := auth.GeneratePswd(email)
    if sErr != nil {
        return User{},err
    }
    hashedPswd, err := auth.GeneratePswd(userSalt + pswd)
    if err != nil {
        return User{},err
    }
    user = User{
    	Id:   ID,
    	Email: email,
    	Password: string(hashedPswd),
    	Salt: userSalt,
    }

    // write to DB
    dbStructure.Users[ID] = user
    err = db.writeDB(dbStructure)
    if err != nil {
	return User{}, err
    }
    return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
	return User{}, err
    }

    user, ok := dbStructure.Users[id]
    if !ok {
    	return User{}, ErrNotExist
    }
    return user, nil
}


func (db *DB) GetUserByEmail(email string) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
	return User{}, err
    }

    // find user
    for _, user := range dbStructure.Users {
	if user.Email == email {
	    return user, nil
	}
    }
    return User{}, ErrNotExist
}


func (db *DB) GetUsers() ([]User, error) {
    data, loadErr := db.loadDB()
    if loadErr != nil {
        return nil, loadErr
    }
    var users []User
    for _, user := range data.Users{
        users = append(users, user)
    }
    return users, nil
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
		Users: map[int]User{},
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
	//fmt.Println(dbStructure)

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


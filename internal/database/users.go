package database

import (
    "log"
    "github.com/Denis-Kuso/chirpy_p/internal/auth"
)


// TODO consider making user a private type
type User struct {
    Email string `json:"email"`
    Password string `json:"password"`
    Id int `json:"id"`
    Salt string
    IsRed bool `json:"is_chirpy_red"`
}


func (db *DB) MakeUserRed(userID int) error {
//find user
    dbStructure, err := db.loadDB()
    if err != nil {
	log.Print(err)
	return ErrReadingDB 
    }
    user, err := db.GetUser(userID)
    if err != nil {
        log.Printf("Err fetching user : %v\n",err)//TODO REPLACE with logs	
        return ErrReadingDB 
    }
    user.IsRed = true
    dbStructure.Users[userID] = user
    // Write to db
    err = db.writeDB(dbStructure)
    if err != nil {
	return err
    }
    return nil
}

// Create user creates a new user and saves to disk
func (db *DB) CreateUser(body string, pswd string) (User, error){
    dbStructure, err := db.loadDB()
	if err != nil {
	    log.Printf("Err during db.loadDB(): %v\n",err)//TODO REPLACE with logs	
	    return User{}, ErrReadingDB
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
		IsRed: false,
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
	log.Print(err)
	return User{}, ErrReadingDB 
    }
    // does user exist?
    //FIND BY ID user, err := db.GetUserByEmail(email)
    user, err := db.GetUser(ID)
    // find user
    if err != nil {
	log.Printf("Cannot update user:%v. ERR:%v\n",email,err)
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
	IsRed: false,
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

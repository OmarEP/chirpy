package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotExit = errors.New("resource does not exit")

type DB struct {
	path string
	mux *sync.RWMutex
}

type Chirp struct {
	ID 		int `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exit 
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux: &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}


// CreateChirp creates a new chirp and saves it to disk 
func (db *DB) CreateChirp(body string)(Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err 
	}

	id := len(dbStructure.Chirps) + 1 
	chirp := Chirp{
		ID: id,
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
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err 
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
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

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err 
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok{
		return Chirp{}, ErrNotExit
	}
	return chirp, nil
}

// ensureDB creates a new database file if it doesn't exit
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err 
}


// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbsStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbsStructure, err 
	}
	err = json.Unmarshal(data, &dbsStructure)
	if err != nil {
		return dbsStructure, err 
	}
	return dbsStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, dat, 0666)
	if err != nil {
		return err
	}
	return nil
}



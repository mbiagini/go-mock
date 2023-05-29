package db

import (
	"encoding/json"
	"fmt"
	"go-mock/model"
	"io/ioutil"
	"os"
	"sync"

	"github.com/mbiagini/go-server-utils/gslog"
)

type EndpointDB interface {
	FindAll()              []model.Endpoint
	FindById(id int)       (model.Endpoint, bool)
	Save(e model.Endpoint) (model.Endpoint, error)
	Delete(id int)
}

// Thread-safe in-memory map of endpoints.
type endpointDB struct {
	sync.RWMutex
	m            map[int]model.Endpoint
	seq          int
}

// The one and only db instance.
var DB EndpointDB

func LoadDB() error {
	db := &endpointDB{
		m: make(map[int]model.Endpoint),
	}

	file, err := os.Open("./resources/db.json")
	if err != nil {
		if os.IsNotExist(err) {
			gslog.Server("File /resources/db.json not found. Creating new DB")
			err = db.saveToFile()
			DB = db
			return err
		} else {
			gslog.Server(fmt.Sprintf("Error reading file /resources/db.json: %s", err.Error()))
			return err
		}
	}

	var endpoints []model.Endpoint
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&endpoints)
	if err != nil {
		gslog.Server(fmt.Sprintf("Error parsing JSON from db file: %s", err.Error()))
		return err
	}
	
	for _, endpoint := range endpoints {
		db.m[endpoint.Id] = endpoint
	}
	DB = db

	return nil
}

// FindAll retrieves all endpoints.
func (db *endpointDB) FindAll() []model.Endpoint {
	db.RLock()
	defer db.RUnlock()

	if len(db.m) == 0 {
		return make([]model.Endpoint, 0)
	}

	copy := make([]model.Endpoint, len(db.m))
	i := 0
	for _, endpoint := range db.m {
		copy[i] = endpoint
		i++
	}

	return copy
}

// FindById returns the endpoint with the given id and a boolean indicating
// if it is present. If no endpoint is found, returns the zero value of Endpoint
// and a false.
func (db *endpointDB) FindById(id int) (model.Endpoint, bool) {
	db.RLock()
	defer db.RUnlock()

	endpoint, ok := db.m[id]
	return endpoint, ok
}

// Save saves a new endpoint using the internal sequence to generate the next id
// and returns this id.
func (db *endpointDB) Save(e model.Endpoint) (model.Endpoint, error) {
	db.Lock()
	defer db.Unlock()

	if e.Id != 0 {
		db.m[e.Id] = e
		err := db.saveToFile()
		return e, err
	}

	db.seq++
	e.Id = db.seq

	db.m[e.Id] = e
	err := db.saveToFile()
	return e, err
}

// Delete receives the id of an endpoint to be removed and deletes it from the
// internal map.
func (db *endpointDB) Delete(id int) {
	db.Lock()
	defer db.Unlock()

	delete(db.m, id)
}

func (db *endpointDB) saveToFile() error {

	endpoints := make([]model.Endpoint, 0, len(db.m))

	for _, endpoint := range db.m {
		endpoints = append(endpoints, endpoint)
	}

	// Convert endpoints to JSON.
	jsonData, err := json.Marshal(endpoints)
	if err != nil {
		return fmt.Errorf("error marshaling db to json: %s", err.Error())
	}

	// Save JSON to file.
	err = ioutil.WriteFile("./resources/db.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error saving JSON to file: %s", err.Error())
	}

	return nil
}
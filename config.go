package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
)

type Condition struct {
	Value      string `json:"value"`
	ResponseId int    `json:"response_id"`
}

type Discriminator struct {
	Parameter  string      `json:"parameter"`
	Location   string      `json:"location"`
	Conditions []Condition `json:"conditions"`
}

type Response struct {
	Id 			 int    `json:"id"`
	Code         int    `json:"code"`
	ContentType  string `json:"content_type"`
	BodyFilename string `json:"body_filename"`
	Delay        int    `json:"delay"`
}

type Endpoint struct {
	Method            string        `json:"method"`
	Path              string        `json:"path"`
	DefaultResponseId int           `json:"default_response_id"`
	HasDiscriminator  bool          `json:"has_discriminator"`
	Discriminator 	  Discriminator `json:"discriminator"`
	Responses         []Response    `json:"responses"`
}

type Config struct {
	Port      int        `json:"port"`
	Basepath  string     `json:"basepath"`
	Endpoints []Endpoint `json:"endpoints"`
}

var Conf *Config

func init() {

	// Default values.
	Conf = &Config{
		Port:     8080,
		Basepath: "/",
	}

	// Configuration filename.
	filename := "./resources/config.json"

	// Open config file.
	configFile, err := os.Open(filename)
	if (err != nil) {
		log.Println("Could not open config file. Make sure there is a config.json file in ./resources directory")
		log.Panicln(err)
	}

	// Parse config file content as JSON.
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(Conf)
	if (err != nil) {
		log.Println("Could not parse config file as JSON. Make sure the config.json file is JSON valid")
		log.Panicln(err)
	}

	// Validate config
	if (!Conf.isValid()) {
		log.Panicln("Validation failed for given config")
	}
}

func (c Config) isValid() bool {
	if (len(c.Endpoints) == 0) {
		log.Println("No endpoints received in configuration. A mock with no endpoints is not a mock to waste resources on")
		return false
	}
	for _, ep := range c.Endpoints {
		if (!validateEndpoint(ep)) {
			return false
		}
	}
	return true
}

func validateEndpoint(ep Endpoint) bool {

	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	validLocations := []string{"PATH", "QUERY", "HEADER"}

	if (!containsCaseInsensitive(validMethods, ep.Method)) {
		log.Println("Invalid method found in endpoint", ep.Path, ":", ep.Method)
		return false
	}

	if (len(ep.Responses) == 0) {
		log.Println("No response found for endpoint", ep.Path)
		return false
	}

	if (ep.HasDiscriminator) {
		d := ep.Discriminator
		if (!containsCaseInsensitive(validLocations, d.Location)) {
			log.Println("Invalid discriminator location in endpoint", ep.Path, ":", d.Location)
			return false
		}
		if (len(d.Conditions) == 0) {
			log.Println("No conditions found for discriminator of endpoint", ep.Path)
			return false
		}
		for _, c := range d.Conditions {
			if (!validateCondition(c, ep.Responses)) {
				log.Println("Validation of condition failed for endpoint", ep.Path)
				return false
			}
		}
	}
	
	for _, r := range ep.Responses {
		if (!validateResponse(r)) {
			log.Println("Invalid response found for endpoint", ep.Path)
			return false
		}
	}

	return true
}

func validateCondition(c Condition, rs []Response) bool {
	for _, r := range rs {
		if (c.ResponseId == r.Id) {
			return true
		}
	}
	log.Println("Could not find response for condition associated with response id", c.ResponseId, "and value", c.Value)
	return false
}

func validateResponse(r Response) bool {
	if (r.BodyFilename == "") {
		return true
	}
	if _, err := os.Stat(r.BodyFilename); errors.Is(err, os.ErrNotExist) {
		log.Println("Filename", r.BodyFilename, "not found for response with id", r.Id)
		return false
	}
	if (r.Delay < 0) {
		log.Println("Invalid delay value", r.Delay, "for response with id", r.Id)
		return false
	}
	return true
}

func containsCaseInsensitive(s []string, str string) bool {
	for _, v := range s {
		if (strings.EqualFold(v, str)) {
			return true
		}
	}
	return false
}
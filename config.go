package main

import (
	"encoding/json"
	"fmt"
	"go-mock/db"
	"os"

	"github.com/mbiagini/go-server-utils/gslog"
)

type Config struct {
	Ip            string              `json:"ip"`
	Port          int                 `json:"port"`
	Basepath      string              `json:"basepath"`
	UploadMaxSize int                 `json:"file_max_size"` 
	Logger        *gslog.LoggerConfig `json:"logger,omitempty"`
	LogFile       gslog.LogFileConfig `json:"log_file"`
}

var Conf Config

func (c *Config) ToString() string {
	json, _ := json.Marshal(c)
	return string(json)
}

func LoadConfiguration(f string) error {

	// Default values.
	var config = Config {
		Ip: "",
		Port:     8080,
		Basepath: "/",
		UploadMaxSize: 1,
	}

	// Open config file.
	file, err := os.Open(f)
	if (err != nil) {
		return fmt.Errorf("error opening configuration file %s: %s", f, err.Error())
	}

	// Parse config file content as JSON.
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&config)
	if (err != nil) {
		return fmt.Errorf("error parsing config file as JSON: %s", err.Error())
	}

	// Configure dependencies.
	if config.Logger != nil {
		gslog.ConfigureLog(*config.Logger)
	}
	gslog.ConfigureLogFile(config.LogFile)

	// Set global variable.
	Conf = config

	err = db.LoadDB()
	if err != nil {
		return fmt.Errorf("error loading DB: %s", err.Error())
	}

	gslog.Server(fmt.Sprintf("config: %s", config.ToString()))

	return nil
}
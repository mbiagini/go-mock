package main

import (
	"encoding/json"
	"fmt"
	"go-mock/model"
	"os"

	"github.com/mbiagini/go-server-utils/gslog"
)

type Config struct {
	Ip        string               `json:"ip"`
	Port      int                  `json:"port"`
	Basepath  string               `json:"basepath"`
	Logger    *gslog.LoggerConfig  `json:"logger"`
	LogFile   gslog.LogFileConfig `json:"log_file"`
	Endpoints []model.Endpoint     `json:"endpoints"`
}

var Conf Config

func LoadConfiguration(f string) error {

	// Default values.
	var config = Config {
		Ip: "",
		Port:     8080,
		Basepath: "/",
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

	return nil
}
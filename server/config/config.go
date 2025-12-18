package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	GuessingTime int `json:"guessingTime"`
	ResultTime   int `json:"resultTime"`
}

var AppConfig Config

func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return err
	}
	return nil
}

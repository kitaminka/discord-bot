package cfg

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var Config Configuration

func init() {
	LoadConfig()
}

func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panicf("Error opening file: %v", err)
	}

	byteValue, _ := io.ReadAll(file)
	err = file.Close()
	if err != nil {
		log.Printf("Error closing file: %v", err)
	}

	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		log.Panicf("Error unmarshalling config: %v", err)
	}
}

package utils

import (
	"encoding/json"
	"log"
)

func LogStructAsJSON(label string, v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error encoding %s: %v", label, err)
		return
	}
	log.Printf("%s: %s", label, string(jsonData))
}

func LogMapAsJSON(label string, v map[string]interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error encoding %s: %v", label, err)
		return
	}
	log.Printf("%s: %s", label, string(jsonData))
}

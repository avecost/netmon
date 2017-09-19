package main

import (
	"encoding/json"
	"os"
)

type AppUser struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Url      []string `json:"url"`
}

func loadUsers(file string) ([]AppUser, error) {
	var users []AppUser
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&users)
	return users, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type AppUser struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Url      []string `json:"url"`
}

func LoadUsers(file string) []AppUser {
	var users []AppUser
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&users)
	return users
}

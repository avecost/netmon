package main

import (
	"io/ioutil"
	"encoding/json"
	//"os"
)

// Franchisee struct contains name and Outlets
type Franchisee struct {
	Operator	string 			`json:"operator"`
	Outlets		[]struct {
		Name		string		`json:"name"`
		Terminals 	[]struct {
			Account 	string 		`json:"account"`
			Status		int 		`json:"status"`
			Privateip	string 		`json:"privateIp"`
			Publicip	string  	`json:"publicIp"`
			Os 			string		`json:"os"`
			Online		string		`json:"online"`
			Lastupdate	string		`json:"lastUpdate"`
		}	`json:"terminals"`
	}	`json:"outlets"`
}

func loadOutlet(file string) ([]Franchisee, error) {
	//pwd, _ := os.Getwd()
	//configJson, err := ioutil.ReadFile(pwd+file)
	configJson, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var instawin []Franchisee
	err = json.Unmarshal(configJson, &instawin)
	if err != nil {
		return nil, err
	}

	return instawin, nil
}

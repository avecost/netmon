package main

import (
	"net/http"
	"log"
)

func IsLoggedIn(r *http.Request) bool {
	session, err := AppSess.Get(r, "session")
	if err != nil {
		log.Println(err)
	}
	if session.Values["loggedIn"] == true {
		return true
	}
	return false
}

func Username(r *http.Request) string {
	session, err := AppSess.Get(r, "session")
	if err != nil {
		log.Println(err)
	}
	return session.Values["username"].(string)
}

package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func HashPassword(sPassword string) (string, error) {
	bPass, err := bcrypt.GenerateFromPassword([]byte(sPassword), 14)
	return string(bPass), err
}

func CheckPasswordHash(sPassword, sHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(sHash), []byte(sPassword))
	return err == nil
}

func ValidUser(sUser, sPassword string, slcUsers []AppUser) bool {
	for _, v := range slcUsers {
		if sUser == v.User && CheckPasswordHash(sPassword, v.Password) {
			return true
		}
	}
	return false
}

func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsLoggedIn(r) {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler(w, r)
	}
}

func UserAllowedURL(slcUsers []AppUser, sUser, sUrl string) bool {
	for _, v := range slcUsers {
		if v.User == sUser {
			for _, url := range v.Url {
				if url == sUrl {
					return true
				}
			}
		}
	}
	return false
}

func HomeUrl(users []AppUser, user string) string {
	if user == "ADMIN" {
		return "/dashboard"
	}
	for _, u := range users {
		if u.User == user {
			return u.Url[0] + "?o=" + user
		}
	}
	return ""
}
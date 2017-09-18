package main

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/avecost/authexample/config"
)

func HashPassword(sPassword string) (string, error) {
	bPass, err := bcrypt.GenerateFromPassword([]byte(sPassword), 14)
	return string(bPass), err
}

func CheckPasswordHash(sPassword, sHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(sHash), []byte(sPassword))
	return err == nil
}

func ValidUser(sUser, sPassword string, slcUsers []config.AppUser) bool {
	for _, v := range slcUsers {
		if sUser == v.User && CheckPasswordHash(sPassword, v.Password) {
			return true
		}
	}
	return false
}

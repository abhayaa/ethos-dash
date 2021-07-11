package main

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/keygen"
)

func main() {

	username := "abhaya#1149"

	user := db.User{
		UserId:   username,
		EthosKey: keygen.Keygen(username),
	}

	db.InsertUser(user)
}

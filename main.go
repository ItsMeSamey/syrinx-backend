package main

import (
	"fmt"
	
	"ccs.ctf/DB"
)

func main() {
	db, err := DB.OpenDb("2024.ctf")
	defer db.Close()

	newUser := DB.User{
		Username: "johnDoe",
		Password: "PasswordHash",
	}
	DB.CreateUser(db, newUser)

	username := "johnDoe"
	password := "PasswordHash"
	authenticatedUser, err := DB.Authenticate(db, username, password)
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}
	fmt.Printf("User authenticated: %s\n", authenticatedUser.Username)
}



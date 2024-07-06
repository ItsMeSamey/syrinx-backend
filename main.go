package main

import (
	"fmt"
	"log"

	"ccs.ctf/DB"

	// "github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := DB.InitDB("mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
	}

	// Create a new user
	err = DB.CreateUser(&DB.User{"a", "mail", "pass", "team id", "discord",})
	if err != nil {
		fmt.Println("Error creating user:", err)
	}

	user, err := DB.UserAuthenticate("a", "pass")
	if err != nil {
		fmt.Println("Auth error", err)
	} else{
		fmt.Printf("TeamId: %d, Discord: %s\n", user.TeamID, user.DiscordID)
	}

}


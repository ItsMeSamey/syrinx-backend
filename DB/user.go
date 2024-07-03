package DB

import (
	"encoding/json"
	"fmt"

	// "github.com/dgrijalva/jwt-go"
	"github.com/boltdb/bolt"
)

// User struct to store user information
type User struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	TeamID   int `json:"tID"`
}

const (
	userBucket = "users"
	teamBucket = "teams"
)

func CreateUser(db *bolt.DB, user User) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := b.Put([]byte(user.Username), data); err != nil {
			return err
		}
		return nil
	})
}

func Authenticate(db *bolt.DB, username, password string) (*User, error) {
	var user User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		val := b.Get([]byte(username))
		if val == nil {
			return fmt.Errorf("User not found")
		}
		return json.Unmarshal(val, &user)
	})

	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("Invalid Password")
	}

	return &user, nil
}


package DB

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	// Used for genearting JWT tokens
	secretKey = "3ut3slEsuGvWxEh4R/OdW+hpSmVrOD8gHMAxMlQXo5CfQhmZvmaH+npQgCb4LFkfW0r9zxNYlrsaY5w/dwsYKw=="
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

func Authenticate(db *bolt.DB, username, password string) (*User, string, error) {
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
		return nil, "", err
	}

	if user.Password != password {
		return nil, "", fmt.Errorf("Invalid Password")
	}

	// Generate JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(), // Token expiry time
	})
	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}


package DB

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// User struct to store user information
type User struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	TeamID   int    `json:"teamID"`
	SessionID string `json:"sesisonID"`
}

const (
	userBucket = "users"
	teamBucket = "teams"
)

func CreateUser(user User) error {
	return DBInstance.Update(func(tx *bolt.Tx) error {
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

func Authenticate(username, password string) (*User, string, error) {
	var user User
	err := DBInstance.View(func(tx *bolt.Tx) error {
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
	
	id, err := genSessionID(&user)
	if err != nil {
		return nil, "", fmt.Errorf("SessionID Creation Failed")
	}
	return &user, id, nil
}

func genSessionID(user *User) (string, error) {
	// TODO: when user reauthanticates, old one should be deleted and a new token must be generated
	_ = user
	return "", nil
}

func GetUserFromSessionID(sessionID string) *User {
	// TODO: lookup the sesisonID table
	return nil
}

func GetSessionIDFromUser(user *User) {
	// TODO: Implement creation of user's sessionID 
	// NOTE: One user must have only 1 session ID,
}

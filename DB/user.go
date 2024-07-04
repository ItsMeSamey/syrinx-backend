package DB

import (
	"fmt"
	"errors"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

const (
	userBucket = "users"
	teamBucket = "teams"
)

// User struct to store user information
type User struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	UserId string	  `json:"userID"`
	TeamID   int    `json:"teamID"`
	SessionID string `json:"sesisonID"`
}

func setSessionID(user *User) (error) {
	// Set the user's ssession id to a <unique> and random base64 encoded string
	// also make a bucket to hold session keys and respective user names
	// TODO: when user reauthanticates, old one should be deleted and a new token must be generated
	_ = user
	return nil
}

func UserExists(username string) (bool, error) {
	// Implement this
	return true, nil
}

func CreateUser(user *User) error {
	// Return error if user is already present
	exists, err := UserExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("CreateUser: User Exists")
	}
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

func Authenticate(username, password string) (*User, error) {
	var user User
	err := DBInstance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		val := b.Get([]byte(username))
		if val == nil {
			return fmt.Errorf("Authenticate: User not found")
		}
		return json.Unmarshal(val, &user)
	})

	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("Authenticate: Invalid Password")
	}
	
	err = setSessionID(&user)
	if err != nil {
		return nil, fmt.Errorf("Authenticate: SessionID Creation Failed")
	}
	return &user, nil
}

func GetUserFromSessionID(sessionID string) (*User, error) {
	// TODO: lookup the sesisonID table
	return nil, nil
}

func GetSessionIDFromUser(user *User) (string, error) {
	// TODO: Implement creation of user's sessionID 
	// NOTE: One user must have only 1 session ID,
	return "", nil
}


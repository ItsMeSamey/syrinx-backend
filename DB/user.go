package DB

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

const (
	userBucket = "users"
	teamBucket = "teams"
	sessionBucket = "sessions"
)

// User struct to store user information
type User struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	UserId string	  `json:"userID"`
	TeamID   int    `json:"teamID"`
	SessionID string `json:"sesisonID"`
}

func genSessionID() (string, error) {
	bytes := make([]byte, 6*16)
	if _, err := rand.Read(bytes); 
	err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func addToBucket(tx *bolt.Tx, bucket string, key []byte, val []byte) error {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return errors.New("addToBucket: bucket is nil")
	}
	if err := b.Put(key, val); err != nil {
		return err
	}
	return nil
}
func getFromBucket(tx *bolt.Tx, bucket string, key []byte) ([]byte, error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, errors.New("getFromBucket: bucket is nil")
	}
	val := b.Get(key)
	if val == nil {
		return nil, fmt.Errorf("getFromBucket: key not found")
	}
	return val, nil
}


func (user *User) setSessionID() error {
	sessionID, err := genSessionID()
	if err != nil {
		return err
	}
	// Set the user's ssession id to a <unique> and random base64 encoded string
	// also make a bucket to hold session keys and respective user names
	// TODO: when user reauthanticates, old one should be deleted and a new token must be generated
	return DBInstance.Update(func(tx *bolt.Tx) error {
		return addToBucket(tx, sessionBucket, []byte(sessionID), []byte(user.Username))
	})
}

func deleteExistingSession(tx *bolt.Tx, userID string) error {
	return errors.New("Not Implemented")
}

func UserExists(username string) (bool, error) {
	isPresent := true
	return isPresent, DBInstance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		if b == nil {
			return errors.New("getFromBucket: bucket is nil")
		}
		val := b.Get([]byte(username))
		isPresent = val == nil
		return nil
	})
}

func (user *User) Create() error {
	// Return error if user is already present
	exists, err := UserExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("CreateUser: User Exists")
	}
	return DBInstance.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return addToBucket(tx, userBucket, []byte(user.Username), data)
	})
}

func (user *User) Authenticate() error {
	var tempUser User
	err := DBInstance.View(func(tx *bolt.Tx) error {
		val, err := getFromBucket(tx, userBucket, []byte(user.Username))
		if err != nil {
			return err
		}
		return json.Unmarshal(val, &tempUser)
	})
	if err != nil {
		return err
	}

	if tempUser.Password != user.Password {
		return fmt.Errorf("Authenticate: Invalid Password")
	}
	user = &tempUser
	return nil
}

func (user *User) GetUserFromSessionID() error {
	if user.SessionID == "" {
		return errors.New("GetUserFromSessionID: SessionID not given")
	}
	return DBInstance.View(func(tx *bolt.Tx) error {
		val, err := getFromBucket(tx, sessionBucket, []byte(user.SessionID))
		if err != nil {
			return err
		}
		return json.Unmarshal(val, &user)
	})
}



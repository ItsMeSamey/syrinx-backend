package DB

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	userBucket = "users"
	teamBucket = "teams"
	sessionBucket = "sessions"
)

// User struct to store user information
type User struct {
	Username  string `json:"user"`
	Password  string `json:"pass"`
	UserId    int    `json:"userID"`
	TeamID    int    `json:"teamID"`
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

func (user *User) Create() error {
	// Return error if user is already present
	exists, err := UserDB.DoesExist(userBucket, []byte(user.Username))
	if err != nil {
		return err
	}
	if exists {
		return errors.New("User.Create: User Exists")
	}

	tries := 0
start:
	sessionID, err := genSessionID()
	if err != nil {
		return err
	}

	exists, err = UserDB.DoesExist(sessionBucket, []byte(sessionID))
	if err != nil {
		return err
	}
	if exists {
		tries += 1
		if tries > 1024*1024 {
			return errors.New("User.Create: 418 I'm A Teapot")
		}
		goto start
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return UserDB.addToBucket(userBucket, []byte(user.Username), data)
}

func UserAuthenticate(username, password string) (*User, error) {
	var user User
	val, err := UserDB.getFromBucket(userBucket, []byte(username))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &user)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("UserAuthenticate: Invalid Password")
	}
	return &user, nil
}

func UserFromSessionID(SessionID []byte) (*User, error) {
	val, err := UserDB.getFromBucket(sessionBucket, SessionID)
	if err != nil {
		return nil, err
	}
	var user User
	return &user, json.Unmarshal(val, &user)
}

func (user *User) UserInLobby(lobbyID int) (bool, error) {
	// Implement
	return false, nil
}
